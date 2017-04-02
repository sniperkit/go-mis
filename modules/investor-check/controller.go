package investorCheck

import (
	"fmt"
	"strconv"

	"bitbucket.org/go-mis/modules/cif"
	email "bitbucket.org/go-mis/modules/email"
	"bitbucket.org/go-mis/modules/r"
	va "bitbucket.org/go-mis/modules/virtual-account"
	"bitbucket.org/go-mis/services"
	iris "gopkg.in/kataras/iris.v4"
)

type totalData struct {
	TotalRows int64 `gorm:"column:totalRows" json:"totalRows"`
}

// FetchDatatables -  fetch data based on parameters sent by datatables
func FetchDatatables(ctx *iris.Context) {
	investors := []InvestorCheck{}
	totalData := totalData{}

	query := "SELECT cif.\"name\", cif.\"phoneNo\", cif.\"idCardNo\", \"bankAccountName\", "
	query += "cif.\"taxCardNo\", cif.\"idCardFilename\", cif.\"taxCardFilename\", cif.\"idCardNo\", "
	query += "cif.\"taxCardNo\", array_to_string(array_agg(virtual_account.\"bankName\"),',') as \"virtualAccountBankName\", "
	query += "array_to_string(array_agg(virtual_account.\"virtualAccountNo\"),',') as \"virtualAccountNumber\" "
	query += "FROM investor "
	query += "LEFT JOIN r_investor_virtual_account ON r_investor_virtual_account.\"investorId\" = investor.id "
	query += "LEFT JOIN virtual_account ON virtual_account.id = r_investor_virtual_account.\"vaId\" "
	query += "JOIN r_cif_investor ON r_cif_investor.\"investorId\" = investor.id "
	query += "JOIN cif ON cif.id = r_cif_investor.\"cifId\" "
	query += "AND cif.\"deletedAt\" IS null AND virtual_account.\"deletedAt\" IS null "
	query += "where cif.\"isValidated\" = false and cif.name ~* '[a-z]+' "
	query += "group by cif.\"name\", cif.\"phoneNo\", cif.\"idCardNo\", \"bankAccountName\", cif.\"taxCardNo\", cif.\"idCardNo\", cif.\"taxCardNo\", cif.\"idCardFilename\", cif.\"taxCardFilename\" "

	queryTotalData := "SELECT count(cif.*) as \"totalRows\" "
	queryTotalData += "FROM cif "
	queryTotalData += "WHERE \"isValidated\" = false "
	queryTotalData += "AND \"deletedAt\" IS NULL "

	if ctx.URLParam("search") != "" {
		query += "AND name ~* '" + ctx.URLParam("search") + "' "
		queryTotalData += "AND name ~* '" + ctx.URLParam("search") + "' "
	}

	services.DBCPsql.Raw(query).Scan(&investors)
	services.DBCPsql.Raw(queryTotalData).Scan(&totalData)

	if ctx.URLParam("limit") != "" {
		query += "LIMIT " + ctx.URLParam("limit") + " "
	} else {
		query += "LIMIT 10 "
	}

	if ctx.URLParam("page") != "" {
		query += "OFFSET " + ctx.URLParam("page")
	} else {
		query += "OFFSET 0 "
	}

	services.DBCPsql.Raw(query).Scan(&investors)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status":    "success",
		"totalRows": totalData.TotalRows,
		"data":      investors,
	})
}

