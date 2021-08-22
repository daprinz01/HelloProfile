package controllers

import (
	"authengine/models"
	"authengine/persistence/orm/authdb"
	"authengine/util"
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
)

// GetUserLanguages is used to retreive languages set by the user
func (env *Env) GetUserLanguages(c echo.Context) (err error) {

	// file, fileHeader, err := r.FormFile("request.AttachmentName[i]")

	// file, err := os.Create(fmt.Sprintf("%s%s", attachmentPath, request.AttachmentName[i].FileName))
	// file.WriteString()

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
	log.WithFields(fields).Info("Get user languages Request received")
	username := c.Param("username")
	log.WithFields(fields).Info(fmt.Sprintf("Username: %s", username))
	if username == "" {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Username not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	languages, err := env.AuthDb.GetUserLanguages(context.Background(), sql.NullString{String: strings.ToLower(username), Valid: true})
	if err != nil {
		errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
		errorResponse.ErrorMessage = util.NO_RECORD_FOUND_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error(fmt.Sprintf("User languages not found: %s", err))

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.WithFields(fields).Info(fmt.Sprintf("Successfully retrieved user languages: %v", languages))
	userLanguages := make([]models.UserLanguage, len(languages))
	for index, value := range languages {
		userLanguages[index].Language = value.Name
		userLanguages[index].Proficiency = value.Proficiency.String
	}
	languageResponse := &models.SuccessResponse{
		ResponseCode:    util.SUCCESS_RESPONSE_CODE,
		ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
		ResponseDetails: userLanguages,
	}
	c.JSON(http.StatusOK, languageResponse)
	return err
}

// AddUserLanguage is used to add languages to a users account
func (env *Env) AddUserLanguage(c echo.Context) (err error) {

	// file, fileHeader, err := r.FormFile("request.AttachmentName[i]")

	// file, err := os.Create(fmt.Sprintf("%s%s", attachmentPath, request.AttachmentName[i].FileName))
	// file.WriteString()

	username := c.Param("username")
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
	log.WithFields(fields).Info("Add user languages Request received")
	log.WithFields(fields).Info(fmt.Sprintf("Username: %s", username))
	if err != nil {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Username not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	language := c.Param("language")

	log.WithFields(fields).Info(fmt.Sprintf("Language: %s", language))
	if err != nil {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Language not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	proficiency := c.Param("proficiency")

	log.WithFields(fields).Info(fmt.Sprintf("Proficiency: %s", proficiency))
	if err != nil {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Proficiency not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	languages, err := env.AuthDb.AddUserLanguage(context.Background(), authdb.AddUserLanguageParams{
		Username:    sql.NullString{String: strings.ToLower(username), Valid: true},
		Name:        strings.ToLower(language),
		Proficiency: sql.NullString{String: strings.ToLower(proficiency), Valid: true},
	})
	if err != nil {
		errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
		errorResponse.ErrorMessage = util.NO_RECORD_FOUND_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while adding user langauge")
		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.WithFields(fields).Info(fmt.Sprintf("Successfully added user languages: %v", languages))

	languageResponse := &models.SuccessResponse{
		ResponseCode:    util.SUCCESS_RESPONSE_CODE,
		ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
	}
	c.JSON(http.StatusOK, languageResponse)
	return err
}

// DeleteUserLanguages is used to retreive languages set by the user
func (env *Env) DeleteUserLanguages(c echo.Context) (err error) {

	// file, fileHeader, err := r.FormFile("request.AttachmentName[i]")

	// file, err := os.Create(fmt.Sprintf("%s%s", attachmentPath, request.AttachmentName[i].FileName))
	// file.WriteString()

	username := c.Param("username")
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
	log.WithFields(fields).Info("Get user languages Request received")
	log.WithFields(fields).Info(fmt.Sprintf("Username: %s", username))
	if err != nil {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Username not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err

	}
	language := c.Param("language")

	if err != nil {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Language not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err

	}
	log.WithFields(fields).Info(fmt.Sprintf("Language: %s", language))
	err = env.AuthDb.DeleteUserLanguage(context.Background(),
		authdb.DeleteUserLanguageParams{
			Username: sql.NullString{String: strings.ToLower(username), Valid: true},
			Name:     strings.ToLower(language),
		})
	if err != nil {
		errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
		errorResponse.ErrorMessage = util.NO_RECORD_FOUND_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error(fmt.Sprintf("Cannot delete language: %s", err))

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.WithFields(fields).Info("Successfully delete user languages")

	languageResponse := &models.SuccessResponse{
		ResponseCode:    util.SUCCESS_RESPONSE_CODE,
		ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
	}
	c.JSON(http.StatusOK, languageResponse)
	return err
}

// GetLanguages is used get languages
func (env *Env) GetLanguages(c echo.Context) (err error) {

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
	log.WithFields(fields).Info("Get languages request received...")
	languages, err := env.AuthDb.GetLanguages(context.Background())
	if err != nil {
		errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
		errorResponse.ErrorMessage = util.NO_RECORD_FOUND_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("languages not found")

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.WithFields(fields).Info(fmt.Sprintf("Successfully retrieved user languages: %v", languages))
	languagesResponse := make([]string, len(languages))
	for index, value := range languages {
		languagesResponse[index] = value.Name
	}
	languageResponse := &models.SuccessResponse{
		ResponseCode:    util.SUCCESS_RESPONSE_CODE,
		ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
		ResponseDetails: languagesResponse,
	}
	c.JSON(http.StatusOK, languageResponse)
	return err
}

// GetLanguage is used get languages
func (env *Env) GetLanguage(c echo.Context) (err error) {

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
	log.WithFields(fields).Info("Get languages request received...")
	language := c.Param("language")

	if err != nil {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Language not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err

	}
	log.WithFields(fields).Info(fmt.Sprintf("Language: %s", language))
	languages, err := env.AuthDb.GetLanguage(context.Background(), strings.ToLower(language))
	if err != nil {
		errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
		errorResponse.ErrorMessage = util.NO_RECORD_FOUND_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("languages not found")

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.WithFields(fields).Info(fmt.Sprintf("Successfully retrieved user languages: %v", languages))
	languagesResponse := languages.Name

	languageResponse := &models.SuccessResponse{
		ResponseCode:    util.SUCCESS_RESPONSE_CODE,
		ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
		ResponseDetails: languagesResponse,
	}
	c.JSON(http.StatusOK, languageResponse)
	return err
}

// AddLanguage is used add languages
func (env *Env) AddLanguage(c echo.Context) (err error) {

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
	log.WithFields(fields).Info("Add languages request received...")
	language := c.Param("language")

	if err != nil {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Language not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	log.WithFields(fields).Info(fmt.Sprintf("Language: %s", language))
	languages, err := env.AuthDb.CreateLanguage(context.Background(), strings.ToLower(language))
	if err != nil {
		errorResponse.Errorcode = util.DUPLICATE_RECORD_ERROR_CODE
		errorResponse.ErrorMessage = util.DUPLICATE_RECORD_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured adding new language")

		c.JSON(http.StatusNotModified, errorResponse)
		return err
	}
	log.WithFields(fields).Info(fmt.Sprintf("Successfully added language: %v", languages))

	languageResponse := &models.SuccessResponse{
		ResponseCode:    util.SUCCESS_RESPONSE_CODE,
		ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
	}
	c.JSON(http.StatusOK, languageResponse)
	return err
}

// UpdateLanguage is used add languages
func (env *Env) UpdateLanguage(c echo.Context) (err error) {

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
	log.WithFields(fields).Info("Update languages request received...")
	language := c.Param("language")

	if err != nil {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Language not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	log.WithFields(fields).Info(fmt.Sprintf("Language: %s", language))

	newLanguage := c.Param("newLanguage")

	log.WithFields(fields).Info(fmt.Sprintf("New Language: %s", strings.ToLower(language)))
	if err != nil {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("New Language not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	languages, err := env.AuthDb.UpdateLanguage(context.Background(), authdb.UpdateLanguageParams{
		Name:   strings.ToLower(newLanguage),
		Name_2: strings.ToLower(language),
	})
	if err != nil {
		errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
		errorResponse.ErrorMessage = util.NO_RECORD_FOUND_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured updating new language")

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.WithFields(fields).Info(fmt.Sprintf("Successfully updated language: %v", languages))
	languagesResponse := languages.Name

	languageResponse := &models.SuccessResponse{
		ResponseCode:    util.SUCCESS_RESPONSE_CODE,
		ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
		ResponseDetails: languagesResponse,
	}
	c.JSON(http.StatusOK, languageResponse)
	return err
}

// DeleteLanguage is used add languages
func (env *Env) DeleteLanguage(c echo.Context) (err error) {

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
	log.WithFields(fields).Info("Delete languages request received...")
	language := c.Param("language")

	if err != nil {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Language not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	log.WithFields(fields).Info(fmt.Sprintf("Language: %s", language))
	err = env.AuthDb.DeleteLanguage(context.Background(), strings.ToLower(language))
	if err != nil {
		errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
		errorResponse.ErrorMessage = util.NO_RECORD_FOUND_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error(fmt.Sprintf("Error occured deleting  language: %s", err))

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.WithFields(fields).Info("Successfully deleted language")

	languageResponse := &models.SuccessResponse{
		ResponseCode:    util.SUCCESS_RESPONSE_CODE,
		ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
	}
	c.JSON(http.StatusOK, languageResponse)
	return err
}
