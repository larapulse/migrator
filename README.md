# MySQL database migrator

<img align="right" width="159px" src="https://github.com/larapulse/migrator/blob/master/logo.png">

[![Build Status](https://travis-ci.org/larapulse/migrator.svg)](https://travis-ci.org/larapulse/migrator)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg)](LICENSE.md)
[![codecov](https://codecov.io/gh/larapulse/migrator/branch/master/graph/badge.svg)](https://codecov.io/gh/larapulse/migrator)
[![Go Report Card](https://goreportcard.com/badge/github.com/larapulse/migrator)](https://goreportcard.com/report/github.com/larapulse/migrator)
[![GoDoc](https://godoc.org/github.com/larapulse/migrator?status.svg)](https://pkg.go.dev/github.com/larapulse/migrator?tab=doc)
[![Release](https://img.shields.io/github/release/larapulse/migrator.svg)](https://github.com/larapulse/migrator/releases)
[![TODOs](https://badgen.net/https/api.tickgit.com/badgen/github.com/larapulse/migrator)](https://www.tickgit.com/browse?repo=github.com/larapulse/migrator)

MySQL database migrator designed to run migrations to your features and manage database schema update with intuitive go code. It is compatible with the latest MySQL v8.

## Installation

To install `migrator` package, you need to install Go and set your Go workspace first.

1. The first need [Go](https://golang.org/) installed (**version 1.13+ is required**), then you can use the below Go command to install `migrator`.

```sh
$ go get -u github.com/larapulse/migrator
```

2. Import it in your code:

```go
import "github.com/larapulse/migrator"
```

## Quick start

Initialize migrator with migration entries:

```go
var migrations = []migrator.Migration{
	{
		Name: "19700101_0001_create_posts_table",
		Up: func() migrator.Schema {
			var s migrator.Schema
			posts := migrator.Table{Name: "posts"}

			posts.UniqueID("id")
			posts.Column("title", migrator.String{Precision: 64})
			posts.Column("content", migrator.Text{})
			posts.Timestamps()

			s.CreateTable(posts)

			return s
		},
		Down: func() migrator.Schema {
			var s migrator.Schema

			s.DropTableIfExists("posts")

			return s
		},
	},
	{
		Name: "19700101_0002_create_comments_table",
		Up: func() migrator.Schema {
			var s migrator.Schema
			comments := migrator.Table{Name: "comments"}

			comments.UniqueID("id")
			comments.UUID("post_id", "", false)
			comments.Column("title", migrator.String{Precision: 64})
			comments.Column("content", migrator.Text{})
			comments.Timestamps()

			comments.Foreign("post_id", "id", "posts", "RESTRICT", "RESTRICT")

			s.CreateTable(comments)

			return s
		},
		Down: func() migrator.Schema {
			var s migrator.Schema

			s.DropTableIfExists("comments")

			return s
		},
	},
}

m := migrator.Migrator{Pool: migrations}
migrated, err = m.Migrate(db)

if err != nil {
	log.Errorf("Could not migrate: %v", err)
	os.Exit(1)
}

if len(migrated) == 0 {
	log.Print("Nothing were migrated.")
}

for _, m := range migrated {
	log.Printf("Migration: %s was migrated ✅", m)
}

log.Print("Migration did run successfully")
```

After the first migration run, `migrations` table will be created:

```
+----+-------------------------------------+-------+---------------------+
| id | name                                | batch | applied_at          |
+----+-------------------------------------+-------+---------------------+
|  1 | 19700101_0001_create_posts_table    |     1 | 2020-06-27 00:00:00 |
|  2 | 19700101_0002_create_comments_table |     1 | 2020-06-27 00:00:00 |
+----+-------------------------------------+-------+---------------------+
```

If you want to use another name for migration table, change it `Migrator` before running migrations:

```go
m := migrator.Migrator{TableName: "_my_app_migrations"}
```

### Transactional migration

In case you have multiple commands within one migration and you want to be sure it is migrated properly, you might enable transactional execution per migration:

```go
var migration = migrator.Migration{
	Name: "19700101_0001_create_posts_and_users_tables",
	Up: func() migrator.Schema {
		var s migrator.Schema
		posts := migrator.Table{Name: "posts"}
		posts.UniqueID("id")
		posts.Timestamps()

		users := migrator.Table{Name: "users"}
		users.UniqueID("id")
		users.Timestamps()

		s.CreateTable(posts)
		s.CreateTable(users)

		return s
	},
	Down: func() migrator.Schema {
		var s migrator.Schema

		s.DropTableIfExists("users")
		s.DropTableIfExists("posts")

		return s
	},
	Transaction: true,
}
```

### Rollback and revert

In case you need to revert your deploy and DB, you can revert last migrated batch:

```go
m := migrator.Migrator{Pool: migrations}
reverted, err := m.Rollback(db)

if err != nil {
	log.Errorf("Could not roll back migrations: %v", err)
	os.Exit(1)
}

if len(reverted) == 0 {
	log.Print("Nothing were rolled back.")
}

for _, m := range reverted {
	log.Printf("Migration: %s was rolled back ✅", m)
}
```

To revert all migrated items back, you have to call `Revert()` on your `migrator`:

```go
m := migrator.Migrator{Pool: migrations}
reverted, err := m.Revert(db)
```

## Customize queries

You may add any column definition to the database on your own, just be sure you implement `columnType` interface:

```go
type customType string

func (ct customType) buildRow() string {
	return string(ct)
}

posts := migrator.Table{Name: "posts"}
posts.UniqueID("id")
posts.Column("data", customType("json not null"))
posts.Timestamps()
```

The same logic is for adding custom commands to the Schema to be migrated or reverted, just be sure you implement `command` interface:

```go
type customCommand string

func (cc customCommand) toSQL() string {
	return string(cc)
}

var s migrator.Schema

c := customCommand("DROP PROCEDURE abc")
s.CustomCommand(c)
```
