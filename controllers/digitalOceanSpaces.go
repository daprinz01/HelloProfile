package controllers

import (
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	log "github.com/sirupsen/logrus"
)

func ConfigureClient(fields log.Fields) (*s3.S3, error) {
	key := os.Getenv("SPACES_KEY")
	secret := os.Getenv("SPACES_SECRET")

	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:    aws.String(os.Getenv("DIGITAL_OCEAN_SPACES_ENDPOINT")),
		Region:      aws.String("us-east-1"),
	}

	newSession, err := session.NewSession(s3Config)
	if err != nil {
		log.WithFields(fields).WithError(err).Error("Error occured while initializing ditigal ocean spaces")
		return nil, err
	} else {
		s3Client := s3.New(newSession)
		return s3Client, nil
	}
}

func UploadFile(fields log.Fields, fileName string, file io.ReadSeeker) error {
	object := s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("DIGITAL_OCEAN_SPACE_NAME")),
		Key:    aws.String(fileName),
		Body:   file,
		ACL:    aws.String("public-read"),
		Metadata: map[string]*string{
			"x-amz-meta-my-key": aws.String(os.Getenv("SPACES_KEY")),
		},
	}
	s3Client, err := ConfigureClient(fields)
	if err != nil {
		return err
	}
	_, err = s3Client.PutObject(&object)
	if err != nil {
		log.WithFields(fields).WithError(err).Error("Error occured while uploading file")
		return err
	}
	return nil
}
