package storageconnector

// StorageConnector interface should be implemented by any service responsible to get template content from a storage manager (FS, Bucket, Redis... etc).
type StorageConnector interface {
	GetTemplateContent(templateName string) string
}
