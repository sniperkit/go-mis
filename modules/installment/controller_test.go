package installment

import (
	"testing"
	iris "gopkg.in/kataras/iris.v4"
	"github.com/valyala/fasthttp"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os/exec"
	"bitbucket.org/go-mis/services"
)


func TestSubmitInstallmentByInstallmentIDWithStatus(t *testing.T) {
	framework := iris.New()

	ctx := framework.AcquireCtx(&fasthttp.RequestCtx{});
	ctx.Set("installment_id", "978625")
	ctx.Set("status", "success")

	SubmitInstallmentByInstallmentIDWithStatus(ctx);
}

func TestStoreInstallment(t *testing.T) {
	db := services.DBCPsql.Begin()
	StoreInstallment(db, 978626, "success")
	db.Commit()
}

func TestUpdateLoanStage(t *testing.T) {

	// initial database
	if err := exec.Command("sh", "db.sh").Run(); err != nil {
		t.Error(err)
	}

	// create db for mock
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=egon dbname=amartha_test password=nakal23baik sslmode=disable")
	if err != nil {
		t.Error(err)
	}


	loan := LoanSchema{}
	if err := db.Table("loan").First(&loan).Error; err != nil {
		t.Error(err)
	}

	if loan.Stage != "INSTALLMENT" {
		t.Error("Loan is not installment");
	}

	// create installment
	installment := Installment{}
	if err := db.Table("installment").Where("id = 50").First(&installment).Error; err != nil {
		t.Error(err)
	}

	if err := UpdateLoanStage(installment, 1, db); err != nil {
		t.Error(err)
	}

	if err := db.Table("loan").First(&loan).Error; err != nil {
		t.Error(err)
	}
	if loan.Stage != "END" {
		t.Error("Loan is not change to END");
	}


	if err := db.Table("loan").Where("loan.id = 2").Scan(&loan).Error; err != nil {
		t.Error(err)
	}

	if loan.Stage != "INSTALLMENT" {
		t.Error("Loan is not installment");
	}

	// create installment
	if err := db.Table("installment").Where("id = 73").Scan(&installment).Error; err != nil {
		t.Error(err)
	}

	if err := UpdateLoanStage(installment, 2, db); err != nil {
		t.Error(err)
	}

	if err := db.Table("loan").First(&loan).Error; err != nil {
		t.Error(err)
	}

	if loan.Stage != "END-EARLY" {
		t.Error("Loan is not change to END-EARLY");
	}

	var count int32
	if err := db.Table("loan_history").Where("remark = 'Automatic update stage END loanId = 1'").Count(&count).Error; err != nil {
		t.Error(err)
	}
	
	if count != 1 {
		t.Error("loan History loanId 1 not exists")
	}

	if err := db.Table("loan_history").Where("remark = 'Automatic update stage END-EARLY loanId = 2'").Count(&count).Error; err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("loan History loanId 2 not exists")
	}




	// create installment
	if err := db.Table("installment").Where("id = 99").Scan(&installment).Error; err != nil {
		t.Error(err)
	}

	if err := UpdateLoanStage(installment, 3, db); err == nil {
		t.Error("It should be error")
	}

	if err := db.Table("loan").Where("id = 3").Scan(&loan).Error; err != nil {
		t.Error(err)
	}

	if loan.Stage != "END-PENDING" {
		t.Error("Loan is not change to END-PENDING");
	}

	if err := db.Table("loan_history").Where("remark = 'Automatic update stage END loanId = 3'").Count(&count).Error; err != nil {
		t.Error(err)
	}

	// Update Loan Stage Meninggal
	if err := db.Table("installment").Where("id = 100").Scan(&installment).Error; err != nil {
		t.Error(err)
	}

	if err := UpdateLoanStage(installment, 4, db); err != nil {
		t.Error(err)
		t.Error("Update Loan Stage 4 error")
	}

	if err := db.Table("loan").Where("id = 4").Scan(&loan).Error; err != nil {
		t.Error(err)
	}

	if loan.Stage != "MENINGGAL" {
		t.Error("Loan 4 is not change to MENINGGAL");
	}

	if err := db.Table("loan_history").Where("remark = 'Automatic update stage MENINGGAL loanId = 4'").Count(&count).Error; err != nil {
		t.Error(err)
	}
}
