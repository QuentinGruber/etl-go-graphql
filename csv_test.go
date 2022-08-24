package main

import (
	"testing"
)

func TestGetSeparator(t *testing.T) {
	separator := getSeparator("id;student_id;cursus;module;grade\n")
	if separator != ";" {
		t.Errorf("getSeparator should detect ';' as separator")
	}
}

func TestStringToLines(t *testing.T) {
	lines, err := StringToLines("line1;\nline2")
	if err != nil {
		t.Error(err)
	}
	if len(lines) != 2 {
		t.Errorf("StringToLines should return 2 lines")
	}
}
