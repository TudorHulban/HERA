package pghera

/*
File concentrates interfaces that are satisfied by Hera.
*/

// IHDDL Concentrates Hera DDL actions that are performed in database specified in constructor.
type IHDDL interface {
	// CreateTable Creates table based on model passed as pointer. In case of simulation it would return the table name and nil.
	CreateTable(model interface{}, simulateOnly bool) (string, string, error)
	// TableExists Checks if table exists. Returns nil if table exists.
	TableExists(tableName string) error
	// DropTable Drops table with cascade option. Any error is returned.
	DropTable(tableName string, withCascade bool) error
}

// IHSQL Concentrates Hera SQL actions that are performed in database specified in constructor.
type IHSQL interface {
	// InsertModel Inserts data presented in a model instance.
	InsertModel(modelData interface{}) error
}
