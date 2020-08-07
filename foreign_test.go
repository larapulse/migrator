package migrator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForeigns(t *testing.T) {
	t.Run("it returns empty on empty keys", func(t *testing.T) {
		f := foreigns{Foreign{}}

		assert.Equal(t, "", f.render())
	})

	t.Run("it renders row from one foreign", func(t *testing.T) {
		f := foreigns{Foreign{Key: "idx_foreign", Column: "test_id", Reference: "id", On: "tests"}}

		assert.Equal(t, "CONSTRAINT `idx_foreign` FOREIGN KEY (`test_id`) REFERENCES `tests` (`id`)", f.render())
	})

	t.Run("it renders row from multiple foreigns", func(t *testing.T) {
		f := foreigns{
			Foreign{Key: "idx_foreign", Column: "test_id", Reference: "id", On: "tests"},
			Foreign{Key: "foreign_idx", Column: "random_id", Reference: "id", On: "randoms"},
		}

		assert.Equal(
			t,
			"CONSTRAINT `idx_foreign` FOREIGN KEY (`test_id`) REFERENCES `tests` (`id`), CONSTRAINT `foreign_idx` FOREIGN KEY (`random_id`) REFERENCES `randoms` (`id`)",
			f.render(),
		)
	})
}

func TestForeign(t *testing.T) {
	t.Run("it builds base constraint", func(t *testing.T) {
		f := Foreign{Key: "foreign_idx", Column: "test_id", Reference: "id", On: "tests"}

		assert.Equal(t, "CONSTRAINT `foreign_idx` FOREIGN KEY (`test_id`) REFERENCES `tests` (`id`)", f.render())
	})

	t.Run("it builds contraint with on_update", func(t *testing.T) {
		f := Foreign{Key: "foreign_idx", Column: "test_id", Reference: "id", On: "tests", OnUpdate: "no action"}

		assert.Equal(t, "CONSTRAINT `foreign_idx` FOREIGN KEY (`test_id`) REFERENCES `tests` (`id`) ON UPDATE NO ACTION", f.render())
	})

	t.Run("it builds contraint without invalid on_update", func(t *testing.T) {
		f := Foreign{Key: "foreign_idx", Column: "test_id", Reference: "id", On: "tests", OnUpdate: "null"}

		assert.Equal(t, "CONSTRAINT `foreign_idx` FOREIGN KEY (`test_id`) REFERENCES `tests` (`id`)", f.render())
	})

	t.Run("it builds contraint with on_update", func(t *testing.T) {
		f := Foreign{Key: "foreign_idx", Column: "test_id", Reference: "id", On: "tests", OnDelete: "set default"}

		assert.Equal(t, "CONSTRAINT `foreign_idx` FOREIGN KEY (`test_id`) REFERENCES `tests` (`id`) ON DELETE SET DEFAULT", f.render())
	})

	t.Run("it builds contraint without invalid on_update", func(t *testing.T) {
		f := Foreign{Key: "foreign_idx", Column: "test_id", Reference: "id", On: "tests", OnDelete: "default"}

		assert.Equal(t, "CONSTRAINT `foreign_idx` FOREIGN KEY (`test_id`) REFERENCES `tests` (`id`)", f.render())
	})

	t.Run("it builds full contraint", func(t *testing.T) {
		f := Foreign{Key: "foreign_idx", Column: "test_id", Reference: "id", On: "tests", OnUpdate: "cascade", OnDelete: "restrict"}

		assert.Equal(t, "CONSTRAINT `foreign_idx` FOREIGN KEY (`test_id`) REFERENCES `tests` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE", f.render())
	})
}

func TestBuildForeignIndexNameOnTable(t *testing.T) {
	assert.Equal(t, "table_test_foreign", BuildForeignNameOnTable("table", "test"))
}
