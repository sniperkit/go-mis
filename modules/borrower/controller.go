package borrower

import (
	"bitbucket.org/go-mis/services"
	//"strconv"
	"fmt"
	iris "gopkg.in/kataras/iris.v4"
	//"encoding/json"
	"bitbucket.org/go-mis/modules/cif"

)

func Init() {
	services.DBCPsql.AutoMigrate(&Borrower{})
	services.BaseCrudInit(Borrower{}, []Borrower{})
}

// GetByID agent by id
func Approve(ctx *iris.Context) {

	// map the payload 
	payload := make(map[string]interface{}, 0)
	if err := ctx.ReadJSON(&payload); err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ktp := payload["client_ktp"].(string)

	// get CIF with with id = id
	id := ctx.Param("id")
	cifData := cif.Cif{}
	services.DBCPsql.Where("\"deletedAt\" IS NULL AND id = ?", id).Scan(&cifData)
	if ktp != "" && (cifData.IdCardNo == ktp)  {
		fmt.Println("user already here")
	} else {
		fmt.Println("new user detected")
	}

	return
}
