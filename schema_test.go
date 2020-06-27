package migrator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchemaCreateTable(t *testing.T) {
	assert := assert.New(t)

	s := Schema{}
	assert.Len(s.pool, 0)

	tb := Table{Name: "test"}
	s.CreateTable(tb)

	assert.Len(s.pool, 1)
	assert.Equal(createTableCommand{tb}, s.pool[0])
}

func TestSchemaDropTable(t *testing.T) {
	assert := assert.New(t)

	s := Schema{}
	assert.Len(s.pool, 0)

	s.DropTable("test", false, "")

	assert.Len(s.pool, 1)
	assert.Equal(dropTableCommand{"test", false, ""}, s.pool[0])
}

func TestSchemaRenameTable(t *testing.T) {
	assert := assert.New(t)

	s := Schema{}
	assert.Len(s.pool, 0)

	s.RenameTable("from", "to")

	assert.Len(s.pool, 1)
	assert.Equal(renameTableCommand{"from", "to"}, s.pool[0])
}

func TestSchemaAlterTable(t *testing.T) {
	assert := assert.New(t)

	s := Schema{}
	assert.Len(s.pool, 0)

	s.AlterTable("table", TableCommands{})

	assert.Len(s.pool, 1)
	assert.Equal(alterTableCommand{"table", TableCommands{}}, s.pool[0])
}

func TestSchemaCustomCommand(t *testing.T) {
	assert := assert.New(t)
	c := testDummyCommand("DROP PROCEDURE abc")

	s := Schema{}
	assert.Len(s.pool, 0)

	s.CustomCommand(c)

	assert.Len(s.pool, 1)
	assert.Equal(c, s.pool[0])
}
