package controllers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/stroiman/go-automapper"
	"helloprofile.com/models"
	"helloprofile.com/util"
)

func (env *Env) GetCallToActions(c echo.Context) (err error) {

	errorResponse := new(models.Errormessage)

	fields := log.Fields{"microservice": "helloprofile.service", "application": "backend", "function": "GetCallToActions"}
	log.WithFields(fields).Info("Get call to action request received...")
	dbCallToAction, err := env.HelloProfileDb.GetCallToActions(context.Background())
	if err != nil {
		errorResponse.Errorcode = util.SQL_ERROR_CODE
		errorResponse.ErrorMessage = util.SQL_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while adding fetch call to action")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	callToActions := make([]models.CallToAction, len(dbCallToAction))
	automapper.MapLoose(dbCallToAction, callToActions)
	log.WithFields(fields).Info("Successfully fetched all call to actions")

	response := &models.SuccessResponse{
		ResponseCode:    util.SUCCESS_RESPONSE_CODE,
		ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
		ResponseDetails: callToActions,
	}
	c.JSON(http.StatusOK, response)
	return err
}

func (env *Env) GetContents(c echo.Context) (err error) {

	errorResponse := new(models.Errormessage)

	fields := log.Fields{"microservice": "helloprofile.service", "application": "backend", "function": "GetContents"}
	log.WithFields(fields).Info("Get contents request received...")
	dbContentTypes, err := env.HelloProfileDb.GetAllContentTypes(context.Background())
	if err != nil {
		errorResponse.Errorcode = util.SQL_ERROR_CODE
		errorResponse.ErrorMessage = util.SQL_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while adding fetch content types")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	callToActions := make([]models.ContentType, len(dbContentTypes))
	automapper.MapLoose(dbContentTypes, callToActions)
	log.WithFields(fields).Info("Successfully fetched all content types")

	response := &models.SuccessResponse{
		ResponseCode:    util.SUCCESS_RESPONSE_CODE,
		ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
		ResponseDetails: callToActions,
	}
	c.JSON(http.StatusOK, response)
	return err
}
