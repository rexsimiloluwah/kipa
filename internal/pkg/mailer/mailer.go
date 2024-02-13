package mailer

import (
	"bytes"
	"html/template"
	"keeper/internal/config"
	"log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sirupsen/logrus"
)

type Mailer struct {
	Client      *sendgrid.Client
	SenderEmail string
	SenderName  string
}

func NewMailer(cfg *config.Config) *Mailer {
	sendgridClient := sendgrid.NewSendClient(cfg.SendgridAPIKey)
	return &Mailer{
		Client:      sendgridClient,
		SenderEmail: cfg.EmailSenderAddr,
		SenderName:  cfg.EmailSenderName,
	}
}

func (m *Mailer) SendEmail(receiverEmailAddr string, receiverName string, subject string, body string) error {
	from := mail.NewEmail(m.SenderName, m.SenderEmail)
	to := mail.NewEmail(receiverName, receiverEmailAddr)
	message := mail.NewSingleEmail(from, subject, to, body, body)
	response, err := m.Client.Send(message)
	if err != nil {
		return err
	}
	logrus.Debug("email response: ", response)
	return nil
}

// Send a mail containing the email verification link
func (m *Mailer) SendEmailVerificationMail(receiverEmailAddr string, receiverName string, subject string, templateData interface{}) error {
	body, err := m.ParseTemplate(
		"./templates/verify-email.html",
		templateData,
	)
	if err != nil {
		logrus.WithError(err).Error(err.Error())
		return err
	}
	err = m.SendEmail(
		receiverEmailAddr,
		receiverName,
		subject,
		body,
	)
	return err
}

// Send a mail containing the reset password link
func (m *Mailer) SendResetPasswordMail(receiverEmailAddr string, receiverName string, subject string, templateData interface{}) error {
	body, err := m.ParseTemplate(
		"./templates/reset-password.html",
		templateData,
	)
	if err != nil {
		logrus.WithError(err).Error(err.Error())
		return err
	}
	err = m.SendEmail(
		receiverEmailAddr,
		receiverName,
		subject,
		body,
	)
	return err
}

func (m *Mailer) ParseTemplate(templateName string, templateData interface{}) (string, error) {
	t, err := template.New("").ParseFiles(templateName, "./templates/base.html")
	if err != nil {
		log.Print(err)
		return "", err
	}

	buf := new(bytes.Buffer)
	if err = t.ExecuteTemplate(buf, "base", templateData); err != nil {
		return "", err
	}

	return buf.String(), nil
}
