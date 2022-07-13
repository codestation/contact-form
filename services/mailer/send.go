package mailer

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
	"golang.org/x/text/message"
	"megpoid.dev/go/contact-form/app/i18n"
	"megpoid.dev/go/contact-form/config"
	"megpoid.dev/go/contact-form/model"
	"megpoid.dev/go/contact-form/templates"
)

type Mailer struct {
	registryTemplate *template.Template
	clientTemplate   *template.Template
	smtpServer       *mail.SMTPServer
	emailFrom        string
	emailTo          []string
	replyTo          string
	appName          string
}

func NewMailer(cfg *config.Config) *Mailer {
	m := &Mailer{
		appName:   cfg.GeneralSettings.AppName,
		emailFrom: cfg.SmtpSettings.EmailFrom,
		emailTo:   cfg.GeneralSettings.EmailTo,
		replyTo:   cfg.GeneralSettings.ReplyTo,
	}
	f := openTemplate("registry.tmpl.html", cfg.GeneralSettings.TemplatesPath)
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	registryTmpl := template.Must(template.New("registry").Parse(string(data)))

	f = openTemplate("client.tmpl.html", cfg.GeneralSettings.TemplatesPath)
	defer f.Close()

	data, err = ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	clientTmpl := template.Must(template.New("registry").Parse(string(data)))

	m.clientTemplate = clientTmpl
	m.registryTemplate = registryTmpl

	server := mail.NewSMTPClient()
	server.Host = cfg.SmtpSettings.SMTPHost
	server.Port = cfg.SmtpSettings.SMTPPort
	server.Username = cfg.SmtpSettings.SMTPUsername
	server.Password = cfg.SmtpSettings.SMTPPassword

	switch cfg.SmtpSettings.SMTPEncryption {
	case "tls":
		server.Encryption = mail.EncryptionSSLTLS
	case "none":
		server.Encryption = mail.EncryptionNone
	case "starttls":
		server.Encryption = mail.EncryptionSTARTTLS
	}

	switch cfg.SmtpSettings.SMTPAuth {
	case "login":
		server.Authentication = mail.AuthLogin
	case "plain":
		server.Authentication = mail.AuthPlain
	case "crammd5":
		server.Authentication = mail.AuthCRAMMD5
	case "none":
		server.Authentication = mail.AuthNone
	}

	server.ConnectTimeout = 30 * time.Second
	server.SendTimeout = 30 * time.Second
	server.Authentication = mail.AuthPlain
	server.KeepAlive = true

	if cfg.SmtpSettings.SMTPSkipVerify {
		server.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	m.smtpServer = server

	return m
}

type templateData struct {
	AppName   string
	FirstName string
	LastName  string
	Email     string
	Phone     string
	Company   string
	Subject   string
	Message   string
}

const templateLang = "es"

func openTemplate(name, templateDir string) io.ReadCloser {
	if templateDir != "" {
		externalFile, err := os.Open(path.Join(templateDir, name))
		if err == nil {
			return externalFile
		}
		log.Printf("Cannot find external template %s, using built-in", name)
	}

	fsDir, err := fs.Sub(templates.Assets(), "email/"+templateLang)
	if err != nil {
		panic(err)
	}
	internalFile, err := fsDir.Open(name)
	if err != nil {
		panic(err)
	}
	return internalFile
}

func (m *Mailer) Send(ctx context.Context, contact *model.Contact) error {
	if m.emailFrom == "" {
		return errors.New("no email-from configured")
	}

	data := templateData{
		AppName:   m.appName,
		FirstName: contact.FirstName,
		LastName:  contact.LastName,
		Email:     contact.Email,
		Phone:     contact.Phone,
		Company:   contact.Company,
		Subject:   contact.Subject,
		Message:   contact.Message,
	}

	client, err := m.smtpServer.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to smtp server: %w", err)
	}

	defer func(client *mail.SMTPClient) {
		err := client.Close()
		if err != nil {
			log.Printf("Failed to clone SMTP connection: %s", err.Error())
		}
	}(client)

	t := message.NewPrinter(i18n.GetLanguageTagsContext(ctx))

	subject := t.Sprintf("[%s] - New contact", m.appName)

	// send registry email to staff
	msg := mail.NewMSG()
	msg.SetFrom(m.emailFrom)
	msg.AddTo(m.emailTo[0])
	if len(m.emailTo) > 1 {
		msg.AddCc(m.emailTo[1:]...)
	}
	if m.replyTo != "" {
		msg.SetReplyTo(m.replyTo)
	}

	msg.SetSubject(subject)

	var registryDoc bytes.Buffer
	err = m.registryTemplate.Execute(&registryDoc, data)
	if err != nil {
		return fmt.Errorf("failed to process registry template: %w", err)
	}

	msg.SetBody(mail.TextHTML, registryDoc.String())

	err = msg.Send(client)
	if err != nil {
		return fmt.Errorf("failed to send email to staff: %w", err)
	}

	subject = t.Sprintf("Thanks for contacting us")

	// send client email
	msg = mail.NewMSG()
	msg.SetFrom(m.emailFrom)
	msg.AddTo(contact.Email)
	msg.SetSubject(subject)

	var clientDoc bytes.Buffer
	err = m.clientTemplate.Execute(&clientDoc, data)
	if err != nil {
		return fmt.Errorf("failed to process client template: %w", err)
	}

	msg.SetBody(mail.TextHTML, clientDoc.String())

	err = msg.Send(client)
	if err != nil {
		return fmt.Errorf("failed to send email to client: %w", err)
	}

	return nil
}
