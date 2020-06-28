package migrator

import (
	"fmt"
	"strings"
)

type command interface {
	toSQL() string
}

type createTableCommand struct {
	t Table
}

func (c createTableCommand) toSQL() string {
	if c.t.Name == "" {
		return ""
	}

	context := c.t.columns.render()
	if context == "" {
		context = "`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT"
	}

	if res := c.t.indexes.render(); res != "" {
		context += ", " + res
	}

	if res := c.t.foreigns.render(); res != "" {
		context += ", " + res
	}

	engine := c.t.Engine
	if engine == "" {
		engine = "InnoDB"
	}

	charset := c.t.Charset
	collation := c.t.Collation
	if charset == "" && collation == "" {
		charset = "utf8mb4"
		collation = "utf8mb4_unicode_ci"
	}
	if charset == "" && collation != "" {
		parts := strings.Split(collation, "_")
		charset = parts[0]
	}
	if charset != "" && collation == "" {
		collation = charset + "_unicode_ci"
	}

	return fmt.Sprintf(
		"CREATE TABLE `%s` (%s) ENGINE=%s DEFAULT CHARSET=%s COLLATE=%s",
		c.t.Name,
		context,
		engine,
		charset,
		collation,
	)
}

type dropTableCommand struct {
	table  string
	soft   bool
	option string
}

func (c dropTableCommand) toSQL() string {
	sql := "DROP TABLE"

	if c.soft {
		sql += " IF EXISTS"
	}

	sql += fmt.Sprintf(" `%s`", c.table)

	var validOptions = list{"RESTRICT", "CASCADE"}
	if validOptions.has(strings.ToUpper(c.option)) {
		sql += " " + strings.ToUpper(c.option)
	}

	return sql
}

type renameTableCommand struct {
	old string
	new string
}

func (c renameTableCommand) toSQL() string {
	return fmt.Sprintf("RENAME TABLE `%s` TO `%s`", c.old, c.new)
}

type alterTableCommand struct {
	name string
	pool TableCommands
}

func (c alterTableCommand) toSQL() string {
	if c.name == "" || len(c.pool) == 0 {
		return ""
	}

	return "ALTER TABLE `" + c.name + "` " + c.poolToSQL()
}

func (c alterTableCommand) poolToSQL() string {
	var sql []string

	for _, tc := range c.pool {
		sql = append(sql, tc.toSQL())
	}

	return strings.Join(sql, ", ")
}
