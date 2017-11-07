package plottingBorrower

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"bitbucket.org/go-mis/config"
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
		err := errors.New("ERROR: No Plotting Params were found in the request body.")
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

// this function toggles activation of plottingparams
func TogglePlottingParamsActivation(ctx *iris.Context) {

	payload := struct {
		InvestorId uint64 `json:investorId`
		Activate   bool   `json:activate`
	}{}

	if err := ctx.ReadJSON(&payload); err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	inv := investor.Investor{
		ID: payload.InvestorId,
	}

	services.DBCPsql.Model(&inv).Update("isBorrowerCriteriaActive", payload.Activate)
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
	})

}

// FindRecomendedLoanByInvestorCriteria - this functoin fetch all loan by criteria
func FindRecomendedLoanByInvestorCriteria(ctx *iris.Context) {
	investorID := ctx.Param("investorId")
	disFrom := ctx.URLParam("disFrom")
	disTo := ctx.URLParam("disTo")
	resultGoloan := make([]RecomendedLoan, 0)

	redisLoan, err := RetriveRecomendedLoanFromRedis(investorID)
	if err != nil {
		log.Println("[ERROR] ", err)
	}
	if len(redisLoan) > 0 {
		fmt.Println("Data from redis")
		fmt.Printf("Redis Loan: %+v", redisLoan)
		ctx.JSON(http.StatusOK, iris.Map{
			"status": "Success",
			"data":   redisLoan,
		})
		return
	}
	resultGoloan, err = RetrieveRecomendedLoanFromLoanService(disFrom, disTo, investorID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, iris.Map{
			"status":  "Error",
			"message": "Internal Server Error",
		})
		return
	}
	fmt.Println("Result go loan: ", resultGoloan)
	ctx.JSON(http.StatusOK, iris.Map{
		"status": "Success",
		"data":   resultGoloan,
	})
}

// RetrieveRecomendedLoanFromLoanService - get all data recomended loan from loan service
func RetrieveRecomendedLoanFromLoanService(disFrom, disTo, investorID string) ([]RecomendedLoan, error) {
	var goloanResp GOLoanSuccessResponse
	goLoanURI := config.Configuration.GoLoanPath + "/" + "loan/plotting-borrower/recomended-loan-investor/" + investorID + "?disFrom=" + disFrom + "&disTo=" + disTo
	fmt.Println("GOLOAN URI: ", goLoanURI)
	body, err := services.CircuitBreaker.Get(goLoanURI)
	if err != nil {
		log.Println("[ERROR] ", err)
		return nil, err
	}
	err = json.Unmarshal(body, &goloanResp)
	if err != nil {
		log.Println("[ERROR] ", err)
		return nil, err
	}
	if goloanResp.Code != 200 && strings.ToUpper(goloanResp.Message) != "SUCCESS" {
		return nil, errors.New("Unable to get recomended loan data from go loan service")
	}
	go func(data []RecomendedLoan) {
		b, errMarshall := json.Marshal(&data)
		redisClient, errRed := services.NewClientRedis()
		if errMarshall != nil || errRed != nil {
			log.Println("[ERROR]", errMarshall)
			log.Println("[ERROR]", errRed)
		} else {
			err := redisClient.SaveRecomendedLoan(investorID, b)
			if err != nil {
				log.Println("[ERROR] ", err)
			}
		}
	}(goloanResp.Data)
	return goloanResp.Data, nil
}

// RetriveRecomendedLoanFromRedis - get data recomened loan from redis
// wrapped wheter is all or specifig by investor id
func RetriveRecomendedLoanFromRedis(investorID string) ([]RecomendedLoan, error) {
	var err error
	loanRedis := make([]RecomendedLoan, 0)
	switch strings.ToUpper(strings.TrimSpace(investorID)) {
	case "ALL":
		loanRedis, err = FindAllRecomendedLoanFromRedis()
	default:
		loanRedis, err = FindRecomendedLoanFromRedis(investorID)
	}
	return loanRedis, err
}

// FindAllRecomendedLoanFromRedis - find all recomended loan from redis
func FindAllRecomendedLoanFromRedis() ([]RecomendedLoan, error) {
	loans := make([]RecomendedLoan, 0)
	redisClient, err := services.NewClientRedis()
	if err != nil {
		return nil, err
	}
	strData, err := redisClient.GetAllRecomendedLoan()
	if err != nil {
		return nil, err
	}
	for i := range strData {
		var recLoan RecomendedLoan
		err = json.Unmarshal([]byte(strData[i]), &recLoan)
		if err != nil {
			return nil, err
		}
		loans = append(loans, recLoan)
	}
	return loans, nil
}

// FindRecomendedLoanFromRedis - find all recomended loan from redis by investor id
func FindRecomendedLoanFromRedis(investorID string) ([]RecomendedLoan, error) {
	loanRedis := make([]RecomendedLoan, 0)
	redisClient, err := services.NewClientRedis()
	if err != nil {
		return nil, err
	}
	key, err := redisClient.GetPRecomendedLoanKey(investorID)
	if err != nil {
		return nil, err
	}
	b, err := redisClient.GetRecomendedLoan(key)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &loanRedis)
	if err != nil {
		return nil, err
	}
	return loanRedis, nil
}
