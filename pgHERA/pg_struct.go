package pghera

// --------- DDL

// ColumnDef provides DDL for table column.
type ColumnDef struct {
	Name       string
	Type       string
	PrimaryKey bool
	NotNull    bool
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

// ----------- INSERT

// RowData is to be used for insert single row
type RowData struct {
	TableName   string
	ColumnNames string
	Values      []string
}

// BulkValues to be used when working with data from several rows from a specific table.
type BulkValues struct {
	TableName   string
	ColumnNames string
	Values      [][]string
}

// CellValue - for select
type CellValue struct {
	ColumnName string
	CellData   interface{}
}

// RowValues contains data from given row.
type RowValues struct {
	ColumnNames []string
	Values      []interface{}
}

// TableData contains data from given table.
type TableData struct {
	ColumnNames []string
	Rows        [][]interface{}
}
