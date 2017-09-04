package mitramanagement

import (
	"errors"
	"strings"

	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

func Init() {
	services.DBCPsql.AutoMigrate(&Status{})
	services.DBCPsql.AutoMigrate(&Reason{})
}

const (
	baseQuery = ` select l.id,
					borrower."borrowerNo",
					cif."name" as "borrowerName",
					"group"."name" as "groupName",
					reason.description as "reason"
				from loan l
					join r_loan_borrower rlb on rlb."loanId" = l.id
					join borrower on borrower.id = rlb."borrowerId"
					join r_loan_installment rli on rli."loanId" = l.id
					join installment on installment.id = rli."installmentId"
					join r_group_borrower rgb on rgb."borrowerId" = borrower.id
					join "group" on "group".id = rgb."groupId"
					join r_cif_borrower rcb on rcb."borrowerId" = borrower.id
					join r_loan_branch rlb on rlb.loanId = l.id
					join branch on branch.id = rlb.branchId
					left join status on status.id = installment.status_id
					left join reason on reason.id = installment.reason_id
				where upper(status.type) = 'MITRA_MANAGEMENT' `
)

func GetBorrowerByInstallmentType(ctx *iris.Context) {
	var borrList []MMBorrower
	branchID := ctx.Get("BRANCH_ID")
	installmentType := ctx.URLParam("status")
	date := ctx.URLParam("date")
	_, err := FindBorrowerByInstallmentType(&borrList, branchID, installmentType, date)
	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":       "error",
			"errorMessage": "Internal Server Error",
		})
	}
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   borrList,
	})
}

func GetBorrowerStatusDescription(ctx *iris.Context) {
	var reasons []Reason
	statusID := ctx.URLParam("statusId")
	_, err := FindBorrowerStatusDescription(&reasons, statusID)
	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":       "error",
			"errorMessage": "Internal Server Error",
		})
	}
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "Success",
		"data":   reasons,
	})
}

func FindBorrowerStatusDescription(reasons *[]Reason, statusID string) (int, error) {
	query := ` select id, description from reason where "statusId" = ? `
	err := services.DBCPsql.Raw(query, statusID).Scan(reasons).Error
	if err != nil {
		return 0, errors.New("Unable to retrive reason data")
	}
	return len(*reasons), nil
}

func FindBorrowerByInstallmentType(borrowers *[]MMBorrower, branchID interface{}, installmentType, date string) (int, error) {
	var expQuery string
	switch strings.ToUpper(installmentType) {
	case "PAR":
		expQuery = baseQuery + ` and upper(installment."type") = 'PAR' `
	case "TR":
		expQuery = baseQuery + ` and (upper(installment."type") = 'TR1' OR upper(installment."type") = 'TR2' `
	case "DO":
		expQuery = baseQuery + ` and upper(installment."type") = 'DROPOUT' `
	default:
		return 0, errors.New("Invalid Installment Type")

	}
	expQuery += ` and installment."createdAt" = ? and branch.id = ? group by l.id order by borrower.name ASC `
	db := services.DBCPsql.Raw(expQuery, date, branchID).Scan(borrowers)
	if db.Error != nil {
		return 0, errors.New("Unable to retrive mitra management borrower")
	}
	return len(*borrowers), nil
}
