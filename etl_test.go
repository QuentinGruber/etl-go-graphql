package main

import (
	"testing"
)

func TestConvertToCamelCase(t *testing.T) {
	convertedName := convertToCamelCase("cou_cou")
	if convertedName != "couCou" {
		t.Errorf("convertToCamelCase should return 'couCou'")
	}
}
