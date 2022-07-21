package migrator

import (
	"fmt"
	"strings"
)

// TableCommands is a pool of commands to be executed on the table.
// https://dev.mysql.com/doc/refman/8.0/en/alter-table.html
type TableCommands []Command

func (tc TableCommands) ToSQL() string {
	rows := []string{}

	for _, c := range tc {
		rows = append(rows, c.ToSQL())
	}

	return strings.Join(rows, ", ")
}

// AddColumnCommand is a command to add the column to the table.
type AddColumnCommand struct {
	Name   string
	Column ColumnType
	After  string
	First  bool
}

func (c AddColumnCommand) ToSQL() string {
	if c.Column == nil {
		return ""
	}

	definition := c.Column.BuildRow()
	if c.Name == "" || definition == "" {
		return ""
	}

	sql := "ADD COLUMN `" + c.Name + "` " + definition

	if c.After != "" {
		sql += " AFTER " + c.After
	} else if c.First {
		sql += " FIRST"
	}

	return sql
}

// RenameColumnCommand is a command to rename a column in the table.
// Warning ⚠️ BC incompatible!
//
// Info ℹ️ extension for Oracle compatibility.
type RenameColumnCommand struct {
	Old string
	New string
}

func (c RenameColumnCommand) ToSQL() string {
	if c.Old == "" || c.New == "" {
		return ""
	}

	return fmt.Sprintf("RENAME COLUMN `%s` TO `%s`", c.Old, c.New)
}

// ModifyColumnCommand is a command to modify column type.
// Warning ⚠️ BC incompatible!
//
// Info ℹ️ extension for Oracle compatibility.
type ModifyColumnCommand struct {
	Name   string
	Column ColumnType
}

func (c ModifyColumnCommand) ToSQL() string {
	if c.Column == nil {
		return ""
	}

	definition := c.Column.BuildRow()
	if c.Name == "" || definition == "" {
		return ""
	}

	return fmt.Sprintf("MODIFY `%s` %s", c.Name, definition)
}

// ChangeColumnCommand is a default command to change column.
// Warning ⚠️ BC incompatible!
type ChangeColumnCommand struct {
	From   string
	To     string
	Column ColumnType
}

func (c ChangeColumnCommand) ToSQL() string {
	if c.Column == nil {
		return ""
	}

	definition := c.Column.BuildRow()
	if c.From == "" || c.To == "" || definition == "" {
		return ""
	}

	return fmt.Sprintf("CHANGE `%s` `%s` %s", c.From, c.To, c.Column.BuildRow())
}

// DropColumnCommand is a command to drop a column from the table.
// Warning ⚠️ BC incompatible!
type DropColumnCommand string

// Info ℹ️ campatible with Oracle
func (c DropColumnCommand) ToSQL() string {
	if c == "" {
		return ""
	}

	return fmt.Sprintf("DROP COLUMN `%s`", c)
}

// AddIndexCommand adds a key to the table.
type AddIndexCommand struct {
	Name    string
	Columns []string
}

func (c AddIndexCommand) ToSQL() string {
	if c.Name == "" || len(c.Columns) == 0 {
		return ""
	}

	return fmt.Sprintf("ADD KEY `%s` (`%s`)", c.Name, strings.Join(c.Columns, "`, `"))
}

// DropIndexCommand removes the key from the table.
type DropIndexCommand string

func (c DropIndexCommand) ToSQL() string {
	if c == "" {
		return ""
	}

	return fmt.Sprintf("DROP KEY `%s`", c)
}

// AddForeignCommand adds the foreign key constraint to the table.
type AddForeignCommand struct {
	Foreign Foreign
}

func (c AddForeignCommand) ToSQL() string {
	if c.Foreign.render() == "" {
		return ""
	}

	return "ADD " + c.Foreign.render()
}

// DropForeignCommand is a command to remove a foreign key constraint.
type DropForeignCommand string

func (c DropForeignCommand) ToSQL() string {
	if c == "" {
		return ""
	}

	return fmt.Sprintf("DROP FOREIGN KEY `%s`", c)
}

// AddUniqueIndexCommand is a command to add a unique key to the table on some columns.
type AddUniqueIndexCommand struct {
	Key     string
	Columns []string
}

func (c AddUniqueIndexCommand) ToSQL() string {
	if c.Key == "" || len(c.Columns) == 0 {
		return ""
	}

	return fmt.Sprintf("ADD UNIQUE KEY `%s` (`%s`)", c.Key, strings.Join(c.Columns, "`, `"))
}

// AddPrimaryIndexCommand is a command to add a primary key.
type AddPrimaryIndexCommand string

func (c AddPrimaryIndexCommand) ToSQL() string {
	if c == "" {
		return ""
	}

	return fmt.Sprintf("ADD PRIMARY KEY (`%s`)", c)
}

// DropPrimaryIndexCommand is a command to remove the primary key from the table.
type DropPrimaryIndexCommand struct{}

func (c DropPrimaryIndexCommand) ToSQL() string {
	return "DROP PRIMARY KEY"
}

// ADD {FULLTEXT | SPATIAL} [INDEX | KEY] [index_name] (key_part,...) [index_option] ...
// DROP {CHECK | CONSTRAINT} symbol
// RENAME {INDEX | KEY} old_index_name TO new_index_name
