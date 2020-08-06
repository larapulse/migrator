package migrator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testColumnType string

func (c testColumnType) buildRow() string {
	return string(c)
}

func TestColumnRender(t *testing.T) {
	t.Run("it renders row from one column", func(t *testing.T) {
		c := columns{column{"test", testColumnType("run")}}

		assert.Equal(t, "`test` run", c.render())
	})

	t.Run("it renders row from multiple columns", func(t *testing.T) {
		c := columns{
			column{"test", testColumnType("run")},
			column{"again", testColumnType("me")},
		}

		assert.Equal(t, "`test` run, `again` me", c.render())
	})
}

func TestInteger(t *testing.T) {
	t.Run("it builds basic column type", func(t *testing.T) {
		c := Integer{}
		assert.Equal(t, "int NOT NULL", c.buildRow())
	})

	t.Run("it build with prefix", func(t *testing.T) {
		c := Integer{Prefix: "super"}
		assert.Equal(t, "superint NOT NULL", c.buildRow())
	})

	t.Run("it builds with precision", func(t *testing.T) {
		c := Integer{Precision: 20}
		assert.Equal(t, "int(20) NOT NULL", c.buildRow())
	})

	t.Run("it builds unsigned", func(t *testing.T) {
		c := Integer{Unsigned: true}
		assert.Equal(t, "int unsigned NOT NULL", c.buildRow())
	})

	t.Run("it builds nullable column type", func(t *testing.T) {
		c := Integer{Nullable: true}
		assert.Equal(t, "int NULL", c.buildRow())
	})

	t.Run("it builds with default value", func(t *testing.T) {
		c := Integer{Default: "0"}
		assert.Equal(t, "int NOT NULL DEFAULT 0", c.buildRow())
	})

	t.Run("it builds with autoincrement", func(t *testing.T) {
		c := Integer{Autoincrement: true}
		assert.Equal(t, "int NOT NULL AUTO_INCREMENT", c.buildRow())
	})

	t.Run("it builds with on_update setter", func(t *testing.T) {
		c := Integer{OnUpdate: "set null"}
		assert.Equal(t, "int NOT NULL ON UPDATE set null", c.buildRow())
	})

	t.Run("it builds with comment", func(t *testing.T) {
		c := Integer{Comment: "test"}
		assert.Equal(t, "int NOT NULL COMMENT 'test'", c.buildRow())
	})

	t.Run("it builds string with all parameters", func(t *testing.T) {
		c := Integer{
			Prefix:        "big",
			Precision:     10,
			Unsigned:      true,
			Nullable:      true,
			Default:       "100",
			Autoincrement: true,
			OnUpdate:      "set null",
			Comment:       "test",
		}

		assert.Equal(
			t,
			"bigint(10) unsigned NULL DEFAULT 100 AUTO_INCREMENT ON UPDATE set null COMMENT 'test'",
			c.buildRow(),
		)
	})
}

