package migrator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeys(t *testing.T) {
	t.Run("it returns empty on empty keys", func(t *testing.T) {
		k := keys{key{}}

		assert.Equal(t, "", k.render())
	})

	t.Run("it renders row from one key", func(t *testing.T) {
		k := keys{key{columns: []string{"test_id"}}}

		assert.Equal(t, "KEY (`test_id`)", k.render())
	})

	t.Run("it renders row from multiple keys", func(t *testing.T) {
		k := keys{
			key{columns: []string{"test_id"}},
			key{columns: []string{"random_id"}},
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
		k := key{}

		assert.Equal(t, "", k.render())
	})

	t.Run("it skips type if it is not in valid list", func(t *testing.T) {
		k := key{typ: "random", columns: []string{"test_id"}}

		assert.Equal(t, "KEY (`test_id`)", k.render())
	})

	t.Run("it renders with type", func(t *testing.T) {
		k := key{typ: "primary", columns: []string{"test_id"}}

		assert.Equal(t, "PRIMARY KEY (`test_id`)", k.render())
	})

	t.Run("it renders with multiple columns", func(t *testing.T) {
		k := key{typ: "unique", columns: []string{"test_id", "random_id"}}

		assert.Equal(t, "UNIQUE KEY (`test_id`, `random_id`)", k.render())
	})

	t.Run("it renders with name", func(t *testing.T) {
		k := key{name: "random_idx", columns: []string{"test_id"}}

		assert.Equal(t, "KEY `random_idx` (`test_id`)", k.render())
	})
}
