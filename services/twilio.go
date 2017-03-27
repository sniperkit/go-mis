package services

import (
	"fmt"

	"github.com/parnurzeal/gorequest"
)

// Twilio - Init Class
type Twilio struct {
	ACCOUNT_SID string
	AUTH_TOKEN  string
	FROM_NUMBER string

	param struct {
		To   string
		From string
		Body string
	}
}

// Nes Instance
func InitTwilio() *Twilio {
	return &Twilio{ACCOUNT_SID: "AC8b453eb30d605054c3a59e06cb88634b", AUTH_TOKEN: "d84ed6d8e25d654b75efb7418c90077c", FROM_NUMBER: "+16317063401"}
}

// SetMarketplaceParam - set value marketPlace Notif
func (t *Twilio) SetParam(to string, body string) {
	t.param.To = to
	t.param.From = t.FROM_NUMBER
	t.param.Body = body
}

// SendNotifMarketplace - send Notif
func (t Twilio) SendSMS() {

	type ParamTest struct {
		To   string
		From string
		Body string
	}

	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + t.ACCOUNT_SID + "/Messages.json"

	request := gorequest.New()
	_, body, _ := request.Post(urlStr).
		SetBasicAuth(t.ACCOUNT_SID, t.AUTH_TOKEN).
		Set("Accept", "application/json").
		Set("Content-Type", "application/x-www-form-urlencoded").
		Send(t.param).
		End()

	fmt.Println(body)
}
