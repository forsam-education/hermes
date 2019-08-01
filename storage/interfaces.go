package storage

import "io"

// TemplateConnector interface should be implemented by any service responsible to get template content from a storage manager (FS, S3 TemplateBucket, Redis... etc).
type TemplateConnector interface {
	// GetTemplateContent should return the content of the template as string
	GetTemplateContent(templateName string) (string, error)
}

// AttachementWriter interface should be implemented by any service responsible to get attachement files from a storage manager (FS, S3 TemplateBucket, Redis... etc).
type AttachementWriter interface {
	// Write should, as expected, write the attachement file to the provided io.Writer.
	Write(attachementPath string, writer io.Writer) error
}
