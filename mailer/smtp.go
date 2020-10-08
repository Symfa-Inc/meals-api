package mailer

import (
	"net/smtp"
	"os"
	"strings"

	"github.com/Aiscom-LLC/meals-api/domain"
)

var auth smtp.Auth

// SendEmail sends registration email on provided email
// returns error
func SendEmail(user domain.User, password string, url string) error {
	auth = smtp.PlainAuth("", os.Getenv("SMTP_EMAIL"), os.Getenv("SMTP_PASSWORD"), "smtp.gmail.com")

	r := NewRequest([]string{user.Email},
		"Добро пожаловать в TastyOffice",
		"Здравствуйте, "+user.FirstName+"\n"+
			"TastyOffice приветствует Вас.\n"+
			"Для завершения регистрации пожалуйста перейдите по ссылке ниже и авторизуйтесь в системе (после первой авторизации Ваш аккаунт будет считаться активным).\n"+
			"Ссылка на приложение: "+url+"\n"+
			"Логин: "+user.Email+"\n"+
			"Пароль: "+password+"\n"+
			"Желаем Вам приятного аппетита и хорошего дня!")

	if err := r.SendEmail(); err != nil {
		return err
	}

	return nil
}

// RecoveryPassword sends email with new random generate password
// return error
func RecoveryPassword(user domain.User, password string, url string) error {
	auth = smtp.PlainAuth("", os.Getenv("SMTP_EMAIL"), os.Getenv("SMTP_PASSWORD"), "smtp.gmail.com")

	r := NewRequest([]string{user.Email},
		"TastyOffice востановление пароля",
		"Здравствуйте,\n"+
			user.FirstName+"\n"+
			"Вас приветствует система TastyOffice. Ваш пароль был успешно заменен на новый\n"+
			"Для доступа к приложению используйте ваш логин и новый пароль:\n"+
			"Логин: "+user.Email+"\n"+
			"Пароль: "+password+"\n"+
			"Войти: "+url)

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
