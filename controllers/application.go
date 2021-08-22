package controllers

import (
	"authengine/models"
	"authengine/persistence/orm/authdb"
	"authengine/util"
	"context"
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
)

// GetApplications is used get countries
func (env *Env) GetApplications(c echo.Context) (err error) {

	errorResponse := new(models.Errormessage)
	application := c.Param("application")
	if application == "" {
		errorResponse.Errorcode = util.APPLICATION_NOT_SPECIFIED_ERROR_CODE
		errorResponse.ErrorMessage = util.APPLICATION_NOT_SPECIFIED_ERROR_MESSAGE
		log.WithField("microservice", "persian.black.authengine.service").Error(util.APPLICATION_NOT_SPECIFIED_ERROR_MESSAGE)
		c.JSON(http.StatusBadRequest, errorResponse)
		return nil
	}
	fields := log.Fields{"microservice": "persian.black.authengine.service", "application": application}
	log.WithFields(fields).Info("Get applications request received...")
	applications, err := env.AuthDb.GetApplications(context.Background())
	if err != nil {
		errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
		errorResponse.ErrorMessage = util.INVALID_APPLICATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error(util.INVALID_APPLICATION_ERROR_MESSAGE)

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}

	log.WithFields(fields).Info("Successfully retrieved application...")
	applicationResponse := make([]models.Application, len(applications))
	for index, value := range applications {
		application := models.Application{
			Application: value.Name,
			Description: value.Description,
		}
		applicationResponse[index] = application
	}
	response := &models.SuccessResponse{
		ResponseCode:    util.SUCCESS_RESPONSE_CODE,
		ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
		ResponseDetails: applicationResponse,
	}
	c.JSON(http.StatusOK, response)
	return err
}

// GetApplication is used get application
func (env *Env) GetApplication(c echo.Context) (err error) {

	errorResponse := new(models.Errormessage)

	application := c.Param("application")
	if application == "" {
		errorResponse.Errorcode = util.APPLICATION_NOT_SPECIFIED_ERROR_CODE
		errorResponse.ErrorMessage = util.APPLICATION_NOT_SPECIFIED_ERROR_MESSAGE
		log.WithField("microservice", "persian.black.authengine.service").Error(util.APPLICATION_NOT_SPECIFIED_ERROR_MESSAGE)
		c.JSON(http.StatusBadRequest, errorResponse)
		return nil
	}
	fields := log.Fields{"microservice": "persian.black.authengine.service", "application": application}
	log.WithFields(fields).Info("Get application request received...")
	dbApplication, err := env.AuthDb.GetApplication(context.Background(), strings.ToLower(application))
	if err != nil {
		errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
		errorResponse.ErrorMessage = util.INVALID_APPLICATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).Info(util.INVALID_APPLICATION_ERROR_MESSAGE)
		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.WithFields(fields).Info(fmt.Sprintf("Successfully retrieved application: %v", dbApplication))
	applicationResponse := models.Application{
		Application: dbApplication.Name,
		Description: dbApplication.Description,
	}

	response := &models.SuccessResponse{
		ResponseCode:    util.SUCCESS_RESPONSE_CODE,
		ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
		ResponseDetails: applicationResponse,
	}
	c.JSON(http.StatusOK, response)
	return err
}

// AddApplication is used add application
func (env *Env) AddApplication(c echo.Context) (err error) {

	errorResponse := new(models.Errormessage)
	application := c.Param("application")
	if application == "" {
		errorResponse.Errorcode = util.APPLICATION_NOT_SPECIFIED_ERROR_CODE
		errorResponse.ErrorMessage = util.APPLICATION_NOT_SPECIFIED_ERROR_MESSAGE
		log.WithField("microservice", "persian.black.authengine.service").Error("Calling application not specified")
		c.JSON(http.StatusBadRequest, errorResponse)
		return nil
	}
	fields := log.Fields{"microservice": "persian.black.authengine.service", "application": application}
	log.WithFields(fields).Info("Add application request received...")
	request := new(models.Application)
	if err = c.Bind(request); err != nil {

		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to marshal request")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	dbApplication, err := env.AuthDb.CreateApplication(context.Background(), authdb.CreateApplicationParams{
		Name:        strings.ToLower(request.Application),
		Description: strings.ToLower(request.Description),
	})
	if err != nil {
		errorResponse.Errorcode = util.DUPLICATE_RECORD_ERROR_CODE
		errorResponse.ErrorMessage = util.DUPLICATE_RECORD_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured adding new application")

		c.JSON(http.StatusNotModified, errorResponse)
		return err
	}
	log.WithFields(fields).Info(fmt.Sprintf("Successfully added application: %v", dbApplication))

	response := &models.SuccessResponse{
		ResponseCode:    util.SUCCESS_RESPONSE_CODE,
		ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
	}
	c.JSON(http.StatusOK, response)
	return err
}

// UpdateApplication is used add application
func (env *Env) UpdateApplication(c echo.Context) (err error) {

	errorResponse := new(models.Errormessage)

	application := c.Param("application")
	if application == "" {
		errorResponse.Errorcode = util.APPLICATION_NOT_SPECIFIED_ERROR_CODE
		errorResponse.ErrorMessage = util.APPLICATION_NOT_SPECIFIED_ERROR_MESSAGE
		log.WithField("microservice", "persian.black.authengine.service").Error("Calling application not specified")
		c.JSON(http.StatusBadRequest, errorResponse)
		return nil
	}
	fields := log.Fields{"microservice": "persian.black.authengine.service", "application": application}
	log.WithFields(fields).Info("Update application request received...")
	request := new(models.Application)
	if err = c.Bind(request); err != nil {

		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to marshal request")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	_, err = env.AuthDb.UpdateApplication(context.Background(), authdb.UpdateApplicationParams{
		Name:        strings.ToLower(request.Application),
		Description: strings.ToLower(request.Description),
		Name_2:      strings.ToLower(application),
	})
	if err != nil {
		errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
		errorResponse.ErrorMessage = util.NO_RECORD_FOUND_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured updating new application")

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.WithFields(fields).Info("Successfully updated application")

	response := &models.SuccessResponse{
		ResponseCode:    util.SUCCESS_RESPONSE_CODE,
		ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
	}
	c.JSON(http.StatusOK, response)
	return err
}

// DeleteApplication is used add application
func (env *Env) DeleteApplication(c echo.Context) (err error) {

	errorResponse := new(models.Errormessage)

	application := c.Param("application")
	if application == "" {
		errorResponse.Errorcode = util.APPLICATION_NOT_SPECIFIED_ERROR_CODE
		errorResponse.ErrorMessage = util.APPLICATION_NOT_SPECIFIED_ERROR_MESSAGE
		log.WithField("microservice", "persian.black.authengine.service").Error("Calling application not specified")
		c.JSON(http.StatusBadRequest, errorResponse)
		return nil
	}
	fields := log.Fields{"microservice": "persian.black.authengine.service", "application": application}
	log.WithFields(fields).Info("Delete application request received...")
	err = env.AuthDb.DeleteApplication(context.Background(), strings.ToLower(application))
	if err != nil {
		errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
		errorResponse.ErrorMessage = util.NO_RECORD_FOUND_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured deleting  application")

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.WithFields(fields).Info("Successfully deleted application")

	response := &models.SuccessResponse{
		ResponseCode:    util.SUCCESS_RESPONSE_CODE,
		ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
	}
	c.JSON(http.StatusOK, response)
	return err
}
