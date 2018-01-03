package loan

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"errors"

	"bitbucket.org/go-mis/modules/account"
	accountTransactionCredit "bitbucket.org/go-mis/modules/account-transaction-credit"
	accountTransactionDebit "bitbucket.org/go-mis/modules/account-transaction-debit"
	cif "bitbucket.org/go-mis/modules/cif"
	"bitbucket.org/go-mis/modules/installment"
	loanHistory "bitbucket.org/go-mis/modules/loan-history"
	"bitbucket.org/go-mis/modules/r"
	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
	"bitbucket.org/go-mis/modules/product-pricing"
	"github.com/jinzhu/gorm"
)

func Init() {
	services.DBCPsql.AutoMigrate(&Loan{})
	services.BaseCrudInit(Loan{}, []Loan{})
}

type TotalData struct {
	TotalRows int64 `gorm:"column:totalRows" json:"totalRows"`
}

// FetchAll - fetchAll Loan data
func FetchAll(ctx *iris.Context) {
	branchID := ctx.Get("BRANCH_ID")

	totalData := TotalData{}
	// queryTotalData := "SELECT DISTINCT COUNT(loan.*) AS \"totalRows\" "
	// queryTotalData += "FROM loan "
	// queryTotalData += "LEFT JOIN r_loan_sector ON r_loan_sector.\"loanId\" = loan.\"id\" "
	// queryTotalData += "LEFT JOIN sector ON r_loan_sector.\"sectorId\" = sector.\"id\" "
	// queryTotalData += "LEFT JOIN r_loan_borrower ON r_loan_borrower.\"borrowerId\" = loan.\"id\" "
	// queryTotalData += "LEFT JOIN borrower ON r_loan_borrower.\"borrowerId\" = borrower.\"id\" "
	// queryTotalData += "LEFT JOIN r_cif_borrower ON r_cif_borrower.\"borrowerId\" = borrower.\"id\""
	// queryTotalData += "LEFT JOIN cif ON r_cif_borrower.\"cifId\" = cif.\"id\" LEFT JOIN r_loan_group ON r_loan_group.\"loanId\" = loan.\"id\" "
	// queryTotalData += "LEFT JOIN \"group\" ON r_loan_group.\"groupId\" = \"group\".\"id\" "
	// queryTotalData += "LEFT JOIN r_loan_branch ON r_loan_branch.\"loanId\" = loan.\"id\" "
	// queryTotalData += "LEFT JOIN branch ON r_loan_branch.\"branchId\" = branch.\"id\" "
	// queryTotalData += "LEFT JOIN r_loan_disbursement ON r_loan_disbursement.\"loanId\" = loan.\"id\" "
	// queryTotalData += "LEFT JOIN disbursement ON disbursement.\"id\" = r_loan_disbursement.\"disbursementId\" "
	// queryTotalData += "WHERE branch.id = ? AND loan.\"deletedAt\" IS NULL AND loan.\"stage\" NOT IN ('END', 'END-EARLY') "

	queryTotalData := "SELECT COUNT(loan.*) AS \"totalRows\" FROM loan "
	queryTotalData += "JOIN r_loan_group ON r_loan_group.\"loanId\" = loan.id "
	queryTotalData += "JOIN r_loan_branch ON r_loan_branch.\"loanId\" = loan.id "
	queryTotalData += "JOIN r_loan_borrower ON r_loan_borrower.\"loanId\" = loan.id "
	queryTotalData += "JOIN r_loan_disbursement ON r_loan_disbursement.\"loanId\" = loan.id "
	queryTotalData += "JOIN r_cif_borrower ON r_cif_borrower.\"borrowerId\" = r_loan_borrower.\"borrowerId\" "
	queryTotalData += "JOIN cif ON cif.id = r_cif_borrower.\"cifId\" "
	queryTotalData += "JOIN \"group\" ON \"group\".id = r_loan_group.\"groupId\" "
	queryTotalData += "JOIN disbursement ON disbursement.id = r_loan_disbursement.\"disbursementId\" "
	queryTotalData += "WHERE r_loan_branch.\"branchId\" = ? "

	if ctx.URLParam("search") != "" {
		queryTotalData += "AND cif.\"name\" ~* '" + ctx.URLParam("search") + "' "
	}

	services.DBCPsql.Raw(queryTotalData, branchID).Find(&totalData)

	loans := []LoanDatatable{}

	var limitPagination int64 = 10
	var offset int64 = 0

	// query := "SELECT DISTINCT loan.*, "
	// query += "sector.\"name\" as \"sector\", "
	// query += "cif.\"name\" as \"borrower\", "
	// query += "\"group\".\"name\" as \"group\", "
	// query += "branch.\"name\" as \"branch\",  "
	// query += "disbursement.\"disbursementDate\" "
	// query += "FROM loan "
	// query += "LEFT JOIN r_loan_sector ON r_loan_sector.\"loanId\" = loan.\"id\" "
	// query += "LEFT JOIN sector ON r_loan_sector.\"sectorId\" = sector.\"id\" "
	// query += "LEFT JOIN r_loan_borrower ON r_loan_borrower.\"borrowerId\" = loan.\"id\" "
	// query += "LEFT JOIN borrower ON r_loan_borrower.\"borrowerId\" = borrower.\"id\" "
	// query += "LEFT JOIN r_cif_borrower ON r_cif_borrower.\"borrowerId\" = borrower.\"id\" "
	// query += "LEFT JOIN cif ON r_cif_borrower.\"cifId\" = cif.\"id\" "
	// query += "LEFT JOIN r_loan_group ON r_loan_group.\"loanId\" = loan.\"id\" "
	// query += "LEFT JOIN \"group\" ON r_loan_group.\"groupId\" = \"group\".\"id\" "
	// query += "LEFT JOIN r_loan_branch ON r_loan_branch.\"loanId\" = loan.\"id\" "
	// query += "LEFT JOIN branch ON r_loan_branch.\"branchId\" = branch.\"id\" "
	// query += "LEFT JOIN r_loan_disbursement ON r_loan_disbursement.\"loanId\" = loan.\"id\" "
	// query += "LEFT JOIN disbursement ON disbursement.\"id\" = r_loan_disbursement.\"disbursementId\" "
	// query += "WHERE branch.id = ? AND loan.\"deletedAt\" IS NULL AND loan.\"stage\" NOT IN ('END', 'END-EARLY') "

	query := "SELECT loan.id as \"loanId\", cif.name AS \"borrower\", borrower.\"borrowerNo\", \"group\".\"name\" AS \"group\", loan.\"submittedLoanDate\", disbursement.\"disbursementDate\", loan.plafond, loan.tenor, loan.rate, loan.stage, loan.id, r_investor_product_pricing_loan.\"investorId\" as \"investorId\" FROM loan "
	query += "JOIN r_loan_group ON r_loan_group.\"loanId\" = loan.id "
	query += "JOIN r_loan_branch ON r_loan_branch.\"loanId\" = loan.id "
	query += "JOIN r_loan_borrower ON r_loan_borrower.\"loanId\" = loan.id "
	query += "JOIN r_investor_product_pricing_loan ON r_investor_product_pricing_loan.\"loanId\" = loan.id "
	query += "JOIN borrower ON r_loan_borrower.\"borrowerId\" = borrower.id "
	query += "JOIN r_loan_disbursement ON r_loan_disbursement.\"loanId\" = loan.id "
	query += "JOIN r_cif_borrower ON r_cif_borrower.\"borrowerId\" = r_loan_borrower.\"borrowerId\" "
	query += "JOIN cif ON cif.id = r_cif_borrower.\"cifId\" "
	query += "JOIN \"group\" ON \"group\".id = r_loan_group.\"groupId\" "
	query += "JOIN disbursement ON disbursement.id = r_loan_disbursement.\"disbursementId\" "
	query += "WHERE r_loan_branch.\"branchId\" = ? "

	if ctx.URLParam("search") != "" {
		query += "AND (cif.\"name\" ~* '" + ctx.URLParam("search") + "' "
		query += "OR \"group\".\"name\" ~* '" + ctx.URLParam("search") + "' "
		query += "OR loan.stage ~* '" + ctx.URLParam("search") + "' "

		if _, err := strconv.Atoi(ctx.URLParam("search")); err == nil {
			query += "OR cast(loan.id as text) like '" + ctx.URLParam("search") + "' "
			query += "OR cast(borrower.\"borrowerNo\" as text) ~* '" + ctx.URLParam("search") + "')"
		} else {
			query += ")"
		}
	}

	if ctx.URLParam("limit") != "" {
		query += "LIMIT " + ctx.URLParam("limit") + " "
	} else {
		query += "LIMIT " + strconv.FormatInt(limitPagination, 10) + " "
	}

	if ctx.URLParam("page") != "" {
		offset, _ = strconv.ParseInt(ctx.URLParam("page"), 10, 64)
		query += "OFFSET " + strconv.FormatInt(offset, 10)
	} else {
		query += "OFFSET 0"
	}

	services.DBCPsql.Raw(query, branchID).Find(&loans)
	ctx.JSON(iris.StatusOK, iris.Map{
		"status":    "success",
		"totalRows": totalData.TotalRows,
		"data":      loans,
	})
}

