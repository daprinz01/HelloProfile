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
		profile, err := env.HelloProfileDb.GetProfile(context.Background(), profileId)
		if err != nil {
			errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
			errorResponse.ErrorMessage = util.NO_RECORD_FOUND_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Profile was not found while trying to add basic block")
			c.JSON(http.StatusBadRequest, errorResponse)
			return err
		}
		if profile.BasicBlockID.UUID != uuid.Nil {
			errorResponse.Errorcode = util.BASIC_BLOCK_EXIST_ERROR_CODE
			errorResponse.ErrorMessage = util.BASIC_BLOCK_EXIST_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Basic block already exists for profile")
			c.JSON(http.StatusBadRequest, errorResponse)
			return err
		}
		log.WithFields(fields).Info(fmt.Sprintf("Basic block to add to profile %s : %v", profileId, request))

		dbAddBasicResult, err := env.HelloProfileDb.AddBasicBlock(context.Background(), helloprofiledb.AddBasicBlockParams{
			Bio:             request.Bio,
			Fullname:        request.Fullname,
			Title:           request.Title,
			CoverColour:     sql.NullString{String: request.CoverColour, Valid: true},
			CoverPhotoUrl:   sql.NullString{String: request.CoverPhotoUrl, Valid: true},
			ProfilePhotoUrl: sql.NullString{String: request.ProfilePhotoUrl, Valid: true},
		})
		if err != nil {
			errorResponse.Errorcode = util.SQL_ERROR_CODE
			errorResponse.ErrorMessage = util.SQL_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while adding basic block for profile ", profileId)
			c.JSON(http.StatusBadRequest, errorResponse)
			return err
		}
		err = env.HelloProfileDb.UpdateProfileWithBasicBlockId(context.Background(), helloprofiledb.UpdateProfileWithBasicBlockIdParams{
			BasicBlockID: uuid.NullUUID{UUID: dbAddBasicResult.ID, Valid: true},
			ID:           profileId,
		})
		log.WithFields(fields).Info("Successfully added basic block for profile ", profileId)

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
	dbBasic, err := env.HelloProfileDb.GetBasicBlock(context.Background(), request.ID)
	if err != nil {
		errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
		errorResponse.ErrorMessage = util.NO_RECORD_FOUND_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Basic block update failed. Basic block not found")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	err = env.HelloProfileDb.UpdateBasicBlock(context.Background(), helloprofiledb.UpdateBasicBlockParams{
		ProfilePhotoUrl: sql.NullString{String: env.GetValue(request.ProfilePhotoUrl, dbBasic.ProfilePhotoUrl.String), Valid: true},
		CoverPhotoUrl:   sql.NullString{String: env.GetValue(request.CoverPhotoUrl, dbBasic.CoverPhotoUrl.String), Valid: true},
		CoverColour:     sql.NullString{String: env.GetValue(request.CoverColour, dbBasic.CoverColour.String), Valid: true},
		Fullname:        env.GetValue(request.Fullname, dbBasic.Fullname),
		Title:           env.GetValue(request.Title, dbBasic.Title),
		Bio:             env.GetValue(request.Bio, dbBasic.Bio),
		ID:              request.ID,
	})
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
	if c.QueryParam("id") != "" && c.QueryParam("profileId") != "" {
		id, err := uuid.Parse(c.QueryParam("id"))
		if err != nil {
			errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
			errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Incorrect id format passed for delete basic block")
			c.JSON(http.StatusBadRequest, errorResponse)
			return err
		}
		profileId, err := uuid.Parse(c.QueryParam("profileId"))
		if err != nil {
			errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
			errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Incorrect id format passed for profileId in delete contact block")
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
			//update the profile to remove the contact ID from the profile
			profile, _ := env.HelloProfileDb.GetProfile(context.Background(), profileId)
			profile.BasicBlockID = uuid.NullUUID{UUID: uuid.Nil, Valid: true}
			_ = env.HelloProfileDb.UpdateProfile(context.Background(), helloprofiledb.UpdateProfileParams{
				UserID:         profile.UserID,
				Status:         profile.Status,
				ProfileName:    profile.ProfileName,
				BasicBlockID:   profile.BasicBlockID,
				ContactBlockID: profile.ContactBlockID,
				PageColor:      profile.PageColor,
				Font:           profile.Font,
				IsDefault:      profile.IsDefault,
				ID:             profile.ID,
			})
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
