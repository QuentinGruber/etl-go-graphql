package main

import "strconv"

func getTypeFromString(typeString string) string {
	if _, err := strconv.ParseInt(typeString, 10, 64); err == nil {
		return "Int"
	}
	if _, err := strconv.ParseFloat(typeString, 64); err == nil {
		return "Float"
	}
	if _, err := strconv.ParseBool(typeString); err == nil {
		return "Boolean"
	}
	return "String"
}
