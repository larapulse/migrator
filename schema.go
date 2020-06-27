// Package migrator represents MySQL database migrator
package migrator

// Schema allows to add commands on schema.
// It should be used within migration to add migration commands.
type Schema struct {
	pool []command
}

// CreateTable allows to create table in schema
//
// Example:
//		var s migrator.Schema
//		t := migrator.Table{Name: "test"}
//
//		s.CreateTable(t)
func (s *Schema) CreateTable(t Table) {
	s.pool = append(s.pool, createTableCommand{t})
}

// DropTable removes table from schema
// Warning ⚠️ BC incompatible
//
// Example:
//		var s migrator.Schema
//		s.DropTable("test", false, "")
//
// Soft delete (drop if exists)
//		s.DropTable("test", true, "")
func (s *Schema) DropTable(name string, soft bool, option string) {
	s.pool = append(s.pool, dropTableCommand{name, soft, option})
}

// RenameTable executes command to rename table
// Warning ⚠️ BC incompatible
//
// Example:
//		var s migrator.Schema
//		s.RenameTable("old", "new")
func (s *Schema) RenameTable(old string, new string) {
	s.pool = append(s.pool, renameTableCommand{old: old, new: new})
}

// AlterTable makes changes on table level
//
// Example:
//		var s migrator.Schema
//		var c TableCommands
//		s.AlterTable("test", c)
func (s *Schema) AlterTable(name string, c TableCommands) {
	s.pool = append(s.pool, alterTableCommand{name, c})
}

// CustomCommand allows to add custom command to the Schema
//
// Example:
//		type customCommand string
//
//		func (c customCommand) toSQL() string {
//			return string(c)
//		}
//
//		c := customCommand("DROP PROCEDURE abc")
//		var s migrator.Schema
//		s.CustomCommand(c)
func (s *Schema) CustomCommand(c command) {
	s.pool = append(s.pool, c)
}
