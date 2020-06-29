package migrator

import "database/sql"

type executableSQL interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// Migration represents migration entity
//
// Name 		should be a unique name to specify migration. It is up to you to choose the name you like
// Up() 		should return Schema with prepared commands to be migrated
// Down()		should return Schema with prepared commands to be reverted
// Transaction	optinal flag to enable transaction for migration
//
// Example:
//		var migration = migrator.Migration{
//			Name: "19700101_0001_create_posts_table",
//			Up: func() migrator.Schema {
//				var s migrator.Schema
//				posts := migrator.Table{Name: "posts"}
//
//				posts.UniqueID("id")
//				posts.Column("title", migrator.String{Precision: 64})
//				posts.Column("content", migrator.Text{})
//				posts.Timestamps()
//
//				s.CreateTable(posts)
//
//				return s
//			},
//			Down: func() migrator.Schema {
//				var s migrator.Schema
//
//				s.DropTableIfExists("posts")
//
//				return s
//			},
//		}
type Migration struct {
	Name        string
	Up          func() Schema
	Down        func() Schema
	Transaction bool
}

func (m Migration) exec(db *sql.DB, commands ...command) error {
	if m.Transaction {
		return runInTransaction(db, commands...)
	}

	return run(db, commands...)
}

func runInTransaction(db *sql.DB, commands ...command) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	err = run(tx, commands...)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func run(db executableSQL, commands ...command) error {
	for _, command := range commands {
		sql := command.toSQL()
		if sql == "" {
			return ErrNoSQLCommandsToRun
		}
		if _, err := db.Exec(sql); err != nil {
			return err
		}
	}

	return nil
}
