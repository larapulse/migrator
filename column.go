package migrator

import (
	"fmt"
	"strconv"
	"strings"
)

type columns []column

func (c columns) render() string {
	rows := []string{}

	for _, item := range c {
		rows = append(rows, fmt.Sprintf("`%s` %s", item.field, item.definition.buildRow()))
	}

	return strings.Join(rows, ", ")
}

type column struct {
	field      string
	definition columnType
}

type columnType interface {
	buildRow() string
}

// Integer represents an integer value in DB: {tiny,small,medium,big}int
//
// Default migrator.Integer will build a sql row: `int NOT NULL`
//
// Examples:
//		tinyint		➡️ migrator.Integer{Prefix: "tiny", Unsigned: true, Precision: 1, Default: "0"}
//			↪️ tinyint(1) unsigned NOT NULL DEFAULT 0
//		int			➡️ migrator.Integer{Nullable: true, OnUpdate: "set null", Comment: "nullable counter"}
//			↪️ int NULL ON UPDATE set null COMMENT 'nullable counter'
//		mediumint	➡️ migrator.Integer{Prefix: "medium", Precision: "255"}
//			↪️ mediumint(255) NOT NULL
//		bigint		➡️ migrator.Integer{Prefix: "big", Unsigned: true, Precision: "255", Autoincrement: true}
//			↪️ bigint(255) unsigned NOT NULL AUTO_INCREMENT
type Integer struct {
	Default  string
	Nullable bool
	Comment  string
	OnUpdate string

	Prefix        string // tiny, small, medium, big
	Unsigned      bool
	Precision     uint16
	Autoincrement bool
}

func (i Integer) buildRow() string {
	sql := i.Prefix + "int"
	if i.Precision > 0 {
		sql += fmt.Sprintf("(%s)", strconv.Itoa(int(i.Precision)))
	}

	if i.Unsigned {
		sql += " unsigned"
	}

	if i.Nullable {
		sql += " NULL"
	} else {
		sql += " NOT NULL"
	}

	if i.Default != "" {
		sql += " DEFAULT " + i.Default
	}

	if i.Autoincrement {
		sql += " AUTO_INCREMENT"
	}

	if i.OnUpdate != "" {
		sql += " ON UPDATE " + i.OnUpdate
	}

	if i.Comment != "" {
		sql += fmt.Sprintf(" COMMENT '%s'", i.Comment)
	}

	return sql
}

// Floatable represents a number with a floating point in DB:
// `float`, `double` or `decimal`
//
// Default migrator.Floatable will build a sql row: `float NOT NULL`
//
// Examples:
//		float	➡️ migrator.Floatable{Precision: 2, Nullable: true}
//			↪️ float(2) NULL
//		real	➡️ migrator.Floatable{Type: "real", Precision: 5, Scale: 2}
//			↪️ real(5,2) NOT NULL
//		double	➡️ migrator.Floatable{Type: "double", Scale: 2, Unsigned: true}
//			↪️ double(0,2) unsigned NOT NULL
//		decimal	➡️ migrator.Floatable{Type: "decimal", Precision: 15, Scale: 2, OnUpdate: "0.0", Comment: "money"}
//			↪️ decimal(15,2) NOT NULL ON UPDATE 0.0 COMMENT 'money'
//		numeric	➡️ migrator.Floatable{Type: "numeric", Default: "0.0"}
//			↪️ numeric NOT NULL DEFAULT 0.0
type Floatable struct {
	Default  string
	Nullable bool
	Comment  string
	OnUpdate string

	Type      string // float, real, double, decimal, numeric
	Unsigned  bool
	Precision uint16
	Scale     uint16
}

