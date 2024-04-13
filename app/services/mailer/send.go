// Copyright 2024 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package mailer

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"time"

	"golang.org/x/text/language"
	"megpoid.dev/go/contact-form/app/model"

	mail "github.com/xhit/go-simple-mail/v2"
	"golang.org/x/text/message"
	"megpoid.dev/go/contact-form/config"
	"megpoid.dev/go/contact-form/templates"
)

type Config struct {
	SmtpSettings    config.SMTPSettings
	GeneralSettings config.GeneralSettings
}

type Mailer struct {
	registryTemplate *template.Template
	clientTemplate   *template.Template
	smtpServer       *mail.SMTPServer
	emailFrom        string
	emailTo          []string
	replyTo          string
	appName          string
	subjectStaff     string
	subjectClient    string
}

func NewMailer(cfg Config) *Mailer {
	m := &Mailer{
		appName:   cfg.GeneralSettings.SenderName,
		emailFrom: cfg.SmtpSettings.EmailFrom,
		emailTo:   cfg.GeneralSettings.EmailTo,
		replyTo:   cfg.GeneralSettings.ReplyTo,
	}
	f := openTemplate("registry.tmpl.html", cfg.GeneralSettings.TemplatesPath, cfg.GeneralSettings.DefaultLanguage)
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	registryTmpl := template.Must(template.New("registry").Parse(string(data)))

	f = openTemplate("client.tmpl.html", cfg.GeneralSettings.TemplatesPath, cfg.GeneralSettings.DefaultLanguage)
	defer f.Close()

	data, err = io.ReadAll(f)
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

	var tag language.Tag
	switch cfg.GeneralSettings.DefaultLanguage {
	case "es":
		tag = language.Spanish
	case "en":
		tag = language.English
	default:
		tag = language.English
	}

	t := message.NewPrinter(tag)

	m.subjectStaff = t.Sprintf("[%s] - New contact", m.appName)
	m.subjectClient = t.Sprintf("Thanks for contacting us")

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

func openTemplate(name, templateDir, lang string) io.ReadCloser {
	if templateDir != "" {
		externalFile, err := os.Open(path.Join(templateDir, name))
		if err == nil {
			return externalFile
		}
		log.Printf("Cannot find external template %s, using built-in", name)
	}

	fsDir, err := fs.Sub(templates.Assets(), "email/"+lang)
	if err != nil {
		panic(err)
	}
	internalFile, err := fsDir.Open(name)
	if err != nil {
		panic(err)
	}
	return internalFile
}

func (m *Mailer) Send(contact *model.Contact) error {
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

	msg.SetSubject(m.subjectStaff)

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

	// send client email
	msg = mail.NewMSG()
	msg.SetFrom(m.emailFrom)
	msg.AddTo(contact.Email)
	msg.SetSubject(m.subjectClient)

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
