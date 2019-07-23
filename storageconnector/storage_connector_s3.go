package storageconnector

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/forsam-education/simplelogger"
)

// The S3 Storage Connector handles getting template content from AWS S3 buckets.
type S3 struct {
	bucket   string
	s3Client *s3.S3
}

// GetTemplateContent fetches template content by it's name from the S3 Bucket.
func (storageConnector S3) GetTemplateContent(templateName string) (string, error) {
	templateS3Object, err := storageConnector.s3Client.GetObject(&s3.GetObjectInput{Bucket: aws.String(storageConnector.bucket), Key: &templateName})
	if err != nil {
		return "", fmt.Errorf("unable to get item in bucket %q, reason %s", storageConnector.bucket, err.Error())
	}
	buf := new(bytes.Buffer)

	simplelogger.GlobalLogger.Info(fmt.Sprintf("Downloaded template %s from S3 storage", templateName), nil)

	_, err = buf.ReadFrom(templateS3Object.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read from object %v, reason %s", templateS3Object, err.Error())
	}

	return buf.String(), nil
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
