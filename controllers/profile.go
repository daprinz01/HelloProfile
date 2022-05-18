package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
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
		profilesChan := make(chan []models.Profile)
		env.getProfiles(user.ID, profilesChan, fields)
		profiles := <-profilesChan
		if len(profiles) <= 0 {
			errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
			errorResponse.ErrorMessage = util.NO_RECORD_FOUND_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Profiles was not found for user")
			c.JSON(http.StatusNotFound, errorResponse)
			return err
		}
		if c.QueryParam("type") == "default" {
			profile := new(models.Profile)
			log.WithFields(fields).Info(`All profiles %v`, profiles)
			for _, value := range profiles {
				if value.IsDefault {
					log.WithFields(fields).Info(`Default profile %v`, &value)
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
		dbProfile.BasicBlockID = uuid.NullUUID{UUID: request.Basic.ID, Valid: true}
		dbProfile.ContactBlockID = uuid.NullUUID{UUID: request.ContactBlock.ID, Valid: true}
		dbProfile.Font = request.Font
		dbProfile.IsDefault = request.IsDefault
		dbProfile.PageColor = request.PageColor
		dbProfile.ProfileName = request.ProfileName
		dbProfile.Status = request.Status
		dbProfile.UserID = user.ID
		dbProfileAddResult, err := env.HelloProfileDb.AddProfile(context.Background(), *dbProfile)
		if err != nil {
			errorResponse.Errorcode = util.SQL_ERROR_CODE
			errorResponse.ErrorMessage = util.SQL_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while adding profile for user ", user.Email)
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
	dbProfile, err := env.HelloProfileDb.GetProfile(context.Background(), request.ID)
	if err != nil {
		errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
		errorResponse.ErrorMessage = util.NO_RECORD_FOUND_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Profile update failed, profile does not exist")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	dbProfile.Font = request.Font
	dbProfile.IsDefault = request.IsDefault
	dbProfile.PageColor = request.PageColor
	dbProfile.ProfileName = request.ProfileName
	dbProfile.Status = request.Status
	dbProfile.ID = request.ID
	err = env.HelloProfileDb.UpdateProfile(context.Background(), helloprofiledb.UpdateProfileParams{
		UserID:         dbProfile.UserID,
		Status:         request.Status,
		ProfileName:    request.ProfileName,
		BasicBlockID:   dbProfile.BasicBlockID,
		ContactBlockID: dbProfile.ContactBlockID,
		PageColor:      request.PageColor,
		Font:           request.Font,
		IsDefault:      request.IsDefault,
		ID:             dbProfile.ID,
	})
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

// DeleteProfile deletes a profile
func (env *Env) DeleteProfile(c echo.Context) (err error) {

	errorResponse := new(models.Errormessage)

	fields := log.Fields{"microservice": "helloprofile.service", "application": "backend", "function": "DeleteSocialBlock"}
	log.WithFields(fields).Info("Delete social block request received...")
	if c.QueryParam("id") != "" {
		id, err := uuid.Parse(c.QueryParam("id"))
		if err != nil {
			errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
			errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Incorrect id format passed for delete profile")
			c.JSON(http.StatusBadRequest, errorResponse)
			return err
		}
		err = env.HelloProfileDb.DeleteProfile(context.Background(), id)
		if err != nil {
			errorResponse.Errorcode = util.SQL_ERROR_CODE
			errorResponse.ErrorMessage = util.SQL_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Profile deletion failed")
			c.JSON(http.StatusBadRequest, errorResponse)
			return err
		} else {
			log.WithFields(fields).Info("Successfully deleted profile")
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
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Id not passed for delete profile request")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
}