func TestFloatable(t *testing.T) {
	t.Run("it builds with default type", func(t *testing.T) {
		c := Floatable{}
		assert.Equal(t, "float NOT NULL", c.buildRow())
	})

	t.Run("it builds basic column type", func(t *testing.T) {
		c := Floatable{Type: "real"}
		assert.Equal(t, "real NOT NULL", c.buildRow())
	})

	t.Run("it builds with precision", func(t *testing.T) {
		c := Floatable{Type: "double", Precision: 20}
		assert.Equal(t, "double(20) NOT NULL", c.buildRow())
	})

	t.Run("it builds with precision and scale", func(t *testing.T) {
		c := Floatable{Type: "decimal", Precision: 10, Scale: 2}
		assert.Equal(t, "decimal(10,2) NOT NULL", c.buildRow())
	})

	t.Run("it builds unsigned", func(t *testing.T) {
		c := Floatable{Unsigned: true}
		assert.Equal(t, "float unsigned NOT NULL", c.buildRow())
	})

	t.Run("it builds nullable column type", func(t *testing.T) {
		c := Floatable{Nullable: true}
		assert.Equal(t, "float NULL", c.buildRow())
	})

	t.Run("it builds with default value", func(t *testing.T) {
		c := Floatable{Default: "0.0"}
		assert.Equal(t, "float NOT NULL DEFAULT 0.0", c.buildRow())
	})

	t.Run("it builds with on_update setter", func(t *testing.T) {
		c := Floatable{OnUpdate: "set null"}
		assert.Equal(t, "float NOT NULL ON UPDATE set null", c.buildRow())
	})

	t.Run("it builds with comment", func(t *testing.T) {
		c := Floatable{Comment: "test"}
		assert.Equal(t, "float NOT NULL COMMENT 'test'", c.buildRow())
	})

	t.Run("it builds string with all parameters", func(t *testing.T) {
		c := Floatable{
			Type:      "decimal",
			Precision: 10,
			Scale:     2,
			Unsigned:  true,
			Nullable:  true,
			Default:   "100.0",
			OnUpdate:  "set null",
			Comment:   "test",
		}

		assert.Equal(
			t,
			"decimal(10,2) unsigned NULL DEFAULT 100.0 ON UPDATE set null COMMENT 'test'",
			c.buildRow(),
		)
	})
}

func TestTimeable(t *testing.T) {
	t.Run("it builds with default type", func(t *testing.T) {
		c := Timable{}
		assert.Equal(t, "timestamp NOT NULL", c.buildRow())
	})

	t.Run("it builds basic column type", func(t *testing.T) {
		c := Timable{Type: "datetime"}
		assert.Equal(t, "datetime NOT NULL", c.buildRow())
	})

	t.Run("it does not set precision for invalid column type", func(t *testing.T) {
		c := Timable{Type: "date", Precision: 3}
		assert.Equal(t, "date NOT NULL", c.buildRow())
	})

	t.Run("it does not set zero precision", func(t *testing.T) {
		c := Timable{Type: "timestamp", Precision: 0}
		assert.Equal(t, "timestamp NOT NULL", c.buildRow())
	})

	t.Run("it does not set invalid precision", func(t *testing.T) {
		c := Timable{Type: "timestamp", Precision: 7}
		assert.Equal(t, "timestamp NOT NULL", c.buildRow())
	})

	t.Run("it builds with precision", func(t *testing.T) {
		c := Timable{Type: "TIMESTAMP", Precision: 6}
		assert.Equal(t, "TIMESTAMP(6) NOT NULL", c.buildRow())
	})

	t.Run("it builds nullable column type", func(t *testing.T) {
		c := Timable{Nullable: true}
		assert.Equal(t, "timestamp NULL", c.buildRow())
	})

	t.Run("it builds with default value", func(t *testing.T) {
		c := Timable{Default: "CURRENT_TIMESTAMP"}
		assert.Equal(t, "timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP", c.buildRow())
	})

	t.Run("it builds with on_update setter", func(t *testing.T) {
		c := Timable{OnUpdate: "set null"}
		assert.Equal(t, "timestamp NOT NULL ON UPDATE set null", c.buildRow())
	})

	t.Run("it builds with comment", func(t *testing.T) {
		c := Timable{Comment: "test"}
		assert.Equal(t, "timestamp NOT NULL COMMENT 'test'", c.buildRow())
	})

	t.Run("it builds string with all parameters", func(t *testing.T) {
		c := Timable{
			Type:     "datetime",
			Nullable: true,
			Default:  "CURRENT_TIMESTAMP",
			OnUpdate: "CURRENT_TIMESTAMP",
			Comment:  "test",
		}

		assert.Equal(
			t,
			"datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'test'",
			c.buildRow(),
		)
	})
}

