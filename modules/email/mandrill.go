package email

import (
	"fmt"
	"log"

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
		Bucket    bool                   `json:"bucket"`
		Subs      map[string]interface{} `json:"subs" `
	}
}

// SetFrom - SetFrom Email
func (m *Mandrill) SetFrom(email string) {
	m.emailParam.From = email
}

// SetBCC - SetBCC Email
func (m *Mandrill) SetBCC(email string) {
	m.emailParam.Bcc = email
}

/* Bucket - for choose email template
   if u bucket true  it use from CK
               false it use from your local*/
func (m *Mandrill) SetBucket(bucket bool) {
	m.emailParam.Bucket = bucket
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
	log.Println("node uploader link: ", config.UploaderApiPath+"/email/send/mandrill")
	m.emailParam.SecretKey = "n0de-U>lo4d3r"
	request := gorequest.New()
	_, body, _ := request.Post(config.UploaderApiPath + "/email/send/mandrill").
		Send(m.emailParam).
		End()

	fmt.Println(body)
}
