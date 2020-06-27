// Package migrator represents MySQL database migrator
package migrator

import "strings"

// Table is an entity to create table
//
// Name			table name
// Engine		default: InnoDB
// Charset		default: utf8mb4 or first part of collation (if set)
// Collation	default: utf8mb4_unicode_ci or charset with `_unicode_ci` suffix
// Comment		optional comment on table
type Table struct {
	Name      string
	columns   columns
	indexes   keys
	foreigns  foreigns
	Engine    string
	Charset   string
	Collation string
	Comment   string
}

// Column adds column to the table
func (t *Table) Column(name string, c columnType) {
	t.columns = append(t.columns, column{field: name, definition: c})
}

// ID adds bigint `id` column that is primary key
func (t *Table) ID(name string) {
	t.Column(name, Integer{
		Prefix:        "big",
		Unsigned:      true,
		Autoincrement: true,
	})
	t.Primary(name)
}

// UniqueID adds unique id column (represented as UUID) that is primary key
func (t *Table) UniqueID(name string) {
	t.UUID(name, "(UUID())", false)
	t.Primary(name)
}

// Boolean represented in DB as tinyint
func (t *Table) Boolean(name string, def string) {
	// tinyint(1)
	t.Column(name, Integer{
		Prefix:    "tiny",
		Unsigned:  true,
		Precision: 1,
		Default:   def,
	})
}

// UUID adds char(36) column
func (t *Table) UUID(name string, def string, nullable bool) {
	// char(36)
	t.Column(name, String{
		Fixed:     true,
		Precision: 36,
		Default:   def,
		Nullable:  nullable,
	})
}

// Timestamps adds default timestamps: `created_at` and `updated_at`
func (t *Table) Timestamps() {
	// created_at not null default CURRENT_TIMESTAMP
	t.Column("created_at", Timable{
		Type:    "timestamp",
		Default: "CURRENT_TIMESTAMP",
	})
	// updated_at not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	t.Column("updated_at", Timable{
		Type:     "timestamp",
		Default:  "CURRENT_TIMESTAMP",
		OnUpdate: "CURRENT_TIMESTAMP",
	})
}

// Primary adds primary key
func (t *Table) Primary(columns ...string) {
	if len(columns) == 0 {
		return
	}

	t.indexes = append(t.indexes, key{
		typ:     "primary",
		columns: columns,
	})
}

// Unique adds unique key on selected columns
func (t *Table) Unique(columns ...string) {
	if len(columns) == 0 {
		return
	}

	t.indexes = append(t.indexes, key{
		name:    t.buildUniqueKeyName(columns...),
		typ:     "unique",
		columns: columns,
	})
}

// Index adds index (key) on selected columns
func (t *Table) Index(name string, columns ...string) {
	if len(columns) == 0 {
		return
	}

	t.indexes = append(t.indexes, key{name: name, columns: columns})
}

// Foreign adds foreign key constraints
func (t *Table) Foreign(column string, reference string, on string, onUpdate string, onDelete string) {
	name := t.buildForeignKeyName(column)
	t.indexes = append(t.indexes, key{
		name:    name,
		columns: []string{column},
	})
	t.foreigns = append(t.foreigns, foreign{
		key:       name,
		column:    column,
		reference: reference,
		on:        on,
		onUpdate:  onUpdate,
		onDelete:  onDelete,
	})
}

func (t *Table) buildUniqueKeyName(columns ...string) string {
	return t.Name + "_" + strings.Join(columns, "_") + "_unique"
}

func (t *Table) buildForeignKeyName(column string) string {
	return t.Name + "_" + column + "_foreign"
}
