package migrator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeys(t *testing.T) {
	t.Run("it returns empty on empty keys", func(t *testing.T) {
		k := keys{Key{}}

		assert.Equal(t, "", k.render())
	})

	t.Run("it renders row from one key", func(t *testing.T) {
		k := keys{Key{Columns: []string{"test_id"}}}

		assert.Equal(t, "KEY (`test_id`)", k.render())
	})

	t.Run("it renders row from multiple keys", func(t *testing.T) {
		k := keys{
			Key{Columns: []string{"test_id"}},
			Key{Columns: []string{"random_id"}},
		}

		assert.Equal(
			t,
			"KEY (`test_id`), KEY (`random_id`)",
			k.render(),
		)
	})
}

func TestKey(t *testing.T) {
	t.Run("it returns empty on empty keys", func(t *testing.T) {
		k := Key{}

		assert.Equal(t, "", k.render())
	})

	t.Run("it skips type if it is not in valid list", func(t *testing.T) {
		k := Key{Type: "random", Columns: []string{"test_id"}}

		assert.Equal(t, "KEY (`test_id`)", k.render())
	})

	t.Run("it renders with type", func(t *testing.T) {
		k := Key{Type: "primary", Columns: []string{"test_id"}}

		assert.Equal(t, "PRIMARY KEY (`test_id`)", k.render())
	})

	t.Run("it renders with multiple columns", func(t *testing.T) {
		k := Key{Type: "unique", Columns: []string{"test_id", "random_id"}}

		assert.Equal(t, "UNIQUE KEY (`test_id`, `random_id`)", k.render())
	})

	t.Run("it renders with name", func(t *testing.T) {
		k := Key{Name: "random_idx", Columns: []string{"test_id"}}

		assert.Equal(t, "KEY `random_idx` (`test_id`)", k.render())
	})
}

func TestBuildUniqueIndexName(t *testing.T) {
	t.Run("It builds name from one column", func(t *testing.T) {
		assert.Equal(t, "table_test_unique", BuildUniqueKeyNameOnTable("table", "test"))
	})

	t.Run("it builds name from multiple columns", func(t *testing.T) {
		assert.Equal(t, "table_test_again_unique", BuildUniqueKeyNameOnTable("table", "test", "again"))
	})
}
