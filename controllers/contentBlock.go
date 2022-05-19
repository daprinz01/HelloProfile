package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"helloprofile.com/models"
	"helloprofile.com/persistence/orm/helloprofiledb"
	"helloprofile.com/util"
)

// AddContentBlock is used create a new content block
func (env *Env) AddContentBlock(c echo.Context) (err error) {

	errorResponse := new(models.Errormessage)

	fields := log.Fields{"microservice": "helloprofile.service", "application": "backend", "function": "AddContentBlock"}
	log.WithFields(fields).Info("Add content request received...")

	if c.Param("profileId") != "" {
		request := new(models.Content)
		if err = c.Bind(request); err != nil {
			errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
			errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to marshal request")
			c.JSON(http.StatusBadRequest, errorResponse)
			return err
		}
		profileId, err := uuid.Parse(c.Param("profileId"))
		if err != nil {
			errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
			errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Profile id passed is invalid and not of type uuid")
			c.JSON(http.StatusBadRequest, errorResponse)
			return err
		}
		isProfileExist, err := env.HelloProfileDb.IsProfileExist(context.Background(), profileId)
		if err != nil || !isProfileExist {
			errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
			errorResponse.ErrorMessage = util.NO_RECORD_FOUND_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Profile was not found while trying to add content")
			c.JSON(http.StatusBadRequest, errorResponse)
			return err
		}

		log.WithFields(fields).Info(fmt.Sprintf("Content block to add to profile %s : %v", profileId, request))
		dbContent := new(helloprofiledb.AddProfileContentParams)
		dbContent.CallToActionID = request.CallToActionID
		dbContent.ContentID = request.ContentID
		dbContent.Description = request.Description
		dbContent.DisplayTitle = request.Title
		dbContent.Order = request.Order
		dbContent.Title = request.Title
		dbContent.Url = request.Url
		dbContent.ProfileID = profileId
		dbAddContentResult, err := env.HelloProfileDb.AddProfileContent(context.Background(), *dbContent)
		if err != nil {
			errorResponse.Errorcode = util.SQL_ERROR_CODE
			errorResponse.ErrorMessage = util.SQL_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while adding content block for profile ", profileId)
			c.JSON(http.StatusBadRequest, errorResponse)
			return err
		}
		log.WithFields(fields).Info("Successfully added content block for profile ", profileId)

		response := &models.SuccessResponse{
			ResponseCode:    util.SUCCESS_RESPONSE_CODE,
			ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
			ResponseDetails: dbAddContentResult.ID,
		}
		c.JSON(http.StatusOK, response)
		return err
	} else {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("ProfileID was not passed in the url params")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
}

// UpdateContentBlock is used udpate the content block of a profile
func (env *Env) UpdateContentBlock(c echo.Context) (err error) {

	errorResponse := new(models.Errormessage)

	fields := log.Fields{"microservice": "helloprofile.service", "application": "backend", "function": "UpdateContentBlock"}
	log.WithFields(fields).Info("Update content block request received...")

	request := new(models.Content)
	if err = c.Bind(request); err != nil {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to marshal request")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	dbContent, err := env.HelloProfileDb.GetProfileContent(context.Background(), request.ID)
	if err != nil {
		errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
		errorResponse.ErrorMessage = util.NO_RECORD_FOUND_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("content block update failed. content block not found")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	err = env.HelloProfileDb.UpdateProfileContent(context.Background(), helloprofiledb.UpdateProfileContentParams{
		Title:          env.GetValue(strings.ToLower(request.Title), dbContent.Title),
		DisplayTitle:   env.GetValue(request.Title, dbContent.Title),
		Description:    env.GetValue(request.Description, dbContent.Description),
		Url:            env.GetValue(request.Url, dbContent.Url),
		CallToActionID: env.GetUUIDValue(request.CallToActionID, dbContent.CallToActionID),
		Order:          env.GetIntValue(request.Order, dbContent.Order),
		ID:             request.ID,
	})
	if err != nil {
		errorResponse.Errorcode = util.SQL_ERROR_CODE
		errorResponse.ErrorMessage = util.SQL_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Content block update failed")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	} else {
		log.WithFields(fields).Info("Successfully updated Content block")
		response := &models.SuccessResponse{
			ResponseCode:    util.SUCCESS_RESPONSE_CODE,
			ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
		}
		c.JSON(http.StatusOK, response)
		return err
	}
}

// DeleteContentBlock deletes a content block
func (env *Env) DeleteContentBlock(c echo.Context) (err error) {

	errorResponse := new(models.Errormessage)

	fields := log.Fields{"microservice": "helloprofile.service", "application": "backend", "function": "DeleteContentBlock"}
	log.WithFields(fields).Info("Delete content block request received...")
	if c.QueryParam("id") != "" {
		id, err := uuid.Parse(c.QueryParam("id"))
		if err != nil {
			errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
			errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Incorrect id format passed for delete content block")
			c.JSON(http.StatusBadRequest, errorResponse)
			return err
		}
		err = env.HelloProfileDb.DeleteProfileContent(context.Background(), id)
		if err != nil {
			errorResponse.Errorcode = util.SQL_ERROR_CODE
			errorResponse.ErrorMessage = util.SQL_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Content block deletion failed")
			c.JSON(http.StatusBadRequest, errorResponse)
			return err
		} else {
			log.WithFields(fields).Info("Successfully deleted content block")
			response := &models.SuccessResponse{
				ResponseCode:    util.SUCCESS_RESPONSE_CODE,
				ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
			}
			c.JSON(http.StatusOK, response)
			return err
		}
	} else {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Id not passed for delete content request")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
}
