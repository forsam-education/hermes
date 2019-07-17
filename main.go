package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/caarlos0/env/v6"
	"github.com/forsam-education/hermes/mailmessage"
	"github.com/forsam-education/hermes/storageconnector"
	"github.com/forsam-education/redriver"
	"gopkg.in/gomail.v2"
)

type config struct {
	Bucket       string `env:"TEMPLATE_BUCKET"`
	SMTPHost     string `env:"SMTP_HOST"`
	SMTPPort     int    `env:"SMTP_PORT" envDefault:"465"`
	SMTPUserName string `env:"SMTP_USER"`
	SMTPPassword string `env:"SMTP_PASS"`
	AWSRegion    string `env:"AWS_REGION_CODE"`
	QueueURL     string `env:"SQS_QUEUE"`
}

// HandleRequest is the main handler function used by the lambda runtime for the incomming event.
func HandleRequest(_ context.Context, event events.SQSEvent) error {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		return fmt.Errorf("unable to parse configuration: %s", err.Error())
	}
	smtpTransport := gomail.NewDialer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUserName, cfg.SMTPPassword)
	storageConnector, err := storageconnector.NewS3(cfg.Bucket, cfg.AWSRegion)
	if err != nil {
		return fmt.Errorf("unable to instantiate S3 storage connector: %s", err.Error())
	}

	messageRedriver := redriver.Redriver{Retries: 3, ConsumedQueueURL: cfg.QueueURL}

	err = messageRedriver.HandleMessages(event.Records, func(event events.SQSMessage) error {
		return mailmessage.SendMail(storageConnector, smtpTransport, event.Body)
	})

	return err
}

func main() {
	lambda.Start(HandleRequest)
}
