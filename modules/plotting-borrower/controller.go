package plottingBorrower

import (
	"encoding/json"
	"errors"

	"bitbucket.org/go-mis/modules/investor"
	"bitbucket.org/go-mis/services"

	"gopkg.in/kataras/iris.v4"
	"strconv"
	"fmt"
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

func FindPlottingBorrower(ctx *iris.Context) {
	
	stageParam := ctx.Param("stage")
	investorIdParams := ctx.URLParam("investorId")
	investorId := 0

	stage:=""
	if(stageParam == "investor") {
		investorId, err := strconv.Atoi(investorIdParams)
		if investorIdParams == "" || err != nil {
			ctx.JSON(iris.StatusBadRequest, iris.Map{
				"message":      "Bad Request",
				"errorMessage": "Invalid User ID",
			})
			return	
		}
		stage = "PRIVATE-INVESTOR"
	} else if(stageParam == "marketplace") {
		stage = "PRIVATE-MARKETPLACE"
	}else {
		ctx.JSON(iris.StatusBadRequest, iris.Map{
			"message":      "Bad Request",
			"errorMessage": "Invalid Stage",
		})
		return	
	}

	loans := RecommendedLoan{}

	query :=`select loan.id as "loanId",cif."name" as "borrowerName","group"."name" as "group",
	branch."name" as "branch",disbursement."disbursementDate"::date as "disbursementDate",loan.plafond,loan.rate,loan.tenor,loan."creditScoreGrade",loan.purpose from loan
	join r_loan_group rlg on rlg."loanId"=loan.id
	join "group" on "group".id = rlg."groupId"
	join r_loan_borrower rlb on rlb."loanId"=loan.id
	join r_cif_borrower rcb on rcb."borrowerId"=rlb."borrowerId"
	join cif on cif.id=rcb."cifId"
	join r_loan_branch rlbr on rlbr."loanId"=loan.id
	join branch on branch.id=rlbr."branchId"
	join r_loan_disbursement rld on rld."loanId"=loan.id
	join disbursement on disbursement.id=rld."disbursementId"
	join r_area_branch rab on rab."branchId"=branch.id
	join r_loan_sector rls on rls."loanId"=loan.id`

	if investorId != "" {
		query+= `
		join r_cif_investor rcfi on rcfi."cifId" = cif.id 
		join investor on investor.id = rcfi."investorId"
		where loan.stage=? and loan."deletedAt" is null and investor.id =? limit 3`;
		services.DBCPsql.Raw(query, stage, investorId).Scan(&loans)
	}else {
		query+= `where loan.stage=? and loan."deletedAt" is null`;
		services.DBCPsql.Raw(query, stage).Scan(&loans)
	}
	
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   loans,
	})

}

//TODO move to go-loan
// this functoin fetch all loan by criteria
func FindRecomendedLoanByInvestorCriteria(ctx *iris.Context) {
	investorId := ctx.Param("investorId")
	disFrom := ctx.URLParam("disFrom")
	disTo := ctx.URLParam("disTo")

	loans := []struct {
		LoanId                   uint64                 `gorm:"column:loanId" json:"loanId"`
		BorrowerName             string                 `gorm:"column:borrowerName" json:"borrowerName"`
		Group		             string                 `gorm:"column:group" json:"group"`
		Branch		             string                 `gorm:"column:branch" json:"branch"`
		DisbursementDate		 string                 `gorm:"column:disbursementDate" json:"disbursementDate"`
		Plafond              	 float64    			`gorm:"column:plafond" json:"plafond"`
		Rate                 	 float64    			`gorm:"column:rate" json:"rate"`
		Tenor                	 uint64     			`gorm:"column:tenor" json:"tenor"`
		CreditScoreGrade     	 string     			`gorm:"column:creditScoreGrade" json:"creditScoreGrade"`
		Purpose              	 string     			`gorm:"column:purpose" json:"purpose"`
	}{}

	queryLoan :=`select loan.id as "loanId",cif."name" as "borrowerName","group"."name" as "group",
	branch."name" as "branch",disbursement."disbursementDate"::date as "disbursementDate",loan.plafond,loan.rate,loan.tenor,loan."creditScoreGrade",loan.purpose from loan
	join r_loan_group rlg on rlg."loanId"=loan.id
	join "group" on "group".id = rlg."groupId"
	join r_loan_borrower rlb on rlb."loanId"=loan.id
	join r_cif_borrower rcb on rcb."borrowerId"=rlb."borrowerId"
	join cif on cif.id=rcb."cifId"
	join r_loan_branch rlbr on rlbr."loanId"=loan.id
	join branch on branch.id=rlbr."branchId"
	join r_loan_disbursement rld on rld."loanId"=loan.id
	join disbursement on disbursement.id=rld."disbursementId"
	join r_area_branch rab on rab."branchId"=branch.id
	join r_loan_sector rls on rls."loanId"=loan.id
	where loan.stage='PRIVATE' and loan."deletedAt" is null`

	if disFrom != ""{
		queryLoan+=` and disbursement."disbursementDate"::date >= '`+disFrom+`' `
	}

	if disTo != ""{
		queryLoan+=` and disbursement."disbursementDate"::date <= '`+disTo+`' `
	}
	if investorId != "-1" {
		plottingParams := struct {
			ID                       uint64           `gorm:"column:id" json:"id"`
			InvestorNo               uint64           `gorm:"column:investorNo" json:"investorNo"`
			InvestorName             string           `gorm:"column:name" json:"investorName"`
			IsBorrowerCriteriaActive *bool            `gorm:"column:isBorrowerCriteriaActive" json:"isBorrowerCriteriaActive"`
			BorrowerCriteriaJSONB    string           `gorm:"column:borrowerCriteria" sql:"type:JSONB NOT NULL DEFAULT '{}'::JSONB" json:"-"`
		}{}

		query := `select investor.id, investor."investorNo", cif."name", investor."isBorrowerCriteriaActive", investor."borrowerCriteria"
        from investor
        join r_cif_investor rci on rci."investorId" = investor.Id
        join cif on cif.id = rci."cifId"
        where investor.id = ?`
		services.DBCPsql.Raw(query, investorId).Scan(&plottingParams)

		// prepare json for borrowerCriteria
		in := []byte(plottingParams.BorrowerCriteriaJSONB)
		var criteria BorrowerCriteria
		json.Unmarshal(in, &criteria)

		// Filter Tenor
		filterTenor(criteria, &queryLoan)
		//Filter Area
		filterArea(criteria, &queryLoan)
		//Filter CreditScoreGrade
		filterCreditScoreGrade(criteria, &queryLoan)
		//Filter Sector
		filterSector(criteria, &queryLoan)
		//Filter Plafond
		filterPlafon(criteria, &queryLoan)
		//Filter Rate
		filterRate(criteria, &queryLoan)
	}
	services.DBCPsql.Raw(queryLoan).Scan(&loans)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   loans,
	})
}

