package controllers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"helloprofile.com/models"
	"helloprofile.com/util"
)

//UploadFile uploads a file to the configured file server and returns the url to the file
func (env *Env) UploadFile(c echo.Context) (err error) {
	fields := log.Fields{"microservice": "helloprofile.service", "application": "helloProfile", "method": "UploadFile"}
	log.WithFields(fields).Info("File upload Request received")
	errorResponse := new(models.Errormessage)

	file, err := c.FormFile("file")
	if err != nil {
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to marshal request")
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	blobFile, err := file.Open()
	if err != nil {
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to read multipart file")
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	var size int64 = file.Size
	buffer := make([]byte, size)
	_, err = blobFile.Read(buffer)
	if err != nil {
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to upload file")
		errorResponse.Errorcode = util.FILE_UPLOAD_ERROR_CODE
		errorResponse.ErrorMessage = util.FILE_UPLOAD_ERROR_MESSAGE
		c.JSON(http.StatusInternalServerError, errorResponse)
		return err
	}
	err = UploadFile(fields, file.Filename, bytes.NewReader(buffer))
	if err != nil {
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to upload file")
		errorResponse.Errorcode = util.FILE_UPLOAD_ERROR_CODE
		errorResponse.ErrorMessage = util.FILE_UPLOAD_ERROR_MESSAGE
		c.JSON(http.StatusInternalServerError, errorResponse)
		return err
	}
	log.WithFields(fields).Info("Successfully uploaded file")

	response := &models.SuccessResponse{
		ResponseCode:    util.SUCCESS_RESPONSE_CODE,
		ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
		ResponseDetails: fmt.Sprintf("%s/%s", os.Getenv("DIGITAL_OCEAN_SPACES_CDN"), file.Filename),
	}
	c.JSON(http.StatusOK, response)
	return err
}

// UploadFile uploads an object
func (c *ClientUploader) UploadFile(file multipart.File, object string) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := c.Cl.Bucket(c.BucketName).Object(c.UploadPath + object).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	return nil
}
