package migrator

import (
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestMigrate(t *testing.T) {
	t.Run("it fails when migration pool is empty", func(t *testing.T) {
		m := Migrator{}
		db, _, resetDB := testDBConnection(t)
		defer resetDB()

		migrated, err := m.Migrate(db)

		assert.Len(t, migrated, 0)
		assert.Error(t, err)
		assert.Equal(t, ErrNoMigrationDefined, err)
	})

	t.Run("it fails when there is invalid item in the migration pool", func(t *testing.T) {
		migration := Migration{}
		m := Migrator{Pool: []Migration{migration}}
		db, _, resetDB := testDBConnection(t)
		defer resetDB()

		migrated, err := m.Migrate(db)

		assert.Len(t, migrated, 0)
		assert.Error(t, err)
		assert.Equal(t, ErrMissingMigrationName, err)
	})

	t.Run("it fails when migration table creation failed", func(t *testing.T) {
		migration := Migration{Name: "test"}
		m := Migrator{Pool: []Migration{migration}}
		db, mock, resetDB := testDBConnection(t)
		defer resetDB()

		mock.ExpectQuery("SELECT").WillReturnRows().WillReturnError(errTestDBQueryFailed)
		mock.ExpectExec("CREATE").WillReturnError(errTestDBExecFailed)

		migrated, err := m.Migrate(db)

		assert.Len(t, migrated, 0)
		assert.Error(t, err)
		assert.Equal(t, fmt.Errorf("Migration table failed to be created: %v", errTestDBExecFailed), err)
	})

	t.Run("it fails while fetching executed list", func(t *testing.T) {
		migration := Migration{Name: "test"}
		m := Migrator{Pool: []Migration{migration}}
		db, mock, resetDB := testDBConnection(t)
		defer resetDB()

		mock.ExpectQuery("SELECT").WillReturnRows()
		mock.ExpectQuery("SELECT id, name, batch, applied_at FROM migrations").WillReturnError(errTestDBExecFailed)

		migrated, err := m.Migrate(db)

		assert.Len(t, migrated, 0)
		assert.Error(t, err)
		assert.Equal(t, errTestDBExecFailed, err)
	})

	t.Run("it skips execution when it was already executed", func(t *testing.T) {
		migration := Migration{Name: "test"}
		m := Migrator{Pool: []Migration{migration}}
		db, mock, resetDB := testDBConnection(t)
		defer resetDB()

		rows := sqlmock.NewRows([]string{"id", "name", "batch", "applied_at"}).AddRow(1, "test", 1, time.Now())

		mock.ExpectQuery("SELECT").WillReturnRows()
		mock.ExpectQuery("SELECT id, name, batch, applied_at FROM migrations").WillReturnRows(rows)

		migrated, err := m.Migrate(db)

		assert.Len(t, migrated, 0)
		assert.Nil(t, err)
	})

	t.Run("it fails executing empty list of migrations", func(t *testing.T) {
		migration := Migration{Name: "test", Up: func() Schema {
			var s Schema
			return s
		}}
		m := Migrator{Pool: []Migration{migration}}
		db, mock, resetDB := testDBConnection(t)
		defer resetDB()

		rows := sqlmock.NewRows([]string{"id", "name", "batch", "applied_at"}).AddRow(1, "new", 1, time.Now())

		mock.ExpectQuery("SELECT").WillReturnRows()
		mock.ExpectQuery("SELECT id, name, batch, applied_at FROM migrations").WillReturnRows(rows)

		migrated, err := m.Migrate(db)

		assert.Len(t, migrated, 0)
		assert.Error(t, err)
		assert.Equal(t, ErrNoSQLCommandsToRun, err)
	})

	t.Run("it fails executing migration commands", func(t *testing.T) {
		migration := Migration{Name: "test", Up: func() Schema {
			var s Schema
			s.pool = append(s.pool, testDummyCommand(""))
			return s
		}}
		m := Migrator{Pool: []Migration{migration}}
		db, mock, resetDB := testDBConnection(t)
		defer resetDB()

		rows := sqlmock.NewRows([]string{"id", "name", "batch", "applied_at"}).AddRow(1, "new", 1, time.Now())

		mock.ExpectQuery("SELECT").WillReturnRows()
		mock.ExpectQuery("SELECT id, name, batch, applied_at FROM migrations").WillReturnRows(rows)

		migrated, err := m.Migrate(db)

		assert.Len(t, migrated, 0)
		assert.Error(t, err)
		assert.Equal(t, ErrNoSQLCommandsToRun, err)
	})

	t.Run("it fails while storing executed migration info", func(t *testing.T) {
		migration := Migration{Name: "test", Up: func() Schema {
			var s Schema
			s.DropTable("test", false, "")
			return s
		}}
		m := Migrator{Pool: []Migration{migration}}
		db, mock, resetDB := testDBConnection(t)
		defer resetDB()

		rows := sqlmock.NewRows([]string{"id", "name", "batch", "applied_at"}).AddRow(1, "new", 1, time.Now())

		mock.ExpectQuery("SELECT").WillReturnRows()
		mock.ExpectQuery("SELECT id, name, batch, applied_at FROM migrations").WillReturnRows(rows)
		mock.ExpectExec("DROP").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT").WillReturnError(errTestDBExecFailed)

		migrated, err := m.Migrate(db)

		assert.Len(t, migrated, 0)
		assert.Error(t, err)
		assert.Equal(t, errTestDBExecFailed, err)
	})

	t.Run("it executes migrations and returns list of migrated items", func(t *testing.T) {
		migration := Migration{Name: "test", Up: func() Schema {
			var s Schema
			s.DropTable("test", false, "")
			return s
		}}
		m := Migrator{Pool: []Migration{migration}}
		db, mock, resetDB := testDBConnection(t)
		defer resetDB()

		rows := sqlmock.NewRows([]string{"id", "name", "batch", "applied_at"}).AddRow(1, "new", 4, time.Now())

		mock.ExpectQuery("SELECT").WillReturnRows()
		mock.ExpectQuery("SELECT id, name, batch, applied_at FROM migrations").WillReturnRows(rows)
		mock.ExpectExec("DROP").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`INSERT .* VALUES \("test", 5\)`).WillReturnResult(sqlmock.NewResult(1, 1))

		migrated, err := m.Migrate(db)

		assert.Len(t, migrated, 1)
		assert.Equal(t, migrated[0], "test")
		assert.Nil(t, err)
	})
}

