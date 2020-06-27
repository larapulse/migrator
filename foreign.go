// Package migrator represents MySQL database migrator
package migrator

import (
	"fmt"
	"strings"
)

type foreigns []foreign

func (f foreigns) render() string {
	values := []string{}

	for _, foreign := range f {
		values = append(values, foreign.render())
	}

	return strings.Join(values, ", ")
}

type foreign struct {
	key       string
	column    string
	reference string // reference field
	on        string // reference table
	onUpdate  string
	onDelete  string
}

func (f foreign) render() string {
	if f.key == "" || f.column == "" || f.on == "" || f.reference == "" {
		return ""
	}

	sql := fmt.Sprintf("CONSTRAINT `%s` FOREIGN KEY (`%s`) REFERENCES `%s` (`%s`)", f.key, f.column, f.on, f.reference)
	if referenceOptions.has(strings.ToUpper(f.onDelete)) {
		sql += " ON DELETE " + strings.ToUpper(f.onDelete)
	}
	if referenceOptions.has(strings.ToUpper(f.onUpdate)) {
		sql += " ON UPDATE " + strings.ToUpper(f.onUpdate)
	}

	return sql
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
