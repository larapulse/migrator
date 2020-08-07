package migrator

import (
	"fmt"
	"strings"
)

type foreigns []Foreign

func (f foreigns) render() string {
	values := []string{}

	for _, foreign := range f {
		values = append(values, foreign.render())
	}

	return strings.Join(values, ", ")
}

// Foreign represents an instance to handle foreign key interactions
type Foreign struct {
	Key       string
	Column    string
	Reference string // reference field
	On        string // reference table
	OnUpdate  string
	OnDelete  string
}

func (f Foreign) render() string {
	if f.Key == "" || f.Column == "" || f.On == "" || f.Reference == "" {
		return ""
	}

	sql := fmt.Sprintf("CONSTRAINT `%s` FOREIGN KEY (`%s`) REFERENCES `%s` (`%s`)", f.Key, f.Column, f.On, f.Reference)
	if referenceOptions.has(strings.ToUpper(f.OnDelete)) {
		sql += " ON DELETE " + strings.ToUpper(f.OnDelete)
	}
	if referenceOptions.has(strings.ToUpper(f.OnUpdate)) {
		sql += " ON UPDATE " + strings.ToUpper(f.OnUpdate)
	}

	return sql
}

// BuildForeignNameOnTable builds a name for the foreign key on the table
func BuildForeignNameOnTable(table string, column string) string {
	return table + "_" + column + "_foreign"
}

var referenceOptions = list{"SET NULL", "CASCADE", "RESTRICT", "NO ACTION", "SET DEFAULT"}

type list []string

func (l list) has(value string) bool {
	for _, item := range l {
		if item == value {
			return true
		}
	}

	return false
}