func TestRollback(t *testing.T) {
	t.Run("it fails when migration pool is empty", func(t *testing.T) {
		m := Migrator{}
		db, _, resetDB := testDBConnection(t)
		defer resetDB()

		reverted, err := m.Rollback(db)

		assert.Len(t, reverted, 0)
		assert.Error(t, err)
		assert.Equal(t, ErrNoMigrationDefined, err)
	})

	t.Run("it fails when there is invalid item in the migration pool", func(t *testing.T) {
		migration := Migration{}
		m := Migrator{Pool: []Migration{migration}}
		db, _, resetDB := testDBConnection(t)
		defer resetDB()

		reverted, err := m.Rollback(db)

		assert.Len(t, reverted, 0)
		assert.Error(t, err)
		assert.Equal(t, ErrMissingMigrationName, err)
	})

	t.Run("it fails when migration table missing", func(t *testing.T) {
		migration := Migration{Name: "test"}
		m := Migrator{Pool: []Migration{migration}}
		db, mock, resetDB := testDBConnection(t)
		defer resetDB()

		mock.ExpectQuery("SELECT").WillReturnRows().WillReturnError(errTestDBQueryFailed)

		reverted, err := m.Rollback(db)

		assert.Len(t, reverted, 0)
		assert.Error(t, err)
		assert.Equal(t, ErrTableNotExists, err)
	})

	t.Run("it fails while fetching executed list", func(t *testing.T) {
		migration := Migration{Name: "test"}
		m := Migrator{Pool: []Migration{migration}}
		db, mock, resetDB := testDBConnection(t)
		defer resetDB()

		mock.ExpectQuery("SELECT").WillReturnRows()
		mock.ExpectQuery("SELECT id, name, batch, applied_at FROM migrations").WillReturnError(errTestDBExecFailed)

		reverted, err := m.Rollback(db)

		assert.Len(t, reverted, 0)
		assert.Error(t, err)
		assert.Equal(t, errTestDBExecFailed, err)
	})

	t.Run("it exits when executed list is empty", func(t *testing.T) {
		migration := Migration{Name: "test"}
		m := Migrator{Pool: []Migration{migration}}
		db, mock, resetDB := testDBConnection(t)
		defer resetDB()

		mock.ExpectQuery("SELECT").WillReturnRows()
		mock.ExpectQuery("SELECT id, name, batch, applied_at FROM migrations").WillReturnRows(sqlmock.NewRows([]string{}))

		reverted, err := m.Rollback(db)

		assert.Len(t, reverted, 0)
		assert.Error(t, err)
		assert.Equal(t, ErrEmptyRollbackStack, err)
	})

	t.Run("it does nothing when executed migration not in the migration pool", func(t *testing.T) {
		migration := Migration{Name: "test", Down: func() Schema {
			var s Schema
			s.pool = append(s.pool, testDummyCommand(""))
			return s
		}}
		m := Migrator{Pool: []Migration{migration}}
		db, mock, resetDB := testDBConnection(t)
		defer resetDB()

		rows := sqlmock.NewRows([]string{"id", "name", "batch", "applied_at"}).AddRow(1, "new", 1, time.Now())

		mock.ExpectQuery("SELECT").WillReturnRows()
		mock.ExpectQuery("SELECT id, name, batch, applied_at FROM migrations").WillReturnRows(rows)

		reverted, err := m.Rollback(db)

		assert.Len(t, reverted, 0)
		assert.Nil(t, err)
	})

	t.Run("it fails executing empty list of commands", func(t *testing.T) {
		migration := Migration{Name: "test", Down: func() Schema {
			var s Schema
			return s
		}}
		m := Migrator{Pool: []Migration{migration}}
		db, mock, resetDB := testDBConnection(t)
		defer resetDB()

		rows := sqlmock.NewRows([]string{"id", "name", "batch", "applied_at"}).AddRow(1, "test", 1, time.Now())

		mock.ExpectQuery("SELECT").WillReturnRows()
		mock.ExpectQuery("SELECT id, name, batch, applied_at FROM migrations").WillReturnRows(rows)

		reverted, err := m.Rollback(db)

		assert.Len(t, reverted, 0)
		assert.Error(t, err)
		assert.Equal(t, ErrNoSQLCommandsToRun, err)
	})

	t.Run("it fails executing migration commands", func(t *testing.T) {
		migration := Migration{Name: "test", Down: func() Schema {
			var s Schema
			s.pool = append(s.pool, testDummyCommand(""))
			return s
		}}
		m := Migrator{Pool: []Migration{migration}}
		db, mock, resetDB := testDBConnection(t)
		defer resetDB()

		rows := sqlmock.NewRows([]string{"id", "name", "batch", "applied_at"}).AddRow(1, "test", 1, time.Now())

		mock.ExpectQuery("SELECT").WillReturnRows()
		mock.ExpectQuery("SELECT id, name, batch, applied_at FROM migrations").WillReturnRows(rows)

		reverted, err := m.Rollback(db)

		assert.Len(t, reverted, 0)
		assert.Error(t, err)
		assert.Equal(t, ErrNoSQLCommandsToRun, err)
	})

	t.Run("it fails while removing executed migration info", func(t *testing.T) {
		migration := Migration{Name: "test", Down: func() Schema {
			var s Schema
			s.DropTable("test", false, "")
			return s
		}}
		m := Migrator{Pool: []Migration{migration}}
		db, mock, resetDB := testDBConnection(t)
		defer resetDB()

		rows := sqlmock.NewRows([]string{"id", "name", "batch", "applied_at"}).AddRow(1, "test", 1, time.Now())

		mock.ExpectQuery("SELECT").WillReturnRows()
		mock.ExpectQuery("SELECT id, name, batch, applied_at FROM migrations").WillReturnRows(rows)
		mock.ExpectExec("DROP").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("DELETE").WillReturnError(errTestDBExecFailed)

		reverted, err := m.Rollback(db)

		assert.Len(t, reverted, 0)
		assert.Error(t, err)
		assert.Equal(t, errTestDBExecFailed, err)
	})

	t.Run("it roll back migrations and returns list of reverted items", func(t *testing.T) {
		migration := Migration{Name: "test", Down: func() Schema {
			var s Schema
			s.DropTable("test", false, "")
			return s
		}}
		m := Migrator{Pool: []Migration{migration, {Name: "new"}}}
		db, mock, resetDB := testDBConnection(t)
		defer resetDB()

		rows := sqlmock.NewRows([]string{"id", "name", "batch", "applied_at"}).
			AddRow(1, "test", 4, time.Now()).
			AddRow(2, "new", 3, time.Now())

		mock.ExpectQuery("SELECT").WillReturnRows()
		mock.ExpectQuery("SELECT id, name, batch, applied_at FROM migrations").WillReturnRows(rows)
		mock.ExpectExec("DROP").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("DELETE FROM migrations WHERE id = ?").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

		reverted, err := m.Rollback(db)

		assert.Len(t, reverted, 1)
		assert.Equal(t, reverted[0], "test")
		assert.Nil(t, err)
	})
}

