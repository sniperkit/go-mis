package cif

import (
	"bitbucket.org/go-mis/services"
	iris "gopkg.in/kataras/iris.v4"
)

func Init() {
	services.DBCPsql.AutoMigrate(&Cif{})
	services.BaseCrudInit(Cif{}, []Cif{})
}

// FetchAll - fetchAll agent data
func FetchAll(ctx *iris.Context) {
	cifs := []CifFragment{}

	query := "SELECT cif.\"id\", cif.\"name\", cif.\"isActivated\", cif.\"isValidated\", "
	query += "r_cif_borrower.\"borrowerId\" as \"borrowerId\", r_cif_investor.\"investorId\" as \"investorId\", "
	query += "r_cif_borrower.\"borrowerId\" IS NOT NULL as \"isBorrower\", r_cif_investor.\"investorId\" IS NOT NULL as \"isInvestor\" "
	query += "FROM cif "
	query += "LEFT JOIN r_cif_borrower ON r_cif_borrower.\"cifId\" = cif.\"id\" "
	query += "LEFT JOIN r_cif_investor ON r_cif_investor.\"cifId\" = cif.\"id\" "

	services.DBCPsql.Raw(query).Find(&cifs)
	ctx.JSON(iris.StatusOK, iris.Map{"data": cifs})
}

// GetByID agent by id
func GetByID(ctx *iris.Context) {
	cif := Cif{}
	services.DBCPsql.Where("\"deletedAt\" IS NULL AND id = ?", ctx.Param("id")).First(&cif)
	if cif == (Cif{}) {
		ctx.JSON(iris.StatusOK, iris.Map{"data": iris.Map{}})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"data": cif})
	}
}
