package controllers

import (
	"cloud.google.com/go/storage"
	"helloprofile.com/persistence/orm/helloprofiledb"
)

// Env is used to declare public variable accessible to the controllers
type Env struct {
	HelloProfileDb *helloprofiledb.Queries
	Uploader       *ClientUploader
}

type ClientUploader struct {
	Cl         *storage.Client
	ProjectID  string
	BucketName string
	UploadPath string
}