// FetchDropping - Fetch loans which in ARCHIVE or DISBURSEMENT-FAILED stage
func FetchDropping(ctx *iris.Context) {
	loanData := []LoanDropping{}

	// ref: dropping.sql
	query := "SELECT loan.id, stage, borrower.\"borrowerNo\", cif_borrower.\"name\" AS borrower, \"group\".\"name\" AS \"group\", investor.id as \"investorId\", cif_investor.name AS investor "
	query += "FROM loan "
	query += "JOIN r_loan_borrower ON r_loan_borrower.\"loanId\" = loan.id "
	query += "JOIN borrower ON borrower.id = r_loan_borrower.\"borrowerId\" "
	query += "JOIN r_cif_borrower ON r_cif_borrower.\"borrowerId\" = borrower.id "
	query += "JOIN (SELECT * FROM cif WHERE \"deletedAt\" IS NULL) AS cif_borrower ON cif_borrower.id = r_cif_borrower.\"cifId\" "
	query += "JOIN r_loan_group ON r_loan_group.\"loanId\" = loan.id "
	query += "JOIN \"group\" ON \"group\".id = r_loan_group.\"groupId\" "
	query += "LEFT JOIN r_investor_product_pricing_loan ON r_investor_product_pricing_loan.\"loanId\"= loan.id "
	query += "LEFT JOIN investor ON investor.id = r_investor_product_pricing_loan.\"investorId\" "
	query += "LEFT JOIN r_cif_investor ON r_cif_investor.\"investorId\" = investor.id "
	query += "LEFT JOIN (SELECT * FROM cif WHERE \"deletedAt\" IS NULL) AS cif_investor ON cif_investor.id = r_cif_investor.\"cifId\" "
	query += "WHERE loan.\"deletedAt\" IS NULL AND (loan.stage = 'ARCHIVE' OR loan.stage = 'DISBURSEMENT-FAILED')"

	services.DBCPsql.Raw(query).Find(&loanData)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   loanData,
	})
}

