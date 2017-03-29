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
