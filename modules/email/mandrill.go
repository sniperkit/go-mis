package email

import (
	"fmt"

	"bitbucket.org/go-mis/config"

	"github.com/parnurzeal/gorequest"
)

// Mandrill - Init Class
type Mandrill struct {
	emailParam struct {
		From      string                 `json:"from" `
		To        string                 `json:"to" `
		Bcc       string                 `json:"bcc" `
		Template  string                 `json:"template" `
		Subject   string                 `json:"subject" `
		SecretKey string                 `json:"secretKey" `
		Subs      map[string]interface{} `json:"subs" `
	}
}

// SetFrom - SetFrom Email
func (m *Mandrill) SetFrom(email string) {
	m.emailParam.From = email
}

// SetTo - SetTo Email
func (m *Mandrill) SetTo(email string) {
	m.emailParam.To = email
}

// SetTo - SetTo Email
func (m *Mandrill) SetBcc(email string) {
	m.emailParam.Bcc = email
}

// Subject - Subject Email
func (m *Mandrill) SetSubject(subject string) {
	m.emailParam.Subject = subject
}

func (m *Mandrill) SetTemplateAndRawBody(template string, raw map[string]interface{}) {
	m.emailParam.Template = template
	m.emailParam.Subs = raw
}

// SendEmail - send Notif
func (m Mandrill) SendEmail() {
	m.emailParam.SecretKey = "n0de-U>lo4d3r"
	request := gorequest.New()
	_, body, _ := request.Post(config.UploaderApiPath + "email/send/mandrill").
		Send(m.emailParam).
		End()

	fmt.Println(body)
}
