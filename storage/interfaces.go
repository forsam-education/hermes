package storage

import "io"

// TemplateFetcher interface should be implemented by any service responsible to get template content from a storage manager (FS, S3 TemplateBucket, Redis... etc).
type TemplateFetcher interface {
	// Fetch should return the content of the template as string
	Fetch(templateName string) (string, error)
}

// AttachementCopier interface should be implemented by any service responsible to get attachement files from a storage manager (FS, S3 TemplateBucket, Redis... etc).
type AttachementCopier interface {
	// Copy should, as expected, copy the attachement file to the provided io.Writer.
	Copy(attachementPath string, writer io.Writer) error
}
