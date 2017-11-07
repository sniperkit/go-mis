package validationTeller

import (
	"log"
	"testing"
)

func TestFindInstallmentData(t *testing.T) {
	var branchID uint64 = 5
	dateParam := "2017-08-24"
	data, err := FindInstallmentData(branchID, dateParam, false)
	if err != nil {
		t.Fatal("Unable to get data installment", err)
	}
	log.Println(data)
	t.Log("Pass")
}

func TestFindDataTransfer(t *testing.T) {
	type args struct {
		branchID uint64
		date     string
	}
	tests := []struct {
		name    string
		args    args
		want    DataTransfer
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "positive test case",
			args: args{
				branchID: 5,
				date:     "2017-08-24",
			},
			want: DataTransfer{
				ID: 14,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindDataTransfer(tt.args.branchID, tt.args.date)
			log.Println(got)
			if err != nil {
				t.Fatal(err)
			}
			t.Log("Pass")
		})
	}
}

func TestFindVTDetailByGroupAndDate(t *testing.T) {
	type args struct {
		groupID uint64
		date    string
	}
	tests := []struct {
		name string
		args args
		want []RawInstallmentDetail
	}{
		{
			name: "Positive Test Case",
			args: args{
				groupID: 1544,
				date:    "2017-08-28",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindVTDetailByGroupAndDate(tt.args.groupID, tt.args.date)
			if err != nil {
				t.Errorf("FindVTDetailByGroupAndDate() error = %v", err)
				return
			}
			t.Log("CASH ON HAND NOTE: ", got[0].CashOnHandNote)
			t.Log("CASH ON RESERVE NOTE: ", got[0].CashOnReserveNote)
		})
	}
}
