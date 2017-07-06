package location

import (
	"strings"

	"bitbucket.org/go-mis/services"
)

const (
	pz   = ".00.000.0000"
	cz   = ".000.0000"
	kecz = ".0000"
)

// location extractor
type SingleExtractor struct {
	Province     string `json:"province"`
	City         string `json:"city"`
	Kecamatan    string `json:"kecamatan"`
	Kelurahan    string `json:"kelurahan"`
	LocationCode string `json:"locationCode"`
}

func (s *SingleExtractor) extract(locCode string) {
	str := strings.Split(locCode, ".")

	p := str[0] + pz
	q := "select name as province from \"inf_location\" where \"locationCode\" = ?"
	services.DBCPsql.Raw(q, p).Scan(&s)

	c := str[0] + "." + str[1] + cz
	q = "select name as city from \"inf_location\" where \"locationCode\" = ?"
	services.DBCPsql.Raw(q, c).Scan(&s)

	kel := str[0] + "." + str[1] + "." + str[2] + kecz
	q = "select name as kelurahan from \"inf_location\" where \"locationCode\" = ?"
	services.DBCPsql.Raw(q, kel).Scan(&s)

	q = "select name as kecamatan from \"inf_location\" where \"locationCode\" = ?"
	services.DBCPsql.Raw(q, locCode).Scan(&s)

	s.LocationCode = locCode

}
