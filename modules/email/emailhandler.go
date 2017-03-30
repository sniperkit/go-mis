package email

func sendEmailVerificationSuccess(email string, name string, va_bca string, va_bca_name string, va_bri string, va_bri_name string) {

	var subs = map[string]interface{}{
		"first_name":  name,
		"va_bca":      va_bca,
		"va_bca_name": va_bca_name,
		"va_bri":      va_bri,
		"va_bri_name": va_bri_name,
	}

	mandrill := new(Mandrill)
	mandrill.SetFrom("hello@amartha.com")
	mandrill.SetTo(email)
	mandrill.SetSubject("Verification Success")
	mandrill.SetTemplateAndRawBody("verification_success", subs)
	mandrill.SendEmail()

}

func sendEmailVerificationFailed(email string, name string) {

	var subs = map[string]interface{}{
		"first_name": name,
	}

	mandrill := new(Mandrill)
	mandrill.SetFrom("hello@amartha.com")
	mandrill.SetTo(email)
	mandrill.SetSubject("Verification Failed")
	mandrill.SetTemplateAndRawBody("verification_failed", subs)
	mandrill.SendEmail()

}

func sendEmailInvestmentSuccess(email string, name string, transferDate string, transferAmount string, transferDestination string) {

	var subs = map[string]interface{}{
		"first_name":           name,
		"transfer_date":        transferDate,
		"transfer_amount":      transferAmount,
		"transfer_destination": transferDestination,
	}

	mandrill := new(Mandrill)
	mandrill.SetFrom("hello@amartha.com")
	mandrill.SetTo(email)
	mandrill.SetSubject("Investment Success")
	mandrill.SetTemplateAndRawBody("investment_success", subs)
	mandrill.SendEmail()

}

func sendEmailInvestmentFailed(email string, name string) {

	var subs = map[string]interface{}{
		"first_name": name,
	}

	mandrill := new(Mandrill)
	mandrill.SetFrom("hello@amartha.com")
	mandrill.SetTo(email)
	mandrill.SetSubject("Investment Failed")
	mandrill.SetTemplateAndRawBody("investment_failed", subs)
	mandrill.SendEmail()

}

func sendEmailDisbursementSucccess(subject string, template string, email string, name string, borrowerName string, purpose string, plafon string, tenor string) {
	sendEmailDisbursement("Disbursement Succcess", "disbursement_success", email, name, borrowerName, purpose, plafon, tenor)
}

func sendEmailDisbursementPending(subject string, template string, email string, name string, borrowerName string, purpose string, plafon string, tenor string) {
	sendEmailDisbursement("Disbursement Pending", "disbursement_pending", email, name, borrowerName, purpose, plafon, tenor)
}

func sendEmailDisbursementPostpone(subject string, template string, email string, name string, borrowerName string, purpose string, plafon string, tenor string) {
	sendEmailDisbursement("Disbursement Failed", "disbursement_postpone", email, name, borrowerName, purpose, plafon, tenor)
}

func sendEmailDisbursement(subject string, template string, email string, name string, borrowerName string, purpose string, plafon string, tenor string) {

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

func sendEmailCashout(email string, name string, cashoutID string) {

	var subs = map[string]interface{}{
		"first_name": name,
		"cashout_id": cashoutID,
	}

	mandrill := new(Mandrill)
	mandrill.SetFrom("hello@amartha.com")
	mandrill.SetTo(email)
	mandrill.SetSubject("Cashout Success")
	mandrill.SetTemplateAndRawBody("cashout_success", subs)
	mandrill.SendEmail()

}
