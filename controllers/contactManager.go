package controllers

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"helloprofile.com/models"
	"helloprofile.com/persistence/orm/helloprofiledb"
	"helloprofile.com/util"
)

// GetProfiles is used get logged in user's contacts
func (env *Env) GetContacts(c echo.Context) (err error) {

	errorResponse := new(models.Errormessage)

	fields := log.Fields{"microservice": "helloprofile.service", "application": "backend", "function": "GetProfiles"}
	log.WithFields(fields).Info("Get contacts belonging to the logged in user request received...")
	email := c.Request().Header.Get("email")
	if email != "" {
		user, err := env.HelloProfileDb.GetUser(context.Background(), sql.NullString{String: email, Valid: true})
		if err != nil {
			errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
			errorResponse.ErrorMessage = util.USER_NOT_FOUND_RESPONSE_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("User was not found")
			c.JSON(http.StatusNotFound, errorResponse)
			return err
		}
		dbContacts, err := env.HelloProfileDb.GetContacts(context.Background(), user.ID)
		var userContacts []models.Contact
		for _, dbContact := range dbContacts {
			profileChan := make(chan models.Profile)
			env.getProfile(dbContact.ProfileID, profileChan, fields)
			userContacts = append(userContacts, models.Contact{
				Profile:    <-profileChan,
				CategoryID: dbContact.ContactCategoryID,
			})
		}

		response := &models.SuccessResponse{
			ResponseCode:    util.SUCCESS_RESPONSE_CODE,
			ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
			ResponseDetails: userContacts,
		}
		c.JSON(http.StatusOK, response)
		return err
	} else {
		errorResponse.Errorcode = util.UNAUTHORIZED_ERROR_CODE
		errorResponse.ErrorMessage = util.UNAUTHORIZED_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("User email not found in header")
		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
}

func (env *Env) AddContact(c echo.Context) (err error) {
	errorResponse := new(models.Errormessage)

	fields := log.Fields{"microservice": "helloprofile.service", "application": "backend", "function": "AddContact"}
	log.WithFields(fields).Info("Add contact request received...")
	email := c.Request().Header.Get("email")
	if email != "" {
		user, err := env.HelloProfileDb.GetUser(context.Background(), sql.NullString{String: email, Valid: true})
		if err != nil {
			errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
			errorResponse.ErrorMessage = util.USER_NOT_FOUND_RESPONSE_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("User was not found")
			c.JSON(http.StatusNotFound, errorResponse)
			return err
		}
		request := new(models.AddContact)
		if err = c.Bind(request); err != nil {
			errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
			errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to marshal request")
			c.JSON(http.StatusBadRequest, errorResponse)
			return err
		}
		addContactResponse, err := env.HelloProfileDb.AddContacts(context.Background(), helloprofiledb.AddContactsParams{
			UserID:            user.ID,
			ProfileID:         request.ProfileID,
			ContactCategoryID: request.CategoryID,
		})
		if err != nil {
			errorResponse.Errorcode = util.SQL_ERROR_CODE
			errorResponse.ErrorMessage = util.SQL_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to add new contact for user ", user.ID)
			c.JSON(http.StatusNotFound, errorResponse)
			return err
		}
		response := &models.SuccessResponse{
			ResponseCode:    util.SUCCESS_RESPONSE_CODE,
			ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
			ResponseDetails: addContactResponse.ID,
		}
		c.JSON(http.StatusOK, response)
		return err
	} else {
		errorResponse.Errorcode = util.UNAUTHORIZED_ERROR_CODE
		errorResponse.ErrorMessage = util.UNAUTHORIZED_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("User email not found in header")
		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
}

// DeleteContact deletes a contact belonging to the logged in user
func (env *Env) DeleteContact(c echo.Context) (err error) {

	errorResponse := new(models.Errormessage)

	fields := log.Fields{"microservice": "helloprofile.service", "application": "backend", "function": "DeleteContact"}
	log.WithFields(fields).Info("Delete contact request received...")
	email := c.Request().Header.Get("email")
	if c.QueryParam("id") != "" && email != "" {
		user, err := env.HelloProfileDb.GetUser(context.Background(), sql.NullString{String: email, Valid: true})
		if err != nil {
			errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
			errorResponse.ErrorMessage = util.USER_NOT_FOUND_RESPONSE_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("User was not found")
			c.JSON(http.StatusNotFound, errorResponse)
			return err
		}
		id, err := uuid.Parse(c.QueryParam("id"))
		if err != nil {
			errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
			errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Incorrect id format passed for delete contact")
			c.JSON(http.StatusBadRequest, errorResponse)
			return err
		}
		err = env.HelloProfileDb.DeleteContact(context.Background(), helloprofiledb.DeleteContactParams{
			UserID:    user.ID,
			ProfileID: id,
		})
		if err != nil {
			errorResponse.Errorcode = util.SQL_ERROR_CODE
			errorResponse.ErrorMessage = util.SQL_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Contact deletion failed")
			c.JSON(http.StatusBadRequest, errorResponse)
			return err
		} else {
			log.WithFields(fields).Info("Successfully deleted contact")
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
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Id not passed for delete contact request")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
}
