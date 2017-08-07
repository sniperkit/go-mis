package installment

import (
	"log"
	"strings"
	"testing"
	"time"
)

func TestGenerateArchiveID(t *testing.T) {
	branchID := "BRANCH-001"
	date := time.Now().Local().Format("2006-01-02")
	archiveID := generateArchiveID(branchID)
	if !strings.Contains(archiveID, branchID) {
		t.Errorf("Expected contains %s actual %s", branchID, archiveID)
	}
	if !strings.Contains(archiveID, date) {
		t.Errorf("Expected contains %s actual %s", date, archiveID)
	}
}

func TestCOH(t *testing.T) {
	coh := []Coh{
		{
			installmentId: 1,
			cash:          50,
		},
		{
			installmentId: 2,
			cash:          100,
		},
		{
			installmentId: 3,
			cash:          150,
		},
	}
	tableTest := []struct {
		Id   uint64
		cash float64
	}{
		{
			Id:   1,
			cash: 50,
		},
		{
			Id:   2,
			cash: 100,
		},
		{
			Id:   3,
			cash: 150,
		},
	}
	for _, v := range tableTest {
		coh := getCOH(v.Id, coh)
		if coh == -1 {
			t.Errorf("Expected %f actual %f\n", v.cash, coh)
		}
	}
	tempCoh := coh[0]
	tempCoh.installmentId = 555555
	res := getCOH(tempCoh.installmentId, coh)
	if res != -1 {
		t.Errorf("Expected %f actual %v\n", res, -1)
	}
}

func TestGetLog(t *testing.T) {
	type Data struct {
		ID         int64
		Amount     float64
		CashOnHand float64
	}
	payload := Data{
		ID:         1,
		Amount:     200,
		CashOnHand: 150,
	}
	branchID := "BRANCH-001"
	dataLog := getLog(branchID, payload)
	if dataLog.Data == nil {
		t.Errorf("Expected %v actual %s\n", payload, dataLog)
	}
	if strings.ToUpper(dataLog.GroupID) != "VALIDASI TELLER" {
		t.Errorf("Expected %s actual %s\n", "VALIDASI TELLER", dataLog.GroupID)
	}
}

func TestPostToLog(t *testing.T) {
	data := struct {
		ID         int64
		Amount     float64
		CashInHand float64
	}{
		ID:         1,
		Amount:     200,
		CashInHand: 200,
	}
	dataLog := Log{
		GroupID:   "VALIDASI TELLER",
		ArchiveID: generateArchiveID("BRANCH-001"),
		Data:      data,
	}
	err := postToLog(dataLog)
	if err != nil {
		log.Println("#ERROR: ", err)
		t.Error("Failed POST loging to GO-LOG APPS")
	}
}
