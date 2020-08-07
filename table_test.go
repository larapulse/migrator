package migrator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTableColumns(t *testing.T) {
	c := testColumnType("test")

	assert := assert.New(t)

	table := Table{}
	assert.Len(table.columns, 0)

	table.Column("test", c)

	assert.Len(table.columns, 1)
	assert.Equal(columns{column{"test", c}}, table.columns)
}

func TestIDColumn(t *testing.T) {
	assert := assert.New(t)
	table := Table{}

	assert.Nil(table.columns)
	assert.Len(table.indexes, 0)

	table.ID("id")

	assert.Len(table.columns, 1)
	assert.Equal("id", table.columns[0].field)
	assert.Equal(Integer{Prefix: "big", Unsigned: true, Autoincrement: true}, table.columns[0].definition)
	assert.Len(table.indexes, 1)
	assert.Equal(Key{Type: "primary", Columns: []string{"id"}}, table.indexes[0])
}

func TestUniqueIDColumn(t *testing.T) {
	assert := assert.New(t)
	table := Table{}

	assert.Nil(table.columns)
	assert.Len(table.indexes, 0)

	table.UniqueID("id")

	assert.Len(table.columns, 1)
	assert.Equal("id", table.columns[0].field)
	assert.Equal(String{Default: "(UUID())", Fixed: true, Precision: 36}, table.columns[0].definition)
	assert.Len(table.indexes, 1)
	assert.Equal(Key{Type: "primary", Columns: []string{"id"}}, table.indexes[0])
}

func TestBinaryID(t *testing.T) {
	assert := assert.New(t)
	table := Table{}

	assert.Nil(table.columns)
	assert.Len(table.indexes, 0)

	table.BinaryID("id")

	assert.Len(table.columns, 1)
	assert.Equal("id", table.columns[0].field)
	assert.Equal(Binary{Default: "(UUID_TO_BIN(UUID()))", Fixed: true, Precision: 16}, table.columns[0].definition)
	assert.Len(table.indexes, 1)
	assert.Equal(Key{Type: "primary", Columns: []string{"id"}}, table.indexes[0])
}

func TestBooleanColumn(t *testing.T) {
	assert := assert.New(t)
	table := Table{}

	assert.Nil(table.columns)

	table.Boolean("flag", "1")

	assert.Len(table.columns, 1)
	assert.Equal("flag", table.columns[0].field)
	assert.Equal(Integer{Prefix: "tiny", Default: "1", Unsigned: true, Precision: 1}, table.columns[0].definition)
}

func TestUUIDColumn(t *testing.T) {
	assert := assert.New(t)
	table := Table{}

	assert.Nil(table.columns)

	table.UUID("uuid", "1111", true)

	assert.Len(table.columns, 1)
	assert.Equal("uuid", table.columns[0].field)
	assert.Equal(String{Default: "1111", Fixed: true, Precision: 36, Nullable: true}, table.columns[0].definition)
}

func TestTimestampsColumn(t *testing.T) {
	assert := assert.New(t)
	table := Table{}

	assert.Nil(table.columns)

	table.Timestamps()

	assert.Len(table.columns, 2)
	assert.Equal("created_at", table.columns[0].field)
	assert.Equal(Timable{Type: "timestamp", Precision: 6, Default: "CURRENT_TIMESTAMP(6)"}, table.columns[0].definition)
	assert.Equal("updated_at", table.columns[1].field)
	assert.Equal(Timable{Type: "timestamp", Precision: 6, Default: "CURRENT_TIMESTAMP(6)", OnUpdate: "CURRENT_TIMESTAMP(6)"}, table.columns[1].definition)
}

func TestIntColumn(t *testing.T) {
	assert := assert.New(t)
	table := Table{}

	assert.Nil(table.columns)

	table.Int("number", 64, true)

	assert.Len(table.columns, 1)
	assert.Equal("number", table.columns[0].field)
	assert.Equal(Integer{Precision: 64, Unsigned: true}, table.columns[0].definition)
}

func TestBigIntColumn(t *testing.T) {
	assert := assert.New(t)
	table := Table{}

	assert.Nil(table.columns)

	table.BigInt("number", 64, true)

	assert.Len(table.columns, 1)
	assert.Equal("number", table.columns[0].field)
	assert.Equal(Integer{Prefix: "big", Precision: 64, Unsigned: true}, table.columns[0].definition)
}