func (f Floatable) buildRow() string {
	sql := f.Type

	if sql == "" {
		sql = "float"
	}

	if f.Scale > 0 {
		sql += fmt.Sprintf("(%s,%s)", strconv.Itoa(int(f.Precision)), strconv.Itoa(int(f.Scale)))
	} else if f.Precision > 0 {
		sql += fmt.Sprintf("(%s)", strconv.Itoa(int(f.Precision)))
	}

	if f.Unsigned {
		sql += " unsigned"
	}

	if f.Nullable {
		sql += " NULL"
	} else {
		sql += " NOT NULL"
	}

	if f.Default != "" {
		sql += " DEFAULT " + f.Default
	}

	if f.OnUpdate != "" {
		sql += " ON UPDATE " + f.OnUpdate
	}

	if f.Comment != "" {
		sql += fmt.Sprintf(" COMMENT '%s'", f.Comment)
	}

	return sql
}

// Timable represents DB representation of timable column type:
// `date`, `datetime`, `timestamp`, `time` or `year`
//
// Default migrator.Timable will build a sql row: `timestamp NOT NULL`.
// Precision from 0 to 6 can be set for `datetime`, `timestamp`, `time`.
//
// Examples:
//		date		➡️ migrator.Timable{Type: "date", Nullable: true}
//			↪️ date NULL
//		datetime	➡️ migrator.Timable{Type: "datetime", Precision: 3, Default: "CURRENT_TIMESTAMP"}
//			↪️ datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP
//		timestamp	➡️ migrator.Timable{Default: "CURRENT_TIMESTAMP(6)", OnUpdate: "CURRENT_TIMESTAMP(6)"}
//			↪️ timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6)
//		time		➡️ migrator.Timable{Type: "time", Comment: "meeting time"}
//			↪️ time NOT NULL COMMENT 'meeting time'
//		year		➡️ migrator.Timable{Type: "year", Nullable: true}
//			↪️ year NULL
type Timable struct {
	Default  string
	Nullable bool
	Comment  string
	OnUpdate string

	Type      string // date, time, datetime, timestamp, year
	Precision uint16
}

func (t Timable) buildRow() string {
	sql := t.Type

	if sql == "" {
		sql = "timestamp"
	}

	validForPrecision := list{"time", "datetime", "timestamp"}
	columnType := strings.ToLower(sql)
	if t.Precision > 0 && t.Precision <= 6 && validForPrecision.has(columnType) {
		sql += fmt.Sprintf("(%s)", strconv.Itoa(int(t.Precision)))
	}

	if t.Nullable {
		sql += " NULL"
	} else {
		sql += " NOT NULL"
	}

	if t.Default != "" {
		sql += " DEFAULT " + t.Default
	}

	if t.OnUpdate != "" {
		sql += " ON UPDATE " + t.OnUpdate
	}

	if t.Comment != "" {
		sql += fmt.Sprintf(" COMMENT '%s'", t.Comment)
	}

	return sql
}

// String represents basic DB string column type: `char` or `varchar`
//
// Default migrator.String will build a sql row: `varchar COLLATE utf8mb4_unicode_ci NOT NULL`
//
// Examples:
//		char	➡️ migrator.String{Fixed: true, Precision: 36, Nullable: true, OnUpdate: "set null", Comment: "uuid"}
//			↪️ char(36) COLLATE utf8mb4_unicode_ci NULL ON UPDATE set null COMMENT 'uuid'
//		varchar	➡️ migrator.String{Precision: 255, Default: "active", Charset: "utf8mb4", Collate: "utf8mb4_general_ci"}
//			↪️ varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'active'
type String struct {
	Default  string
	Nullable bool
	Comment  string
	OnUpdate string

	Charset string
	Collate string

	Fixed     bool // char for fixed, otherwise varchar
	Precision uint16
}

func (s String) buildRow() string {
	sql := ""

	if !s.Fixed {
		sql += "var"
	}

	sql += "char"

	if s.Precision > 0 {
		sql += fmt.Sprintf("(%s)", strconv.Itoa(int(s.Precision)))
	}

	if s.Charset != "" {
		sql += " CHARACTER SET " + s.Charset
	}

	if s.Collate != "" {
		sql += " COLLATE " + s.Collate
	} else if s.Charset == "" {
		// use default
		sql += " COLLATE utf8mb4_unicode_ci"
	}

	if s.Nullable {
		sql += " NULL"
	} else {
		sql += " NOT NULL"
	}

	sql += buildDefaultForString(s.Default)

	if s.OnUpdate != "" {
		sql += " ON UPDATE " + s.OnUpdate
	}

	if s.Comment != "" {
		sql += fmt.Sprintf(" COMMENT '%s'", s.Comment)
	}

	return sql
}

