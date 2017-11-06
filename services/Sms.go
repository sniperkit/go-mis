package services

import (
	"strings"
	"fmt"
)

func SendSMS(phoneNumber string, message string) {
	if strings.HasPrefix(phoneNumber, "0") {
		phoneNumber = strings.Replace(phoneNumber, "0", "+62", 1)
	}
	if strings.HasPrefix(phoneNumber, "+62811") || strings.HasPrefix(phoneNumber, "62811") ||
		strings.HasPrefix(phoneNumber, "+62812") || strings.HasPrefix(phoneNumber, "62812") ||
		strings.HasPrefix(phoneNumber, "+62813") || strings.HasPrefix(phoneNumber, "62813") ||
		strings.HasPrefix(phoneNumber, "+62821") || strings.HasPrefix(phoneNumber, "62821") ||
		strings.HasPrefix(phoneNumber, "+62822") || strings.HasPrefix(phoneNumber, "62822") ||
		strings.HasPrefix(phoneNumber, "+62823") || strings.HasPrefix(phoneNumber, "62823") ||
		strings.HasPrefix(phoneNumber, "+62852") || strings.HasPrefix(phoneNumber, "62852") ||
		strings.HasPrefix(phoneNumber, "+62853") || strings.HasPrefix(phoneNumber, "62853") ||
		strings.HasPrefix(phoneNumber, "+62851") || strings.HasPrefix(phoneNumber, "62851") {
		nextmo := InitNEXMO()
		nextmo.SetParam(phoneNumber, message)
		nextmo.SendSMS()
		return
	}
	fmt.Println("#phoneNumber", phoneNumber, "Sending sms", message)
	twilio := InitTwilio()

	twilio.SetParam(phoneNumber, message)
	twilio.SendSMS()
}