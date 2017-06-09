package cif

import (
	"strconv"
	"bitbucket.org/go-mis/services"
	iris "gopkg.in/kataras/iris.v4"
	"fmt"
)

func Init() {
	services.DBCPsql.AutoMigrate(&Cif{})
	services.BaseCrudInit(Cif{}, []Cif{})
}

type TotalData struct {
	TotalRows int64 `gorm:"column:totalRows" json:"totalRows"`
}

func FetchAll(ctx *iris.Context) {
	totalData := TotalData{}

	queryTotalData := "SELECT count(cif.*) as \"totalRows\" "
	queryTotalData += "FROM cif "
	queryTotalData += "LEFT JOIN r_cif_borrower ON r_cif_borrower.\"cifId\" = cif.\"id\" "
	queryTotalData += "LEFT JOIN r_cif_investor ON r_cif_investor.\"cifId\" = cif.\"id\" "

	if ctx.URLParam("search") != "" {
		// queryTotalData += "WHERE cif.name LIKE '%" + ctx.URLParam("search") + "%' "
		queryTotalData += "WHERE cif.\"name\" ~* '" + ctx.URLParam("search") + "' "
	}

	services.DBCPsql.Raw(queryTotalData).Find(&totalData)

	var limitPagination int64 = 10
	var offset int64 = 0
	cifInvestorBorrower := []CifInvestorBorrower{}

	query := "SELECT cif.\"id\", cif.\"cifNumber\", cif.\"name\", cif.\"isActivated\", cif.\"isValidated\", "
	query += "r_cif_borrower.\"borrowerId\" as \"borrowerId\", r_cif_investor.\"investorId\" as \"investorId\", "
	query += "r_cif_borrower.\"borrowerId\" IS NOT NULL as \"isBorrower\", r_cif_investor.\"investorId\" IS NOT NULL as \"isInvestor\" "
	query += "FROM cif "
	query += "LEFT JOIN r_cif_borrower ON r_cif_borrower.\"cifId\" = cif.\"id\" "
	query += "LEFT JOIN r_cif_investor ON r_cif_investor.\"cifId\" = cif.\"id\" "

	if ctx.URLParam("search") != "" {
		// query += "WHERE cif.name LIKE '%" + ctx.URLParam("search") + "%' "
		query += "WHERE cif.\"name\" ~* '" + ctx.URLParam("search") + "' "
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

	services.DBCPsql.Raw(query).Find(&cifInvestorBorrower)
	ctx.JSON(iris.StatusOK, iris.Map{
		"status":    "success",
		"totalRows": totalData.TotalRows,
		"data":      cifInvestorBorrower,
	})
}

func GetCifSummary(ctx *iris.Context) {
	cifSummary := CifSummary{}

	query := "SELECT (SELECT COUNT(cif.*) AS \"totalRegisteredCif\" FROM cif), (SELECT COUNT(investor.*) AS \"totalInvestor\" FROM investor), (SELECT COUNT(borrower.*) AS \"totalBorrower\" FROM borrower)"

	services.DBCPsql.Raw(query).Find(&cifSummary)
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   cifSummary,
	})
}

type CifInvestorAccount struct {
	CifID        uint64  `gorm:"column:cifId" json:"cifId"`
	InvestorID   uint64  `gorm:"column:investorId" json:"investorId"`
	Name         string  `gorm:"column:name" json:"name"`
	TotalBalance float64 `gorm:"column:totalBalance" json:"totalBalance"`
}

// GetCifInvestorAccount - Get CIF, investor and account data
func GetCifInvestorAccount(ctx *iris.Context) {
	email := ctx.URLParam("email")

	query := "SELECT cif.id AS \"cifId\", r_cif_investor.\"investorId\" AS \"investorId\", cif.name, account.\"totalBalance\" "
	query += "FROM cif "
	query += "JOIN r_cif_investor ON r_cif_investor.\"cifId\" = cif.id "
	query += "JOIN r_account_investor ON r_account_investor.\"investorId\" = r_cif_investor.\"investorId\" "
	query += "JOIN account ON account.id = r_account_investor.\"accountId\" "
	query += "WHERE cif.username = ? "

	cifInvestorAccountObj := new(CifInvestorAccount)

	services.DBCPsql.Raw(query, email).Scan(&cifInvestorAccountObj)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   cifInvestorAccountObj,
	})
}