// Text represents long text column type represented in DB as:
//  - {tiny,medium,long}text
//  - {tiny,medium,long}blob
//
// Default migrator.Text will build a sql row: `text COLLATE utf8mb4_unicode_ci NOT NULL`
//
// Examples:
//		tinytext	➡️ migrator.Text{Prefix: "tiny"}
//			↪️ tinytext COLLATE utf8mb4_unicode_ci NOT NULL
//		text		➡️ migrator.Text{Nullable: true, OnUpdate: "set null", Comment: "write your comment here"}
//			↪️ text COLLATE utf8mb4_unicode_ci NULL ON UPDATE set null COMMENT 'write your comment here'
//		mediumtext	➡️ migrator.Text{Prefix: "medium"}
//			↪️ mediumtext COLLATE utf8mb4_unicode_ci NOT NULL
//		longtext	➡️ migrator.Text{Prefix: "long", Default: "write you text", Charset: "utf8mb4", Collate: "utf8mb4_general_ci"}
//			↪️ longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'write you text'
//		tinyblob	➡️ migrator.Text{Prefix: "tiny", Blob: true}
//			↪️ tinyblob NOT NULL
//		blob		➡️ migrator.Text{Blob: true}
//			↪️ blob NOT NULL
//		mediumblob	➡️ migrator.Text{Prefix: "medium", Blob: true}
//			↪️ mediumblob NOT NULL
//		longblob	➡️ migrator.Text{Prefix: "long", Blob: true}
//			↪️ longblob NOT NULL
type Text struct {
	Default  string
	Nullable bool
	Comment  string
	OnUpdate string

	Charset string
	Collate string

	Prefix string // tiny, medium, long
	Blob   bool   // for binary
}

func (t Text) buildRow() string {
	sql := t.Prefix

	if t.Blob {
		sql += "blob"
	} else {
		sql += "text"
	}

	if t.Charset != "" {
		sql += " CHARACTER SET " + t.Charset
	}

	if t.Collate != "" {
		sql += " COLLATE " + t.Collate
	} else if t.Charset == "" && t.Blob == false {
		// use default
		sql += " COLLATE utf8mb4_unicode_ci"
	}

	if t.Nullable {
		sql += " NULL"
	} else {
		sql += " NOT NULL"
	}

	sql += buildDefaultForString(t.Default)

	if t.OnUpdate != "" {
		sql += " ON UPDATE " + t.OnUpdate
	}

	if t.Comment != "" {
		sql += fmt.Sprintf(" COMMENT '%s'", t.Comment)
	}

	return sql
}

// JSON represents DB column type `json`
//
// Default migrator.JSON will build a sql row: `json NOT NULL`
//
// Examples:
//		➡️ migrator.JSON{Nullable: true, Comment: "user data"}
//			↪️ json NULL COMMENT 'user data'
//		➡️ migrator.JSON{Default: "{}", OnUpdate: "{}"}
//			↪️ json NOT NULL DEFAULT '{}' ON UPDATE {}
type JSON struct {
	Default  string
	Nullable bool
	Comment  string
	OnUpdate string
}

func (j JSON) buildRow() string {
	sql := "json"

	if j.Nullable {
		sql += " NULL"
	} else {
		sql += " NOT NULL"
	}

	sql += buildDefaultForString(j.Default)

	if j.OnUpdate != "" {
		sql += " ON UPDATE " + j.OnUpdate
	}

	if j.Comment != "" {
		sql += fmt.Sprintf(" COMMENT '%s'", j.Comment)
	}

	return sql
}

