package mailmessage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/forsam-education/hermes/storage"
	"github.com/forsam-education/simplelogger"
	"gopkg.in/gomail.v2"
	htemplate "html/template"
	"io"
	ttemplate "text/template"
)

type mailMessage struct {
	FromName        string                 `json:"from_name"`
	FromAddress     string                 `json:"from_address"`
	ToAddress       string                 `json:"to_address"`
	ReplyToAddress  string                 `json:"reply_to"`
	Template        string                 `json:"template_name"`
	Subject         string                 `json:"subject"`
	CC              []string               `json:"cc,omitempty"`
	BCC             []string               `json:"bcc,omitempty"`
	Attachements    []string               `json:"attachements,omitempty"`
	TemplateContext map[string]interface{} `json:"template_context"`
}

func buildMailContent(templateConnector storage.TemplateConnector, attachementWriter storage.AttachementWriter, mailMsg *mailMessage) (*gomail.Message, error) {
	message := gomail.NewMessage()

	htmlTemplateContent, err := templateConnector.GetTemplateContent(fmt.Sprintf("%s.html.template", mailMsg.Template))
	if err != nil {
		return nil, err
	}
	txtTemplateContent, err := templateConnector.GetTemplateContent(fmt.Sprintf("%s.txt.template", mailMsg.Template))
	if err != nil {
		return nil, err
	}
	htmlTmpl, _ := htemplate.New("htmlTemplate").Parse(htmlTemplateContent)
	txtTmpl, _ := ttemplate.New("textTemplate").Parse(txtTemplateContent)

	var htmlTmplBuffer bytes.Buffer
	err = htmlTmpl.Execute(&htmlTmplBuffer, mailMsg.TemplateContext)
	if err != nil {
		return nil, fmt.Errorf("unable to execute HTML template: %s", err.Error())
	}

	var txtTmplBuffer bytes.Buffer
	err = txtTmpl.Execute(&txtTmplBuffer, mailMsg.TemplateContext)
	if err != nil {
		return nil, fmt.Errorf("unable to execute TXT template: %s", err.Error())
	}

	ccAddresses := make([]string, len(mailMsg.CC))
	for i, ccRecipient := range mailMsg.CC {
		ccAddresses[i] = message.FormatAddress(ccRecipient, "")
	}

	bccAddresses := make([]string, len(mailMsg.BCC))
	for i, bccRecipient := range mailMsg.BCC {
		bccAddresses[i] = message.FormatAddress(bccRecipient, "")
	}

	message.SetBody("text/plain", txtTmplBuffer.String())
	message.AddAlternative("text/html", htmlTmplBuffer.String())
	message.SetAddressHeader("From", mailMsg.FromAddress, mailMsg.FromName)
	message.SetHeader("To", mailMsg.ToAddress)
	message.SetHeader("Subject", mailMsg.Subject)
	message.SetHeader("Cc", ccAddresses...)
	message.SetHeader("Bcc", bccAddresses...)
	message.SetHeader("Reply-To", mailMsg.ReplyToAddress)
	for _, att := range mailMsg.Attachements {
		message.Attach(att, gomail.SetCopyFunc(func(writer io.Writer) error {
			return attachementWriter.WriteFile(att, writer)
		}))
	}

	return message, nil
}

// SendMail builds and sends a mail through SMTP transport
func SendMail(templateConnector storage.TemplateConnector, attachementWriter storage.AttachementWriter, smtpTransport *gomail.Dialer, messageBody string) error {
	var mailMsg mailMessage

	err := json.Unmarshal([]byte(messageBody), &mailMsg)
	if err != nil {
		return fmt.Errorf("unable tu unmarshal email: %s", err.Error())
	}

	mail, err := buildMailContent(templateConnector, attachementWriter, &mailMsg)
	if err != nil {
		return err
	}

	if err := smtpTransport.DialAndSend(mail); err != nil {
		return fmt.Errorf("unable to send email through smtp: %s", err.Error())
	}

	simplelogger.GlobalLogger.Info("Sent email message", simplelogger.LogExtraData{"mail": mailMsg})

	return nil
}
