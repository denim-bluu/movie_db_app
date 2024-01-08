package mailer

import (
	"bytes"
	"embed"
	"text/template"
	"time"

	"github.com/go-mail/mail/v2"
)

//go:embed "templates"
var templateFS embed.FS

type Smtp struct {
	Host     string
	Port     int
	Username string
	Password string
	Sender   string
}

type Mailer struct {
	dialer *mail.Dialer
	Smtp   Smtp
}

func New(smtp Smtp) *Mailer {
	dialer := mail.NewDialer(smtp.Host, smtp.Port, smtp.Username, smtp.Password)
	dialer.Timeout = 5 * time.Second

	mailer := &Mailer{
		dialer: dialer,
		Smtp:   smtp,
	}
	return mailer
}

func (m *Mailer) SendWithRetry(msg *mail.Message) error {
	var err error
	for i := 1; i <= 3; i++ {
		err := m.dialer.DialAndSend(msg)
		if nil == err {
			return nil
		}

		time.Sleep(500 * time.Millisecond)
	}
	return err
}

func (m *Mailer) CreateMessage(recipient, subject, htmlBody, plainBody string) *mail.Message {
	msg := mail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.Smtp.Sender)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", htmlBody)
	msg.AddAlternative("text/plain", plainBody)
	return msg
}

func (m *Mailer) parseTemplates(tmpl *template.Template, data any) (subject, plainBody, htmlBody string, err error) {
	subject, err = m.executeTemplate(tmpl, "subject", data)
	if err != nil {
		return
	}

	plainBody, err = m.executeTemplate(tmpl, "plainBody", data)
	if err != nil {
		return
	}

	htmlBody, err = m.executeTemplate(tmpl, "htmlBody", data)
	return
}

func (m *Mailer) executeTemplate(tmpl *template.Template, name string, data any) (string, error) {
	buf := new(bytes.Buffer)
	err := tmpl.ExecuteTemplate(buf, name, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (m *Mailer) Send(recipient, templateFile string, data any) error {
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject, plainBody, htmlBody, err := m.parseTemplates(tmpl, data)
	if err != nil {
		return err
	}

	msg := m.CreateMessage(recipient, subject, htmlBody, plainBody)

	return m.SendWithRetry(msg)
}
