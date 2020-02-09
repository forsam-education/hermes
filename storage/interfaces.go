package storage

import "io"

// TemplateFetcher interface should be implemented by any service responsible to get template content from a storage manager (FS, S3 TemplateBucket, Redis... etc).
type TemplateFetcher interface {
	// Fetch should return the content of the template as string
	Fetch(templateName string) (string, error)
}

// AttachmentCopier interface should be implemented by any service responsible to get attachment files from a storage manager (FS, S3 TemplateBucket, Redis... etc).
type AttachmentCopier interface {
	// Copy should, as expected, copy the attachment file to the provided io.Writer.
	Copy(attachmentPath string, writer io.Writer) error
}
