package plottingBorrower

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/go-mis/config"
	"bitbucket.org/go-mis/modules/investor"
	"bitbucket.org/go-mis/modules/loan-history"
	"bitbucket.org/go-mis/modules/loan-order"
	"bitbucket.org/go-mis/modules/account"
	"bitbucket.org/go-mis/modules/r"
	"bitbucket.org/go-mis/services"

	"gopkg.in/kataras/iris.v4"
	"math/rand"
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

	isActive := ctx.URLParam("isActive")

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
			where "borrowerCriteria" <> '{}' and 
			investor."deletedAt" is null`
	if isActive == "true" {
		query += ` and "isBorrowerCriteriaActive" = true`
	}
	if err := services.DBCPsql.Raw(query).Scan(&investors).Error; err != nil {
		ctx.JSON(iris.StatusOK, iris.Map{
			"status":  "Error",
			"message": "Internal server error",
		})
		return
	}

	lenData := len(investors)

	if lenData > 0 {
		totalRows = lenData
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
	where investor.id=?  and 

	investor."deletedAt" is null`

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

// FindRecommendedLoanByInvestorCriteria - this functoin fetch all loan by criteria
func FindPlottingBorrower(ctx *iris.Context) {

	stageParam := ctx.Param("stage")
	investorIdParams := ctx.URLParam("investorId")
	investorId := 0
	var err error

	stage := ""
	if stageParam == "investor" {
		investorId, err = strconv.Atoi(investorIdParams)
		if investorId <= 0 || investorIdParams == "" || err != nil {
			ctx.JSON(iris.StatusBadRequest, iris.Map{
				"message":      "Bad Request",
				"errorMessage": "Invalid User ID",
			})
			return
		}
		stage = "PRIVATE-INVESTOR"
	} else if stageParam == "marketplace" {
		stage = "PRIVATE-MARKETPLACE"
	} else {
		ctx.JSON(iris.StatusBadRequest, iris.Map{
			"message":      "Bad Request",
			"errorMessage": "Invalid Stage",
		})
		return
	}

	loans := []RecommendedLoan{}

	query := `select loan.id as "loanId",cif."name" as "borrowerName","group"."name" as "group",
    branch."name" as "branch",disbursement."disbursementDate"::date as "disbursementDate",loan.plafond,loan.rate,loan.tenor,loan."creditScoreGrade",loan.purpose 
    from loan
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
    join r_investor_product_pricing_loan rippl on rippl."loanId"=loan.id`

	fmt.Printf("ancok %s %d %s ", stage, investorId, investorIdParams)

	if investorId > 0 {
		query += ` join investor on investor.id = rippl."investorId"
		where loan.stage=? and loan."deletedAt" is null and investor.id =?`
		services.DBCPsql.Raw(query, stage, investorId).Scan(&loans)
		fmt.Printf("bangsed %s %d", stage, investorId)
	} else {
		query += ` where loan.stage=? and loan."deletedAt" is null`
		services.DBCPsql.Raw(query, stage).Scan(&loans)
	}

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   loans,
	})

}

func GetSchedulerHistory(ctx *iris.Context) {
	date := ctx.Param("date")
	loans := []SchedulerLoan{}
	query := `select loan.id as "loanId",cif."name" as "borrowerName","group"."name" as "group",
        branch."name" as "branch",loan.plafond,loan.rate,loan."creditScoreGrade",loan_history."createdAt" as "schedulerTime" from loan
        join r_loan_history on r_loan_history."loanId"=loan.id
        join loan_history on loan_history.id=r_loan_history."loanHistoryId"
        join r_loan_borrower rlb on rlb."loanId"=loan.id
        join r_cif_borrower rcb on rcb."borrowerId"=rlb."borrowerId"
        join cif on cif.id=rcb."cifId"
        join r_loan_group rlg on rlg."loanId"=loan.id
        join "group" on "group".id=rlg."groupId"
        join r_loan_branch rlbr on rlbr."loanId"=loan.id
        join branch on branch.id=rlbr."branchId"
        where loan_history."stageFrom"='PRIVATE-MARKETPLACE' and loan_history."stageTo"='MARKETPLACE' and loan_history."createdAt"::date=?`
	services.DBCPsql.Raw(query, date).Scan(&loans)
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   loans,
	})
}

