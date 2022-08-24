package main

import (
	"testing"
)

func TestSumGetTypeFromString(t *testing.T) {
	intType := getTypeFromString("5")
	if intType != "Int" {
		t.Errorf("Fail to identify int type")
	}
	stringType := getTypeFromString("qzdqz")
	if stringType != "String" {
		t.Errorf("Fail to identify int type")
	}
	floatType := getTypeFromString("5.5")
	if floatType != "Float" {
		t.Errorf("Fail to identify int type")
	}
	booleanType := getTypeFromString("true")
	if booleanType != "Boolean" {
		t.Errorf("Fail to identify int type")
	}
}
