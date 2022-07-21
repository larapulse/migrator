package migrator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTableCommands(t *testing.T) {
	t.Run("it returns empty on empty commands list", func(t *testing.T) {
		c := TableCommands{}
		assert.Equal(t, "", c.ToSQL())
	})

	t.Run("it renders row from one command", func(t *testing.T) {
		c := TableCommands{testCommand("test")}
		assert.Equal(t, "Do action on test", c.ToSQL())
	})

	t.Run("it renders row from multiple commands", func(t *testing.T) {
		c := TableCommands{testCommand("test"), testCommand("bang")}
		assert.Equal(t, "Do action on test, Do action on bang", c.ToSQL())
	})
}

func TestAddColumnCommand(t *testing.T) {
	t.Run("it returns an empty string if column definition missing", func(t *testing.T) {
		c := AddColumnCommand{Name: "tests"}
		assert.Equal(t, "", c.ToSQL())
	})

	t.Run("it returns an empty string if column name missing", func(t *testing.T) {
		c := AddColumnCommand{Column: testColumnType("test")}
		assert.Equal(t, "", c.ToSQL())
	})

	t.Run("it returns an empty string if column definition empty", func(t *testing.T) {
		c := AddColumnCommand{Name: "tests", Column: testColumnType("")}
		assert.Equal(t, "", c.ToSQL())
	})

	t.Run("it returns base row", func(t *testing.T) {
		c := AddColumnCommand{Name: "test_id", Column: testColumnType("definition")}
		assert.Equal(t, "ADD COLUMN `test_id` definition", c.ToSQL())
	})

	t.Run("it returns row with after column", func(t *testing.T) {
		c := AddColumnCommand{Name: "test_id", Column: testColumnType("definition"), After: "id"}
		assert.Equal(t, "ADD COLUMN `test_id` definition AFTER id", c.ToSQL())
	})

	t.Run("it returns row with first flag", func(t *testing.T) {
		c := AddColumnCommand{Name: "test_id", Column: testColumnType("definition"), First: true}
		assert.Equal(t, "ADD COLUMN `test_id` definition FIRST", c.ToSQL())
	})
}

func TestRenameColumnCommand(t *testing.T) {
	t.Run("it returns an empty string if old name missing", func(t *testing.T) {
		c := RenameColumnCommand{New: "test"}
		assert.Equal(t, "", c.ToSQL())
	})

	t.Run("it returns an empty string if new name missing", func(t *testing.T) {
		c := RenameColumnCommand{Old: "test"}
		assert.Equal(t, "", c.ToSQL())
	})

	t.Run("it returns a proper row", func(t *testing.T) {
		c := RenameColumnCommand{Old: "from", New: "to"}
		assert.Equal(t, "RENAME COLUMN `from` TO `to`", c.ToSQL())
	})
}

func TestModifyColumnCommand(t *testing.T) {
	t.Run("it returns an empty string if column definition missing", func(t *testing.T) {
		c := ModifyColumnCommand{Name: "tests"}
		assert.Equal(t, "", c.ToSQL())
	})

	t.Run("it returns an empty string if column name missing", func(t *testing.T) {
		c := ModifyColumnCommand{Column: testColumnType("test")}
		assert.Equal(t, "", c.ToSQL())
	})

	t.Run("it returns an empty string if column definition empty", func(t *testing.T) {
		c := ModifyColumnCommand{Name: "tests", Column: testColumnType("")}
		assert.Equal(t, "", c.ToSQL())
	})

	t.Run("it returns a proper row", func(t *testing.T) {
		c := ModifyColumnCommand{Name: "test_id", Column: testColumnType("definition")}
		assert.Equal(t, "MODIFY `test_id` definition", c.ToSQL())
	})
}

func TestChangeColumnCommand(t *testing.T) {
	t.Run("it returns an empty string if column definition missing", func(t *testing.T) {
		c := ChangeColumnCommand{From: "tests", To: "something"}
		assert.Equal(t, "", c.ToSQL())
	})

	t.Run("it returns an empty string if column from name missing", func(t *testing.T) {
		c := ChangeColumnCommand{To: "something", Column: testColumnType("test")}
		assert.Equal(t, "", c.ToSQL())
	})

	t.Run("it returns an empty string if column to name missing", func(t *testing.T) {
		c := ChangeColumnCommand{From: "tests", Column: testColumnType("test")}
		assert.Equal(t, "", c.ToSQL())
	})

	t.Run("it returns an empty string if column definition empty", func(t *testing.T) {
		c := ChangeColumnCommand{From: "tests", To: "something", Column: testColumnType("")}
		assert.Equal(t, "", c.ToSQL())
	})

	t.Run("it returns a proper row", func(t *testing.T) {
		c := ChangeColumnCommand{From: "tests", To: "something", Column: testColumnType("definition")}
		assert.Equal(t, "CHANGE `tests` `something` definition", c.ToSQL())
	})
}