func TestFloatColumn(t *testing.T) {
	assert := assert.New(t)
	table := Table{}

	assert.Nil(table.columns)

	table.Float("number", 15, 2)

	assert.Len(table.columns, 1)
	assert.Equal("number", table.columns[0].field)
	assert.Equal(Floatable{Precision: 15, Scale: 2}, table.columns[0].definition)
}

func TestFixedFloatColumn(t *testing.T) {
	assert := assert.New(t)
	table := Table{}

	assert.Nil(table.columns)

	table.FixedFloat("number", 15, 2)

	assert.Len(table.columns, 1)
	assert.Equal("number", table.columns[0].field)
	assert.Equal(Floatable{Type: "decimal", Precision: 15, Scale: 2}, table.columns[0].definition)
}

func TestDecimalColumn(t *testing.T) {
	assert := assert.New(t)
	table := Table{}

	assert.Nil(table.columns)

	table.Decimal("number", 15, 2)

	assert.Len(table.columns, 1)
	assert.Equal("number", table.columns[0].field)
	assert.Equal(Floatable{Type: "decimal", Precision: 15, Scale: 2}, table.columns[0].definition)
}

func TestVarcharColumn(t *testing.T) {
	assert := assert.New(t)
	table := Table{}

	assert.Nil(table.columns)

	table.Varchar("string", 64)

	assert.Len(table.columns, 1)
	assert.Equal("string", table.columns[0].field)
	assert.Equal(String{Precision: 64}, table.columns[0].definition)
}

func TestCharColumn(t *testing.T) {
	assert := assert.New(t)
	table := Table{}

	assert.Nil(table.columns)

	table.Char("string", 32)

	assert.Len(table.columns, 1)
	assert.Equal("string", table.columns[0].field)
	assert.Equal(String{Fixed: true, Precision: 32}, table.columns[0].definition)
}

func TestTextColumn(t *testing.T) {
	assert := assert.New(t)
	table := Table{}

	assert.Nil(table.columns)

	table.Text("string", true)

	assert.Len(table.columns, 1)
	assert.Equal("string", table.columns[0].field)
	assert.Equal(Text{Nullable: true}, table.columns[0].definition)
}

func TestBlobColumn(t *testing.T) {
	assert := assert.New(t)
	table := Table{}

	assert.Nil(table.columns)

	table.Blob("string", true)

	assert.Len(table.columns, 1)
	assert.Equal("string", table.columns[0].field)
	assert.Equal(Text{Blob: true, Nullable: true}, table.columns[0].definition)
}

func TestJsonColumn(t *testing.T) {
	assert := assert.New(t)
	table := Table{}

	assert.Nil(table.columns)

	table.JSON("data")

	assert.Len(table.columns, 1)
	assert.Equal("data", table.columns[0].field)
	assert.Equal(JSON{}, table.columns[0].definition)
}

func TestTimestampColumn(t *testing.T) {
	assert := assert.New(t)
	table := Table{}

	assert.Nil(table.columns)

	table.Timestamp("date", true, "CURRENT_TIMESTAMP")

	assert.Len(table.columns, 1)
	assert.Equal("date", table.columns[0].field)
	assert.Equal(Timable{Nullable: true, Default: "CURRENT_TIMESTAMP"}, table.columns[0].definition)
}

func TestPreciseTimestampColumn(t *testing.T) {
	assert := assert.New(t)
	table := Table{}

	assert.Nil(table.columns)

	table.PreciseTimestamp("date", 3, true, "CURRENT_TIMESTAMP")

	assert.Len(table.columns, 1)
	assert.Equal("date", table.columns[0].field)
	assert.Equal(Timable{Precision: 3, Nullable: true, Default: "CURRENT_TIMESTAMP"}, table.columns[0].definition)
}

func TestDateColumn(t *testing.T) {
	assert := assert.New(t)
	table := Table{}

	assert.Nil(table.columns)

	table.Date("date", true, "NOW()")

	assert.Len(table.columns, 1)
	assert.Equal("date", table.columns[0].field)
	assert.Equal(Timable{Type: "date", Nullable: true, Default: "NOW()"}, table.columns[0].definition)
}

