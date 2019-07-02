package storageconnector

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
)

// The S3 Storage Connector handles getting template content from AWS S3 buckets.
type S3 struct {
	bucket   string
	s3Client *s3.S3
}

// GetTemplateContent fetches template content by it's name from the S3 Bucket.
func (storageConnector S3) GetTemplateContent(templateName string) string {
	templateS3Object, err := storageConnector.s3Client.GetObject(&s3.GetObjectInput{Bucket: aws.String(storageConnector.bucket), Key: &templateName})
	if err != nil {
		exitErrorf("Unable to get item in bucket %q, %v", storageConnector.bucket, err)
	}
	buf := new(bytes.Buffer)

	_, err = buf.ReadFrom(templateS3Object.Body)
	if err != nil {
		exitErrorf("Unable to read from object %v, %v", templateS3Object, err)
	}

	return buf.String()
}

func exitErrorf(msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

// NewS3 instanciates an S3 with the AWS Session and AWS S3 Client
func NewS3(bucket string) *S3 {
	p := new(S3)
	p.bucket = bucket
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String("eu-west-1"),
		},
	)

	if err != nil {
		exitErrorf("Unable to connect to AWS: $s", err)
	}

	p.s3Client = s3.New(sess)

	return p
}