func TestDropColumnCommand(t *testing.T) {
	t.Run("it returns an empty string if column name missing", func(t *testing.T) {
		c := DropColumnCommand("")
		assert.Equal(t, "", c.ToSQL())
	})

	t.Run("it returns a proper row", func(t *testing.T) {
		c := DropColumnCommand("test_id")
		assert.Equal(t, "DROP COLUMN `test_id`", c.ToSQL())
	})
}

func TestAddIndexCommand(t *testing.T) {
	t.Run("it returns an empty string if index name missing", func(t *testing.T) {
		c := AddIndexCommand{Columns: []string{"test"}}
		assert.Equal(t, "", c.ToSQL())
	})

	t.Run("it returns an empty string if columns list empty", func(t *testing.T) {
		c := AddIndexCommand{Name: "test", Columns: []string{}}
		assert.Equal(t, "", c.ToSQL())
	})

	t.Run("it returns a proper row", func(t *testing.T) {
		c := AddIndexCommand{Name: "test_idx", Columns: []string{"test"}}
		assert.Equal(t, "ADD KEY `test_idx` (`test`)", c.ToSQL())
	})
}

func TestDropIndexCommand(t *testing.T) {
	t.Run("it returns an empty string if index name missing", func(t *testing.T) {
		c := DropIndexCommand("")
		assert.Equal(t, "", c.ToSQL())
	})

	t.Run("it returns a proper row", func(t *testing.T) {
		c := DropIndexCommand("test_idx")
		assert.Equal(t, "DROP KEY `test_idx`", c.ToSQL())
	})
}

func TestAddForeignCommand(t *testing.T) {
	t.Run("it returns an empty string on missing foreign key", func(t *testing.T) {
		c := AddForeignCommand{}
		assert.Equal(t, "", c.ToSQL())
	})

	t.Run("it builds a proper row", func(t *testing.T) {
		c := AddForeignCommand{Foreign{Key: "idx_foreign", Column: "test_id", Reference: "id", On: "tests"}}
		assert.Equal(t, "ADD CONSTRAINT `idx_foreign` FOREIGN KEY (`test_id`) REFERENCES `tests` (`id`)", c.ToSQL())
	})
}

func TestDropForeignCommand(t *testing.T) {
	t.Run("it returns an empty string if index name missing", func(t *testing.T) {
		c := DropForeignCommand("")
		assert.Equal(t, "", c.ToSQL())
	})

	t.Run("it returns a proper row", func(t *testing.T) {
		c := DropForeignCommand("test_idx")
		assert.Equal(t, "DROP FOREIGN KEY `test_idx`", c.ToSQL())
	})
}

func TestAddUniqueIndexCommand(t *testing.T) {
	t.Run("it returns an empty string if index name missing", func(t *testing.T) {
		c := AddUniqueIndexCommand{Columns: []string{"test"}}
		assert.Equal(t, "", c.ToSQL())
	})

	t.Run("it returns an empty string if columns list empty", func(t *testing.T) {
		c := AddUniqueIndexCommand{Key: "test", Columns: []string{}}
		assert.Equal(t, "", c.ToSQL())
	})

	t.Run("it returns a proper row", func(t *testing.T) {
		c := AddUniqueIndexCommand{Key: "test_idx", Columns: []string{"test"}}
		assert.Equal(t, "ADD UNIQUE KEY `test_idx` (`test`)", c.ToSQL())
	})
}

func TestAddPrimaryIndexCommand(t *testing.T) {
	t.Run("it returns an empty string if index name missing", func(t *testing.T) {
		c := AddPrimaryIndexCommand("")
		assert.Equal(t, "", c.ToSQL())
	})

	t.Run("it returns a proper row", func(t *testing.T) {
		c := AddPrimaryIndexCommand("test_idx")
		assert.Equal(t, "ADD PRIMARY KEY (`test_idx`)", c.ToSQL())
	})
}

func TestDropPrimaryIndexCommand(t *testing.T) {
	c := DropPrimaryIndexCommand{}
	assert.Equal(t, "DROP PRIMARY KEY", c.ToSQL())
}