func TestString(t *testing.T) {
	t.Run("it builds with default type", func(t *testing.T) {
		c := String{}
		assert.Equal(t, "varchar COLLATE utf8mb4_unicode_ci NOT NULL", c.buildRow())
	})

	t.Run("it builds fixed", func(t *testing.T) {
		c := String{Fixed: true}
		assert.Equal(t, "char COLLATE utf8mb4_unicode_ci NOT NULL", c.buildRow())
	})

	t.Run("it builds with precision", func(t *testing.T) {
		c := String{Precision: 255}
		assert.Equal(t, "varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL", c.buildRow())
	})

	t.Run("it builds with charset", func(t *testing.T) {
		c := String{Charset: "utf8"}
		assert.Equal(t, "varchar CHARACTER SET utf8 NOT NULL", c.buildRow())
	})

	t.Run("it builds with collate", func(t *testing.T) {
		c := String{Collate: "utf8mb4_general_ci"}
		assert.Equal(t, "varchar COLLATE utf8mb4_general_ci NOT NULL", c.buildRow())
	})

	t.Run("it builds nullable column type", func(t *testing.T) {
		c := String{Nullable: true}
		assert.Equal(t, "varchar COLLATE utf8mb4_unicode_ci NULL", c.buildRow())
	})

	t.Run("it builds with default value", func(t *testing.T) {
		c := String{Default: "done"}
		assert.Equal(t, "varchar COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'done'", c.buildRow())
	})

	t.Run("it builds with on_update setter", func(t *testing.T) {
		c := String{OnUpdate: "set null"}
		assert.Equal(t, "varchar COLLATE utf8mb4_unicode_ci NOT NULL ON UPDATE set null", c.buildRow())
	})

	t.Run("it builds with comment", func(t *testing.T) {
		c := String{Comment: "test"}
		assert.Equal(t, "varchar COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'test'", c.buildRow())
	})

	t.Run("it builds string with all parameters", func(t *testing.T) {
		c := String{
			Fixed:     true,
			Precision: 36,
			Nullable:  true,
			Charset:   "utf8mb4",
			Collate:   "utf8mb4_general_ci",
			Default:   "nice",
			OnUpdate:  "set null",
			Comment:   "test",
		}

		assert.Equal(
			t,
			"char(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT 'nice' ON UPDATE set null COMMENT 'test'",
			c.buildRow(),
		)
	})
}

func TestText(t *testing.T) {
	t.Run("it builds with default type", func(t *testing.T) {
		c := Text{}
		assert.Equal(t, "text COLLATE utf8mb4_unicode_ci NOT NULL", c.buildRow())
	})

	t.Run("it builds with prefix", func(t *testing.T) {
		c := Text{Prefix: "medium"}
		assert.Equal(t, "mediumtext COLLATE utf8mb4_unicode_ci NOT NULL", c.buildRow())
	})

	t.Run("it builds blob", func(t *testing.T) {
		c := Text{Blob: true}
		assert.Equal(t, "blob COLLATE utf8mb4_unicode_ci NOT NULL", c.buildRow())
	})

	t.Run("it builds blob with prefix", func(t *testing.T) {
		c := Text{Prefix: "tiny", Blob: true}
		assert.Equal(t, "tinyblob COLLATE utf8mb4_unicode_ci NOT NULL", c.buildRow())
	})

	t.Run("it builds with charset", func(t *testing.T) {
		c := Text{Charset: "utf8"}
		assert.Equal(t, "text CHARACTER SET utf8 NOT NULL", c.buildRow())
	})

	t.Run("it builds with collate", func(t *testing.T) {
		c := Text{Collate: "utf8mb4_general_ci"}
		assert.Equal(t, "text COLLATE utf8mb4_general_ci NOT NULL", c.buildRow())
	})

	t.Run("it builds nullable column type", func(t *testing.T) {
		c := Text{Nullable: true}
		assert.Equal(t, "text COLLATE utf8mb4_unicode_ci NULL", c.buildRow())
	})

	t.Run("it builds with default value", func(t *testing.T) {
		c := Text{Default: "done"}
		assert.Equal(t, "text COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'done'", c.buildRow())
	})

	t.Run("it builds with on_update setter", func(t *testing.T) {
		c := Text{OnUpdate: "set null"}
		assert.Equal(t, "text COLLATE utf8mb4_unicode_ci NOT NULL ON UPDATE set null", c.buildRow())
	})

	t.Run("it builds with comment", func(t *testing.T) {
		c := Text{Comment: "test"}
		assert.Equal(t, "text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'test'", c.buildRow())
	})

	t.Run("it builds string with all parameters", func(t *testing.T) {
		c := Text{
			Prefix:   "long",
			Blob:     true,
			Nullable: true,
			Charset:  "utf8mb4",
			Collate:  "utf8mb4_general_ci",
			Default:  "nice",
			OnUpdate: "set null",
			Comment:  "test",
		}

		assert.Equal(
			t,
			"longblob CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT 'nice' ON UPDATE set null COMMENT 'test'",
			c.buildRow(),
		)
	})
}

