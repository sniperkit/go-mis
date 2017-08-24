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
