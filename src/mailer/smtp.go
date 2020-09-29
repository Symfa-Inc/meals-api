package mailer

import (
	"net/smtp"
	"os"
	"strings"

	"github.com/Aiscom-LLC/meals-api/src/domain"
)

var auth smtp.Auth

// SendEmail sends registration email on provided email
// returns error
func SendEmail(user domain.User, password string) error {
	auth = smtp.PlainAuth("", os.Getenv("SMTP_EMAIL"), os.Getenv("SMTP_PASSWORD"), "smtp.gmail.com")

	r := NewRequest([]string{user.Email},
		"Добро пожаловать в TastyOffice",
		"Здравствуйте,\n"+
			user.FirstName+"\n"+
			"Вас приветствует система TastyOffice. Для завершения регистрации нажмите на кнопку войти, что бы подтвердить учетную запись\n"+
			"Для доступа к приложению используйте ваш логин и пароль:\n"+
			"Логин: "+user.Email+"\n"+
			"Пароль: "+password+"\n"+
			"Войти: https://meals.d1.aisnovations.com/login")

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
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + r.subject + "!\n"
	to := "To: " + strings.Join(r.to, ", ") + "\n"
	msg := []byte(subject + to + mime + "\n" + r.body)
	addr := "smtp.gmail.com:587"

	if err := smtp.SendMail(addr, auth, os.Getenv("SMTP_EMAIL"), r.to, msg); err != nil {
		return err
	}
	return nil
}
