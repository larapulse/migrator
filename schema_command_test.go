package migrator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCommand string

func (c testCommand) ToSQL() string {
	return "Do action on " + string(c)
}

func TestCreateTableCommand(t *testing.T) {
	t.Run("it returns empty string when table name missing", func(t *testing.T) {
		tb := Table{}
		c := createTableCommand{tb}

		assert.Equal(t, "", c.ToSQL())
	})

	t.Run("it renders default table", func(t *testing.T) {
		tb := Table{Name: "test"}
		c := createTableCommand{tb}

		assert.Equal(
			t,
			"CREATE TABLE `test` (`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci",
			c.ToSQL(),
		)
	})

	t.Run("it renders columns", func(t *testing.T) {
		tb := Table{
			Name: "test",
			columns: []column{
				{"test", testColumnType("random thing")},
				{"random", testColumnType("another thing")},
			},
		}
		c := createTableCommand{tb}

		assert.Equal(
			t,
			"CREATE TABLE `test` (`test` random thing, `random` another thing) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci",
			c.ToSQL(),
		)
	})

	t.Run("it renders indexes", func(t *testing.T) {
		tb := Table{
			Name: "test",
			indexes: []Key{
				{Name: "idx_rand", Columns: []string{"id"}},
				{Columns: []string{"id", "name"}},
			},
		}
		c := createTableCommand{tb}

		assert.Equal(
			t,
			strings.Join([]string{
				"CREATE TABLE `test` (",
				"`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT, ",
				"KEY `idx_rand` (`id`), KEY (`id`, `name`)",
				") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci",
			}, ""),
			c.ToSQL(),
		)
	})

	t.Run("it renders foreigns", func(t *testing.T) {
		tb := Table{
			Name: "test",
			foreigns: []Foreign{
				{Key: "idx_foreign", Column: "test_id", Reference: "id", On: "tests"},
				{Key: "foreign_idx", Column: "random_id", Reference: "id", On: "randoms"},
			},
		}
		c := createTableCommand{tb}

		assert.Equal(
			t,
			strings.Join([]string{
				"CREATE TABLE `test` (",
				"`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT, ",
				"CONSTRAINT `idx_foreign` FOREIGN KEY (`test_id`) REFERENCES `tests` (`id`), ",
				"CONSTRAINT `foreign_idx` FOREIGN KEY (`random_id`) REFERENCES `randoms` (`id`)",
				") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci",
			}, ""),
			c.ToSQL(),
		)
	})

	t.Run("it renders engine", func(t *testing.T) {
		tb := Table{Name: "test", Engine: "MyISAM"}
		c := createTableCommand{tb}

		assert.Equal(
			t,
			"CREATE TABLE `test` (`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci",
			c.ToSQL(),
		)
	})

	t.Run("it renders charset and collation", func(t *testing.T) {
		tb := Table{Name: "test", Charset: "rand", Collation: "random_io"}
		c := createTableCommand{tb}

		assert.Equal(
			t,
			"CREATE TABLE `test` (`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT) ENGINE=InnoDB DEFAULT CHARSET=rand COLLATE=random_io",
			c.ToSQL(),
		)
	})

	t.Run("it renders charset and manually add collation", func(t *testing.T) {
		tb := Table{Name: "test", Charset: "utf8"}
		c := createTableCommand{tb}

		assert.Equal(
			t,
			"CREATE TABLE `test` (`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci",
			c.ToSQL(),
		)
	})

	t.Run("it renders collation and manually add charset", func(t *testing.T) {
		tb := Table{Name: "test", Collation: "utf8_general_ci"}
		c := createTableCommand{tb}

		assert.Equal(
			t,
			"CREATE TABLE `test` (`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci",
			c.ToSQL(),
		)
	})

	t.Run("it renders all together", func(t *testing.T) {
		tb := Table{
			Name: "test",
			columns: []column{
				{"test", testColumnType("random thing")},
				{"random", testColumnType("another thing")},
			},
			indexes: []Key{
				{Name: "idx_rand", Columns: []string{"id"}},
				{Columns: []string{"id", "name"}},
			},
			foreigns: []Foreign{
				{Key: "idx_foreign", Column: "test_id", Reference: "id", On: "tests"},
				{Key: "foreign_idx", Column: "random_id", Reference: "id", On: "randoms"},
			},
			Engine:    "MyISAM",
			Charset:   "rand",
			Collation: "random_io",
		}
		c := createTableCommand{tb}

		assert.Equal(
			t,
			strings.Join([]string{
				"CREATE TABLE `test` (",
				"`test` random thing, `random` another thing, ",
				"KEY `idx_rand` (`id`), KEY (`id`, `name`), ",
				"CONSTRAINT `idx_foreign` FOREIGN KEY (`test_id`) REFERENCES `tests` (`id`), ",
				"CONSTRAINT `foreign_idx` FOREIGN KEY (`random_id`) REFERENCES `randoms` (`id`)",
				") ENGINE=MyISAM DEFAULT CHARSET=rand COLLATE=random_io",
			}, ""),
			c.ToSQL(),
		)
	})
}

func TestDropTableCommand(t *testing.T) {
	t.Run("it drops table", func(t *testing.T) {
		c := dropTableCommand{"test", false, ""}
		assert.Equal(t, "DROP TABLE `test`", c.ToSQL())
	})

	t.Run("it drops table if exists", func(t *testing.T) {
		c := dropTableCommand{"test", true, ""}
		assert.Equal(t, "DROP TABLE IF EXISTS `test`", c.ToSQL())
	})

	t.Run("it drops table with cascade flag", func(t *testing.T) {
		c := dropTableCommand{"test", false, "cascade"}
		assert.Equal(t, "DROP TABLE `test` CASCADE", c.ToSQL())
	})

	t.Run("it drops table if exists with restrict flag", func(t *testing.T) {
		c := dropTableCommand{"test", true, "restrict"}
		assert.Equal(t, "DROP TABLE IF EXISTS `test` RESTRICT", c.ToSQL())
	})
}

func TestRenameTableCommand(t *testing.T) {
	c := renameTableCommand{"from", "to"}

	assert.Equal(t, "RENAME TABLE `from` TO `to`", c.ToSQL())
}

func TestAlterTableCommand(t *testing.T) {
	t.Run("it returns an empty command if table name is missing", func(t *testing.T) {
		c := alterTableCommand{pool: TableCommands{testCommand("test")}}

		assert.Equal(t, "", c.ToSQL())
	})

	t.Run("it returns an empty command if pool is empty", func(t *testing.T) {
		c := alterTableCommand{name: "test"}

		assert.Equal(t, "", c.ToSQL())
	})

	t.Run("it renders command with one alter sub-command", func(t *testing.T) {
		c := alterTableCommand{name: "test", pool: TableCommands{testCommand("test")}}

		assert.Equal(t, "ALTER TABLE `test` Do action on test", c.ToSQL())
	})

	t.Run("it renders command with multiple alter sub-command", func(t *testing.T) {
		c := alterTableCommand{
			name: "test",
			pool: TableCommands{testCommand("test"), testCommand("bang")},
		}

		assert.Equal(t, "ALTER TABLE `test` Do action on test, Do action on bang", c.ToSQL())
	})
}