func TestRevert(t *testing.T) {
	t.Run("it fails when migration pool is empty", func(t *testing.T) {
		m := Migrator{}
		db, _, resetDB := testDBConnection(t)
		defer resetDB()

		reverted, err := m.Revert(db)

		assert.Len(t, reverted, 0)
		assert.Error(t, err)
		assert.Equal(t, ErrNoMigrationDefined, err)
	})

	t.Run("it fails when there is invalid item in the migration pool", func(t *testing.T) {
		migration := Migration{}
		m := Migrator{Pool: []Migration{migration}}
		db, _, resetDB := testDBConnection(t)
		defer resetDB()

		reverted, err := m.Revert(db)

		assert.Len(t, reverted, 0)
		assert.Error(t, err)
		assert.Equal(t, ErrMissingMigrationName, err)
	})

	t.Run("it fails when migration table missing", func(t *testing.T) {
		migration := Migration{Name: "test"}
		m := Migrator{Pool: []Migration{migration}}
		db, mock, resetDB := testDBConnection(t)
		defer resetDB()

		mock.ExpectQuery("SELECT").WillReturnRows().WillReturnError(errTestDBQueryFailed)

		reverted, err := m.Revert(db)

		assert.Len(t, reverted, 0)
		assert.Error(t, err)
		assert.Equal(t, ErrTableNotExists, err)
	})

	t.Run("it fails while fetching executed list", func(t *testing.T) {
		migration := Migration{Name: "test"}
		m := Migrator{Pool: []Migration{migration}}
		db, mock, resetDB := testDBConnection(t)
		defer resetDB()

		mock.ExpectQuery("SELECT").WillReturnRows()
		mock.ExpectQuery("SELECT id, name, batch, applied_at FROM migrations").WillReturnError(errTestDBExecFailed)

		reverted, err := m.Revert(db)

		assert.Len(t, reverted, 0)
		assert.Error(t, err)
		assert.Equal(t, errTestDBExecFailed, err)
	})

	t.Run("it exits when executed list is empty", func(t *testing.T) {
		migration := Migration{Name: "test"}
		m := Migrator{Pool: []Migration{migration}}
		db, mock, resetDB := testDBConnection(t)
		defer resetDB()

		mock.ExpectQuery("SELECT").WillReturnRows()
		mock.ExpectQuery("SELECT id, name, batch, applied_at FROM migrations").WillReturnRows(sqlmock.NewRows([]string{}))

		reverted, err := m.Revert(db)

		assert.Len(t, reverted, 0)
		assert.Error(t, err)
		assert.Equal(t, ErrEmptyRollbackStack, err)
	})

	t.Run("it does nothing when executed migration not in the migration pool", func(t *testing.T) {
		migration := Migration{Name: "test", Down: func() Schema {
			var s Schema
			s.pool = append(s.pool, testDummyCommand(""))
			return s
		}}
		m := Migrator{Pool: []Migration{migration}}
		db, mock, resetDB := testDBConnection(t)
		defer resetDB()

		rows := sqlmock.NewRows([]string{"id", "name", "batch", "applied_at"}).AddRow(1, "new", 1, time.Now())

		mock.ExpectQuery("SELECT").WillReturnRows()
		mock.ExpectQuery("SELECT id, name, batch, applied_at FROM migrations").WillReturnRows(rows)

		reverted, err := m.Revert(db)

		assert.Len(t, reverted, 0)
		assert.Nil(t, err)
	})

	t.Run("it fails executing empty list of commands", func(t *testing.T) {
		migration := Migration{Name: "test", Down: func() Schema {
			var s Schema
			return s
		}}
		m := Migrator{Pool: []Migration{migration}}
		db, mock, resetDB := testDBConnection(t)
		defer resetDB()

		rows := sqlmock.NewRows([]string{"id", "name", "batch", "applied_at"}).AddRow(1, "test", 1, time.Now())

		mock.ExpectQuery("SELECT").WillReturnRows()
		mock.ExpectQuery("SELECT id, name, batch, applied_at FROM migrations").WillReturnRows(rows)

		reverted, err := m.Revert(db)

		assert.Len(t, reverted, 0)
		assert.Error(t, err)
		assert.Equal(t, ErrNoSQLCommandsToRun, err)
	})

	t.Run("it fails executing migration commands", func(t *testing.T) {
		migration := Migration{Name: "test", Down: func() Schema {
			var s Schema
			s.pool = append(s.pool, testDummyCommand(""))
			return s
		}}
		m := Migrator{Pool: []Migration{migration}}
		db, mock, resetDB := testDBConnection(t)
		defer resetDB()

		rows := sqlmock.NewRows([]string{"id", "name", "batch", "applied_at"}).AddRow(1, "test", 1, time.Now())

		mock.ExpectQuery("SELECT").WillReturnRows()
		mock.ExpectQuery("SELECT id, name, batch, applied_at FROM migrations").WillReturnRows(rows)

		reverted, err := m.Revert(db)

		assert.Len(t, reverted, 0)
		assert.Error(t, err)
		assert.Equal(t, ErrNoSQLCommandsToRun, err)
	})

	t.Run("it fails while removing executed migration info", func(t *testing.T) {
		migration := Migration{Name: "test", Down: func() Schema {
			var s Schema
			s.DropTable("test", false, "")
			return s
		}}
		m := Migrator{Pool: []Migration{migration}}
		db, mock, resetDB := testDBConnection(t)
		defer resetDB()

		rows := sqlmock.NewRows([]string{"id", "name", "batch", "applied_at"}).AddRow(1, "test", 1, time.Now())

		mock.ExpectQuery("SELECT").WillReturnRows()
		mock.ExpectQuery("SELECT id, name, batch, applied_at FROM migrations").WillReturnRows(rows)
		mock.ExpectExec("DROP").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("DELETE").WillReturnError(errTestDBExecFailed)

		reverted, err := m.Revert(db)

		assert.Len(t, reverted, 0)
		assert.Error(t, err)
		assert.Equal(t, errTestDBExecFailed, err)
	})

	t.Run("it roll back migrations and returns list of reverted items", func(t *testing.T) {
		m := Migrator{Pool: []Migration{
			{Name: "test", Down: func() Schema {
				var s Schema
				s.DropTable("test", false, "")
				return s
			}},
			{Name: "new", Down: func() Schema {
				var s Schema
				s.DropTable("test", false, "")
				return s
			}},
		}}
		db, mock, resetDB := testDBConnection(t)
		defer resetDB()

		rows := sqlmock.NewRows([]string{"id", "name", "batch", "applied_at"}).
			AddRow(1, "test", 4, time.Now()).
			AddRow(2, "new", 3, time.Now())

		mock.ExpectQuery("SELECT").WillReturnRows()
		mock.ExpectQuery("SELECT id, name, batch, applied_at FROM migrations").WillReturnRows(rows)
		mock.ExpectExec("DROP").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("DELETE FROM migrations WHERE id = ?").WithArgs(2).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("DROP").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("DELETE FROM migrations WHERE id = ?").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

		reverted, err := m.Revert(db)

		assert.Len(t, reverted, 2)
		assert.Equal(t, reverted[0], "new")
		assert.Equal(t, reverted[1], "test")
		assert.Nil(t, err)
	})
}

