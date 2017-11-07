package email

import (
	"fmt"
	"time"

	"bitbucket.org/go-mis/services"
	humanize "github.com/dustin/go-humanize"
)

type ClientInvestor struct {
	ID               uint64    `json:"id" gorm:"column:id"`
	CreditScoreGrade string    `json:"creditScoreGrade" gorm:"column:creditScoreGrade"`
	Plafond          float64   `json:"plafond" gorm:"column:plafond"`
	IsInsurance       bool      `json:"isInsurance" gorm:"column:isInsurance"`
	Tenor            uint64    `json:"tenor" gorm:"column:tenor"`
	Purpose          string    `json:"purpose" gorm:"column:purpose"`
	BorrowerName     string    `json:"borrowerName" gorm:"column:borrowerName"`
	DisbursementDate time.Time `json:"disbursementDate" gorm:"column:disbursementDate"`
	RevenueProjection float64   `json:"revenueProjection" gorm:"column:revenueProjection"`
}

type ClientInvestorEmailBody struct {
	Plafond          string `json:"plafond" `
	Tenor            uint64 `json:"tenor" `
	Purpose          string `json:"purpose" `
	BorrowerName     string `json:"borrowerName" `
	DisbursementDate string `json:"disbursementDate"`
	RevenueProjection string   `json:"revenueProjection" gorm:"column:revenueProjection"`
}

type Donations struct {
	Name   string `json:"name"`
	Total  uint64 `json:"total"`
	Amount string `json:"amount"`
}

//GetClientInvestor get all investor data for email notification
func GetClientInvestor(InvestorID uint64, OrderNo string, Stage string) ([]ClientInvestor, []ClientInvestorEmailBody, float64, float64) {
	var InsuranceRate float64 = 0.015
	clientInvestors := []ClientInvestor{}
	var resultsEmails []ClientInvestorEmailBody
	var total float64
	var totalInsurance, totalInsuranceAmount float64

	query := `select distinct loan.id,
				loan."creditScoreGrade",
				loan.plafond,
				loan.tenor,
				loan.purpose,
				loan."isInsurance",
				cif."name" as "borrowerName",
				disbursement."disbursementDate",
				(loan.plafond + (loan.plafond * loan.rate * (pp."returnOfInvestment"))) "revenueProjection"
			from loan
				join r_loan_borrower ON r_loan_borrower."loanId" = loan.id
				join borrower ON borrower.id = r_loan_borrower."borrowerId"
				join r_cif_borrower ON r_cif_borrower."borrowerId" = borrower.id
				join cif ON cif.id = r_cif_borrower."cifId"
				join r_investor_product_pricing_loan ON r_investor_product_pricing_loan."loanId" = loan.id
				join product_pricing pp on pp.id = r_investor_product_pricing_loan."productPricingId"
				join r_loan_order rlo on rlo."loanId" = loan.id
				join loan_order lo on lo.id  = rlo."loanOrderId"
				join r_loan_disbursement on r_loan_disbursement."loanId" = loan.id
				join disbursement on disbursement.id = r_loan_disbursement."disbursementId"
			where r_investor_product_pricing_loan."investorId" = ? and loan.stage = ?
			and lo."orderNo" = ? and lo."deletedAt" is null
			and r_investor_product_pricing_loan."deletedAt" is null`
	// log.Println(query)
	services.DBCPsql.Raw(query, InvestorID, Stage, OrderNo).Scan(&clientInvestors)

	totalInsurance = 0
	total = 0
	for _, val := range clientInvestors {
		total += val.Plafond
		var emailBody ClientInvestorEmailBody
		emailBody.BorrowerName = val.BorrowerName
		emailBody.Purpose = val.Purpose
		plafondStr := fmt.Sprintf("Rp %s", humanize.Commaf(val.Plafond))
		emailBody.Plafond = plafondStr
		emailBody.Tenor = val.Tenor
		emailBody.DisbursementDate = val.DisbursementDate.Format("02-01-2006")
		emailBody.RevenueProjection = fmt.Sprintf("Rp %s", humanize.Commaf(val.RevenueProjection))
		resultsEmails = append(resultsEmails, emailBody)

		totalInsurance = totalInsurance + 1
	}
	totalInsuranceAmount = (total * (totalInsurance * InsuranceRate))
		total = total + totalInsuranceAmount
	return clientInvestors, resultsEmails, total, totalInsuranceAmount
}

func SendEmailIInvestmentSuccess(name string, toEmail string, OrderNo string, investorId uint64, voucherAmount float64) {
	fmt.Println("#Email send investment success", toEmail)
	_, clientInvestorsEmail, total, totalInsuranceAmount := GetClientInvestor(investorId, OrderNo, "INVESTOR")

	subs := map[string]interface{}{
		"first_name": name,
		"clients":    clientInvestorsEmail,
		"total":      fmt.Sprintf("Rp. %s", humanize.Commaf(total)),
		"donations":  []Donations{},
		"totalInsuranceAmount":  fmt.Sprintf("Rp. %s", humanize.Commaf(totalInsuranceAmount)),
		"voucherAmount":         voucherAmount,
		"voucherAmountHumanize": fmt.Sprintf("Rp. %s", humanize.Commaf(voucherAmount)),
	}

	mandrill := &Mandrill{}
	mandrill.SetFrom("Amartha <hello@amartha.com>")
	mandrill.SetTo(toEmail)
	mandrill.SetBCC("investing@amartha.com")
	mandrill.SetSubject("[Amartha] Selamat, Transaksi Anda Berhasil")
	if totalInsuranceAmount > 0 {
		mandrill.SetTemplateAndRawBody("investment_balance_success_v3_jamkrindo", subs)
	} else {
		mandrill.SetTemplateAndRawBody("investment_balance_success_v3", subs)
	}
	mandrill.SetBucket(true)
	mandrill.SendEmail()
}

