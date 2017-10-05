package services

import (
	"fmt"

	"github.com/parnurzeal/gorequest"
)

// SNS - Init Class
type NEXMO struct {
	API_KEY    string
	API_SECRET string
	param      struct {
		To        string `json:"to" `
		From      string `json:"from" `
		Text      string `json:"text" `
		ApiKey    string `json:"api_key" `
		ApiSecret string `json:"api_secret" `
	}
}

// InitSNS - initial sns third party
func InitNEXMO() *NEXMO {
	return &NEXMO{API_KEY: "aa0b5754", API_SECRET: "16335aed1d56a69c"}
}

func (n *NEXMO) SetParam(to string, body string) {
	n.param.To = to
	n.param.Text = body
	n.param.From = "AMARTHA"
	n.param.ApiKey = n.API_KEY
	n.param.ApiSecret = n.API_SECRET
}

// SendSMS - sending a message using sns
func (n NEXMO) SendSMS() {

	urlStr := "https://rest.nexmo.com/sms/json"
	request := gorequest.New()

	_, body, _ := request.Post(urlStr).
		Send(n.param).
		End()

	fmt.Println(body)
}

