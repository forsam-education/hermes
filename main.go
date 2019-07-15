package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/caarlos0/env/v6"
	"github.com/forsam-education/hermes/mailmessage"
	"github.com/forsam-education/hermes/storageconnector"
	"golang.org/x/sync/errgroup"
	"gopkg.in/gomail.v2"
)

type config struct {
	Bucket       string `env:"TEMPLATE_BUCKET"`
	SMTPHost     string `env:"SMTP_HOST"`
	SMTPPort     int    `env:"SMTP_PORT" envDefault:"465"`
	SMTPUserName string `env:"SMTP_USER"`
	SMTPPassword string `env:"SMTP_PASS"`
}

func processMails(ctx context.Context, connector storageconnector.StorageConnector, dialer *gomail.Dialer, messages []events.SQSMessage) error {
	errs, _ := errgroup.WithContext(ctx)

	for _, message := range messages {
		errs.Go(func() error {
			return mailmessage.SendMail(connector, dialer, message.Body)
		})
	}

	return errs.Wait()
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

	return processMails(ctx, storageConnector, smtpTransport, event.Records)
}

func main() {
	lambda.Start(HandleRequest)
}
