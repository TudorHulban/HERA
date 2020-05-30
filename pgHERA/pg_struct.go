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
