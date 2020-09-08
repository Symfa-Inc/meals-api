package mailer

import (
	"bytes"
	"github.com/Aiscom-LLC/meals-api/interfaces"
	"html/template"
	"net/smtp"
	"os"
	"strings"
)

var auth smtp.Auth

type templateData struct {
	Name     string
	URL      string
	Email    string
	Password string
}

// SendEmail sends registration email on provided email
// returns error
func SendEmail(user interfaces.User, password string) error {

	auth = smtp.PlainAuth("", os.Getenv("SMTP_EMAIL"), os.Getenv("SMTP_PASSWORD"), "smtp.gmail.com")
	templateData := templateData{
		Name:     user.FirstName + " " + user.LastName,
		URL:      "meals.d1.aisnovations.com/login",
		Email:    user.Email,
		Password: password,
	}

	r := NewRequest([]string{user.Email}, "Invitation to AIS Meals", "")
	dir, _ := os.Getwd()

	if err := r.ParseTemplate(dir+"/src/mailer/email_template.html", templateData); err != nil {
		return err
	}

	if err := r.SendEmail(); err != nil {
		return err
	}

	return nil
}

// Request struct
type Request struct {
	to      []string
	subject string
	body    string
}

// NewRequest creates pointer to Request and fills it with
// provided data, returns pointer
func NewRequest(to []string, subject, body string) *Request {
	return &Request{
		to:      to,
		subject: subject,
		body:    body,
	}
}

// SendEmail is method on Request struct
// which sends message, returns error
func (r *Request) SendEmail() error {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + r.subject + "!\n"
	to := "To: " + strings.Join(r.to, ", ") + "\n"
	msg := []byte(subject + to + mime + "\n" + r.body)
	addr := "smtp.gmail.com:587"

	if err := smtp.SendMail(addr, auth, os.Getenv("SMTP_EMAIL"), r.to, msg); err != nil {
		return err
	}
	return nil
}

// ParseTemplate takes two arguments, parses struct and putting it in template
// return error
func (r *Request) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}