//TODO move to go-loan
func filterPlafon(criteria BorrowerCriteria, queryLoan *string) {
	fmt.Println("Plafon",criteria.Plafon)
	if criteria.Plafon.OptionType==4 {
		addedQuery := " and (loan.plafond >"+strconv.Itoa(criteria.Plafon.From)+" and loan.plafond <"+strconv.Itoa(criteria.Plafon.To)+")"
		*queryLoan += addedQuery
	}else {
		operatorReference := struct {
			Operator string `gorm:"column:operator" json:"operator"`
		}{}

		query := `select operator from plotting_borrower_operator_reference where id = ?`
		services.DBCPsql.Raw(query, criteria.Plafon.OptionType).Scan(&operatorReference)
		if operatorReference.Operator != "" {
			addedQuery := " and loan.plafond " + operatorReference.Operator + " " + strconv.Itoa(criteria.Plafon.From)
			*queryLoan += addedQuery
		}
	}
}

//TODO move to go-loan
func filterRate(criteria BorrowerCriteria, queryLoan *string) {
	if criteria.Rate.OptionType==4 {
		addedQuery := " and (loan.rate >"+strconv.FormatFloat(criteria.Rate.From, 'f', -1, 64)+" and loan.rate <"+strconv.FormatFloat(criteria.Rate.To, 'E', -1, 64)+")"
		*queryLoan += addedQuery
	}else {
		operatorReference := struct {
			Operator string `gorm:"column:operator" json:"operator"`
		}{}

		query := `select operator from plotting_borrower_operator_reference where id = ?`
		services.DBCPsql.Raw(query, criteria.Rate.OptionType).Scan(&operatorReference)
		if operatorReference.Operator != "" {
			addedQuery := " and loan.rate " + operatorReference.Operator + " " + strconv.FormatFloat(criteria.Rate.From, 'E', -1, 64)
			*queryLoan += addedQuery
		}
	}
}

//TODO move to go-loan
func filterTenor(criteria BorrowerCriteria, queryLoan *string) {
	if len(criteria.Tenor) > 0 {
		addedQuery := " and loan.tenor in ("
		for i := 0; i < len(criteria.Tenor); i++ {
			if i != 0 {
				addedQuery += ","
			}
			addedQuery += strconv.Itoa(criteria.Tenor[i])
		}
		addedQuery += ")"
		*queryLoan += addedQuery
	}
}

//TODO move to go-loan
func filterArea(criteria BorrowerCriteria, queryLoan *string) {
	if len(criteria.Area) > 0 {
		addedQuery := ` and rab."areaId" in (`
		for i := 0; i < len(criteria.Area); i++ {
			if i != 0 {
				addedQuery += ","
			}
			addedQuery += strconv.Itoa(criteria.Area[i].ID)
		}
		addedQuery += ")"
		*queryLoan += addedQuery
	}
}

//TODO move to go-loan
func filterCreditScoreGrade(criteria BorrowerCriteria, queryLoan *string)  {
	if len(criteria.CreditScoreGrade) > 0 {
		addedQuery := ` and loan."creditScoreGrade" in (`
		for i := 0; i < len(criteria.CreditScoreGrade); i++ {
			if i != 0 {
				addedQuery += ","
			}
			addedQuery += `'`+criteria.CreditScoreGrade[i]+`'`
		}
		addedQuery += ")"
		*queryLoan += addedQuery
	}
}

//TODO move to go-loan
func filterSector(criteria BorrowerCriteria, queryLoan *string)  {
	if len(criteria.Sector) > 0 {
		addedQuery := ` and rls."sectorId" in (`
		for i := 0; i < len(criteria.Sector); i++ {
			if i != 0 {
				addedQuery += ","
			}
			addedQuery += strconv.Itoa(criteria.Sector[i].ID)
		}
		addedQuery += ")"
		*queryLoan += addedQuery
	}
}