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
	assert.Equal(key{typ: "primary", columns: []string{"id"}}, table.indexes[0])
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
	assert.Equal(key{typ: "primary", columns: []string{"id"}}, table.indexes[0])
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
	assert.Equal(Timable{Type: "timestamp", Default: "CURRENT_TIMESTAMP"}, table.columns[0].definition)
	assert.Equal("updated_at", table.columns[1].field)
	assert.Equal(Timable{Type: "timestamp", Default: "CURRENT_TIMESTAMP", OnUpdate: "CURRENT_TIMESTAMP"}, table.columns[1].definition)
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
		assert.Equal(key{typ: "primary", columns: []string{"id", "name"}}, table.indexes[0])
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
		assert.Equal(key{name: "table_id_name_unique", typ: "unique", columns: []string{"id", "name"}}, table.indexes[0])
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
		assert.Equal(key{name: "test_idx", columns: []string{"id", "name"}}, table.indexes[0])
	})
}

func TestTableForeignIndex(t *testing.T) {
	assert := assert.New(t)
	table := Table{Name: "table"}

	assert.Nil(table.indexes)
	assert.Nil(table.foreigns)

	table.Foreign("test_id", "id", "tests", "set null", "cascade")

	assert.Len(table.indexes, 1)
	assert.Equal(key{name: "table_test_id_foreign", columns: []string{"test_id"}}, table.indexes[0])
	assert.Len(table.foreigns, 1)
	assert.Equal(
		foreign{key: "table_test_id_foreign", column: "test_id", reference: "id", on: "tests", onUpdate: "set null", onDelete: "cascade"},
		table.foreigns[0],
	)
}

func TestBuildUniqueIndexName(t *testing.T) {
	t.Run("It builds name from one column", func(t *testing.T) {
		table := Table{Name: "table"}

		assert.Equal(t, "table_test_unique", table.buildUniqueKeyName("test"))
	})

	t.Run("it builds name from multiple columns", func(t *testing.T) {
		table := Table{Name: "table"}

		assert.Equal(t, "table_test_again_unique", table.buildUniqueKeyName("test", "again"))
	})
}

func TestBuildForeignIndexName(t *testing.T) {
	table := Table{Name: "table"}

	assert.Equal(t, "table_test_foreign", table.buildForeignKeyName("test"))
}