func TestJson(t *testing.T) {
	t.Run("it builds with default type", func(t *testing.T) {
		c := JSON{}
		assert.Equal(t, "json NOT NULL", c.buildRow())
	})

	t.Run("it builds nullable column type", func(t *testing.T) {
		c := JSON{Nullable: true}
		assert.Equal(t, "json NULL", c.buildRow())
	})

	t.Run("it builds with default value", func(t *testing.T) {
		c := JSON{Default: "{}"}
		assert.Equal(t, "json NOT NULL DEFAULT '{}'", c.buildRow())
	})

	t.Run("it builds with on_update setter", func(t *testing.T) {
		c := JSON{OnUpdate: "set null"}
		assert.Equal(t, "json NOT NULL ON UPDATE set null", c.buildRow())
	})

	t.Run("it builds with comment", func(t *testing.T) {
		c := JSON{Comment: "test"}
		assert.Equal(t, "json NOT NULL COMMENT 'test'", c.buildRow())
	})

	t.Run("it builds string with all parameters", func(t *testing.T) {
		c := JSON{
			Nullable: true,
			Default:  "{}",
			OnUpdate: "set null",
			Comment:  "test",
		}

		assert.Equal(
			t,
			"json NULL DEFAULT '{}' ON UPDATE set null COMMENT 'test'",
			c.buildRow(),
		)
	})
}

func TestEnum(t *testing.T) {
	t.Run("it builds with default type", func(t *testing.T) {
		c := Enum{}
		assert.Equal(t, "enum('') NOT NULL", c.buildRow())
	})

	t.Run("it builds with multiple flag", func(t *testing.T) {
		c := Enum{Multiple: true}
		assert.Equal(t, "set('') NOT NULL", c.buildRow())
	})

	t.Run("it builds with values", func(t *testing.T) {
		c := Enum{Values: []string{"active", "inactive"}}
		assert.Equal(t, "enum('active', 'inactive') NOT NULL", c.buildRow())
	})

	t.Run("it builds nullable column type", func(t *testing.T) {
		c := Enum{Nullable: true}
		assert.Equal(t, "enum('') NULL", c.buildRow())
	})

	t.Run("it builds with default value", func(t *testing.T) {
		c := Enum{Default: "valid"}
		assert.Equal(t, "enum('') NOT NULL DEFAULT 'valid'", c.buildRow())
	})

	t.Run("it builds with on_update setter", func(t *testing.T) {
		c := Enum{OnUpdate: "set null"}
		assert.Equal(t, "enum('') NOT NULL ON UPDATE set null", c.buildRow())
	})

	t.Run("it builds with comment", func(t *testing.T) {
		c := Enum{Comment: "test"}
		assert.Equal(t, "enum('') NOT NULL COMMENT 'test'", c.buildRow())
	})

	t.Run("it builds string with all parameters", func(t *testing.T) {
		c := Enum{
			Multiple: true,
			Values:   []string{"male", "female", "other"},
			Nullable: true,
			Default:  "male,female",
			OnUpdate: "set null",
			Comment:  "test",
		}

		assert.Equal(
			t,
			"set('male', 'female', 'other') NULL DEFAULT 'male,female' ON UPDATE set null COMMENT 'test'",
			c.buildRow(),
		)
	})
}