func TestTimeColumn(t *testing.T) {
	assert := assert.New(t)
	table := Table{}

	assert.Nil(table.columns)

	table.Time("time", true, "NOW()")

	assert.Len(table.columns, 1)
	assert.Equal("time", table.columns[0].field)
	assert.Equal(Timable{Type: "time", Nullable: true, Default: "NOW()"}, table.columns[0].definition)
}

func TestYearColumn(t *testing.T) {
	assert := assert.New(t)
	table := Table{}

	assert.Nil(table.columns)

	table.Year("year", true, "YEAR(NOW())")

	assert.Len(table.columns, 1)
	assert.Equal("year", table.columns[0].field)
	assert.Equal(Timable{Type: "year", Nullable: true, Default: "YEAR(NOW())"}, table.columns[0].definition)
}

func TestBinaryColumn(t *testing.T) {
	assert := assert.New(t)
	table := Table{}

	assert.Nil(table.columns)

	table.Binary("binary", 36, true)

	assert.Len(table.columns, 1)
	assert.Equal("binary", table.columns[0].field)
	assert.Equal(Binary{Fixed: true, Precision: 36, Nullable: true}, table.columns[0].definition)
}

func TestVarbinaryColumn(t *testing.T) {
	assert := assert.New(t)
	table := Table{}

	assert.Nil(table.columns)

	table.Varbinary("binary", 36, true)

	assert.Len(table.columns, 1)
	assert.Equal("binary", table.columns[0].field)
	assert.Equal(Binary{Precision: 36, Nullable: true}, table.columns[0].definition)
}

func TestTablePrimaryIndex(t *testing.T) {
	t.Run("it skips adding key on empty columns list", func(t *testing.T) {
		assert := assert.New(t)
		table := Table{}

		assert.Nil(table.indexes)

		table.Primary()

		assert.Nil(table.indexes)
	})

	t.Run("it adds primary key", func(t *testing.T) {
		assert := assert.New(t)
		table := Table{}

		assert.Nil(table.indexes)

		table.Primary("id", "name")

		assert.Len(table.indexes, 1)
		assert.Equal(Key{Type: "primary", Columns: []string{"id", "name"}}, table.indexes[0])
	})
}

func TestTableUniqueIndex(t *testing.T) {
	t.Run("it skips adding key on empty columns list", func(t *testing.T) {
		assert := assert.New(t)
		table := Table{}

		assert.Nil(table.indexes)

		table.Unique()

		assert.Nil(table.indexes)
	})

	t.Run("it adds unique key", func(t *testing.T) {
		assert := assert.New(t)
		table := Table{Name: "table"}

		assert.Nil(table.indexes)

		table.Unique("id", "name")

		assert.Len(table.indexes, 1)
		assert.Equal(Key{Name: "table_id_name_unique", Type: "unique", Columns: []string{"id", "name"}}, table.indexes[0])
	})
}

func TestTableIndex(t *testing.T) {
	t.Run("it skips adding key on empty columns list", func(t *testing.T) {
		assert := assert.New(t)
		table := Table{}

		assert.Nil(table.indexes)

		table.Index("test")

		assert.Nil(table.indexes)
	})

	t.Run("it adds unique key", func(t *testing.T) {
		assert := assert.New(t)
		table := Table{Name: "table"}

		assert.Nil(table.indexes)

		table.Index("test_idx", "id", "name")

		assert.Len(table.indexes, 1)
		assert.Equal(Key{Name: "test_idx", Columns: []string{"id", "name"}}, table.indexes[0])
	})
}

func TestTableForeignIndex(t *testing.T) {
	assert := assert.New(t)
	table := Table{Name: "table"}

	assert.Nil(table.indexes)
	assert.Nil(table.foreigns)

	table.Foreign("test_id", "id", "tests", "set null", "cascade")

	assert.Len(table.indexes, 1)
	assert.Equal(Key{Name: "table_test_id_foreign", Columns: []string{"test_id"}}, table.indexes[0])
	assert.Len(table.foreigns, 1)
	assert.Equal(
		Foreign{Key: "table_test_id_foreign", Column: "test_id", Reference: "id", On: "tests", OnUpdate: "set null", OnDelete: "cascade"},
		table.foreigns[0],
	)
}
