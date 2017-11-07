package services

import (
	"log"
	"testing"
)

func TestPostToLog(t *testing.T) {
	data := struct {
		ID         uint64
		Projection float64
		CashOnHand float64
	}{
		ID:         1,
		Projection: 2000000.00,
		CashOnHand: 2000000.00,
	}
	dataLog := Log{
		GroupID:   ConstructNotesGroupId(1000, "2017-08-23"),
		ArchiveID: "VT-NOTES",
		Data:      data,
	}
	type args struct {
		l Log
	}
	var testArgs args
	testArgs.l = dataLog
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Positive test case",
			args: testArgs,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := PostToLog(tt.args.l)
			if err != nil {
				t.Errorf("Error when post to log")
			}
		})
	}
}

func TestGetNotes(t *testing.T) {
	var branchID uint64 = 1000
	dateStr := "2017-08-23"
	groupID := ConstructNotesGroupId(branchID, dateStr)
	notes, err := GetNotes(groupID)
	if err != nil {
		t.Fatal("Unable to retrieve notes from GO-LOG App")
	}
	if len(notes) > 1 {
		t.Fatal("Number of notes must be 1")
	}
	t.Log("Should return 1 notes")
	log.Println("notes: ", notes)
}
