package controllers

import (
	"context"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"helloprofile.com/models"
	"helloprofile.com/persistence/orm/helloprofiledb"
	"helloprofile.com/util"
)

// Save profile is used to save a users profile without being authenticated
func (env *Env) SaveProfile(c echo.Context) (err error) {
	fields := log.Fields{"microservice": "helloprofile.service", "application": c.Param("application"), "method": "SaveProfile"}
	log.WithFields(fields).Info("Save profile Request received")
	errorResponse := new(models.Errormessage)

	request := new(models.SaveProfileRequest)
	if err = c.Bind(request); err != nil {
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to marshal request")
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	if !request.IsTermsAgreed {
		errorResponse.Errorcode = util.TERMS_NOT_AGREED_ERROR_CODE
		errorResponse.ErrorMessage = util.TERMS_NOT_AGREED_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Saving profile failed, terms was not agreed")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	isProfileExist, err := env.HelloProfileDb.IsProfileExist(context.Background(), request.ProfileId)
	if err != nil {
		errorResponse.Errorcode = util.SQL_ERROR_CODE
		errorResponse.ErrorMessage = util.SQL_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while checking if profile exists")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	if isProfileExist {
		errorResponse.Errorcode = util.PROFILE_NOT_FOUND_ERROR_CODE
		errorResponse.ErrorMessage = util.PROFILE_NOT_FOUND_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Profile was not found")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	savedProfile, err := env.HelloProfileDb.CreateSavedProfile(context.Background(), helloprofiledb.CreateSavedProfileParams{
		FirstName: request.Firstname,
		LastName:  request.Lastname,
		Email:     strings.ToLower(request.Email),
		ProfileID: request.ProfileId,
	})
	if err != nil {
		errorResponse.Errorcode = util.SQL_ERROR_CODE
		errorResponse.ErrorMessage = util.SQL_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error saving profile")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	log.WithFields(fields).Info("Successfully saved profile")
	resetResponse := &models.SuccessResponse{
		ResponseCode:    util.SUCCESS_RESPONSE_CODE,
		ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
		ResponseDetails: savedProfile,
	}
	c.JSON(http.StatusOK, resetResponse)
	return err
}
