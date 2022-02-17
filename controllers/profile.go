package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/stroiman/go-automapper"
	"helloprofile.com/models"
	"helloprofile.com/persistence/orm/helloprofiledb"
	"helloprofile.com/util"
)

// GetProfiles is used get profiles
func (env *Env) GetProfiles(c echo.Context) (err error) {

	errorResponse := new(models.Errormessage)

	fields := log.Fields{"microservice": "helloprofile.service", "application": "backend", "function": "GetProfiles"}
	log.WithFields(fields).Info("Get profiles request received...")
	if c.QueryParam("email") != "" {
		user, err := env.HelloProfileDb.GetUser(context.Background(), sql.NullString{String: c.QueryParam("email"), Valid: true})
		if err != nil {
			errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
			errorResponse.ErrorMessage = util.USER_NOT_FOUND_RESPONSE_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("User was not found")
			c.JSON(http.StatusNotFound, errorResponse)
			return err
		}
		profiles := make(chan []models.Profile)
		env.getProfiles(user.ID, profiles, fields)
		if len(profiles) <= 0 {
			errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
			errorResponse.ErrorMessage = util.NO_RECORD_FOUND_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("User was not found")
			c.JSON(http.StatusNotFound, errorResponse)
			return err
		}
		if c.QueryParam("type") == "default" {
			profile := new(models.Profile)
			for _, value := range <-profiles {
				if value.IsDefault {
					profile = &value
				}
			}
			response := &models.SuccessResponse{
				ResponseCode:    util.SUCCESS_RESPONSE_CODE,
				ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
				ResponseDetails: &profile,
			}
			c.JSON(http.StatusOK, response)
		} else {
			response := &models.SuccessResponse{
				ResponseCode:    util.SUCCESS_RESPONSE_CODE,
				ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
				ResponseDetails: profiles,
			}
			c.JSON(http.StatusOK, response)

		}
		return err
	} else {
		profiles, err := env.HelloProfileDb.GetAllProfiles(context.Background())
		if err != nil {
			errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
			errorResponse.ErrorMessage = util.NO_RECORD_FOUND_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Profiles not found")

			c.JSON(http.StatusNotFound, errorResponse)
			return err
		}
		log.WithFields(fields).Info("Successfully retrieved states...")
		profilesResponse := make([]models.Profile, len(profiles))
		for index, value := range profiles {
			profile := models.Profile{
				Status:      value.Status,
				ID:          value.ID,
				ProfileName: value.ProfileName,
				IsDefault:   value.IsDefault,
				PageColor:   value.PageColor,
				Font:        value.Font,
			}
			profilesResponse[index] = profile
		}
		response := &models.SuccessResponse{
			ResponseCode:    util.SUCCESS_RESPONSE_CODE,
			ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
			ResponseDetails: profilesResponse,
		}
		c.JSON(http.StatusOK, response)
	}
	return err
}

// AddProfile is used create a new profile
func (env *Env) AddProfile(c echo.Context) (err error) {

	errorResponse := new(models.Errormessage)

	fields := log.Fields{"microservice": "helloprofile.service", "application": "backend", "function": "AddProfile"}
	log.WithFields(fields).Info("Add profile request received...")

	if c.Param("email") != "" {
		request := new(models.Profile)
		if err = c.Bind(request); err != nil {
			errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
			errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to marshal request")
			c.JSON(http.StatusBadRequest, errorResponse)
			return err
		}
		user, err := env.HelloProfileDb.GetUser(context.Background(), sql.NullString{String: c.Param("email"), Valid: true})
		if err != nil {
			errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
			errorResponse.ErrorMessage = util.USER_NOT_FOUND_RESPONSE_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("User was not found while trying to add profile")
			c.JSON(http.StatusBadRequest, errorResponse)
			return err
		}
		log.WithFields(fields).Info(fmt.Sprintf("Profile to add to user %s : %v", user.Email, request))
		dbProfile := new(helloprofiledb.AddProfileParams)
		automapper.MapLoose(request, dbProfile)
		dbProfile.UserID = user.ID
		dbProfileAddResult, err := env.HelloProfileDb.AddProfile(context.Background(), *dbProfile)
		if err != nil {
			errorResponse.Errorcode = util.SQL_ERROR_CODE
			errorResponse.ErrorMessage = util.SQL_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while adding profile for user %s", user.Email)
			c.JSON(http.StatusBadRequest, errorResponse)
			return err
		}
		log.WithFields(fields).Info("Successfully added profile")

		response := &models.SuccessResponse{
			ResponseCode:    util.SUCCESS_RESPONSE_CODE,
			ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
			ResponseDetails: dbProfileAddResult.ID,
		}
		c.JSON(http.StatusOK, response)
		return err
	} else {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Email was not passed in the url params")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
}

// UpdateProfile is used udpate a profile
func (env *Env) UpdateProfile(c echo.Context) (err error) {

	errorResponse := new(models.Errormessage)

	fields := log.Fields{"microservice": "helloprofile.service", "application": "backend", "function": "UpdateProfile"}
	log.WithFields(fields).Info("Update profile request received...")

	request := new(models.Profile)
	if err = c.Bind(request); err != nil {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to marshal request")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	dbProfile := new(helloprofiledb.UpdateProfileParams)
	automapper.MapLoose(request, dbProfile)
	err = env.HelloProfileDb.UpdateProfile(context.Background(), *dbProfile)
	if err != nil {
		errorResponse.Errorcode = util.SQL_ERROR_CODE
		errorResponse.ErrorMessage = util.SQL_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Profile update failed")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	} else {
		log.WithFields(fields).Info("Successfully updated profile")
		response := &models.SuccessResponse{
			ResponseCode:    util.SUCCESS_RESPONSE_CODE,
			ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
		}
		c.JSON(http.StatusOK, response)
		return err
	}
}
