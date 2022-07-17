package controllers

import (
	"context"
	"path/filepath"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/google/uuid"
	"google.golang.org/api/option"
	"helloprofile.com/persistence/orm/helloprofiledb"
)

// Env is used to declare public variable accessible to the controllers
type Env struct {
	HelloProfileDb *helloprofiledb.Queries
	Uploader       *ClientUploader
	FirebaseClient *auth.Client
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

func SetupFirebase() *auth.Client {
	serviceAccountKeyFilePath, err := filepath.Abs("./serviceAccountKey.json")
	if err != nil {
		panic("Unable to load serviceAccountKeys.json file")
	}
	opt := option.WithCredentialsFile(serviceAccountKeyFilePath)
	//Firebase admin SDK initialization
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic("Firebase load error")
	}
	//Firebase Auth
	auth, err := app.Auth(context.Background())
	if err != nil {
		panic("Firebase load error")
	}
	return auth
}
