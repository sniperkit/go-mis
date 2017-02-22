package borrower

import (
	"bitbucket.org/go-mis/services"
	//"strconv"
	//"fmt"
	iris "gopkg.in/kataras/iris.v4"
	//"encoding/json"
	"bitbucket.org/go-mis/modules/cif"
	"bitbucket.org/go-mis/modules/loan"
	"bitbucket.org/go-mis/modules/r"
	productPricing "bitbucket.org/go-mis/modules/product-pricing"

)

func Init() {
	services.DBCPsql.AutoMigrate(&Borrower{})
	services.BaseCrudInit(Borrower{}, []Borrower{})
}

// Approve prospective borrower 
func Approve(ctx *iris.Context) {

	// map the payload 
	payload := make(map[string]interface{}, 0)
	if err := ctx.ReadJSON(&payload); err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ktp := payload["client_ktp"].(string)

	// get CIF with with idCardNo = ktp
	cifData := cif.Cif{}
	if ktp != "" {
		services.DBCPsql.Table("cif").Select("\"idCardNo\"").Where("\"idCardNo\" = ?", ktp).Scan(&cifData)
		if cifData.IdCardNo != "" {
			// found. use existing cif
			// get borrower id
			borrower := Borrower{}
			services.DBCPsql.Table("r_cif_borrower").Where("\"cifId\" =?", cifData.ID).Scan(&borrower)

			// get loan id
			loan := loan.Loan{}
			services.DBCPsql.Table("r_loan_borrower").Where("\"borrowerId\" =?", borrower.ID).Scan(&loan)

			// which loan pricing would we like to use?
			// get the newest one
			UseProductPricing(0, loan.ID)
		} else {
			// not found. create new CIF
			cifData = CreateCIF(payload)
			services.DBCPsql.Table("cif").Create(&cifData)

			// reserve one row for this new borrower 
			newBorrower := Borrower{}
			services.DBCPsql.Table("borrower").Create(&newBorrower)

			rCifBorrower := r.RCifBorrower{
				CifId:cifData.ID,
				BorrowerId:newBorrower.ID,
			}
			services.DBCPsql.Table("r_cif_borrower").Create(&rCifBorrower)

			// reserve one loan record for this new borrower
			loan := loan.Loan{}
			services.DBCPsql.Table("loan").Create(&loan)

			rLoanBorrower := r.RLoanBorrower{
				LoanId:loan.ID,
				BorrowerId:newBorrower.ID,
			}
			services.DBCPsql.Table("r_loan_borrower").Create(&rLoanBorrower)

			// which loan pricing would we like to use?
			// get the newest one
			UseProductPricing(0, loan.ID)
		}

	}
	return
}


func UseProductPricing(investorId uint64, loanId uint64) {
	pPricing := productPricing.ProductPricing{}
	services.DBCPsql.Table("product_pricing").Last(&pPricing)

	rInvProdPriceLoan := r.RInvestorProductPricingLoan{
		InvestorId:investorId,
		ProductPricingId:pPricing.ID,
		LoanId:loanId,
	}
	services.DBCPsql.Table("r_investor_product_pricing_loan").Create(&rInvProdPriceLoan)
}

func CreateCIF(payload map[string]interface{}) cif.Cif {
	// convert each payload  field into empty string 
	var cpl map[string]string
	cpl = make(map[string]string)
	for k, v := range payload {
		cpl[k] = v.(string)
	}

	// map each payload field to it's respective cif field
	newCif := cif.Cif{}

	// here the only payload that match the cif fields
	newCif.Username						= cpl["client_simplename"]
	newCif.Name								= cpl["client_fullname"]
	newCif.PlaceOfBirth       = cpl["client_birthplace"]
	newCif.DateOfBirth        = cpl["client_birthdate"]
	newCif.IdCardNo           = cpl["client_ktp"]
	newCif.IdCardFilename     = cpl["photo_ktp"]
	newCif.TaxCardNo          = cpl["client_npwp"]
	newCif.MaritalStatus      = cpl["client_marital_status"]
	newCif.MotherName         = cpl["client_ibu_kandung"]
	newCif.Religion           = cpl["client_religion"]
	newCif.Address            = cpl["client_alamat"]
	newCif.Kelurahan          = cpl["client_desa"]
	newCif.Kecamatan          = cpl["kecamatan"]
	return newCif
}