// UpdateStage - Update Stage Loan
func UpdateStage(ctx *iris.Context) {
	// Habib: logic dipindah ke executeUpdateStage supaya bisa di-reuse
	// loanData := Loan{}

	// loanID := ctx.Param("id")
	// stage := ctx.Param("stage")
	// services.DBCPsql.First(&loanData, loanID)
	// if loanData == (Loan{}) {
	// 	ctx.JSON(iris.StatusInternalServerError, iris.Map{
	// 		"status":  "error",
	// 		"message": "Can't find any loan detail.",
	// 	})
	// 	return
	// }

	// loanHistoryData := loanHistory.LoanHistory{StageFrom: loanData.Stage, StageTo: stage, Remark: "loanId=" + fmt.Sprintf("%v", loanData.ID) + " updated stage to " + stage}
	// services.DBCPsql.Table("loan_history").Create(&loanHistoryData)

	// rLoanHistory := r.RLoanHistory{LoanId: loanData.ID, LoanHistoryId: loanHistoryData.ID}
	// services.DBCPsql.Table("r_loan_history").Create(&rLoanHistory)

	// services.DBCPsql.Table("loan").Where("\"id\" = ?", loanData.ID).UpdateColumn("stage", stage)

	// ctx.JSON(iris.StatusOK, iris.Map{
	// 	"status":    "success",
	// 	"stageFrom": loanData.Stage,
	// 	"stageTo":   stage,
	// })

	loanID, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	stage := ctx.Param("stage")

	loanStage, err := executeUpdateStage(loanID, stage)
	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": "Can't find any loan detail.",
		})
		return
	}

	ctx.JSON(iris.StatusOK, iris.Map{
		"status":    "success",
		"stageFrom": loanStage,
		"stageTo":   stage,
	})
	return
}

func executeUpdateStage(id uint64, stage string) (string, error) {
	loanData := Loan{}
	services.DBCPsql.Where("id = ?", id).First(&loanData)
	fmt.Printf("%v", loanData)
	if loanData == (Loan{}) {
		return "", errors.New("loan not found")
	}

	loanHistoryData := loanHistory.LoanHistory{StageFrom: loanData.Stage, StageTo: stage, Remark: "loanId=" + fmt.Sprintf("%v", loanData.ID) + " updated stage to " + stage}
	services.DBCPsql.Table("loan_history").Create(&loanHistoryData)

	rLoanHistory := r.RLoanHistory{LoanId: loanData.ID, LoanHistoryId: loanHistoryData.ID}
	services.DBCPsql.Table("r_loan_history").Create(&rLoanHistory)

	services.DBCPsql.Table("loan").Where("\"id\" = ?", loanData.ID).UpdateColumn("stage", stage)

	if loanData.Stage == "DISBURSEMENT-FAILED" && (stage == "MARKETPLACE" || stage == "PRIVATE") {
		services.DBCPsql.Table("r_investor_product_pricing_loan").Where(" \"loanId\" = ?", loanData.ID).Update("investorId", nil)
		services.DBCPsql.Table("r_investor_product_pricing_loan").Where(" \"loanId\" = ?", loanData.ID).Update("updatedAt", time.Now())
	}

	if loanData.Stage == "ARCHIVE" {
		services.DBCPsql.Table("r_investor_product_pricing_loan").Where(" \"loanId\" = ?", loanData.ID).Update("investorId", nil)
		services.DBCPsql.Table("r_investor_product_pricing_loan").Where(" \"loanId\" = ?", loanData.ID).Update("updatedAt", time.Now())
		// services.DBCPsql.Table("r_investor_product_pricing_loan").Where(" \"loanId\" = ?", loanData.ID).Update("deletedAt", time.Now())
	}
	
	// reset product pricing to retail 
	// if stageTo is `PRIVATE` or `MARKETPLACE`
	if stage == `PRIVATE` || stage == `MARKETPLACE` {
		productPricingSchema := productPricing.ProductPricing{}
		services.DBCPsql.Table(`product_pricing`).Where(`"isInstitutional" = false AND current_timestamp BETWEEN "startDate" AND "endDate"`).Scan(&productPricingSchema)
		services.DBCPsql.Table("r_investor_product_pricing_loan").Where(" \"loanId\" = ?", loanData.ID).Update("productPricingId", productPricingSchema.ID)
	}

	return loanData.Stage, nil
}

