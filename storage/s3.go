package storage

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/forsam-education/simplelogger"
	"io"
)

// S3 handles getting template content from AWS S3 buckets.
type S3 struct {
	bucket   string
	s3Client *s3.S3
}

// GetTemplateContent fetches template content by it's name from the S3 TemplateBucket.
func (s3Connector *S3) GetTemplateContent(templateName string) (string, error) {
	templateS3Object, err := s3Connector.s3Client.GetObject(&s3.GetObjectInput{Bucket: aws.String(s3Connector.bucket), Key: &templateName})
	if err != nil {
		return "", fmt.Errorf("unable to get item in bucket %q: %s", s3Connector.bucket, err.Error())
	}
	buf := new(bytes.Buffer)

	simplelogger.GlobalLogger.Info(fmt.Sprintf("Downloaded template %s from S3 storage", templateName), nil)

	_, err = buf.ReadFrom(templateS3Object.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read from object %v: %s", templateS3Object, err.Error())
	}

	return buf.String(), nil
}

// Write fetches attachement content by it's name from the S3 bucket and writes to attach it to an email.
func (s3Connector *S3) Write(attachementPath string, writer io.Writer) error {
	attachementS3Object, err := s3Connector.s3Client.GetObject(&s3.GetObjectInput{Bucket: aws.String(s3Connector.bucket), Key: &attachementPath})
	if err != nil {
		return fmt.Errorf("unable to get item in bucket %q: %s", s3Connector.bucket, err.Error())
	}
	simplelogger.GlobalLogger.Info(fmt.Sprintf("Downloaded attachement %s from S3 storage", attachementPath), nil)

	_, err = io.Copy(writer, attachementS3Object.Body)
	if err != nil {
		return fmt.Errorf("unable to read from object %v: %s", attachementS3Object, err.Error())
	}

	return nil
}

// NewS3 instanciates an S3 with the AWS Session and AWS S3 Client
func NewS3(bucket string, region string) (*S3, error) {
	p := new(S3)
	p.bucket = bucket
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(region),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to AWS: %s", err.Error())
	}

	p.s3Client = s3.New(sess)

	simplelogger.GlobalLogger.Info("Connected to S3 storage", simplelogger.LogExtraData{"bucket": bucket})

	return p, nil
}
