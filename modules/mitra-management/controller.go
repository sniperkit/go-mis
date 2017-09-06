package mitramanagement

import (
	"errors"
	"log"
	"strings"

	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

func Init() {
	services.DBCPsql.AutoMigrate(&Status{})
	services.DBCPsql.AutoMigrate(&Reason{})
}

const (
	selectBorrowerQuery = ` select l.id as "loanId",
								borrower."borrowerNo" as "borrowerNumber",
								cif."name" as "borrowerName",
								"group"."name" as "groupName",
								reason.description as "reason" `

	selectDetailDO = ` , l.plafond as "plafond", l.tenor as "tenor",
							installment."createdAt" as "doDate",
							installment."type" as "type" `
	selectDetailPAR = `, installment."paidInstallment" as "nominal",
							installment."createdAt" as "parDate",
							agent.fullname as "agent" `
	selectDetailTR = ` , installment."paidInstallment" as "nominal",
							installment."createdAt" as "trDate",
							agent.fullname as "agent",
							installment."type" as "type" `
	fromQuery = ` from loan l
					join r_loan_borrower rlb on rlb."loanId" = l.id
					join borrower on borrower.id = rlb."borrowerId"
					join r_loan_installment rli on rli."loanId" = l.id
					join installment on installment.id = rli."installmentId"
					join r_group_borrower rgb on rgb."borrowerId" = borrower.id
					join "group" on "group".id = rgb."groupId"
					join r_cif_borrower rcb on rcb."borrowerId" = borrower.id
					join cif on cif.id = rcb."cifId"
					join r_loan_branch rlbranch on rlbranch."loanId" = l.id
					join branch on branch.id = rlbranch."branchId"
					left join status on status.id = installment."statusId"
					left join reason on reason.id = installment."reasonId" 
					left join r_group_agent rga on rga."groupId" = "group".id
					left join agent on agent.id = rga."agentId" `
	whereQuery     = ` where upper(status.type) = 'MITRA_MANAGEMENT' and branch.id = ? `
	borrowerQuery  = selectBorrowerQuery + fromQuery + whereQuery
	dODetailQuery  = selectBorrowerQuery + selectDetailDO + fromQuery + whereQuery
	pARDetailQuery = selectBorrowerQuery + selectDetailPAR + fromQuery + whereQuery
	tRDetailQuery  = selectBorrowerQuery + selectDetailTR + fromQuery + whereQuery
)

func GetBorrowerByInstallmentTypeAndDate(ctx *iris.Context) {
	var borrList []MMBorrower
	branchID := ctx.Get("BRANCH_ID")
	installmentType := ctx.Param("borrowerType")
	date := ctx.Param("date")
	log.Println("Installment type: ", installmentType)
	log.Println("date: ", date)
	_, err := FindBorrowerByInstallmentType(&borrList, branchID, installmentType, date)
	if err != nil {
		errorResponse(ctx, err.Error(), iris.StatusInternalServerError)
	}
	successResponse(ctx, borrList)
}

func GetBorrowerDetailByInstallmentTypeAndDate(ctx *iris.Context) {
	branchID := ctx.Get("BRANCH_ID")
	installmentParam, _ := ctx.URLParamInt64("installmentID")
	installmentID := uint64(installmentParam)
	installmentType := ctx.URLParam("status")
	date := ctx.URLParam("date")
	switch strings.ToUpper(installmentType) {
	case "DO":
		var doDetails []MMDOBorrower
		_, err := FindDODetailBorrower(&doDetails, branchID, date, installmentID)
		if err != nil {
			errorResponse(ctx, err.Error(), iris.StatusInternalServerError)
		}
		successResponse(ctx, doDetails)
	case "PAR":
		var parDetails []MMPARBorrower
		_, err := FindPARDetailBorrower(&parDetails, branchID, date, installmentID)
		if err != nil {
			errorResponse(ctx, err.Error(), iris.StatusInternalServerError)
		}
		successResponse(ctx, parDetails)
	case "TR":
		var trDetails []MMTRBorrower
		_, err := FindTRDetailBorrower(&trDetails, branchID, date, installmentID)
		if err != nil {
			errorResponse(ctx, err.Error(), iris.StatusInternalServerError)
		}
		successResponse(ctx, trDetails)
	}
}

func GetBorrowerStatusReason(ctx *iris.Context) {
	var reasons []Reason
	statusID := ctx.Param("statusId")
	_, err := FindBorrowerStatusReason(&reasons, statusID)
	if err != nil {
		errorResponse(ctx, err.Error(), iris.StatusInternalServerError)
	}
	successResponse(ctx, reasons)
}

func FindBorrowerStatusReason(reasons *[]Reason, statusID string) (int, error) {
	query := ` select id, "statusId", description from reason where "statusId" = ? `
	err := services.DBCPsql.Raw(query, statusID).Scan(reasons).Error
	if err != nil {
		return 0, errors.New("Unable to retrive reason data")
	}
	return len(*reasons), nil
}

func FindBorrowerByInstallmentType(borrowers *[]MMBorrower, branchID interface{}, installmentType, date string) (int, error) {
	expQuery, err := findBorrowerQuery(installmentType)
	log.Println(expQuery)
	if err != nil {
		return 0, err
	}
	expQuery += ` and installment."createdAt"::date = ? order by cif."name" ASC `
	log.Println(expQuery)
	db := services.DBCPsql.Raw(expQuery, branchID, date).Scan(borrowers)
	if db.Error != nil {
		return 0, errors.New("Unable to retrive mitra management borrower")
	}
	return len(*borrowers), nil
}

func FindDODetailBorrower(doDetails *[]MMDOBorrower, branchID interface{}, date string, installmentID uint64) (int, error) {
	query, err := findBorrowerDetailQuery("DO", installmentID)
	if err != nil {
		return 0, err
	}
	query += ` and installment."createdAt"::date = ? order by cif."name" ASC `
	err = services.DBCPsql.Raw(query, branchID, date).Scan(doDetails).Error
	if err != nil {
		return 0, errors.New("Unable to retrive data Borrower DO detail")
	}
	return len(*doDetails), err
}

func FindPARDetailBorrower(parDetails *[]MMPARBorrower, branchID interface{}, date string, installmentID uint64) (int, error) {
	query, err := findBorrowerDetailQuery("PAR", installmentID)
	if err != nil {
		return 0, err
	}
	query += ` and installment."createdAt"::date = ? order by cif."name" ASC `
	err = services.DBCPsql.Raw(query, branchID, date).Scan(parDetails).Error
	if err != nil {
		return 0, err
	}
	return len(*parDetails), nil
}

func FindTRDetailBorrower(trDetails *[]MMTRBorrower, branchID interface{}, date string, installmentID uint64) (int, error) {
	query, err := findBorrowerDetailQuery("TR", installmentID)
	if err != nil {
		return 0, err
	}
	err = services.DBCPsql.Raw(query, branchID, date).Scan(trDetails).Error
	if err != nil {
		return 0, err
	}
	return len(*trDetails), nil
}

func findBorrowerQuery(installmentType string) (string, error) {
	switch strings.ToUpper(installmentType) {
	case "PAR":
		return borrowerQuery + ` and upper(installment."type") = 'PAR' `, nil
	case "TR":
		return borrowerQuery + ` and (upper(installment."type") = 'TR1' OR upper(installment."type") = 'TR2' ) `, nil
	case "DO":
		return borrowerQuery + ` and upper(installment."type") = 'DROPOUT' `, nil
	default:
		return "", errors.New("Invalid Installment Type")

	}
}

func findBorrowerDetailQuery(installmentType string, installmentID uint64) (string, error) {
	switch strings.ToUpper(installmentType) {
	case "PAR":
		return pARDetailQuery + ` and upper(installment."type") = 'PAR' and and installment.id = ? `, nil
	case "TR":
		return tRDetailQuery + ` and (upper(installment."type") = 'TR1' OR upper(installment."type") = 'TR2' ) and installment.id = ? `, nil
	case "DO":
		return dODetailQuery + ` and upper(installment."type") = 'DROPOUT' and installment.id = ? `, nil
	default:
		return "", errors.New("Invalid Installment Type")

	}
}

func successResponse(ctx *iris.Context, data interface{}) {
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "Success",
		"data":   data,
	})
	return
}

func errorResponse(ctx *iris.Context, data interface{}, errStatus int) {
	ctx.JSON(errStatus, iris.Map{
		"status": "Success",
		"data":   data,
	})
	return
}
