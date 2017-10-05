package email

func SendEmailVerificationSuccess(email string, name string, va_bca string, va_bca_name string, va_mandiri string, va_mandiri_name string) {

	var subs = map[string]interface{}{
		"first_name":      name,
		"va_bca":          va_bca,
		"va_bca_name":     va_bca_name,
		"va_mandiri":      va_mandiri,
		"va_mandiri_name": va_mandiri_name,
	}

	mandrill := new(Mandrill)
	mandrill.SetFrom("hello@amartha.com")
	mandrill.SetTo(email)
	mandrill.SetBcc("investing@amartha.com")
	mandrill.SetSubject("[Amartha] Selamat Data Anda Sudah Terverifikasi ")
	mandrill.SetTemplateAndRawBody("verification_success_v2", subs)
	mandrill.SendEmail()

}

// SendEmailVerificationFailed - Send email to investor because of his incomplete data
func SendEmailVerificationFailed(email string, name string, reasons []string) {
	var subs = map[string]interface{}{
		"first_name": name,
		"reason":     reasons,
	}
	mandrill := new(Mandrill)
	mandrill.SetFrom("hello@amartha.com")
	mandrill.SetTo(email)
	mandrill.SetSubject("[Amartha] Verifikasi Data Anda Gagal ")
	mandrill.SetTemplateAndRawBody("verification_failed", subs)
	mandrill.SendEmail()
}

func SendEmailInvestmentSuccess(email string, name string, transferDate string, transferAmount string, transferDestination string) {

	var subs = map[string]interface{}{
		"first_name":           name,
		"transfer_date":        transferDate,
		"transfer_amount":      transferAmount,
		"transfer_destination": transferDestination,
	}

	mandrill := new(Mandrill)
	mandrill.SetFrom("hello@amartha.com")
	mandrill.SetTo(email)
	mandrill.SetSubject("[Amartha] Selamat, Transaksi Anda Berhasil")
	mandrill.SetTemplateAndRawBody("investment_success", subs)
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
