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

// AddBasicBlock is used create a new profile
func (env *Env) AddBasicBlock(c echo.Context) (err error) {

	errorResponse := new(models.Errormessage)

	fields := log.Fields{"microservice": "helloprofile.service", "application": "backend", "function": "AddBasicBlock"}
	log.WithFields(fields).Info("Add basic block request received...")

	if c.Param("profileId") != "" {
		request := new(models.Basic)
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
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Profile id passed is invalid")
			c.JSON(http.StatusBadRequest, errorResponse)
			return err
		}
		isProfileExist, err := env.HelloProfileDb.IsProfileExist(context.Background(), profileId)
		if err != nil || !isProfileExist {
			errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
			errorResponse.ErrorMessage = util.NO_RECORD_FOUND_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Profile was not found while trying to add basic block")
			c.JSON(http.StatusBadRequest, errorResponse)
			return err
		}

		log.WithFields(fields).Info(fmt.Sprintf("Basic block to add to profile %s : %v", profileId, request))
		dbBasic := new(helloprofiledb.AddBasicBlockParams)
		dbBasic.Bio = request.Bio
		dbBasic.Fullname = request.Fullname
		dbBasic.Title = request.Title
		dbBasic.CoverColour = sql.NullString{String: request.CoverColour, Valid: true}
		dbBasic.CoverPhotoUrl = sql.NullString{String: request.CoverPhotoUrl, Valid: true}
		dbBasic.ProfilePhotoUrl = sql.NullString{String: request.ProfilePhotoUrl, Valid: true}
		dbAddBasicResult, err := env.HelloProfileDb.AddBasicBlock(context.Background(), *dbBasic)
		if err != nil {
			errorResponse.Errorcode = util.SQL_ERROR_CODE
			errorResponse.ErrorMessage = util.SQL_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while adding basic block for profile %s", profileId)
			c.JSON(http.StatusBadRequest, errorResponse)
			return err
		}
		err = env.HelloProfileDb.UpdateProfileWithBasicBlockId(context.Background(), helloprofiledb.UpdateProfileWithBasicBlockIdParams{
			BasicBlockID: uuid.NullUUID{UUID: dbAddBasicResult.ID, Valid: true},
			ID:           profileId,
		})
		log.WithFields(fields).Info("Successfully added basic block for profile %s", profileId)

		response := &models.SuccessResponse{
			ResponseCode:    util.SUCCESS_RESPONSE_CODE,
			ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
			ResponseDetails: dbAddBasicResult.ID,
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

// UpdateBasicBlock is used udpate the basic block of a profile
func (env *Env) UpdateBasicBlock(c echo.Context) (err error) {

	errorResponse := new(models.Errormessage)

	fields := log.Fields{"microservice": "helloprofile.service", "application": "backend", "function": "UpdateBasicBlock"}
	log.WithFields(fields).Info("Update basic block request received...")

	request := new(models.Basic)
	if err = c.Bind(request); err != nil {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to marshal request")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	dbBasic := new(helloprofiledb.UpdateBasicBlockParams)
	dbBasic.Bio = request.Bio
	dbBasic.Fullname = request.Fullname
	dbBasic.Title = request.Title
	dbBasic.CoverColour = sql.NullString{String: request.CoverColour, Valid: true}
	dbBasic.CoverPhotoUrl = sql.NullString{String: request.CoverPhotoUrl, Valid: true}
	dbBasic.ProfilePhotoUrl = sql.NullString{String: request.ProfilePhotoUrl, Valid: true}
	err = env.HelloProfileDb.UpdateBasicBlock(context.Background(), *dbBasic)
	if err != nil {
		errorResponse.Errorcode = util.SQL_ERROR_CODE
		errorResponse.ErrorMessage = util.SQL_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Basic block update failed")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	} else {
		log.WithFields(fields).Info("Successfully updated basic block")
		response := &models.SuccessResponse{
			ResponseCode:    util.SUCCESS_RESPONSE_CODE,
			ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
		}
		c.JSON(http.StatusOK, response)
		return err
	}
}

// Deletes a basic block
func (env *Env) DeleteBasicBlock(c echo.Context) (err error) {

	errorResponse := new(models.Errormessage)

	fields := log.Fields{"microservice": "helloprofile.service", "application": "backend", "function": "DeleteBasicBlock"}
	log.WithFields(fields).Info("Delete basic block request received...")
	if c.Param("id") != "" {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
			errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Incorrect id format passed for delete basic block")
			c.JSON(http.StatusBadRequest, errorResponse)
			return err
		}
		err = env.HelloProfileDb.DeleteBasicBlock(context.Background(), id)
		if err != nil {
			errorResponse.Errorcode = util.SQL_ERROR_CODE
			errorResponse.ErrorMessage = util.SQL_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Basic block update failed")
			c.JSON(http.StatusBadRequest, errorResponse)
			return err
		} else {
			log.WithFields(fields).Info("Successfully deleted basic block")
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
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to marshal request")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
}