// Enum represents choosable value. In the database represented by: `enum` or `set`
//
// Default migrator.Enum will build a sql row: `enum('') NOT NULL`
//
// Examples:
//		enum	➡️ migrator.Enum{Values: []string{"on", "off"}, Default: "off", Nullable: true, OnUpdate: "set null"}
//			↪️ enum('on', 'off') NULL DEFAULT 'off' ON UPDATE set null
//		set		➡️ migrator.Enum{Values: []string{"1", "2", "3"}, Comment: "options"}
//			↪️ set('1', '2', '3') NOT NULL COMMENT 'options'
type Enum struct {
	Default  string
	Nullable bool
	Comment  string
	OnUpdate string

	Values   []string
	Multiple bool // "set", otherwise "enum"
}

func (e Enum) buildRow() string {
	sql := ""

	if e.Multiple {
		sql += "set"
	} else {
		sql += "enum"
	}

	sql += "('" + strings.Join(e.Values, "', '") + "')"

	if e.Nullable {
		sql += " NULL"
	} else {
		sql += " NOT NULL"
	}

	sql += buildDefaultForString(e.Default)

	if e.OnUpdate != "" {
		sql += " ON UPDATE " + e.OnUpdate
	}

	if e.Comment != "" {
		sql += fmt.Sprintf(" COMMENT '%s'", e.Comment)
	}

	return sql
}

// Bit represents default `bit` column type
//
// Default migrator.Bit will build a sql row: `bit NOT NULL`
//
// Examples:
//		➡️ migrator.Bit{Precision: 8, Default: "1", Comment: "mario game code"}
//			↪️ bit(8) NOT NULL DEFAULT 1 COMMENT 'mario game code'
//		➡️ migrator.Bit{Precision: 64, Nullable: true, OnUpdate: "set null"}
//			↪️ bit(64) NULL ON UPDATE set null
type Bit struct {
	Default  string
	Nullable bool
	Comment  string
	OnUpdate string

	Precision uint16
}

func (b Bit) buildRow() string {
	sql := "bit"

	if b.Precision > 0 {
		sql += "(" + strconv.Itoa(int(b.Precision)) + ")"
	}

	if b.Nullable {
		sql += " NULL"
	} else {
		sql += " NOT NULL"
	}

	if b.Default != "" {
		sql += " DEFAULT " + b.Default
	}

	if b.OnUpdate != "" {
		sql += " ON UPDATE " + b.OnUpdate
	}

	if b.Comment != "" {
		sql += fmt.Sprintf(" COMMENT '%s'", b.Comment)
	}

	return sql
}

// Binary represents binary column type: `binary` or `varbinary`
//
// Default migrator.Binary will build a sql row: `varbinary NOT NULL`
//
// Examples:
//		binary		➡️ migrator.Binary{Fixed: true, Precision: 36, Default: "1", Comment: "uuid"}
//			↪️ binary(36) NOT NULL DEFAULT 1 COMMENT 'uuid'
//		varbinary	➡️ migrator.Binary{Precision: 255, Nullable: true, OnUpdate: "set null"}
//			↪️ varbinary(255) NULL ON UPDATE set null
type Binary struct {
	Default  string
	Nullable bool
	Comment  string
	OnUpdate string

	Fixed     bool // binary for fixed, otherwise varbinary
	Precision uint16
}

func (b Binary) buildRow() string {
	sql := ""

	if !b.Fixed {
		sql += "var"
	}

	sql += "binary"

	if b.Precision > 0 {
		sql += fmt.Sprintf("(%s)", strconv.Itoa(int(b.Precision)))
	}

	if b.Nullable {
		sql += " NULL"
	} else {
		sql += " NOT NULL"
	}

	if b.Default != "" {
		sql += " DEFAULT " + b.Default
	}

	if b.OnUpdate != "" {
		sql += " ON UPDATE " + b.OnUpdate
	}

	if b.Comment != "" {
		sql += fmt.Sprintf(" COMMENT '%s'", b.Comment)
	}

	return sql
}

func buildDefaultForString(v string) string {
	if v == "" {
		return ""
	}

	if v[:1] == "(" && v[len(v)-1:] == ")" {
		return fmt.Sprintf(" DEFAULT %s", v)
	}

	if v == "<empty>" || v == "<nil>" {
		v = ""
	}

	return fmt.Sprintf(" DEFAULT '%s'", v)
}
