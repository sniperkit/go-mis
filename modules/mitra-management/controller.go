package mitramanagement

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jinzhu/gorm"

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
								borrower."id" as "borrowerId",
								cif."name" as "borrowerName",
								"group"."name" as "groupName",
								reason.description as "reason",
								reason.id as "reasonId",
								installment.id as "installmentId" `

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
                    join r_loan_group rlg on rlg."loanId" = rli."loanId"
                    join "group" on "group".id = rlg."groupId"
                    join r_cif_borrower rcb on rcb."borrowerId" = borrower.id
                    join cif on cif.id = rcb."cifId"
                    join r_loan_branch rlbranch on rlbranch."loanId" = l.id
                    join branch on branch.id = rlbranch."branchId"
                    left join status on status.id = installment."statusId"
                    left join reason on reason.id = installment."reasonId"
                    left join r_group_agent rga on rga."groupId" = "group".id
                    left join agent on agent.id = rga."agentId"   `
	whereQuery     = ` where branch.id = ? `
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
		var DODetail MMDOBorrower
		err := FindDODetailBorrower(&DODetail, branchID, date, installmentID)
		if err != nil {
			errorResponse(ctx, err.Error(), iris.StatusInternalServerError)
		}
		successResponse(ctx, DODetail)
	case "PAR":
		var PARDetail MMPARBorrower
		err := FindPARDetailBorrower(&PARDetail, branchID, date, installmentID)
		if err != nil {
			errorResponse(ctx, err.Error(), iris.StatusInternalServerError)
		}
		successResponse(ctx, PARDetail)
	case "TR":
		var TRDetails MMTRBorrower
		err := FindTRDetailBorrower(&TRDetails, branchID, date, installmentID)
		fmt.Println(TRDetails)
		if err != nil {
			errorResponse(ctx, err.Error(), iris.StatusInternalServerError)
		}
		successResponse(ctx, TRDetails)
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
	if err != nil {
		return 0, err
	}
	expQuery += ` and installment."createdAt"::date = ? and installment."deletedAt" is null order by cif."name" ASC `
	log.Println("QUERY TO RUN: ", expQuery)
	db := services.DBCPsql.Raw(expQuery, branchID, date).Scan(borrowers)
	if db.Error != nil {
		return 0, errors.New("Unable to retrive mitra management borrower")
	}
	return len(*borrowers), nil
}

func FindDODetailBorrower(doDetails *MMDOBorrower, branchID interface{}, date string, installmentID uint64) error {
	query, err := findBorrowerDetailQuery("DO", installmentID)
	log.Println("Query to run: ", query)
	if err != nil {
		return err
	}
	query += ` and installment."createdAt"::date = ? order by cif."name" ASC `
	err = services.DBCPsql.Raw(query, branchID, installmentID, date).Scan(doDetails).Error
	if err != nil {
		return errors.New("Unable to retrive data Borrower DO detail")
	}
	return nil
}

func FindPARDetailBorrower(parDetails *MMPARBorrower, branchID interface{}, date string, installmentID uint64) error {
	query, err := findBorrowerDetailQuery("PAR", installmentID)
	if err != nil {
		return err
	}
	query += ` and installment."createdAt"::date = ? order by cif."name" ASC `
	err = services.DBCPsql.Raw(query, branchID, installmentID, date).Scan(parDetails).Error
	if err != nil {
		return errors.New("Unable to retrive data Borrower PAR detail")
	}
	return nil
}

func FindTRDetailBorrower(trDetails *MMTRBorrower, branchID interface{}, date string, installmentID uint64) error {
	query, err := findBorrowerDetailQuery("TR", installmentID)
	if err != nil {
		return err
	}
	query += ` and installment."createdAt"::date = ? order by cif."name" ASC `
	err = services.DBCPsql.Raw(query, branchID, installmentID, date).Scan(trDetails).Error
	if err != nil {
		return errors.New("Unable to retrive data Borrower TR detail")
	}
	return nil
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
		return pARDetailQuery + ` and upper(installment."type") = 'PAR' and installment.id = ? `, nil
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

func GetStatusAll(ctx *iris.Context) {
	s := []Status{}
	err := services.DBCPsql.Where("type = 'mitra_management'").Find(&s).Error
	if err != nil {
		log.Println("[INFO] Params is not valid")
		ctx.JSON(iris.StatusBadRequest, iris.Map{
			"message":      "Bad Request",
			"errorMessage": err.Error(),
		})
		return
	}

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   s,
	})
}

func SubmitReason(ctx *iris.Context) {
	payload := struct {
		InstallmentID uint64 `json:"installmentId"`
		BorrowerID    uint64 `json:"borrowerId"`
		StatusID      uint64 `json:"statusId"`
		ReasonID      uint64 `json:"reasonId"`
	}{}

	err := ctx.ReadJSON(&payload)
	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// time
	t := time.Now().Format("2006-01-02 15:04:05")

	db := services.DBCPsql.Begin()
	// update installment
	q := `update Installment set "statusId" = ?, "reasonId" = ?, "updatedAt"=? where id=?`
	err = db.Exec(q, payload.StatusID, payload.ReasonID, t, payload.InstallmentID).Error
	if err != nil {
		ProcessErrorAndRollback(ctx, db, "Error Update Installment: "+err.Error())
		return
	}

	if payload.StatusID == 1 {
		q = `update borrower set "doDate" = ? where id = ?`
		err = services.DBCPsql.Exec(q, t, payload.BorrowerID).Error
		if err != nil {
			ProcessErrorAndRollback(ctx, db, "Error Update Borrower: "+err.Error())
			return
		}
	}

	db.Commit()

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   "installment data has ben updated",
	})
}

func ProcessErrorAndRollback(ctx *iris.Context, db *gorm.DB, message string) {
	db.Rollback()
	ctx.JSON(iris.StatusInternalServerError, iris.Map{
		"status":  "error",
		"message": message,
	})
}
