package storage

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"log"
)

// S3 handles getting template content from AWS S3 buckets. It implements both AttachmentCopier and TemplateFetcher interfaces.
type S3 struct {
	bucket   string
	s3Client *s3.S3
}

// Fetch the template content by it's name from the S3 TemplateBucket and returns content.
func (s3Connector *S3) Fetch(templateName string) (string, error) {
	templateS3Object, err := s3Connector.s3Client.GetObject(&s3.GetObjectInput{Bucket: aws.String(s3Connector.bucket), Key: &templateName})
	if err != nil {
		return "", fmt.Errorf("unable to get item in bucket %q: %s", s3Connector.bucket, err.Error())
	}
	buf := new(bytes.Buffer)

	log.Printf("Downloaded template %s from S3 storage", templateName)

	_, err = buf.ReadFrom(templateS3Object.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read from object %v: %s", templateS3Object, err.Error())
	}

	return buf.String(), nil
}

// Copy fetches attachment content by it's name from the S3 bucket and copies it to attach it to an email.
func (s3Connector *S3) Copy(attachmentPath string, writer io.Writer) error {
	attachmentS3Object, err := s3Connector.s3Client.GetObject(&s3.GetObjectInput{Bucket: aws.String(s3Connector.bucket), Key: &attachmentPath})
	if err != nil {
		return fmt.Errorf("unable to get item in bucket %q: %s", s3Connector.bucket, err.Error())
	}
	log.Printf("Downloaded attachment %s from S3 storage", attachmentPath)

	_, err = io.Copy(writer, attachmentS3Object.Body)
	if err != nil {
		return fmt.Errorf("unable to read from object %v: %s", attachmentS3Object, err.Error())
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

	fmt.Printf("Connected to S3 storage bucket %s", bucket)

	return p, nil
}
