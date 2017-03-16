package email

import (
	"fmt"

	"bitbucket.org/go-mis/config"
	"github.com/parnurzeal/gorequest"
)

// SendgridEmailTemplate - select email template
func SendgridEmailTemplate(key string) string {
	var template = make(map[string]string, 10)
	template["UNVERIFIED_DATA"] = "088a5e63-a4f6-4ccd-bed4-765ae5021088"
	template["VERIFIED_DATA"] = "3099d51a-294d-4f89-8449-e13d4c3d2f8b"


	return template[key]
}

// Sendgrid - Init Class
type Sendgrid struct {
	emailParam struct {
		From      EmailUser              `json:"from" `
		To        []EmailUser            `json:"to" `
		Template  string                 `json:"template" `
		Subject   string                 `json:"subject" `
		SecretKey string                 `json:"secretKey" `
		Subs      map[string]interface{} `json:"subs" `
	}
}

// EmailUser - Email User Struct
type EmailUser struct {
	Name  string `json:"name" `
	Email string `json:"email" `
}

// SetFrom - SetFrom Email
func (s *Sendgrid) SetFrom(name string, email string) {
	s.emailParam.From = EmailUser{Name: name, Email: email}
}

// SetTo - SetTo Email
func (s *Sendgrid) SetTo(name string, email string) {
	user := EmailUser{Name: name, Email: email}
	s.emailParam.To = []EmailUser{user}
	//	EmailUser{Name: name, Email: email}
}

// Subject - Subject Email
func (s *Sendgrid) SetSubject(subject string) {
	s.emailParam.Subject = subject
}

func (s *Sendgrid) VerifiedBodyEmail(template string, first_name string, email string, vaData map[string]string) {
	s.emailParam.Template = SendgridEmailTemplate(template)
	var subs map[string]interface{} = map[string]interface{}{
		"[%first_name%]": first_name,
		"[%va_bca%]":  vaData["BCA"],
		"[%va_bri%]":  vaData["BRI"],
		"[%va_name%]": vaData["BCA_HOLDER"],
		"[%email%]": email,
	}

	fmt.Println("SUB:",subs)
	s.emailParam.Subs = subs
}

func (s *Sendgrid) SetVerificationBodyEmail(template string, first_name string, full_name string, email string, url string) {
	s.emailParam.Template = SendgridEmailTemplate(template)
	var subs map[string]interface{} = map[string]interface{}{
		"[%first_name%]": first_name,
		"[%full_name%]":  full_name,
		"[%email%]":      email,
		"[%url%]":        url,
	}
	s.emailParam.Subs = subs
}

// SendEmail - send Notif
func (s Sendgrid) SendEmail() {
	s.emailParam.SecretKey = "n0de-U>lo4d3r"
	request := gorequest.New()
	_, body, _ := request.Post(config.UploaderApiPath + "/email/send/sendgrid").
		Send(s.emailParam).
		End()

	//res2B, _ := json.Marshal(s.emailParam)
	//fmt.Println(string(res2B))

	fmt.Println(body)
}
