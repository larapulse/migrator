package migrator

import "strings"

type keys []Key

func (k keys) render() string {
	values := []string{}

	for _, key := range k {
		value := key.render()
		if value != "" {
			values = append(values, value)
		}
	}

	return strings.Join(values, ", ")
}

// Key represents an instance to handle key (index) interactions
type Key struct {
	Name    string
	Type    string // primary, unique
	Columns []string
}

var keyTypes = list{"PRIMARY", "UNIQUE"}

func (k Key) render() string {
	if len(k.Columns) == 0 {
		return ""
	}

	sql := ""
	if keyTypes.has(strings.ToUpper(k.Type)) {
		sql += strings.ToUpper(k.Type) + " "
	}

	sql += "KEY"

	if k.Name != "" {
		sql += " `" + k.Name + "`"
	}

	sql += " (`" + strings.Join(k.Columns, "`, `") + "`)"

	return sql
}

// BuildUniqueKeyNameOnTable builds a name for the foreign key on the table
func BuildUniqueKeyNameOnTable(table string, columns ...string) string {
	return table + "_" + strings.Join(columns, "_") + "_unique"
}
