package plottingBorrower

import (
	"encoding/json"
	"errors"

	"bitbucket.org/go-mis/modules/investor"
	"bitbucket.org/go-mis/services"

	"gopkg.in/kataras/iris.v4"
)

// This function saves potting paramaters as borrower criteria
// into investor table. The data would be saved in json format
func SavePlottingParams(ctx *iris.Context) {
	// convert requestbody to string
	pp := string(ctx.Request.Body())
	if pp == "" {
		err := errors.New("No Plotting Params were found in the request body.")
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	in := []byte(pp)
	var raw map[string]interface{}
	json.Unmarshal(in, &raw)

	// set plottingParrams as borrowerCrtiteria on investor
	var invId uint64
	val, ok := raw["investorId"].(float64)
	if !ok {
		err := errors.New("ERROR: investorId is not a number")
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	invId = uint64(val)

	var isBorrowerCriteriaActive bool
	val2, ok := raw["isBorrowerCriteriaActive"].(bool)
	if !ok {
		err := errors.New("ERROR: IsBorrowerCriteriaActive is not a bool")
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	isBorrowerCriteriaActive = val2

	delete(raw, "investorId")
	delete(raw, "isBorrowerCriteriaActive")
	out, _ := json.Marshal(raw)

	// model to update
	inv := investor.Investor{
		ID: invId,
	}

	// new data
	data := map[string]interface{}{
		"borrowerCriteria":         out,
		"isBorrowerCriteriaActive": isBorrowerCriteriaActive,
	}

	services.DBCPsql.Model(&inv).Updates(data)
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
	})

}

// This function get all investors which their borrowerCriteria is not null
func ListPlottingParams(ctx *iris.Context) {

	totalRows := 0
	investors := []struct {
		ID                       uint64 `gorm:"column:id" json:"id"`
		InvestorNo               uint64 `gorm:"column:investorNo" json:"investorNo"`
		InvestorName             string `gorm:"column:name" json:"investorName"`
		IsBorrowerCriteriaActive *bool  `gorm:"column:isBorrowerCriteriaActive" json:"isBorrowerCriteriaActive"`
	}{}

	query := `select investor.id, investor."investorNo", cif."name", investor."isBorrowerCriteriaActive" 
	from investor
	join r_cif_investor rci on rci."investorId" = investor.Id
	join cif on cif.id = rci."cifId"
	where "borrowerCriteria" <> '{}' and investor."deletedAt" is null`
	services.DBCPsql.Raw(query).Scan(&investors)
	for _, _ = range investors {
		totalRows += 1
	}

	ctx.JSON(iris.StatusOK, iris.Map{
		"status":    "success",
		"data":      investors,
		"totalRows": totalRows,
	})

}

// this functoin fetch list investor that eligble for create borrowerCriteria
func FindEligbleInvestor(ctx *iris.Context) {
	investorId := ctx.Param("investorId")

	query := `select investor.id,cif."name","investorNo","borrowerCriteria" from investor
	join r_cif_investor on r_cif_investor."investorId"=investor.id
	join cif on cif.id=r_cif_investor."cifId"
	where investor.id=?  and investor."deletedAt" is null`

	investor := EligbleInvestor{}
	services.DBCPsql.Raw(query, investorId).Scan(&investor)

	if investor.BorrowerCriteria != "{}" {
		ctx.JSON(iris.StatusOK, iris.Map{
			"status":  "error",
			"message": "Investor already has criteria",
		})
		return
	}

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   investor,
	})

}

// this functoin fetch the detail of the plotting params
func GetPlottingParamsDetail(ctx *iris.Context) {
	investorId := ctx.Param("investorId")

	plottingParams := struct {
		ID                       uint64                 `gorm:"column:id" json:"id"`
		InvestorNo               uint64                 `gorm:"column:investorNo" json:"investorNo"`
		InvestorName             string                 `gorm:"column:name" json:"investorName"`
		IsBorrowerCriteriaActive *bool                  `gorm:"column:isBorrowerCriteriaActive" json:"isBorrowerCriteriaActive"`
		BorrowerCriteriaJSONB    string                 `gorm:"column:borrowerCriteria" sql:"type:JSONB NOT NULL DEFAULT '{}'::JSONB" json:"-"`
		BorrowerCriteria         map[string]interface{} `json:"borrowerCriteria"`
	}{}

	query := `select investor.id, investor."investorNo", cif."name", investor."isBorrowerCriteriaActive", investor."borrowerCriteria" 
        from investor
        join r_cif_investor rci on rci."investorId" = investor.Id
        join cif on cif.id = rci."cifId"
        where investor.id = ?`
	services.DBCPsql.Raw(query, investorId).Scan(&plottingParams)

	// prepare json for borrowerCriteria
	in := []byte(plottingParams.BorrowerCriteriaJSONB)
	var raw map[string]interface{}
	json.Unmarshal(in, &raw)
	plottingParams.BorrowerCriteria = raw

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   plottingParams,
	})
}
