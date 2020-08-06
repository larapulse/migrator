// Package migrator represents MySQL database migrator
//
// MySQL database migrator designed to run migrations to your features and manage database schema update with intuitive go code.
// It is compatible with the latest MySQL v8.
package migrator

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

const migrationTable = "migrations"

var (
	// ErrTableNotExists returns when migration table not found
	ErrTableNotExists = errors.New("Migration table does not exist")

	// ErrNoMigrationDefined returns when no migrations defined in the migrations pool
	ErrNoMigrationDefined = errors.New("No migrations defined")

	// ErrEmptyRollbackStack returns when nothing can be reverted
	ErrEmptyRollbackStack = errors.New("Nothing to rollback, there are no migration executed")

	// ErrMissingMigrationName returns when migration name is missing
	ErrMissingMigrationName = errors.New("Missing migration name")

	// ErrNoSQLCommandsToRun returns when migration is invalid and has no commands in the pool
	ErrNoSQLCommandsToRun = errors.New("There are no commands to be executed")
)

type migrationEntry struct {
	id        uint64
	name      string
	batch     uint64
	appliedAt time.Time
}

// Migrator represents a struct with migrations, that should be executed.
//
// Default migration table name is `migrations`, but it can be re-defined.
// Pool is a list of migrations that should be migrated.
type Migrator struct {
	// Name of the table to track executed migrations
	TableName string
	// stack of migrations
	Pool     []Migration
	executed []migrationEntry
}

// Migrate runs all migrations from pool and stores in migration table executed migration.
func (m Migrator) Migrate(db *sql.DB) (migrated []string, err error) {
	if len(m.Pool) == 0 {
		return migrated, ErrNoMigrationDefined
	}

	if err := m.checkMigrationPool(); err != nil {
		return migrated, err
	}

	if err := m.createMigrationTable(db); err != nil {
		return migrated, fmt.Errorf("Migration table failed to be created: %v", err)
	}

	if err := m.fetchExecuted(db); err != nil {
		return migrated, err
	}

	batch := m.batch() + 1
	table := m.table()

	for _, item := range m.Pool {
		if m.isExecuted(item.Name) {
			continue
		}

		s := item.Up()
		if len(s.pool) == 0 {
			return migrated, ErrNoSQLCommandsToRun
		}
		if err := item.exec(db, s.pool...); err != nil {
			return migrated, err
		}

		entry := migrationEntry{name: item.Name, batch: batch}
		sql := fmt.Sprintf("INSERT INTO `%s` (`name`, `batch`) VALUES (\"%s\", %d)", table, entry.name, entry.batch)

		if _, err := db.Exec(sql); err != nil {
			return migrated, err
		}

		migrated = append(migrated, item.Name)
	}

	return migrated, nil
}

// Rollback reverts last executed batch of migrations.
func (m Migrator) Rollback(db *sql.DB) (reverted []string, err error) {
	if len(m.Pool) == 0 {
		return reverted, ErrNoMigrationDefined
	}

	if err := m.checkMigrationPool(); err != nil {
		return reverted, err
	}

	if !m.hasTable(db) {
		return reverted, ErrTableNotExists
	}

	if err := m.fetchExecuted(db); err != nil {
		return reverted, err
	}

	if len(m.executed) == 0 {
		return reverted, ErrEmptyRollbackStack
	}

	table := m.table()
	revertable := m.lastBatchExecuted()

	for i := len(revertable) - 1; i >= 0; i-- {
		name := revertable[i].name

		for j := len(m.Pool) - 1; j >= 0; j-- {
			item := m.Pool[j]

			if item.Name == name {
				s := item.Down()
				if len(s.pool) == 0 {
					return reverted, ErrNoSQLCommandsToRun
				}
				if err := item.exec(db, s.pool...); err != nil {
					return reverted, err
				}

				if _, err := db.Exec(fmt.Sprintf("DELETE FROM %s WHERE id = ?", table), revertable[i].id); err != nil {
					return reverted, err
				}

				reverted = append(reverted, name)
			}
		}
	}

	return reverted, nil
}

// Revert reverts all executed migration from the pool.
func (m Migrator) Revert(db *sql.DB) (reverted []string, err error) {
	if len(m.Pool) == 0 {
		return reverted, ErrNoMigrationDefined
	}

	if err := m.checkMigrationPool(); err != nil {
		return reverted, err
	}

	if !m.hasTable(db) {
		return reverted, ErrTableNotExists
	}

	if err := m.fetchExecuted(db); err != nil {
		return reverted, err
	}

	if len(m.executed) == 0 {
		return reverted, ErrEmptyRollbackStack
	}

	table := m.table()

	for i := len(m.executed) - 1; i >= 0; i-- {
		name := m.executed[i].name

		for j := len(m.Pool) - 1; j >= 0; j-- {
			item := m.Pool[j]

			if item.Name == name {
				s := item.Down()
				if len(s.pool) == 0 {
					return reverted, ErrNoSQLCommandsToRun
				}
				if err := item.exec(db, s.pool...); err != nil {
					return reverted, err
				}

				if _, err := db.Exec(fmt.Sprintf("DELETE FROM %s WHERE id = ?", table), m.executed[i].id); err != nil {
					return reverted, err
				}

				reverted = append(reverted, name)
			}
		}
	}

	return reverted, nil
}

func (m Migrator) checkMigrationPool() error {
	var names []string

	for _, item := range m.Pool {
		if item.Name == "" {
			return ErrMissingMigrationName
		}

		for _, exist := range names {
			if exist == item.Name {
				return fmt.Errorf(`Migration "%s" is duplicated in the pool`, exist)
			}
		}

		names = append(names, item.Name)
	}

	return nil
}

func (m Migrator) createMigrationTable(db *sql.DB) error {
	if m.hasTable(db) {
		return nil
	}

	sql := fmt.Sprintf(
		"CREATE TABLE %s (%s) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci",
		m.table(),
		strings.Join([]string{
			"id int(10) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY",
			"name varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL",
			"batch int(11) NOT NULL",
			"applied_at timestamp(6) NULL DEFAULT CURRENT_TIMESTAMP(6)",
		}, ", "),
	)

	_, err := db.Exec(sql)

	return err
}

func (m Migrator) hasTable(db *sql.DB) bool {
	_, hasTable := db.Query("SELECT * FROM " + m.table())

	return hasTable == nil
}

func (m Migrator) table() string {
	table := m.TableName
	if table == "" {
		table = migrationTable
	}

	return table
}

func (m Migrator) batch() uint64 {
	var batch uint64

	for _, item := range m.executed {
		if item.batch > batch {
			batch = item.batch
		}
	}

	return batch
}

func (m *Migrator) fetchExecuted(db *sql.DB) error {
	rows, err := db.Query("SELECT id, name, batch, applied_at FROM " + m.table() + " ORDER BY applied_at ASC")
	if err != nil {
		return err
	}
	m.executed = []migrationEntry{}

	for rows.Next() {
		var entry migrationEntry

		if err := rows.Scan(&entry.id, &entry.name, &entry.batch, &entry.appliedAt); err != nil {
			return err
		}

		m.executed = append(m.executed, entry)
	}

	return nil
}

func (m Migrator) isExecuted(name string) bool {
	for _, item := range m.executed {
		if item.name == name {
			return true
		}
	}

	return false
}

func (m Migrator) lastBatchExecuted() []migrationEntry {
	batch := m.batch()
	var result []migrationEntry

	for _, item := range m.executed {
		if item.batch == batch {
			result = append(result, item)
		}
	}

	return result
}
