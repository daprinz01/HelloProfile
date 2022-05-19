package controllers

import (
	"cloud.google.com/go/storage"
	"github.com/google/uuid"
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

func (env *Env) GetValue(request, user string) string {
	if request == "" {
		return user
	}
	return request
}

func (env *Env) GetIntValue(request, user int32) int32 {
	if request == 0 {
		return user
	}
	return request
}
func (env *Env) GetUUIDValue(request, user uuid.UUID) uuid.UUID {
	if request == uuid.Nil {
		return user
	}
	return request
}

func (env *Env) GetBoolValue(request, user bool) bool {
	if !request {
		return user
	}
	return request
}