func TestCheckMigrationPool(t *testing.T) {
	t.Run("it is successful on empty pool", func(t *testing.T) {
		m := Migrator{}
		err := m.checkMigrationPool()

		assert.Nil(t, err)
	})

	t.Run("It is successful for proper pool", func(t *testing.T) {
		m := Migrator{Pool: []Migration{
			{Name: "test"},
			{Name: "random"},
		}}
		err := m.checkMigrationPool()

		assert.Nil(t, err)
	})

	t.Run("it returns an error on missing migration name", func(t *testing.T) {
		m := Migrator{Pool: []Migration{
			{Name: "test"},
			{Name: "random"},
			{Name: ""},
		}}
		err := m.checkMigrationPool()

		assert.Error(t, err)
		assert.Equal(t, ErrMissingMigrationName, err)
	})

	t.Run("it returns an error on duplicated migration name", func(t *testing.T) {
		m := Migrator{Pool: []Migration{
			{Name: "test"},
			{Name: "random"},
			{Name: "again"},
			{Name: "migration"},
			{Name: "again"},
		}}
		err := m.checkMigrationPool()

		assert.NotNil(t, err)
		assert.Equal(t, `Migration "again" is duplicated in the pool`, err.Error())
	})
}

func TestCreateMigrationTable(t *testing.T) {
	t.Run("it ignores creation if table exists", func(t *testing.T) {
		m := Migrator{}
		db, mock, resetDB := testDBConnection(t)
		defer resetDB()

		mock.ExpectQuery(`SELECT \* FROM migrations`).WillReturnRows().WillReturnError(nil)

		err := m.createMigrationTable(db)

		assert.Nil(t, err)
	})

	t.Run("it creates migration table", func(t *testing.T) {
		m := Migrator{}
		db, mock, resetDB := testDBConnection(t)
		defer resetDB()

		mock.ExpectQuery(`SELECT \* FROM migrations`).WillReturnError(errTestDBQueryFailed)
		sql := `CREATE TABLE migrations \(id int\(10\) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY, name varchar\(255\) COLLATE utf8mb4_unicode_ci NOT NULL, batch int\(11\) NOT NULL, applied_at timestamp NULL DEFAULT CURRENT_TIMESTAMP\) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`
		mock.ExpectExec(sql).WillReturnResult(sqlmock.NewResult(1, 1))

		err := m.createMigrationTable(db)

		assert.Nil(t, err)
	})

	t.Run("it fails creating table", func(t *testing.T) {
		m := Migrator{}
		db, mock, resetDB := testDBConnection(t)
		defer resetDB()

		mock.ExpectQuery(`SELECT \* FROM migrations`).WillReturnError(errTestDBQueryFailed)
		sql := `CREATE TABLE migrations \(` +
			`id int\(10\) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY, ` +
			`name varchar\(255\) COLLATE utf8mb4_unicode_ci NOT NULL, ` +
			`batch int\(11\) NOT NULL, applied_at timestamp NULL DEFAULT CURRENT_TIMESTAMP\) ` +
			`ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`
		mock.ExpectExec(sql).WillReturnError(errTestDBExecFailed)

		err := m.createMigrationTable(db)

		assert.Error(t, err)
		assert.Equal(t, errTestDBExecFailed, err)
	})
}