// Verify - verify the selected investor
func Verify(ctx *iris.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	// status type: verified or declined
	status := ctx.Param("status")

	if status == "verified" {
		services.DBCPsql.Table("cif").Where("id = ?", id).Update("isValidated", true)

		cifSchema := cif.Cif{}
		services.DBCPsql.Table("cif").Where("id = ?", id).Scan(&cifSchema)

		// get investor id
		inv := &r.RCifInvestor{}
		services.DBCPsql.Table("r_cif_investor").Where("\"cifId\" = ?", id).Scan(&inv)

		// get virtual account
		rInvVa := []r.RInvestorVirtualAccount{}
		services.DBCPsql.Table("r_investor_virtual_account").Where("\"investorId\" = ?", inv.InvestorId).Scan(&rInvVa)

		vaObj := &va.VirtualAccount{}
		userVa := []va.VirtualAccount{}
		for _, val := range rInvVa {
			services.DBCPsql.Table("virtual_account").Where("\"id\" = ?", val.VirtualAccountId).Scan(&vaObj)
			userVa = append(userVa, *vaObj)
		}

		vaData := make(map[string]string)
		for _, val := range userVa {
			if val.BankName == "BRI" {
				vaData["BRI"] = val.VirtualAccountNo
				vaData["BRI_HOLDER"] = val.VirtualAccountName
			} else if val.BankName == "BCA" {
				vaData["BCA"] = val.VirtualAccountNo
				vaData["BCA_HOLDER"] = val.VirtualAccountName
			}
		}

		if cifSchema.Username != "" {
			fmt.Println("Sending email..")
			sendgrid := email.Sendgrid{}
			sendgrid.SetFrom("Amartha", "no-reply@amartha.com")
			sendgrid.SetTo(cifSchema.Name, cifSchema.Username)
			sendgrid.SetSubject(cifSchema.Name + ", Verifikasi Data Anda Berhasil")
			sendgrid.VerifiedBodyEmail("VERIFIED_DATA", cifSchema.Name, cifSchema.Username, vaData)
			sendgrid.SendEmail()
		}

		if cifSchema.PhoneNo != "" {
			// send sms notification
			fmt.Println("Sending sms ... ")
			twilio := services.InitTwilio()
			message := "Selamat data Anda sudah terverifikasi. Silakan login ke dashboard Anda dan mulai berinvestasi. www.amartha.com"
			twilio.SetParam(cifSchema.PhoneNo, message)
			twilio.SendSMS()
		}

	} else {
		cifSchema := cif.Cif{}
		services.DBCPsql.Table("cif").Where("id = ?", id).Scan(&cifSchema)

		sendgrid := email.Sendgrid{}
		sendgrid.SetFrom("Amartha", "no-reply@amartha.com")
		sendgrid.SetTo(cifSchema.Name, cifSchema.Username)
		sendgrid.SetSubject(cifSchema.Name + ", Verifikasi Data Anda Gagal")
		sendgrid.SetVerificationBodyEmail("UNVERIFIED_DATA", cifSchema.Name, cifSchema.Name, cifSchema.Username, "")
		sendgrid.SendEmail()
	}

	ctx.JSON(iris.StatusOK, iris.Map{
		"status":             "success",
		"verificationStatus": status,
	})
}

// Verified - verify the selected investor
func Verified(ctx *iris.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	// status type: verified or declined
	status := ctx.Param("status")

	if status == "verified" {
		services.DBCPsql.Table("cif").Where("id = ?", id).Update("isVerified", true)

		cifSchema := cif.Cif{}
		services.DBCPsql.Table("cif").Where("id = ?", id).Scan(&cifSchema)

		// get investor id
		inv := &r.RCifInvestor{}
		services.DBCPsql.Table("r_cif_investor").Where("\"cifId\" = ?", id).Scan(&inv)

		// get virtual account
		rInvVa := []r.RInvestorVirtualAccount{}
		services.DBCPsql.Table("r_investor_virtual_account").Where("\"investorId\" = ?", inv.InvestorId).Scan(&rInvVa)

		vaObj := &va.VirtualAccount{}
		userVa := []va.VirtualAccount{}
		for _, val := range rInvVa {
			services.DBCPsql.Table("virtual_account").Where("\"id\" = ?", val.VirtualAccountId).Scan(&vaObj)
			userVa = append(userVa, *vaObj)
		}

		vaData := make(map[string]string)
		for _, val := range userVa {
			if val.BankName == "BRI" {
				vaData["BRI"] = val.VirtualAccountNo
				vaData["BRI_HOLDER"] = val.VirtualAccountName
			} else if val.BankName == "BCA" {
				vaData["BCA"] = val.VirtualAccountNo
				vaData["BCA_HOLDER"] = val.VirtualAccountName
			}
		}

		if cifSchema.Username != "" {
			fmt.Println("Sending email..")
			sendgrid := email.Sendgrid{}
			sendgrid.SetFrom("Amartha", "no-reply@amartha.com")
			sendgrid.SetTo(cifSchema.Name, cifSchema.Username)
			sendgrid.SetSubject(cifSchema.Name + ", Verifikasi Data Anda Berhasil")
			sendgrid.VerifiedBodyEmail("VERIFIED_DATA", cifSchema.Name, cifSchema.Username, vaData)
			sendgrid.SendEmail()
		}

		if cifSchema.PhoneNo != "" {
			// send sms notification
			fmt.Println("Sending sms ... ")
			twilio := services.InitTwilio()
			message := "Selamat data Anda sudah terverifikasi. Silakan login ke dashboard Anda dan mulai berinvestasi. www.amartha.com"
			twilio.SetParam(cifSchema.PhoneNo, message)
			twilio.SendSMS()
		}

	} else {
		cifSchema := cif.Cif{}
		services.DBCPsql.Table("cif").Where("id = ?", id).Scan(&cifSchema)

		sendgrid := email.Sendgrid{}
		sendgrid.SetFrom("Amartha", "no-reply@amartha.com")
		sendgrid.SetTo(cifSchema.Name, cifSchema.Username)
		sendgrid.SetSubject(cifSchema.Name + ", Verifikasi Data Anda Gagal")
		sendgrid.SetVerificationBodyEmail("UNVERIFIED_DATA", cifSchema.Name, cifSchema.Name, cifSchema.Username, "")
		sendgrid.SendEmail()
	}

	ctx.JSON(iris.StatusOK, iris.Map{
		"status":             "success",
		"verificationStatus": status,
	})
}
