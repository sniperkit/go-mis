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
	cifs := []Cif{}
	services.DBCPsql.Where("\"deletedAt\" IS NULL").Order("id asc").Find(&cifs)
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