func TestHasTable(t *testing.T) {
	t.Run("it returns true if table exists", func(t *testing.T) {
		m := Migrator{}
		db, mock, resetDB := testDBConnection(t)
		defer resetDB()

		mock.ExpectQuery(`SELECT \* FROM migrations`).WillReturnRows().WillReturnError(nil)
		got := m.hasTable(db)

		assert.Equal(t, true, got)
	})

	t.Run("it returns false if table does not exist", func(t *testing.T) {
		m := Migrator{}
		db, mock, resetDB := testDBConnection(t)
		defer resetDB()

		mock.ExpectQuery(`SELECT \* FROM migrations`).WillReturnError(errTestDBQueryFailed)
		got := m.hasTable(db)

		assert.Equal(t, false, got)
	})
}

func TestMigrationTable(t *testing.T) {
	t.Run("it returns default table name", func(t *testing.T) {
		m := Migrator{}
		got := m.table()

		assert.Equal(t, "migrations", got)
	})

	t.Run("it returns selected table name", func(t *testing.T) {
		m := Migrator{TableName: "table"}
		got := m.table()

		assert.Equal(t, "table", got)
	})
}

func TestBatch(t *testing.T) {
	t.Run("it returns zero on empty executed list", func(t *testing.T) {
		m := Migrator{}
		got := m.batch()

		assert.Equal(t, uint64(0), got)
	})

	t.Run("it returns zero if migration batch is zero", func(t *testing.T) {
		m := Migrator{
			executed: []migrationEntry{
				{batch: uint64(0)},
			},
		}
		got := m.batch()

		assert.Equal(t, uint64(0), got)
	})

	t.Run("it returns the biggest batch from migration list", func(t *testing.T) {
		m := Migrator{
			executed: []migrationEntry{
				{batch: uint64(6)},
				{batch: uint64(3)},
				{batch: uint64(15)},
				{batch: uint64(12)},
			},
		}
		got := m.batch()

		assert.Equal(t, uint64(15), got)
	})
}

