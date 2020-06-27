package migrator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForeigns(t *testing.T) {
	t.Run("it returns empty on empty keys", func(t *testing.T) {
		f := foreigns{foreign{}}

		assert.Equal(t, "", f.render())
	})

	t.Run("it renders row from one foreign", func(t *testing.T) {
		f := foreigns{foreign{key: "idx_foreign", column: "test_id", reference: "id", on: "tests"}}

		assert.Equal(t, "CONSTRAINT `idx_foreign` FOREIGN KEY (`test_id`) REFERENCES `tests` (`id`)", f.render())
	})

	t.Run("it renders row from multiple foreigns", func(t *testing.T) {
		f := foreigns{
			foreign{key: "idx_foreign", column: "test_id", reference: "id", on: "tests"},
			foreign{key: "foreign_idx", column: "random_id", reference: "id", on: "randoms"},
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
		f := foreign{key: "foreign_idx", column: "test_id", reference: "id", on: "tests"}

		assert.Equal(t, "CONSTRAINT `foreign_idx` FOREIGN KEY (`test_id`) REFERENCES `tests` (`id`)", f.render())
	})

	t.Run("it builds contraint with on_update", func(t *testing.T) {
		f := foreign{key: "foreign_idx", column: "test_id", reference: "id", on: "tests", onUpdate: "no action"}

		assert.Equal(t, "CONSTRAINT `foreign_idx` FOREIGN KEY (`test_id`) REFERENCES `tests` (`id`) ON UPDATE NO ACTION", f.render())
	})

	t.Run("it builds contraint without invalid on_update", func(t *testing.T) {
		f := foreign{key: "foreign_idx", column: "test_id", reference: "id", on: "tests", onUpdate: "null"}

		assert.Equal(t, "CONSTRAINT `foreign_idx` FOREIGN KEY (`test_id`) REFERENCES `tests` (`id`)", f.render())
	})

	t.Run("it builds contraint with on_update", func(t *testing.T) {
		f := foreign{key: "foreign_idx", column: "test_id", reference: "id", on: "tests", onDelete: "set default"}

		assert.Equal(t, "CONSTRAINT `foreign_idx` FOREIGN KEY (`test_id`) REFERENCES `tests` (`id`) ON DELETE SET DEFAULT", f.render())
	})

	t.Run("it builds contraint without invalid on_update", func(t *testing.T) {
		f := foreign{key: "foreign_idx", column: "test_id", reference: "id", on: "tests", onDelete: "default"}

		assert.Equal(t, "CONSTRAINT `foreign_idx` FOREIGN KEY (`test_id`) REFERENCES `tests` (`id`)", f.render())
	})

	t.Run("it builds full contraint", func(t *testing.T) {
		f := foreign{key: "foreign_idx", column: "test_id", reference: "id", on: "tests", onUpdate: "cascade", onDelete: "restrict"}

		assert.Equal(t, "CONSTRAINT `foreign_idx` FOREIGN KEY (`test_id`) REFERENCES `tests` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE", f.render())
	})
}
