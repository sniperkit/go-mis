package validationTeller

import (
	"strconv"

	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

var STAGE map[uint64]string = map[uint64]string{
	1: "AGENT",
	2: "TELLER",
	3: "PENDING",
	4: "REVIEW",
	5: "APPROVE",
}

func GetData(ctx *iris.Context) {
	params := struct {
		BranchId uint64 `json:"branchId"`
		Date     string `json:"date"`
	}{}

	err := ctx.ReadJSON(&params)

	if err != nil {
		ctx.JSON(iris.StatusBadRequest, iris.Map{"status": "error", "message": err.Error()})
		return
	}
	query := `select i.id as "installmentId", a.fullname,g.name, sum(i."paidInstallment") "repayment",sum(i.reserve) "tabungan",sum(i."paidInstallment"+i.reserve) "total", 
				coalesce(sum(case
				when d.stage = 'SUCCESS' then plafond end
				),0) "totalCair",
				coalesce(sum(case
				when d.stage = 'FAILED' then plafond end
				),0) "totalGagalDropping"
				from loan l join r_loan_group rlg on l.id = rlg."loanId"
				join "group" g on g.id = rlg."groupId"
				join r_group_agent rga on g.id = rga."groupId"
				join agent a on a.id = rga."agentId"
				join r_loan_branch rlb on rlb."loanId" = l.id
				join branch b on b.id = rlb."branchId"
				join r_loan_installment rli on rli."loanId" = l.id
				join installment i on i.id = rli."installmentId"
				join r_loan_disbursement rld on rld."loanId" = l.id
				join disbursement d on d.id = rld."disbursementId"
				where l."deletedAt" isnull and b.id= ? and coalesce(i."transactionDate",i."createdAt")::date = ? and l.stage = 'INSTALLMENT'
				group by g.name, a.fullname
				order by a.fullname`

	result := []ValidationTellerData{}
	services.DBCPsql.Raw(query, params.BranchId, params.Date).Scan(&result)
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   result,
	})
}

func UpdateInstallmentStage(ctx *iris.Context) {
	params := struct {
		InstallmentID uint64 `json:"installmentId"`
		Stage         uint64 `json:"stage"`
	}{}

	err := ctx.ReadJSON(&params)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, iris.Map{"status": "error", "message": err.Error()})
		return
	}

	if err := db.Table("installment").Where("id = ?", params.InstallmentID).UpdateColumn("stage", STAGE[stage]).Error; err != nil {
		return err
	}

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   "Installment id:" + strconv.Itoa(params.InstallmentID) + " updated. Stage:" + STAGE[params.Stage],
	})
}