func TestPoolExecuted(t *testing.T) {
	t.Run("it fails executing query", func(t *testing.T) {
		m := Migrator{}
		db, mock, resetDB := testDBConnection(t)
		defer resetDB()

		mock.ExpectQuery("SELECT id, name, batch, applied_at FROM migrations").WillReturnError(errTestDBQueryFailed)

		err := m.fetchExecuted(db)

		assert.Error(t, err)
		assert.Equal(t, errTestDBQueryFailed, err)
		assert.Nil(t, m.executed)
	})

	t.Run("it fails scanning row", func(t *testing.T) {
		m := Migrator{}
		db, mock, resetDB := testDBConnection(t)
		defer resetDB()

		rows := sqlmock.NewRows([]string{"id", "name", "batch", "applied_at"}).
			AddRow(1, "first", 1, time.Now()).
			AddRow(2, "second", 1, "test")

		mock.ExpectQuery("SELECT id, name, batch, applied_at FROM migrations").WillReturnRows(rows)

		got := m.fetchExecuted(db)

		assert.Error(t, got)
		assert.NotNil(t, m.executed)
		assert.Len(t, m.executed, 1)
	})

	t.Run("it returns a list of executed migrations", func(t *testing.T) {
		m := Migrator{}
		db, mock, resetDB := testDBConnection(t)
		defer resetDB()

		rows := sqlmock.NewRows([]string{"id", "name", "batch", "applied_at"}).
			AddRow(1, "first", 1, time.Now()).
			AddRow(2, "second", 1, time.Now())

		mock.ExpectQuery("SELECT id, name, batch, applied_at FROM migrations").WillReturnRows(rows)

		err := m.fetchExecuted(db)

		assert.Nil(t, err)
		assert.NotNil(t, m.executed)
		assert.Len(t, m.executed, 2)
	})
}

