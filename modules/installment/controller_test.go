package installment

import (
	"testing"
	iris "gopkg.in/kataras/iris.v4"
	"github.com/valyala/fasthttp"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os/exec"
)


func TestSubmitInstallmentByInstallmentIDWithStatus(t *testing.T) {
	t.Log("kucing");

	framework := iris.New()

	ctx := framework.AcquireCtx(&fasthttp.RequestCtx{});
	ctx.Set("installment_id", "978625")
	ctx.Set("status", "success")

	SubmitInstallmentByInstallmentIDWithStatus(ctx);
}

func TestStoreInstallment(t *testing.T) {
	StoreInstallment(978625, "success")
}

func TestUpdateLoanStageNormal(t *testing.T) {

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

}