func GetLoanDetail(ctx *iris.Context) {
	loanObj := Loan{}
	borrowerObj := LoanBorrowerProfile{}
	installmentObj := []installment.Installment{}

	loanId, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		fmt.Println(err)
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": "Something went wrong. Please try again later.",
		})
		return
	}

	services.DBCPsql.Table("loan").Where("\"id\" = ?", loanId).First(&loanObj)

	queryBorrowerObj := "SELECT borrower.\"id\" AS \"borrowerId\",cif.\"cifNumber\", cif.\"name\", \"group\".\"name\" AS \"group\", area.\"name\" AS \"area\", branch.\"name\" AS \"branch\" "
	queryBorrowerObj += "FROM loan "
	queryBorrowerObj += "LEFT JOIN r_loan_borrower ON r_loan_borrower.\"loanId\" = loan.\"id\" "
	queryBorrowerObj += "LEFT JOIN borrower ON borrower.\"id\" = r_loan_borrower.\"borrowerId\" "
	queryBorrowerObj += "LEFT JOIN r_cif_borrower ON r_cif_borrower.\"borrowerId\" = borrower.\"id\" "
	queryBorrowerObj += "LEFT JOIN cif ON cif.\"id\" = r_cif_borrower.\"cifId\" "
	queryBorrowerObj += "LEFT JOIN r_loan_group ON r_loan_group.\"loanId\" = loan.\"id\" "
	queryBorrowerObj += "LEFT JOIN \"group\" ON \"group\".\"id\" = r_loan_group.\"groupId\" "
	queryBorrowerObj += "LEFT JOIN r_loan_branch ON r_loan_branch.\"loanId\" = loan.\"id\" "
	queryBorrowerObj += "LEFT JOIN branch ON branch.\"id\" = r_loan_branch.\"branchId\" "
	queryBorrowerObj += "LEFT JOIN r_area_branch ON r_area_branch.\"branchId\" = branch.\"id\" "
	queryBorrowerObj += "LEFT JOIN area ON area.\"id\" = r_area_branch.\"areaId\" "
	queryBorrowerObj += "WHERE loan.\"id\" = ?"

	services.DBCPsql.Raw(queryBorrowerObj, loanId).Scan(&borrowerObj)

	queryInstallmentObj := "SELECT * "
	queryInstallmentObj += "FROM installment "
	queryInstallmentObj += "JOIN r_loan_installment ON r_loan_installment.\"installmentId\" = installment.\"id\" "
	queryInstallmentObj += "JOIN loan ON loan.\"id\" = r_loan_installment.\"loanId\" "
	queryInstallmentObj += "WHERE loan.\"id\" = ? and installment.\"deletedAt\" is null"

	services.DBCPsql.Raw(queryInstallmentObj, loanId).Scan(&installmentObj)

	investorCifObj := cif.Cif{}

	queryInvestorObj := "SELECT * FROM cif "
	queryInvestorObj += "JOIN r_cif_investor ON r_cif_investor.\"cifId\" = cif.id "
	queryInvestorObj += "JOIN r_investor_product_pricing_loan ON r_investor_product_pricing_loan.\"investorId\" = r_cif_investor.\"investorId\" "
	queryInvestorObj += "WHERE r_investor_product_pricing_loan.\"loanId\" = ? "

	services.DBCPsql.Raw(queryInvestorObj, loanId).Scan(&investorCifObj)

	productPricingObj  := productPricing.ProductPricing{}
	queryProductPricing := "SELECT product_pricing.\"id\", product_pricing.\"returnOfInvestment\", product_pricing.\"administrationFee\", "
	queryProductPricing += "product_pricing.\"serviceFee\", product_pricing.\"startDate\", product_pricing.\"endDate\", product_pricing.\"isInstitutional\" "
	queryProductPricing += "FROM product_pricing "
	queryProductPricing += "LEFT JOIN r_investor_product_pricing_loan ON r_investor_product_pricing_loan.\"productPricingId\"=product_pricing.\"id\" WHERE r_investor_product_pricing_loan.\"loanId\" = ?"
	services.DBCPsql.Raw(queryProductPricing,loanId).Scan(&productPricingObj)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data": iris.Map{
			"loan":        loanObj,
			"borrower":    borrowerObj,
			"installment": installmentObj,
			"investor":    investorCifObj,
			"productPricing":    productPricingObj,
		},
	})
}

type BorrowerObj struct {
	Fullname   string `gorm:"column:name" json:"name"`
	BorrowerNo string `gorm:"column:borrowerNo" json:"borrowerNo"`
	Branch 			string `gorm:"column:branch" json:"branch"`
	IdCardNo		string `gorm:"column:idCardNo" json:"idCardNo"`
	Address		 	string `gorm:"column:address" json:"address"`
	Kelurahan		string `gorm:"column:kelurahan" json:"kelurahan"`
	Kecamatan		string `gorm:"column:kecamatan " json:"kecamatan"`
	Group      	string `gorm:"column:group" json:"group"`
	TempatLahir string `gorm:"column:tempatLahir" json:"tempatLahir"`
	TanggalLahir string `gorm:"column:tanggalLahir" json:"tanggalLahir"`
	Desa		string `gorm:"column::desa" json:"desa"`
	Status		string `gorm:"column::status" json:"status"`
	NamaPenanggungJawab string `gorm:"column:nama_pj" json:"nama_pj"`
	TempatLahirPenanggungJawab string `gorm:"column:pj_tempatlahir" json:"pj_tempatlahir"`
	TglLahirPenanggungJawab string `gorm:"column:pj_tgllahir" json:"pj_tgllahir"`
	HubPenanggungJawab string `gorm:"column:hubungan" json:"data_hubungan"`
}

