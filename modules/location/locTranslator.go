package location

import (
	"errors"
	"strings"

	"bitbucket.org/go-mis/services"
)

type LocTranslator interface {
	GetLocation()
}

func TranslateLocation(locs interface{}) (error, LocTranslator) {
	switch v := locs.(type) {
	case string:
		l := location{}
		l.LocationCode = v
		return nil, &l
	case []interface{}:
		var realVal []string

		for _, val := range v {
			realVal = append(realVal, val.(string))
		}

		ls := locations{}
		ls.LocationCodes = realVal
		return nil, &ls
	default:
		return errors.New("Wrong arguments"), nil
	}
}

const (
	pz   = ".00.000.0000"
	cz   = ".000.0000"
	kecz = ".0000"
)

// location extractor
type location struct {
	Province     string `json:"province"`
	City         string `json:"city"`
	Kecamatan    string `json:"kecamatan"`
	Kelurahan    string `json:"kelurahan"`
	LocationCode string `json:"locationCode"`
}

func (l *location) extract(locCode string) {
	str := strings.Split(locCode, ".")

	if str[0] != "00" {
		p := str[0] + pz
		q := "select name as province from \"inf_location\" where \"locationCode\" = ?"
		services.DBCPsql.Raw(q, p).Scan(&l)
	}

	if str[1] != "00" {
		c := str[0] + "." + str[1] + cz
		q := "select name as city from \"inf_location\" where \"locationCode\" = ?"
		services.DBCPsql.Raw(q, c).Scan(&l)
	}

	if str[2] != "000" {
		kel := str[0] + "." + str[1] + "." + str[2] + kecz
		q := "select name as kelurahan from \"inf_location\" where \"locationCode\" = ?"
		services.DBCPsql.Raw(q, kel).Scan(&l)
	}

	if str[3] != "0000" {
		q := "select name as kecamatan from \"inf_location\" where \"locationCode\" = ?"
		services.DBCPsql.Raw(q, locCode).Scan(&l)
	}

	l.LocationCode = locCode

}

func (l *location) GetLocation() {
	l.extract(l.LocationCode)
}

type locations struct {
	Locations     []location
	LocationCodes []string
}

func (ls *locations) extract(locCode []string) {
	locs := []location{}
	for _, val := range locCode {
		l := location{}
		l.extract(val)
		locs = append(locs, l)
	}

	ls.Locations = locs
}

func (ls *locations) GetLocation() {
	ls.extract(ls.LocationCodes)
}
