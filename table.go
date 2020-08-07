package migrator

// Table is an entity to create a table.
//
// - Name		table name
// - Engine		default: InnoDB
// - Charset	default: utf8mb4 or first part of collation (if set)
// - Collation	default: utf8mb4_unicode_ci or charset with `_unicode_ci` suffix
// - Comment	optional comment on table
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

// Column adds a column to the table
func (t *Table) Column(name string, c columnType) {
	t.columns = append(t.columns, column{field: name, definition: c})
}

// ID adds bigint `id` column that is the primary key
func (t *Table) ID(name string) {
	t.Column(name, Integer{
		Prefix:        "big",
		Unsigned:      true,
		Autoincrement: true,
	})
	t.Primary(name)
}

// UniqueID adds unique id column (represented as UUID) that is the primary key
func (t *Table) UniqueID(name string) {
	t.UUID(name, "(UUID())", false)
	t.Primary(name)
}

// BinaryID adds unique binary id column (represented as UUID) that is the primary key
func (t *Table) BinaryID(name string) {
	t.Column(name, Binary{Fixed: true, Precision: 16, Default: "(UUID_TO_BIN(UUID()))"})
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
	// created_at timestamp(6) not null default CURRENT_TIMESTAMP(6)
	t.Column("created_at", Timable{
		Type:      "timestamp",
		Precision: 6,
		Default:   "CURRENT_TIMESTAMP(6)",
	})
	// updated_at timestamp(6) not null default CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6)
	t.Column("updated_at", Timable{
		Type:      "timestamp",
		Precision: 6,
		Default:   "CURRENT_TIMESTAMP(6)",
		OnUpdate:  "CURRENT_TIMESTAMP(6)",
	})
}

// Int adds int(precision) column to the table
func (t *Table) Int(name string, precision uint16, unsigned bool) {
	t.Column(name, Integer{Precision: precision, Unsigned: unsigned})
}

// BigInt adds bigint(precision) column to the table
func (t *Table) BigInt(name string, precision uint16, unsigned bool) {
	t.Column(name, Integer{Prefix: "big", Precision: precision, Unsigned: unsigned})
}

// Float adds float(precision,scale) column to the table
func (t *Table) Float(name string, precision uint16, scale uint16) {
	t.Column(name, Floatable{Precision: precision, Scale: scale})
}

// FixedFloat is an alias to decimal(precision,scale) column
func (t *Table) FixedFloat(name string, precision uint16, scale uint16) {
	t.Decimal(name, precision, scale)
}

// Decimal adds decimal(precision,scale) column to the table
func (t *Table) Decimal(name string, precision uint16, scale uint16) {
	t.Column(name, Floatable{Type: "decimal", Precision: precision, Scale: scale})
}

// Varchar adds varchar(precision) column to the table
func (t *Table) Varchar(name string, precision uint16) {
	t.Column(name, String{Precision: precision})
}

// Char adds char(precision) column to the table
func (t *Table) Char(name string, precision uint16) {
	t.Column(name, String{Fixed: true, Precision: precision})
}

// Text adds text column to the table
func (t *Table) Text(name string, nullable bool) {
	t.Column(name, Text{Nullable: nullable})
}

// Blob adds blob column to the table
func (t *Table) Blob(name string, nullable bool) {
	t.Column(name, Text{Blob: true, Nullable: nullable})
}

// JSON adds json column to the table
func (t *Table) JSON(name string) {
	t.Column(name, JSON{})
}

// Timestamp adds timestamp column to the table
func (t *Table) Timestamp(name string, nullable bool, def string) {
	t.Column(name, Timable{Nullable: nullable, Default: def})
}

// PreciseTimestamp adds timestamp column with precision to the table
func (t *Table) PreciseTimestamp(name string, precision uint16, nullable bool, def string) {
	t.Column(name, Timable{Precision: precision, Nullable: nullable, Default: def})
}

// Date adds date column to the table
func (t *Table) Date(name string, nullable bool, def string) {
	t.Column(name, Timable{Type: "date", Nullable: nullable, Default: def})
}

// Time adds time column to the table
func (t *Table) Time(name string, nullable bool, def string) {
	t.Column(name, Timable{Type: "time", Nullable: nullable, Default: def})
}

// Year adds year column to the table
func (t *Table) Year(name string, nullable bool, def string) {
	t.Column(name, Timable{Type: "year", Nullable: nullable, Default: def})
}

// Binary adds binary(precision) column to the table
func (t *Table) Binary(name string, precision uint16, nullable bool) {
	t.Column(name, Binary{Fixed: true, Precision: precision, Nullable: nullable})
}

// Varbinary adds varbinary(precision) column to the table
func (t *Table) Varbinary(name string, precision uint16, nullable bool) {
	t.Column(name, Binary{Precision: precision, Nullable: nullable})
}

// Primary adds primary key
func (t *Table) Primary(columns ...string) {
	if len(columns) == 0 {
		return
	}

	t.indexes = append(t.indexes, Key{
		Type:    "primary",
		Columns: columns,
	})
}

// Unique adds unique key on selected columns
func (t *Table) Unique(columns ...string) {
	if len(columns) == 0 {
		return
	}

	t.indexes = append(t.indexes, Key{
		Name:    BuildUniqueKeyNameOnTable(t.Name, columns...),
		Type:    "unique",
		Columns: columns,
	})
}

// Index adds index (key) on selected columns
func (t *Table) Index(name string, columns ...string) {
	if len(columns) == 0 {
		return
	}

	t.indexes = append(t.indexes, Key{Name: name, Columns: columns})
}

// Foreign adds foreign key constraints
func (t *Table) Foreign(column string, reference string, on string, onUpdate string, onDelete string) {
	name := BuildForeignNameOnTable(t.Name, column)
	t.indexes = append(t.indexes, Key{
		Name:    name,
		Columns: []string{column},
	})
	t.foreigns = append(t.foreigns, Foreign{
		Key:       name,
		Column:    column,
		Reference: reference,
		On:        on,
		OnUpdate:  onUpdate,
		OnDelete:  onDelete,
	})
}
