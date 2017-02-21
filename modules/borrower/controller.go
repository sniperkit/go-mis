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
	// check each payload value which is interface is not nil
	// convert each value into string (we only deal with string this time) 
	var cpl map[string]string
	cpl = make(map[string]string)
	for k, _ := range payload {
		if payload[k] != nil {
			cpl[k] = payload[k].(string)
		}
	}

	// assign the data to corresppndence cif
	newCif := cif.Cif{}

	newCif.Username						= cpl["client_simplename"]
	newCif.Name								= cpl["client_fullname"]
	newCif.PlaceOfBirth       = cpl["client_birthplace"]
	newCif.DateOfBirth        = cpl["client_birthdate"]
	newCif.IdCardNo           = cpl["client_ktp"]
	//newCif.IdCardValidDate    = payload[""].(string)
	newCif.IdCardFilename     = cpl["photo_ktp"]
	//newCif.TaxCardNo          = payload[""].(string)
	//newCif.TaxCardFilename    = payload[""].(string)
	newCif.MaritalStatus      = cpl["client_marital_status"]
	newCif.MotherName         = cpl["client_ibu_kandung"]
	newCif.Religion           = cpl["client_religion"]
	newCif.Address            = cpl["client_alamat"]
	//newCif.Kelurahan          = payload[""].(string)
	newCif.Kecamatan          = cpl["kecamatan"]
	//newCif.City               = payload[""].(string)
	//newCif.Province           = payload[""].(string)
	//newCif.Nationality        = payload[""].(string)
	//newCif.Zipcode            = payload[""].(string)
	//newCif.PhoneNo            = payload[""].(string)
	//newCif.CompanyName        = payload[""].(string)
	//newCif.CompanyAddress     = payload[""].(string)
	//newCif.Occupation         = payload[""].(string)
	//newCif.Income             = payload[""].(string)
	//newCif.IncomeSourceFund   = payload[""].(string)
	//newCif.IncomeSourceCountry= payload[""].(string)
	//newCif.IsActivated        = payload[""].(string)
	//newCif.IsVAlidated        = payload[""].(string)
	//newCif.CreatedAt          = payload[""].(string)
	//newCif.UpdatedAt          = payload[""].(string)
	//newCif.DeletedAt          = payload[""].(string)
	return newCif
}