// GetAkadData - Get data to be shown in Akad
func GetAkadData(ctx *iris.Context) {
	loanID, _ := strconv.Atoi(ctx.Param("id"))
	data := Akad{}

	query := `SELECT loan.*, disbursement. "disbursementDate", product_pricing."returnOfInvestment", 
	product_pricing."administrationFee", product_pricing."serviceFee", "group"."name" as "group", 
	r_investor_product_pricing_loan."investorId", "orderNo"
	FROM loan 
	JOIN r_investor_product_pricing_loan ON r_investor_product_pricing_loan."loanId" = loan.id 
	JOIN product_pricing ON product_pricing.id = r_investor_product_pricing_loan."productPricingId" 
	JOIN r_loan_disbursement ON r_loan_disbursement."loanId" = loan.id 
	JOIN disbursement ON disbursement.id = r_loan_disbursement."disbursementId" 
	JOIN r_loan_group ON r_loan_group."loanId" = loan.id 
	JOIN "group" ON "group".id = r_loan_group."groupId" 
	JOIN "r_loan_order" ON "loan".id = r_loan_order."loanId" 
	JOIN "loan_order" ON "r_loan_order"."loanOrderId" = loan_order."id" 
	WHERE loan.id = ? AND loan."deletedAt" IS NULL`

	errAkad := services.DBCPsql.Raw(query, loanID).Scan(&data).Error

	if errAkad != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "Error",
			"message": "Can't Find Akad Based On Loan ID given. Error: " + errAkad.Error(),
		})
		return
	}

	queryGetInvestor := `SELECT cif.*
	FROM r_investor_product_pricing_loan
	JOIN r_cif_investor ON r_cif_investor."investorId" = r_investor_product_pricing_loan."investorId"
	JOIN cif ON cif.id = r_cif_investor."cifId"
	WHERE r_investor_product_pricing_loan."loanId" = ? AND r_investor_product_pricing_loan."deletedAt" IS NULL LIMIT 1`

	investorData := cif.Cif{}

	errCif := services.DBCPsql.Raw(queryGetInvestor, loanID).Scan(&investorData).Error

	if errCif != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "Error",
			"message": "Can't Find CIF Based On Loan ID given. Error: " + errAkad.Error(),
		})
		return
	}

	queryGetBorrower := `SELECT cif."name", cif."address" as "address", cif."kelurahan" as "kelurahan", cif."kecamatan" as kecamatan, 
	cif."idCardNo" as "idCardNo" ,borrower."borrowerNo", "group"."name" as "group", branch."name" as "branch", lr._raw::json ->> 'client_birthplace' as "tempatLahir",
	lr._raw::json ->> 'client_birthdate' as "tanggalLahir", lr._raw::json ->> 'client_desa' as "desa", lr._raw::json ->> 'client_maritalstatus' as "status",
	lr._raw::json ->> 'data_suami' as "nama_pj", lr._raw::json ->> 'data_suami_tempatlahir' as "pj_tempatlahir", lr._raw::json ->> 'data_suami_tgllahir' as "pj_tgllahir", 
	lr._raw::json ->> 'data_hubungan' as "hubungan"
	FROM loan
	JOIN r_loan_borrower on r_loan_borrower."loanId" = loan.id
	JOIN borrower ON borrower.id = r_loan_borrower."borrowerId"
	JOIN r_cif_borrower ON r_cif_borrower."borrowerId" = borrower.id
	JOIN cif ON cif.id = r_cif_borrower."cifId"
	JOIN loan_raw lr on lr."loanId" = loan.id
	JOIN r_loan_group on r_loan_group."loanId" = loan.id
	JOIN "group" on "group".id = r_loan_group."groupId"
	JOIN r_group_branch ON r_group_branch."groupId" = "group"."id"
	JOIN branch on branch."id" = r_group_branch."branchId"
	WHERE loan.id = ? AND loan."deletedAt" IS NULL LIMIT 1`

	borrowerData := BorrowerObj{}

	errBorrower := services.DBCPsql.Raw(queryGetBorrower, loanID).Scan(&borrowerData).Error

	if errBorrower != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "Error",
			"message": "Can't Find Borrower Based On Loan ID given. Error: " + errAkad.Error(),
		})
		return
	}
	
	floatTenor := float64(data.Tenor)
	weeklyBase := Round(data.Plafond/floatTenor, 2)
	weeklyMargin := Round(data.Rate*data.Plafond*data.ReturnOfInvestment/floatTenor, 2)
	weeklyFeeBorrower := Round(data.Rate*data.Plafond*data.AdminitrationFee/floatTenor, 2)
	weeklyFeeInvestor := Round(data.Rate*data.Plafond*data.ServiceFee/floatTenor, 2)

	noReserveTime, _ := time.Parse("2006-01-02", "2017-04-03")
	augustTime, _ := time.Parse("2006-01-02", "2016-08-29")
	submittedLoanTime, _ := time.Parse("2006-01-02T15:04:05-07:00", data.SubmittedLoanDate)

	var reserve uint64
	if submittedLoanTime.After(noReserveTime) {
		reserve = 0
	} else {
		if submittedLoanTime.After(augustTime) {
			switch {
			case data.Plafond <= 3000100:
				reserve = 3000
			case data.Plafond <= 5000100:
				reserve = 4000
			case data.Plafond <= 7000100:
				reserve = 5000
			case data.Plafond <= 9000100:
				reserve = 6000
			case data.Plafond <= 11000100:
				reserve = 7000
			default:
				reserve = 8000
			}
		} else {
			switch {
			case data.Plafond < 1500001:
				reserve = 2000
			case data.Plafond < 2500001:
				reserve = 3000
			case data.Plafond < 3500001:
				reserve = 4000
			case data.Plafond < 4500001:
				reserve = 5000
			case data.Plafond < 5000001:
				reserve = 6000
			default:
				reserve = 7000
			}
		}
	}

	var sentAgreementType string
	if data.AgreementType == "" {
		sentAgreementType = "MBA"
	} else {
		sentAgreementType = data.AgreementType
	}

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data": iris.Map{
			"_id":               data.ID,
			"adminFee":			 data.AdminitrationFee,
			"serviceFee":		 data.ServiceFee,
			"disbursementDate":  data.DisbursementDate,
			"agreementType":     sentAgreementType,
			"purpose":           data.Purpose,
			"plafond":           data.Plafond,
			"tenor":             data.Tenor,
			"installment":       data.Installment,
			"rate": 			 data.Rate,
			"returnOfInvestment": data.ReturnOfInvestment,
			"orderNo":			 data.OrderNo,
			"weeklyBase":        weeklyBase,
			"weeklyMargin":      weeklyMargin,
			"weeklyFeeBorrower": weeklyFeeBorrower,
			"weeklyFeeInvestor": weeklyFeeInvestor,
			"reserve":           reserve,
			"borrower":          borrowerData,
			"investorId":        data.InvestorID,
			"investor":          investorData,
		},
	})
}

