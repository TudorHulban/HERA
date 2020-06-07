package pghera

// getRDBMSType Helper returns RDBMS type based on Go type.
// If PK switch to auto increment field.
func getRDBMSType(goType string, isPK bool) string {
	if isPK {
		switch goType {
		case "int":
			return "serial"
		case "int64":
			return "bigserial"
		}
		return ""
	}

	switch goType {
	case "string":
		return "text"
	case "*string":
		return "text"
	case "int":
		return "bigint"
	case "int64":
		return "bigint"
	case "float64":
		return "numeric"
	case "*float64":
		return "numeric"
	case "bool":
		return "boolean"
	case "*bool":
		return "boolean"
	}
	return ""
}
