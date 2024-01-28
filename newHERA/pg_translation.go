package main

func reflectToPG(reflectType string, isPK bool) string {
	if isPK {
		switch reflectType {
		case "int", "int64":
			return "bigserial"

		default:
			return "serial"
		}
	}

	switch reflectType {
	case "string", "*string":
		return "text"

	case "int", "int64":
		return "bigint"

	case "float64", "*float64":
		return "numeric"

	case "bool", "*bool":
		return "boolean"
	}

	return ""
}
