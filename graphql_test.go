package main

import (
	"testing"
)

func TestExtractQueriedDataset(t *testing.T) {
	datasetsResult1 := extractQueriedDataset("{query:{staff{email,Campus}}}")
	if len(datasetsResult1) != 1 {
		t.Errorf("extractQueriedDataset should return 1 dataset")
	}
	if datasetsResult1[0] != "staff" {
		t.Errorf("extractQueriedDataset should return staff dataset")
	}
	datasetsResult2 := extractQueriedDataset("{query:{students{last_name,first_name,gender,student_id},grades{credit,cursus,student_id}}}")
	if len(datasetsResult2) != 2 {
		t.Errorf("extractQueriedDataset should return 2 dataset")
	}
	if datasetsResult2[0] != "students" {
		t.Errorf("Fail to extract students dataset")
	}
	if datasetsResult2[1] != "grades" {
		t.Errorf("Fail to extract grades dataset")
	}

	datasetsResult3 := extractQueriedDataset("{query:{students(last_name:paul){last_name,first_name,gender,student_id},grades(credit:12){credit,cursus,student_id}}}")
	if len(datasetsResult3) != 2 {
		t.Errorf("extractQueriedDataset should return 2 dataset")
	}
	if datasetsResult3[0] != "students" {
		t.Errorf("Fail to extract students dataset")
	}
	if datasetsResult3[1] != "grades" {
		t.Errorf("Fail to extract grades dataset")
	}
}