// RefundAndChangeStageTo - refund investor balance and change loan stage
func RefundAndChangeStageTo(ctx *iris.Context) {
	loanID, _ := strconv.ParseUint(ctx.Param("loan_id"), 10, 64)
	stage := ctx.Param("stage")

	// get loan_id, investor_id, account_id, plafond
	refundBase := RefundBase{}
	// ref: refund-base.sql
	queryRefundBase := `SELECT loan.id AS loan_id, loan."isInsurance", investor.id AS investor_id, account.id AS account_id, loan.plafond `
	queryRefundBase += "FROM loan "
	queryRefundBase += "JOIN r_investor_product_pricing_loan ON r_investor_product_pricing_loan.\"loanId\" = loan.id "
	queryRefundBase += "JOIN investor ON investor.id = r_investor_product_pricing_loan.\"investorId\" "
	queryRefundBase += "JOIN r_account_investor ON r_account_investor.\"investorId\" = investor.id "
	queryRefundBase += "JOIN account ON account.id = r_account_investor.\"accountId\" "
	queryRefundBase += "WHERE loan.\"deletedAt\" IS NULL AND loan.id = ? "
	services.DBCPsql.Raw(queryRefundBase, loanID).First(&refundBase)

	// add new account_transaction_debit entry
	transaction := accountTransactionDebit.AccountTransactionDebit{
		Type:            "REFUND",
		TransactionDate: time.Now(),
		Amount:          refundBase.Plafond,
		Remark:          "",
	}
	services.DBCPsql.Table("account_transaction_debit").Create(&transaction)

	// add new account_transaction_debit_loan entry
	transactionLoan := accountTransactionDebit.AccountTransactionDebitLoan{
		AccountTransactionDebitID: transaction.ID,
		LoanID: refundBase.LoanID,
	}
	services.DBCPsql.Table("r_account_transaction_debit_loan").Create(&transactionLoan)

	// connect the entry to investor account
	rTransaction := r.RAccountTransactionDebit{
		AccountId:                 refundBase.AccountID,
		AccountTransactionDebitId: transaction.ID,
	}
	services.DBCPsql.Table("r_account_transaction_debit").Create(&rTransaction)
	
	// Check whether loan included insurance or NOT
	// if yes, then refund the amount based on plafond * insurance-rate 
	// insurance-rate = 1.5% = 0.015
	if refundBase.IsInsurance {
		insuranceRate := 0.015
		
		transactionInsurance := accountTransactionDebit.AccountTransactionDebit{
			Type:            "REFUND-INSURANCE",
			TransactionDate: time.Now(),
			Amount:          refundBase.Plafond * insuranceRate,
			Remark:          "Refund insurance",
		}
		services.DBCPsql.Table("account_transaction_debit").Create(&transactionInsurance)
	
		// add new account_transaction_debit_loan entry
		transactionLoanInsurance := accountTransactionDebit.AccountTransactionDebitLoan{
			AccountTransactionDebitID: transactionInsurance.ID,
			LoanID: refundBase.LoanID,
		}
		services.DBCPsql.Table("r_account_transaction_debit_loan").Create(&transactionLoanInsurance)
	
		// connect the entry to investor account
		rTransactionInsurance := r.RAccountTransactionDebit{
			AccountId:                 refundBase.AccountID,
			AccountTransactionDebitId: transactionInsurance.ID,
		}
		services.DBCPsql.Table("r_account_transaction_debit").Create(&rTransactionInsurance)
	}

	// calculate account balance and save it to account
	queryTotalDebit := "SELECT SUM(account_transaction_debit.amount) "
	queryTotalDebit += "FROM account_transaction_debit "
	queryTotalDebit += "JOIN r_account_transaction_debit ON 	r_account_transaction_debit.\"accountTransactionDebitId\" = 					account_transaction_debit.id "
	queryTotalDebit += "WHERE r_account_transaction_debit.\"accountId\" = ? and account_transaction_debit.\"deletedAt\" is null"

	queryTotalCredit := "SELECT SUM(account_transaction_credit.amount) "
	queryTotalCredit += "FROM account_transaction_credit "
	queryTotalCredit += "JOIN r_account_transaction_credit ON 	r_account_transaction_credit.\"accountTransactionCreditId\" = 					account_transaction_credit.id "
	queryTotalCredit += "WHERE r_account_transaction_credit.\"accountId\" = ? and account_transaction_credit.\"deletedAt\" is null "

	debit := AccountSum{}
	credit := AccountSum{}
	services.DBCPsql.Raw(queryTotalDebit, refundBase.AccountID).Scan(&debit)
	services.DBCPsql.Raw(queryTotalCredit, refundBase.AccountID).Scan(&credit)
	totalBalance := debit.Sum - credit.Sum

	services.DBCPsql.Table("account").Where("id = ?", refundBase.AccountID).Updates(account.Account{TotalDebit: debit.Sum, TotalCredit: credit.Sum, TotalBalance: totalBalance})

	loanStage, err := executeUpdateStage(loanID, stage)
	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": "Can't find any loan detail.",
		})
		return
	}

	ctx.JSON(iris.StatusOK, iris.Map{
		"status":    "success",
		"stageFrom": loanStage,
		"stageTo":   stage,
	})
}

