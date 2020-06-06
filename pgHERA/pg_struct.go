package pghera

// User Type used for developing. TODO: Move to test or documentation.
type User struct {
	tableName   struct{}    `hera:"theTableName"`                        // nolint
	ID          int64       `hera:"pk"`                                  // nolint
	name        string      `hera:"default:xx, index"`                   // nolint
	age         int         `hera:"required, column-name:theage, index"` // nolint
	isConnected bool        `hera:"default:true"`                        // nolint
	comment     string      `hera:"-"`                                   // nolint
	toSkip      interface{} // nolint
}

// --------- DDL

// ColumnDef provides DDL for table column.
type ColumnDef struct {
	PrimaryKey bool
	NotNull    bool
	Name       string
	Type       string
}

// TableDDL is nutshell for database table DDL.
type TableDDL struct {
	Name        string
	TableFields []ColumnDef
}

// SchemaDDL is nutshell for database schema.
type SchemaDDL struct {
	Tables []TableDDL
}

// ----------- INSERT. Data contains values and names.

// TableColumnNames Type for CSV of column names.
type TableColumnNames struct {
	ColumnNames []string
}

// RowValues contains data from given row.
type RowValues struct {
	Values []interface{}
}

// RowData contains data from given table.
type RowData struct {
	TableName string
	TableColumnNames
	Data []RowValues
}

// CellData - for select
type CellData struct {
	ColumnName string
	CellValue  interface{}
}