func GetCifBorrower (ctx *iris.Context){
	borrowerId := ctx.Get("id")

	query := "SELECT borrower.\"id\" as \"borrowerId\", "
	query += " borrower.\"borrowerNo\" as \"borrowerNo\", "
	query += " borrower.\"isCheckedTerm\" as \"isCheckedTerm\", "
	query += " borrower.\"isCheckedPrivacy\" as \"isCheckedPrivacy\", "
	query += " borrower.\"id\" as \"village\", "
	query += " cif.\"cifNumber\" as \"cifNumber\", "
	query += " cif.\"username\" as \"username\", "
	query += " cif.\"password\" as \"password\", "
	query += " cif.\"name\" as \"name\", "
	query += " cif.\"gender\" as \"gender\", "
	query += " cif.\"placeOfBirth\" as \"placeOfBirth\", "
	query += " cif.\"dateOfBirth\" as \"dateOfBirth\", "
	query += " cif.\"idCardNo\" as \"idCardNo\", "
	query += " cif.\"idCardValidDate\" as \"idCardValidDate\", "
	query += " cif.\"idCardFilename\" as \"idCardFilename\", "
	query += " cif.\"taxCardNo\" as \"taxCardNo\", "
	query += " cif.\"taxCardFilename\" as \"taxCardFilename\", "
	query += " cif.\"maritalStatus\" as \"maritalStatus\", "
	query += " cif.\"mothersName\" as \"mothersName\", "
	query += " cif.\"religion\" as \"religion\", "
	query += " cif.\"address\" as \"address\", "
	query += " cif.\"rt\" as \"rt\", "
	query += " cif.\"rw\" as \"rw\", "
	query += " cif.\"kelurahan\" as \"kelurahan\", "
	query += " cif.\"kecamatan\" as \"kecamatan\", "
	query += " cif.\"province\" as \"province\", "
	query += " cif.\"nationality\" as \"nationality\", "
	query += " cif.\"zipcode\" as \"zipcode\", "
	query += " cif.\"phoneNo\" as \"phoneNo\", "
	query += " cif.\"companyName\" as \"companyName\", "
	query += " cif.\"companyAddress\" as \"companyAddress\", "
	query += " cif.\"occupation\" as \"occupation\", "
	query += " cif.\"income\" as \"income\", "
	query += " cif.\"incomeSourceFund\" as \"incomeSourceFund\", "
	query += " cif.\"incomeSourceCountry\" as \"incomeSourceCountry\", "
	query += " cif.\"isActivated\" as \"isActivated\", "
	query += " cif.\"isValidated\" as \"isValidated\", "
	query += " cif.\"isVerified\" as \"isVerified\" "
	query += "FROM borrower "
	query += "LEFT JOIN r_cif_borrower rcb ON rcb.\"borrowerId\" = borrower.\"id\" "
	query += "LEFT JOIN cif ON cif.\"id\" = rcb.\"cifId\" "
	query += "WHERE borrower.\"id\" = ? "

	borrowerAll := CifBorrower{}

	services.DBCPsql.Raw(query, borrowerId).Scan(&borrowerAll)
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   borrowerAll,
	})
}


func GetCifInvestor (ctx *iris.Context){
	investorId := ctx.Get("id")

	query := "SELECT investor.\"id\" as \"investorId\", "
	query += " investor.\"isCheckedTerm\" as \"isCheckedTerm\", "
	query += " investor.\"isCheckedPrivacy\" as \"isCheckedPrivacy\", "
	query += " investor.\"investorNo\" as \"investorNo\", "
	query += " investor.\"isInstitutional\" as \"isInstitutional\", "
	query += " investor.\"bankName\" as \"bankName\", "
	query += " investor.\"bankBranch\" as \"bankBranch\", "
	query += " investor.\"bankAccountName\" as \"bankAccountName\", "
	query += " investor.\"bankAccountNo\" as \"bankAccountNo\", "
	query += " cif.\"id\" as \"cifId\", "
	query += " cif.\"cifNumber\" as \"cifNumber\", "
	query += " cif.\"username\" as \"username\", "
	query += " cif.\"password\" as \"password\", "
	query += " cif.\"name\" as \"name\", "
	query += " cif.\"gender\" as \"gender\", "
	query += " cif.\"placeOfBirth\" as \"placeOfBirth\", "
	query += " cif.\"dateOfBirth\" as \"dateOfBirth\", "
	query += " cif.\"idCardNo\" as \"idCardNo\", "
	query += " cif.\"idCardValidDate\" as \"idCardValidDate\", "
	query += " cif.\"idCardFilename\" as \"idCardFilename\", "
	query += " cif.\"taxCardNo\" as \"taxCardNo\", "
	query += " cif.\"taxCardFilename\" as \"taxCardFilename\", "
	query += " cif.\"maritalStatus\" as \"maritalStatus\", "
	query += " cif.\"mothersName\" as \"mothersName\", "
	query += " cif.\"religion\" as \"religion\", "
	query += " cif.\"address\" as \"address\", "
	query += " cif.\"rt\" as \"rt\", "
	query += " cif.\"rw\" as \"rw\", "
	query += " cif.\"kelurahan\" as \"kelurahan\", "
	query += " cif.\"kecamatan\" as \"kecamatan\", "
	query += " cif.\"province\" as \"province\", "
	query += " cif.\"nationality\" as \"nationality\", "
	query += " cif.\"zipcode\" as \"zipcode\", "
	query += " cif.\"phoneNo\" as \"phoneNo\", "
	query += " cif.\"companyName\" as \"companyName\", "
	query += " cif.\"companyAddress\" as \"companyAddress\", "
	query += " cif.\"occupation\" as \"occupation\", "
	query += " cif.\"income\" as \"income\", "
	query += " cif.\"incomeSourceFund\" as \"incomeSourceFund\", "
	query += " cif.\"incomeSourceCountry\" as \"incomeSourceCountry\", "
	query += " cif.\"isActivated\" as \"isActivated\", "
	query += " cif.\"isValidated\" as \"isValidated\", "
	query += " cif.\"isVerified\" as \"isVerified\" "
	query += "FROM investor "
	query += "LEFT JOIN r_cif_investor rcb ON rcb.\"investorId\" = investor.\"id\" "
	query += "LEFT JOIN cif ON cif.\"id\" = rcb.\"cifId\" "
	query += "WHERE investor.\"id\" = ? "

	investorAll := CifInvestor{}

	services.DBCPsql.Raw(query, investorId).Scan(&investorAll)
	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   investorAll,
	})
}

func UpdateInvestorCif (ctx *iris.Context){
	investorId := ctx.Get("investorId")
	cifId := ctx.Get("cifId")

	loanData := Loan{}

	fmt.Println(investorId)
	fmt.Println(cifId)

	// services.DBCPsql.Table("investor").Where(" \"investorId\" = ?", investorId).Update("investorId", nil)
	// services.DBCPsql.Table("cif").Where(" \"cifId\" = ?", cif).Update("cifId", nil)




}