// Round - round the nearest value to nearest integer
func Round(val float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= 0.5 {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

// GetLoanStageHistory - get loan stage history by loan id
func GetLoanStageHistory(ctx *iris.Context) {
	loanID, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)

	query := "SELECT r_loan_history.\"createdAt\", loan_history.id, loan_history.\"stageFrom\", loan_history.\"stageTo\", loan_history.remark "
	query += "FROM loan_history "
	query += "JOIN r_loan_history ON r_loan_history.\"loanHistoryId\" = loan_history.id "
	query += "WHERE loan_history.\"deletedAt\" is null and r_loan_history.\"loanId\" = ? "

	historyData := []LoanStageHistory{}
	services.DBCPsql.Raw(query, loanID).Scan(&historyData)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   historyData,
	})
}

// AssignInvestorToLoan - assign investor to loan
func AssignInvestorToLoan(ctx *iris.Context) {
	cifID := ctx.URLParam("cifId")
	loanID := ctx.URLParam("loanId")
	investorId := ctx.URLParam("investorId")

	rCifInvestorSchema := r.RCifInvestor{}
	services.DBCPsql.Table("r_cif_investor").Where("\"cifId\" = ?", cifID).Scan(&rCifInvestorSchema)

	loanSchema := Loan{}
	services.DBCPsql.Table("loan").Where("\"id\" = ?", loanID).Scan(&loanSchema)

	loanHistorySchema := &loanHistory.LoanHistory{StageFrom: loanSchema.Stage, StageTo: "INVESTOR", Remark: "MANUAL ASSIGN loanId=" + fmt.Sprintf("%v", loanSchema.ID) + " investorId=" + fmt.Sprintf("%v", rCifInvestorSchema.InvestorId)}

	transac := services.DBCPsql.Begin()
	err:=transac.Table("loan_history").Create(loanHistorySchema).Error
	if err!=nil {
		processErrorAndRollback(ctx,transac,err,"CREATE-LOAN-HISTORY")
		return
	}

	fmt.Println("Insert r_loan_history")
	rLoanHistorySchema := &r.RLoanHistory{LoanId: loanSchema.ID, LoanHistoryId: loanHistorySchema.ID}
	err=transac.Table("r_loan_history").Create(rLoanHistorySchema).Error
	fmt.Println("Insert r_loan_history_finish")
	if err!=nil {
		processErrorAndRollback(ctx,transac,err,"CREATE-R-LOAN-HISTORY")
		return
	}

	fmt.Println("UpdateLoan",loanID)
	err=transac.Exec(`update loan set stage='INVESTOR' where loan.id=?`,loanID).Error
	fmt.Println("Finish update Loan")
	if err!=nil {
		processErrorAndRollback(ctx,transac,err,"UPDATE-STAGE-INVESTOR")
		return
	}
	// Checking Function if investorId has a product pricing.
	rInvestorProductPricing := r.RInvestorProductPricing{}
	transac.Table("r_investor_product_pricing").Where("\"investorId\" = ?", investorId).Scan(&rInvestorProductPricing)
	fmt.Println("ProductPricing")
	if rInvestorProductPricing.ID == 0{
		err=transac.Table("r_investor_product_pricing_loan").Where("\"loanId\" = ?", loanSchema.ID).UpdateColumn("investorId", rCifInvestorSchema.InvestorId).Error
		if err!=nil {
			processErrorAndRollback(ctx,transac,err,"UPDATE-R-CIF-INVESTOR")
			return
		}
	}else{
		err=transac.Table("r_investor_product_pricing_loan").Where("\"loanId\" = ?", loanSchema.ID).UpdateColumn("investorId", rCifInvestorSchema.InvestorId).Error
		if err!=nil {
			processErrorAndRollback(ctx,transac,err,"UPDATE-R-INVESTOR-PRODUCT-PRICING-LOAN")
			return
		}
		err=transac.Table("r_investor_product_pricing_loan").Where("\"loanId\" = ?", loanSchema.ID).UpdateColumn("productPricingId", rInvestorProductPricing.ProductPricingId).Error
		if err!=nil {
			processErrorAndRollback(ctx,transac,err,"UPDATE-R-INVESTOR-PRODUCT-PRICING-LOAN")
			return
		}
	}
	fmt.Println("ProductPricingFinish")

	rAccountInvestorSchema := r.RAccountInvestor{}
	transac.Table("r_account_investor").Where("\"investorId\" = ?", rCifInvestorSchema.InvestorId).Scan(&rAccountInvestorSchema)

	accountTransactionCreditSchema := &accountTransactionCredit.AccountTransactionCredit{Type: "INVEST", TransactionDate: time.Now(), Amount: loanSchema.Plafond, Remark: "MANUAL ASSIGN"}
	err=transac.Table("account_transaction_credit").Create(accountTransactionCreditSchema).Error
	if err!=nil {
		processErrorAndRollback(ctx,transac,err,"UPDATE-ACCOUNT-TRANSACTION-CREDIT")
		return
	}

	rAccountTransactionCreditSchema := &r.RAccountTransactionCredit{AccountId: rAccountInvestorSchema.AccountId, AccountTransactionCreditId: accountTransactionCreditSchema.ID}
	err=transac.Table("r_account_transaction_credit").Create(rAccountTransactionCreditSchema).Error
	if err!=nil {
		processErrorAndRollback(ctx,transac,err,"UPDATE-R-ACCOUNT-TRANSACTION-CREDIT")
		return
	}
	rAccounTransactionCreditLoanSchema := &r.RAccountTransactionCreditLoan{AccountTransactionCreditId: accountTransactionCreditSchema.ID, LoanId: loanSchema.ID}
	err=transac.Table("r_account_transaction_credit_loan").Create(rAccounTransactionCreditLoanSchema).Error
	if err!=nil {
		processErrorAndRollback(ctx,transac,err,"r_account_transaction_credit_loan")
		return
	}

	totalDebit := accountTransactionDebit.GetTotalAccountTransactionDebitByTransac(transac,rAccountInvestorSchema.AccountId)
	totalCredit := accountTransactionCredit.GetTotalAccountTransactionCreditByTransac(transac,rAccountInvestorSchema.AccountId)

	totalBalance := totalDebit - totalCredit

	err=transac.Exec(`UPDATE "account" SET "totalDebit" = ?, "totalCredit" = ?, "totalBalance" = ?  WHERE (id = ?)`,totalDebit,totalCredit,totalBalance,rAccountInvestorSchema.AccountId).Error

	if err!=nil {
		processErrorAndRollback(ctx,transac,err,"ACCOUNT")
		return
	}

	queryReferal:=`with A as (
select id, "inviterInvestorId" from referral where "inviteeInvestorId" = ` + strconv.FormatUint(rCifInvestorSchema.InvestorId,10) + ` and "inviterGetTimestamp" isnull and "deletedAt" isnull
),
B as (
    update referral set "inviterGetTimestamp" = current_timestamp, "inviteeUseTimestamp" = current_timestamp
    from A where A.id = referral.id
    returning referral.id,concat('SUCCESSFUL REFERRAL INVITEE = ',"inviteeInvestorId",' INVITER = ',referral."inviterInvestorId",' REFERALL ID = ',referral.id) "remark",
    referral."inviterInvestorId","inviteeInvestorId","inviterGetTimestamp","inviteeGetTimestamp","createdAt","updatedAt","deletedAt"
),
C as (
    insert into account_transaction_debit ("type", amount, remark, "transactionDate","createdAt", "updatedAt")
    select 'REFERRAL-INVITER',100000,"remark",current_timestamp,current_timestamp,current_timestamp from B
    returning id, remark
),
D as (
    insert into r_account_transaction_debit_referral ("accountTransactionDebitId","referralId","createdAt")
    select C.id, split_part(remark,' ',12)::int, current_timestamp from C),
I as (
    insert into r_account_transaction_debit ("accountTransactionDebitId","accountId","createdAt")
    select C.id, rai."accountId", current_timestamp from C
    join r_account_investor rai on split_part(C.remark,' ',8)::int = rai."investorId"
    ),
J as(
    update account set threshold = coalesce(threshold,0) + total * 100000, "totalDebit" = coalesce("totalDebit",0) + total * 100000, "totalBalance" = coalesce("totalBalance",0) + total * 100000
    from (select rai."accountId" id, count(1) "total" from r_account_investor rai join B on rai."investorId" = B."inviterInvestorId" group by "accountId" ) foo where account.id = foo.id
    returning account.id
)
select * from B`
	err=transac.Exec(queryReferal).Error
	if err!=nil {
		processErrorAndRollback(ctx,transac,err,"REFERAL")
		return
	}
	transac.Commit()

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   iris.Map{},
	})
}

func processErrorAndRollback(ctx *iris.Context, db *gorm.DB, err error, process string) {
	fmt.Println("#ERROR#MANUALASSIGN",process,err)
	db.Rollback()
	ctx.JSON(iris.StatusInternalServerError, iris.Map{"status": "error", "error": "Error on " + process + " " + err.Error()})
}