func SendEmailVerificationSuccess(email string, name string, va_bca string, va_bca_name string, va_mandiri string, va_mandiri_name string) {
	var subs = map[string]interface{}{
		"first_name":      name,
		"va_bca":          va_bca,
		"va_bca_name":     va_bca_name,
		"va_mandiri":      va_mandiri,
		"va_mandiri_name": va_mandiri_name,
	}
	fmt.Println("Subs email: ", subs)
	mandrill := new(Mandrill)
	mandrill.SetFrom("hello@amartha.com")
	mandrill.SetTo(email)
	mandrill.SetBcc("investing@amartha.com")
	mandrill.SetSubject("[Amartha] Selamat Data Anda Sudah Terverifikasi ")
	mandrill.SetTemplateAndRawBody("verification_success_v2", subs)
	mandrill.SetBucket(true)
	mandrill.SendEmail()
}

// SendEmailVerificationFailed - Send email to investor because of his incomplete data
func SendEmailVerificationFailed(email string, name string, reasons []string) {
	var subs = map[string]interface{}{
		"investor_name": name,
		"reasons":       reasons,
	}
	mandrill := new(Mandrill)
	mandrill.SetFrom("hello@amartha.com")
	mandrill.SetTo(email)
	mandrill.SetSubject("[Amartha] Verifikasi Data Anda Gagal ")
	mandrill.SetTemplateAndRawBody("decline_v1", subs)
	mandrill.SetBucket(true)
	mandrill.SendEmail()
}

func SendEmailInvestmentFailed(email string, name string) {

	var subs = map[string]interface{}{
		"first_name": name,
	}

	mandrill := new(Mandrill)
	mandrill.SetFrom("hello@amartha.com")
	mandrill.SetTo(email)
	mandrill.SetSubject("[Amartha] Transaksi Anda Gagal")
	mandrill.SetTemplateAndRawBody("investment_failed", subs)
	mandrill.SendEmail()

}

func SendEmailDisbursementSucccess(subject string, template string, email string, name string, borrowerName string, purpose string, plafon string, tenor string, totalPeople string, totalFund string, disbursementDate string) {
	var subs = map[string]interface{}{
		"first_name":        name,
		"disbursement_date": disbursementDate,
		"total_people":      totalPeople,
		"total_fund":        totalFund,
		"borrower_name":     borrowerName,
		"purpose":           purpose,
		"plafon":            plafon,
		"tenor":             tenor,
	}

	mandrill := new(Mandrill)
	mandrill.SetFrom("hello@amartha.com")
	mandrill.SetTo(email)
	mandrill.SetSubject("[Amartha] Pencairan Mitra Usaha Anda Telah Berhasil")
	mandrill.SetTemplateAndRawBody("disbursement_success", subs)
	mandrill.SendEmail()
}

func SendEmailDisbursementPending(subject string, template string, email string, name string, borrowerName string, purpose string, plafon string, tenor string, totalPeople string, totalFund string) {
	SendEmailDisbursement("[Amartha] Pencairan Mitra Usaha Anda Gagal", "disbursement_pending", email, name, borrowerName, purpose, plafon, tenor, totalPeople, totalFund)
}

func SendEmailDisbursementPostpone(subject string, template string, email string, name string, borrowerName string, purpose string, plafon string, tenor string, totalPeople string, totalFund string) {
	SendEmailDisbursement("[Amartha] Pencairan Mitra Usaha Anda Gagal", "disbursement_postpone", email, name, borrowerName, purpose, plafon, tenor, totalPeople, totalFund)
}

func SendEmailDisbursement(subject string, template string, email string, name string, borrowerName string, purpose string, plafon string, tenor string, totalPeople string, totalFund string) {

	var subs = map[string]interface{}{
		"first_name":    name,
		"borrower_name": borrowerName,
		"purpose":       purpose,
		"plafon":        plafon,
		"tenor":         tenor,
	}

	mandrill := new(Mandrill)
	mandrill.SetFrom("hello@amartha.com")
	mandrill.SetTo(email)
	mandrill.SetSubject(subject)
	mandrill.SetTemplateAndRawBody(template, subs)
	mandrill.SendEmail()

}

func SendEmailCashout(email string, name string, cashoutID string) {

	var subs = map[string]interface{}{
		"first_name": name,
		"cashout_id": cashoutID,
	}

	mandrill := new(Mandrill)
	mandrill.SetFrom("hello@amartha.com")
	mandrill.SetTo(email)
	mandrill.SetSubject("[Amartha] Cashout Request Anda Telah Kami Terima")
	mandrill.SetTemplateAndRawBody("cashout_success", subs)
	mandrill.SendEmail()

}

func SendEmailUpkDelay(email string, name string, totalPeople string, totalFund string, borrowerName string, purpose string, plafon string, tenor string) {

	var subs = map[string]interface{}{
		"first_name":    name,
		"total_people":  totalPeople,
		"total_fund":    totalFund,
		"borrower_name": borrowerName,
		"purpose":       purpose,
		"plafon":        plafon,
		"tenor":         tenor,
	}

	mandrill := new(Mandrill)
	mandrill.SetFrom("hello@amartha.com")
	mandrill.SetTo(email)
	mandrill.SetSubject("[Amartha] Penundaan Pencairan")
	mandrill.SetTemplateAndRawBody("upk_delay", subs)
	mandrill.SendEmail()

}
