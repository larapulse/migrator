// Package migrator represents MySQL database migrator
package migrator

import "strings"

type keys []key

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

type key struct {
	name    string
	typ     string // primary, unique
	columns []string
}

var keyTypes = list{"PRIMARY", "UNIQUE"}

func (k key) render() string {
	if len(k.columns) == 0 {
		return ""
	}

	sql := ""
	if keyTypes.has(strings.ToUpper(k.typ)) {
		sql += strings.ToUpper(k.typ) + " "
	}

	sql += "KEY"

	if k.name != "" {
		sql += " `" + k.name + "`"
	}

	sql += " (`" + strings.Join(k.columns, "`, `") + "`)"

	return sql
}
