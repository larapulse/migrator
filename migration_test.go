package migrator

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var (
	errTestDBExecFailed        = errors.New("DB exec command failed")
	errTestDBQueryFailed       = errors.New("DB query command failed")
	errTestDBTransactionFailed = errors.New("DB transaction failed")
	errTestLastInsertID        = errors.New("Failed to get last insert ID")
	errTestAffectedRows        = errors.New("Failed to amount of affected rows")
)

type testDummyCommand string

func (c testDummyCommand) ToSQL() string {
	return string(c)
}

func testDBConnection(t *testing.T) (db *sql.DB, mock sqlmock.Sqlmock, resetDB func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	resetDB = func() {
		defer db.Close()

		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}

	return
}

func TestMigrationExec(t *testing.T) {
	t.Run("it executes migration in transaction", func(t *testing.T) {
		m := Migration{Transaction: true}

		db, mock, resetDB := testDBConnection(t)
		defer resetDB()

		commands := []Command{
			testCommand("test"),
			testDummyCommand("test"),
		}

		mock.ExpectBegin()
		mock.ExpectExec(commands[0].ToSQL()).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(commands[1].ToSQL()).WillReturnResult(sqlmock.NewResult(2, 1))
		mock.ExpectCommit()

		// now we execute our method
		if err := m.exec(db, nil, commands...); err != nil {
			t.Errorf("error was not expected while running query: %s", err)
		}
	})

	t.Run("it executes general transaction", func(t *testing.T) {
		m := Migration{}
		db, mock, resetDB := testDBConnection(t)
		defer resetDB()

		commands := []Command{
			testCommand("test"),
			testDummyCommand("test"),
		}
		mock.ExpectExec(commands[0].ToSQL()).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(commands[1].ToSQL()).WillReturnResult(sqlmock.NewResult(2, 1))

		// now we execute our method
		if err := m.exec(db, nil, commands...); err != nil {
			t.Errorf("error was not expected while running query: %s", err)
		}
	})
}

func TestRunInTransaction(t *testing.T) {
	t.Run("it returns an error if transaction wasn't started", func(t *testing.T) {
		db, mock, resetDB := testDBConnection(t)
		defer resetDB()

		commands := []Command{}
		want := sqlmock.ErrCancelled
		mock.ExpectBegin().WillReturnError(want)

		// now we execute our method
		got := runInTransaction(db, nil, commands...)
		assert.Equal(t, want, got)
	})

	t.Run("it rolled back transaction in case of error", func(t *testing.T) {
		db, mock, resetDB := testDBConnection(t)
		defer resetDB()

		commands := []Command{testDummyCommand("run")}
		want := sqlmock.ErrCancelled

		mock.ExpectBegin()
		mock.ExpectExec(commands[0].ToSQL()).WillReturnError(want)
		mock.ExpectRollback()

		// now we execute our method
		got := runInTransaction(db, nil, commands...)
		assert.Equal(t, want, got)
	})

	t.Run("it returns an error if committing transaction was unsuccessful", func(t *testing.T) {
		db, mock, resetDB := testDBConnection(t)
		defer resetDB()

		commands := []Command{}
		want := sqlmock.ErrCancelled

		mock.ExpectBegin()
		mock.ExpectCommit().WillReturnError(want)

		// now we execute our method
		got := runInTransaction(db, nil, commands...)
		assert.Equal(t, want, got)
	})

	t.Run("it executes all commands", func(t *testing.T) {
		db, mock, resetDB := testDBConnection(t)
		defer resetDB()

		commands := []Command{
			testCommand("test"),
			testDummyCommand("test"),
		}

		mock.ExpectBegin()
		mock.ExpectExec(commands[0].ToSQL()).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(commands[1].ToSQL()).WillReturnResult(sqlmock.NewResult(2, 1))
		mock.ExpectCommit()

		// now we execute our method
		if err := runInTransaction(db, nil, commands...); err != nil {
			t.Errorf("error was not expected while running query: %s", err)
		}
	})
}

func TestRun(t *testing.T) {
	t.Run("it returns an error on invalid command", func(t *testing.T) {
		db, mock, resetDB := testDBConnection(t)
		defer resetDB()

		commands := []Command{
			testCommand("test"),
			testDummyCommand(""),
		}

		mock.ExpectExec(commands[0].ToSQL()).WillReturnResult(sqlmock.NewResult(1, 1))

		err := run(db, nil, commands...)

		assert.Error(t, err)
		assert.Equal(t, ErrNoSQLCommandsToRun, err)
	})

	t.Run("it returns an error on DB command execution", func(t *testing.T) {
		db, mock, resetDB := testDBConnection(t)
		defer resetDB()

		commands := []Command{
			testCommand("test"),
			testDummyCommand("dead"),
		}

		mock.ExpectExec(commands[0].ToSQL()).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(commands[1].ToSQL()).WillReturnError(errTestDBExecFailed)

		err := run(db, nil, commands...)

		assert.Error(t, err)
		assert.Equal(t, errTestDBExecFailed, err)
	})

	t.Run("it executes all commands", func(t *testing.T) {
		db, mock, resetDB := testDBConnection(t)
		defer resetDB()

		commands := []Command{
			testCommand("test"),
			testDummyCommand("test"),
		}

		mock.ExpectExec(commands[0].ToSQL()).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(commands[1].ToSQL()).WillReturnResult(sqlmock.NewResult(2, 1))

		err := run(db, nil, commands...)

		assert.Nil(t, err)
	})
}
