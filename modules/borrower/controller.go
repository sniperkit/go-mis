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
		// use CIF data
		fmt.Println("use CIF data")
	} else {
		// create new CIF
		// populate new CIF data
		CreateCIF(payload)
		// save to db
		services.DBCPsql.Table("cif").Create(&rUserMisAccessToken)

	}
	return
}

func CreateCIF(payload map[string]interface{}) {
	// check each payload value which is interface is not nil
	// convert each value into string (we only deal with string this time) 
	var cpl map[string]string
	cpl = make(map[string]string)
	for k, _ := range payload {
		if payload[k] != nil {
			cpl[k] = payload[k].(string)
		}
	}

	// assign the data to corresppndence cif
	newCif := cif.Cif{}

	newCif.Username						= cpl["client_simplename"]
	newCif.Name								= cpl["client_fullname"]
	newCif.PlaceOfBirth       = cpl["client_birthplace"]
	newCif.DateOfBirth        = cpl["client_birthdate"]
	newCif.IdCardNo           = cpl["client_ktp"]
	//newCif.IdCardValidDate    = payload[""].(string)
	newCif.IdCardFilename     = cpl["photo_ktp"]
	//newCif.TaxCardNo          = payload[""].(string)
	//newCif.TaxCardFilename    = payload[""].(string)
	newCif.MaritalStatus      = cpl["client_marital_status"]
	newCif.MotherName         = cpl["client_ibu_kandung"]
	newCif.Religion           = cpl["client_religion"]
	newCif.Address            = cpl["client_alamat"]
	//newCif.Kelurahan          = payload[""].(string)
	newCif.Kecamatan          = cpl["kecamatan"]
	//newCif.City               = payload[""].(string)
	//newCif.Province           = payload[""].(string)
	//newCif.Nationality        = payload[""].(string)
	//newCif.Zipcode            = payload[""].(string)
	//newCif.PhoneNo            = payload[""].(string)
	//newCif.CompanyName        = payload[""].(string)
	//newCif.CompanyAddress     = payload[""].(string)
	//newCif.Occupation         = payload[""].(string)
	//newCif.Income             = payload[""].(string)
	//newCif.IncomeSourceFund   = payload[""].(string)
	//newCif.IncomeSourceCountry= payload[""].(string)
	//newCif.IsActivated        = payload[""].(string)
	//newCif.IsVAlidated        = payload[""].(string)
	//newCif.CreatedAt          = payload[""].(string)
	//newCif.UpdatedAt          = payload[""].(string)
	//newCif.DeletedAt          = payload[""].(string)
	fmt.Println(newCif)
}
