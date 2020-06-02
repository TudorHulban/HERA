package pghera

// translationTable Type used for providing persistence for the created translation table.
type translationTable map[string]string

func newTranslationTable() *translationTable {
	t := make(translationTable)
	t["string"] = "text"
	t["*string"] = "text"
	t["int"] = "bigint"
	t["int64"] = "bigint"
	t["float64"] = "numeric"
	t["*float64"] = "numeric"
	t["bool"] = "boolean"
	t["*bool"] = "boolean"

	return &t
}