// UpdateLoanStageHandler - update loan stage handler
func UpdateLoanStageHandler(ctx *iris.Context) {
	loanResponse := struct {
		Status  int64  `json:"status"`
		Message string `json:"message"`
	}{}
	payload := UpdateStageRequest{}
	if err := ctx.ReadJSON(&payload); err != nil {
		log.Println("Error",err)
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "Error",
			"message": "Invalid json",
		})
		return
	}
	log.Println("Payload: ", payload)
	if !isValidStage(payload.StageFrom) {
		ctx.JSON(iris.StatusBadRequest, iris.Map{
			"status":  "Error",
			"message": "Invalid stage from",
		})
		return
	}
	if !isValidStage(payload.StageTo) {
		ctx.JSON(iris.StatusBadRequest, iris.Map{
			"status":  "Error",
			"message": "Invalid stage to",
		})
		return
	}
	if payload.StageTo=="INVESTOR"{
		if err:=updateToInvestor(payload);err!=nil{
			ctx.JSON(iris.StatusBadRequest, iris.Map{
				"status":  "Error",
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(iris.StatusOK, iris.Map{
			"status":  "Success",
			"message": "Succeed: update loan stage from:" + payload.StageFrom + " to:" + payload.StageTo,
		})
		return
	}
	if !isUsingInvestorID(payload.StageTo) {
		payload.InvestorId = 0
	}
	goLoanURI := config.Configuration.GoLoanPath + "/" + "loan/plotting-borrower/update/loan-stage/"
	bPayload, err := json.Marshal(&payload)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, iris.Map{
			"status":  "Internal server error",
			"message": "Invalid stage from",
		})
		return
	}
	loanPayload := strings.NewReader(string(bPayload))
	respBytes, err := services.CircuitBreaker.Put(goLoanURI, loanPayload)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, iris.Map{
			"status":  "Internal server error",
			"message": "Invalid stage from",
		})
		return
	}
	err = json.Unmarshal(respBytes, &loanResponse)
	log.Printf("Loan response %+v", loanResponse)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, iris.Map{
			"status":  "Internal server error",
			"message": "Invalid stage to",
		})
		return
	}
	ctx.JSON(iris.StatusOK, iris.Map{
		"status":  "Success",
		"message": "Succeed: update loan stage from:" + payload.StageFrom + " to:" + payload.StageTo,
	})
}
func updateToInvestor(payload UpdateStageRequest) error{
	remarkFlag := "PENDING"
	currentBalance := GetCurrentBalance(payload.InvestorId)
	if currentBalance < payload.Amount {
		fmt.Println("#ERROR#Balance#notenough",payload)
		return errors.New("Balance is not enough")
	}
	orderNo := generateOrderNumber(payload.InvestorId)
	db := services.DBCPsql.Begin()
	order := &loanOrder.LoanOrder{Remark: remarkFlag, OrderNo: orderNo}
	db.Table("loan_order").Create(order)
	updateStageQuery := `UPDATE loan SET "stage" = 'ORDERED' WHERE loan."id" IN ( ` + getStrLoanId(payload.LoanId) + ` )`
	if err := db.Exec(updateStageQuery); err != nil && err.RowsAffected == 0 {
		db.Rollback()
		return errors.New("Error update loan stage to ordered")
	}
	for _,loanId := range payload.LoanId {
		// INSERT TO R LOAN ORDER
		insertRloQuery := `INSERT INTO r_loan_order("loanId","loanOrderId", "createdAt", "updatedAt") VALUES(?,?,current_timestamp,current_timestamp)`
		if err := db.Exec(insertRloQuery, loanId, order.ID); err != nil && err.RowsAffected == 0 {
			db.Rollback()
			return errors.New("Error insert r_loan_order")
		}
		// INSERT TO LOAN HISTORY
		investorID := strconv.FormatUint(payload.InvestorId, 10)
		loanHistorySchema := &loanHistory.LoanHistory{StageFrom: "CART", StageTo: "ORDERED", Remark: "ORDERED loanId=" + fmt.Sprintf("%v", loanId) + " investorId=" + investorID}
		db.Table("loan_history").Create(loanHistorySchema)

		// INSERT TO R_LOAN HISTORY
		rLoanHistorySchema := &r.RLoanHistory{LoanId: uint64(loanId), LoanHistoryId: loanHistorySchema.ID}
		db.Table("r_loan_history").Create(rLoanHistorySchema)
	}
	db.Commit()
	//make new transaction so order still created no matter accept failed or success
	newDb:=services.DBCPsql.Begin()
	if err:=loanOrder.AcceptOrder(orderNo,false,newDb);err!=nil{
		db.Rollback()
		fmt.Println("ERROR-ACCEPT-ORDER",err)
		return err
	}
	newDb.Commit()
	return nil
}