func TestBit(t *testing.T) {
	t.Run("it builds basic column type", func(t *testing.T) {
		c := Bit{}
		assert.Equal(t, "bit NOT NULL", c.buildRow())
	})

	t.Run("it builds with precision", func(t *testing.T) {
		c := Bit{Precision: 20}
		assert.Equal(t, "bit(20) NOT NULL", c.buildRow())
	})

	t.Run("it builds nullable column type", func(t *testing.T) {
		c := Bit{Nullable: true}
		assert.Equal(t, "bit NULL", c.buildRow())
	})

	t.Run("it builds with default value", func(t *testing.T) {
		c := Bit{Default: "1"}
		assert.Equal(t, "bit NOT NULL DEFAULT 1", c.buildRow())
	})

	t.Run("it builds with on_update setter", func(t *testing.T) {
		c := Bit{OnUpdate: "set null"}
		assert.Equal(t, "bit NOT NULL ON UPDATE set null", c.buildRow())
	})

	t.Run("it builds with comment", func(t *testing.T) {
		c := Bit{Comment: "test"}
		assert.Equal(t, "bit NOT NULL COMMENT 'test'", c.buildRow())
	})

	t.Run("it builds string with all parameters", func(t *testing.T) {
		c := Bit{
			Precision: 10,
			Nullable:  true,
			Default:   "0",
			OnUpdate:  "set null",
			Comment:   "test",
		}

		assert.Equal(
			t,
			"bit(10) NULL DEFAULT 0 ON UPDATE set null COMMENT 'test'",
			c.buildRow(),
		)
	})
}

func TestBinary(t *testing.T) {
	t.Run("it builds with default type", func(t *testing.T) {
		c := Binary{}
		assert.Equal(t, "varbinary NOT NULL", c.buildRow())
	})

	t.Run("it builds fixed", func(t *testing.T) {
		c := Binary{Fixed: true}
		assert.Equal(t, "binary NOT NULL", c.buildRow())
	})

	t.Run("it builds with precision", func(t *testing.T) {
		c := Binary{Precision: 255}
		assert.Equal(t, "varbinary(255) NOT NULL", c.buildRow())
	})

	t.Run("it builds nullable column type", func(t *testing.T) {
		c := Binary{Nullable: true}
		assert.Equal(t, "varbinary NULL", c.buildRow())
	})

	t.Run("it builds with default value", func(t *testing.T) {
		c := Binary{Default: "1"}
		assert.Equal(t, "varbinary NOT NULL DEFAULT 1", c.buildRow())
	})

	t.Run("it builds with on_update setter", func(t *testing.T) {
		c := Binary{OnUpdate: "set null"}
		assert.Equal(t, "varbinary NOT NULL ON UPDATE set null", c.buildRow())
	})

	t.Run("it builds with comment", func(t *testing.T) {
		c := Binary{Comment: "test"}
		assert.Equal(t, "varbinary NOT NULL COMMENT 'test'", c.buildRow())
	})

	t.Run("it builds string with all parameters", func(t *testing.T) {
		c := Binary{
			Fixed:     true,
			Precision: 36,
			Nullable:  true,
			Default:   "1",
			OnUpdate:  "set null",
			Comment:   "test",
		}

		assert.Equal(
			t,
			"binary(36) NULL DEFAULT 1 ON UPDATE set null COMMENT 'test'",
			c.buildRow(),
		)
	})
}

func TestBuildDefaultForString(t *testing.T) {
	t.Run("it returns an empty string if default value is missing", func(t *testing.T) {
		got := buildDefaultForString("")

		assert.Equal(t, "", got)
	})

	t.Run("it builds default with expression", func(t *testing.T) {
		got := buildDefaultForString("(UUID())")
		want := " DEFAULT (UUID())"

		assert.Equal(t, want, got)
	})

	t.Run("it builds normal default", func(t *testing.T) {
		got := buildDefaultForString("value")
		want := " DEFAULT 'value'"

		assert.Equal(t, want, got)
	})
}
