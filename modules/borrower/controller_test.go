package borrower

import (
	"testing"
	"bitbucket.org/go-mis/services"
)

func TestGetOrCreateBorrowerId(t *testing.T) {

	var e error
	if e != nil {
		t.Error(e)
	}

	input := map[string]interface{}{
		"client_ktp": "3271066002690001",
	}

	db := services.DBCPsql.Begin()

	i, e := GetOrCreateBorrowerId(input, db)
	if i != 26170 {
		t.Error("Got wrong ID")
	}

	input = map[string]interface{}{
		"client_ktp":            "3271066002690002",
		"client_simplename":     "Egon",
		"client_fullname":       "Egon Firman",
		"client_birthplace":     "Bandung",
		"client_birthdate":      "2017-02-01aaa",
		"photo_ktp":             "",
		"client_npwp":           "",
		"client_marital_status": "Menikah",
		"client_ibu_kandung":    "Sri",
		"client_religion":       "Islam",
		"client_alamat":         "Jalan Sukawrana",
		"client_desa":           "Somethign",
		"client_kecamatan":      "S",
		"client_rt":             "001",
		"client_rw":             "003",
		"data_pendapatan_istri": "300000",
	}

	i, e = GetOrCreateBorrowerId(input, db)
	if e == nil {
		t.Error("It should be error")

	}
	db.Rollback()
}
