package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/caarlos0/env/v6"
	"github.com/forsam-education/hermes/storageconnector"
	"golang.org/x/sync/errgroup"
	"gopkg.in/gomail.v2"
	htemplate "html/template"
	ttemplate "text/template"
)

type config struct {
	Bucket       string `env:"TEMPLATE_BUCKET"`
	SMTPHost     string `env:"SMTP_HOST"`
	SMTPPort     int    `env:"SMTP_PORT" envDefault:"465"`
	SMTPUserName string `env:"SMTP_USER"`
	SMTPPassword string `env:"SMTP_PASS"`
}

type mailMessage struct {
	FromName        string                 `json:"from_name"`
	FromAddress     string                 `json:"from_address"`
	ToAddress       string                 `json:"to_address"`
	ReplyToAddress  string                 `json:"reply_to"`
	Template        string                 `json:"template_name"`
	Subject         string                 `json:"subject"`
	CC              []string               `json:"cc,omitempty"`
	BCC             []string               `json:"bcc,omitempty"`
	TemplateContext map[string]interface{} `json:"template_context"`
}

func buildMailContent(storageConnector storageconnector.StorageConnector, mailMsg *mailMessage) (*gomail.Message, error) {
	message := gomail.NewMessage()

	htmlTemplateContent, err := storageConnector.GetTemplateContent(fmt.Sprintf("%s.html.template", mailMsg.Template))
	if err != nil {
		return nil, err
	}
	txtTemplateContent, err := storageConnector.GetTemplateContent(fmt.Sprintf("%s.txt.template", mailMsg.Template))
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

	return message, nil
}

func sendMail(storageConnector storageconnector.StorageConnector, smtpTransport *gomail.Dialer, messageBody string) error {
	var mailMsg mailMessage

	err := json.Unmarshal([]byte(messageBody), &mailMsg)
	if err != nil {
		return fmt.Errorf("unable tu unmarshal email: %s", err.Error())
	}

	mail, err := buildMailContent(storageConnector, &mailMsg)
	if err != nil {
		return err
	}

	if err := smtpTransport.DialAndSend(mail); err != nil {
		return fmt.Errorf("unable to send email through smtp: %s", err.Error())
	}

	return nil
}

// HandleRequest is the main handler function used by the lambda runtime for the incomming event.
func HandleRequest(ctx context.Context, event events.SQSEvent) error {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		return fmt.Errorf("unable to parse configuration: %s", err.Error())
	}
	smtpTransport := gomail.NewDialer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUserName, cfg.SMTPPassword)
	storageConnector, err := storageconnector.NewS3(cfg.Bucket)
	if err != nil {
		return fmt.Errorf("unable to instantiate S3 storage connector: %s", err.Error())
	}

	errs, ctx := errgroup.WithContext(ctx)

	for _, message := range event.Records {
		errs.Go(func() error {
			return sendMail(storageConnector, smtpTransport, message.Body)
		})
	}

	return errs.Wait()
}

func main() {
	lambda.Start(HandleRequest)
}