func TestIsExecuted(t *testing.T) {
	t.Run("it returns false on empty executed list", func(t *testing.T) {
		m := Migrator{}
		got := m.isExecuted("test")

		assert.Equal(t, false, got)
	})

	t.Run("it returns false if migration wasn't executed yet", func(t *testing.T) {
		m := Migrator{
			executed: []migrationEntry{
				{name: "test"},
				{name: "random"},
				{name: "lorem"},
				{name: "ipsum"},
			},
		}
		got := m.isExecuted("")

		assert.Equal(t, false, got)
	})

	t.Run("it returns true if migration was executed", func(t *testing.T) {
		m := Migrator{
			executed: []migrationEntry{
				{name: "test"},
				{name: "random"},
				{name: "lorem"},
				{name: "ipsum"},
			},
		}
		got := m.isExecuted("random")

		assert.Equal(t, true, got)
	})
}

func TestLastExecutedForBatch(t *testing.T) {
	t.Run("it returns an empty list if nothing found for biggest batch", func(t *testing.T) {
		m := Migrator{}
		got := m.lastBatchExecuted()

		assert.Len(t, got, 0)
	})

	t.Run("", func(t *testing.T) {
		m := Migrator{
			executed: []migrationEntry{
				{name: "test", batch: 1},
				{name: "again", batch: 3},
				{name: "random", batch: 2},
				{name: "lorem", batch: 3},
				{name: "ipsum", batch: 3},
			},
		}
		got := m.lastBatchExecuted()

		assert.Len(t, got, 3)
	})
}