func getStrLoanId(LoanId []uint64) string{
	result:=""
	for i:=0;i<len(LoanId);i++{
		if i != 0 {
			result += ","
		}
		result += strconv.FormatUint(LoanId[i],10)
	}
	return result
}

//Generate OrderId by [YMDHIS][XXX][INVESTORID]
func generateOrderNumber(investorID uint64) string {
	now := time.Now()
	numRand := rand.Intn(999)
	rand.Seed(time.Now().UTC().UnixNano())
	timestamp := fmt.Sprintf("%d%d%d%d%d%d%d%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), numRand, investorID)
	return timestamp[2:]
}

func GetCurrentBalance(InvestorID uint64) float64 {
	account := account.Account{}
	accountDb := services.DBCPsql.Table("account").Select("*")
	accountDb = accountDb.Joins(`JOIN r_account_investor rai on rai."accountId" = account.id`)
	accountDb = accountDb.Where(`rai."investorId" = ?`, InvestorID)
	accountDb.First(&account)

	totalBalance := account.TotalBalance
	return totalBalance

}

func isValidStage(stage string) bool {
	log.Println("Stage: ", stage)
	stages := []string{"PRIVATE", "PRIVATE-INVESTOR", "INVESTOR", "MARKETPLACE", "PRIVATE-MARKETPLACE"}
	for i := range stages {
		if strings.ToUpper(stages[i]) == strings.ToUpper(stage) {
			return true
		}
	}
	return false
}

func isUsingInvestorID(stage string) bool {
	fmt.Println("STAGE",stage)
	invStages := []string{"PRIVATE-INVESTOR", "INVESTOR"}
	for i := range invStages {
		if strings.ToUpper(stage) == invStages[i] {
			return true
		}
	}
	return false
}

// this functoin fetch all loan by criteria
func FindRecommendedLoanByInvestorCriteria(ctx *iris.Context) {
	investorID := ctx.Param("investorId")
	disFrom := ctx.URLParam("disFrom")
	disTo := ctx.URLParam("disTo")
	resultGoloan := make([]RecommendedLoan, 0)

	//redisLoan, err := RetriveRecommendedLoanFromRedis(investorID)
	//if err != nil {
	//	log.Println("[ERROR] ", err)
	//}
	//log.Println("Investor ID: ", investorID)
	//if len(redisLoan) > 0 {
	//
	//	if investorID != "-1" {
	//		disToDate, errToDate := time.Parse("2006-01-02", disTo)
	//		disFromDate, errFromDate := time.Parse("2006-01-02", disFrom)
	//
	//		if errToDate != nil {
	//			log.Println("Error disto: ", errToDate)
	//			ctx.JSON(http.StatusInternalServerError, iris.Map{
	//				"status":  "Error",
	//				"message": errToDate.Error(),
	//			})
	//			return
	//		}
	//
	//		if errFromDate != nil {
	//			log.Println("Error disfrom: ", errFromDate)
	//			ctx.JSON(http.StatusInternalServerError, iris.Map{
	//				"status":  "Error",
	//				"message": errFromDate.Error(),
	//			})
	//			return
	//		}
	//
	//		resultGoloanByDate := make([]RecommendedLoan, 0, len(redisLoan))
	//
	//		for _, loan := range redisLoan {
	//			loandDate, errLoanDate := time.Parse("2006-01-02T15:04:05Z", loan.DisbursementDate)
	//			if errLoanDate != nil {
	//				log.Println("errLoanDate: ", errLoanDate)
	//				ctx.JSON(http.StatusInternalServerError, iris.Map{
	//					"status":  "Error",
	//					"message": errLoanDate.Error(),
	//				})
	//				return
	//			}
	//
	//			if loandDate.After(disFromDate) && loandDate.Before(disToDate) {
	//				log.Printf("disfrom %s to disto %s and loand date: %s , result: %t ", disFrom, disTo, loandDate, (loandDate.After(disFromDate) && loandDate.Before(disToDate)))
	//				if loan.LoanId > 0 {
	//					log.Println("Loan ID: ", loan.LoanId)
	//					resultGoloanByDate = append(resultGoloanByDate, loan)
	//				}
	//			}
	//
	//			ctx.JSON(http.StatusOK, iris.Map{
	//				"status": "Success",
	//				"data":   resultGoloanByDate,
	//			})
	//			return
	//		}
	//	} else {
	//		ctx.JSON(http.StatusOK, iris.Map{
	//			"status": "Success",
	//			"data":   redisLoan,
	//		})
	//		return
	//	}
	//}

	resultGoloan, err := RetrieveRecommendedLoanFromLoanService(disFrom, disTo, investorID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, iris.Map{
			"status":  "Error",
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, iris.Map{
		"status": "Success",
		"data":   resultGoloan,
	})
}

// RetrieveRecommendedLoanFromLoanService - get all data recomended loan from loan service
func RetrieveRecommendedLoanFromLoanService(disFrom, disTo, investorID string) ([]RecommendedLoan, error) {
	var goloanResp GOLoanSuccessResponse
	goLoanURI := config.Configuration.GoLoanPath + "/" + "loan/plotting-borrower/recomended-loan-investor/" + investorID + "?disFrom=" + disFrom + "&disTo=" + disTo
	fmt.Println("GOLOAN URI: ", goLoanURI)
	body, err := services.CircuitBreaker.Get(goLoanURI)
	if err != nil {
		log.Println("[ERROR] Calling API to go-loan service", err)
		return nil, err
	}
	err = json.Unmarshal(body, &goloanResp)
	if err != nil {
		log.Println("[ERROR] Json unmarshall", err)
		return nil, err
	}
	if goloanResp.Code != 200 && strings.ToUpper(goloanResp.Message) != "SUCCESS" {
		return nil, errors.New("Unable to get recomended loan data from go loan service")
	}
	go func(data []RecommendedLoan) {
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

// RetriveRecommendedLoanFromRedis - get data recomened loan from redis
// wrapped wheter is all or specifig by investor id
func RetriveRecommendedLoanFromRedis(investorID string) ([]RecommendedLoan, error) {
	var err error
	loanRedis := make([]RecommendedLoan, 0)
	switch strings.ToUpper(strings.TrimSpace(investorID)) {
	case "-1":
		loanRedis, err = FindAllRecommendedLoanFromRedis()
	default:
		loanRedis, err = FindRecommendedLoanFromRedis(investorID)
	}
	return loanRedis, err
}

// FindAllRecommendedLoanFromRedis - find all recomended loan from redis
func FindAllRecommendedLoanFromRedis() ([]RecommendedLoan, error) {
	loans := make([]RecommendedLoan, 0)
	redisClient, err := services.NewClientRedis()
	if err != nil {
		return nil, err
	}
	strData, err := redisClient.GetAllRecomendedLoan()
	if err != nil {
		return nil, err
	}
	for i := range strData {
		var recLoan RecommendedLoan
		err = json.Unmarshal([]byte(strData[i]), &recLoan)
		if err != nil {
			return nil, err
		}
		loans = append(loans, recLoan)
	}
	return loans, nil
}

// FindRecommendedLoanFromRedis - find all recomended loan from redis by investor id
func FindRecommendedLoanFromRedis(investorID string) ([]RecommendedLoan, error) {
	loanRedis := make([]RecommendedLoan, 0)
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
