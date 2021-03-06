package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"helloprofile.com/models"
	"helloprofile.com/persistence/orm/helloprofiledb"
	"helloprofile.com/util"
)

// AddProfile is used create a new profile
func (env *Env) AddContactBlock(c echo.Context) (err error) {

	errorResponse := new(models.Errormessage)

	fields := log.Fields{"microservice": "helloprofile.service", "application": "backend", "function": "AddBasicBlock"}
	log.WithFields(fields).Info("Add contact block request received...")

	if c.Param("profileId") != "" {
		request := new(models.ContactBlock)
		if err = c.Bind(request); err != nil {
			errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
			errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to marshal request")
			c.JSON(http.StatusBadRequest, errorResponse)
			return err
		}
		profileId, err := uuid.Parse(c.Param("profileId"))
		if err != nil || profileId == uuid.Nil {
			errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
			errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to marshal request")
			c.JSON(http.StatusBadRequest, errorResponse)
			return err
		}
		profile, err := env.HelloProfileDb.GetProfile(context.Background(), profileId)
		if err != nil {
			errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
			errorResponse.ErrorMessage = util.NO_RECORD_FOUND_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Profile was not found while trying to add contact block")
			c.JSON(http.StatusBadRequest, errorResponse)
			return err
		}
		if profile.ContactBlockID.UUID != uuid.Nil {
			errorResponse.Errorcode = util.CONTACT_BLOCK_EXIST_ERROR_CODE
			errorResponse.ErrorMessage = util.CONTACT_BLOCK_EXIST_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Contact block already exists for profile")
			c.JSON(http.StatusBadRequest, errorResponse)
			return err
		}

		log.WithFields(fields).Info(fmt.Sprintf("Contact block to add to profile %s : %v", profileId, request))

		dbAddContactResult, err := env.HelloProfileDb.AddContactBlock(context.Background(), helloprofiledb.AddContactBlockParams{
			Address: request.Address,
			Email:   request.Email,
			Phone:   request.Phone,
			Website: request.Website,
		})
		if err != nil {
			errorResponse.Errorcode = util.SQL_ERROR_CODE
			errorResponse.ErrorMessage = util.SQL_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while adding contact block for profile ", profileId)
			c.JSON(http.StatusBadRequest, errorResponse)
			return err
		}

		err = env.HelloProfileDb.UpdateProfileWithContactBlockId(context.Background(), helloprofiledb.UpdateProfileWithContactBlockIdParams{
			ContactBlockID: uuid.NullUUID{UUID: dbAddContactResult.ID, Valid: true},
			ID:             profileId,
		})
		log.WithFields(fields).Info("Successfully added contact block to profile ", profileId)

		response := &models.SuccessResponse{
			ResponseCode:    util.SUCCESS_RESPONSE_CODE,
			ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
			ResponseDetails: dbAddContactResult.ID,
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

// UpdateContactBlock is used udpate the contact block of a profile
func (env *Env) UpdateContactBlock(c echo.Context) (err error) {

	errorResponse := new(models.Errormessage)

	fields := log.Fields{"microservice": "helloprofile.service", "application": "backend", "function": "UpdateContactBlock"}
	log.WithFields(fields).Info("Update contact block request received...")

	request := new(models.ContactBlock)
	if err = c.Bind(request); err != nil {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to marshal request")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	dbContact, err := env.HelloProfileDb.GetContactBlock(context.Background(), request.ID)
	if err != nil {
		errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
		errorResponse.ErrorMessage = util.NO_RECORD_FOUND_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Contact block update failed. Contact block not found")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	err = env.HelloProfileDb.UpdateContactBlock(context.Background(), helloprofiledb.UpdateContactBlockParams{
		Phone:   env.GetValue(request.Phone, dbContact.Phone),
		Email:   env.GetValue(request.Email, dbContact.Email),
		Address: env.GetValue(request.Address, dbContact.Address),
		Website: env.GetValue(request.Website, dbContact.Website),
		ID:      request.ID,
	})
	if err != nil {
		errorResponse.Errorcode = util.SQL_ERROR_CODE
		errorResponse.ErrorMessage = util.SQL_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Contact block update failed")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	} else {
		log.WithFields(fields).Info("Successfully updated contact block")
		response := &models.SuccessResponse{
			ResponseCode:    util.SUCCESS_RESPONSE_CODE,
			ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
		}
		c.JSON(http.StatusOK, response)
		return err
	}
}

// DeleteContactBlock deletes a contact block
func (env *Env) DeleteContactBlock(c echo.Context) (err error) {

	errorResponse := new(models.Errormessage)

	fields := log.Fields{"microservice": "helloprofile.service", "application": "backend", "function": "DeleteContactBlock"}
	log.WithFields(fields).Info("Delete contact block request received...")
	if c.QueryParam("id") != "" && c.QueryParam("profileId") != "" {
		id, err := uuid.Parse(c.QueryParam("id"))
		if err != nil {
			errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
			errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Incorrect id format passed for delete contact block")
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
		err = env.HelloProfileDb.DeleteContactBlock(context.Background(), id)
		if err != nil {
			errorResponse.Errorcode = util.SQL_ERROR_CODE
			errorResponse.ErrorMessage = util.SQL_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Contact block update failed")
			c.JSON(http.StatusBadRequest, errorResponse)
			return err
		} else {
			//update the profile to remove the contact ID from the profile
			profile, _ := env.HelloProfileDb.GetProfile(context.Background(), profileId)
			profile.ContactBlockID = uuid.NullUUID{UUID: uuid.Nil, Valid: true}
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
			log.WithFields(fields).Info("Successfully deleted contact block")
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
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to read contact block id")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
}

// //
// func (env *Env) SaveContact(c echo.Context) (err error){

// }
