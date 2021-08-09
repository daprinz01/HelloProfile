package controllers

import (
	"authengine/models"
	"authengine/persistence/orm/authdb"
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
)

// GetLanguageProficiencies is used get proficiencies
func (env *Env) GetLanguageProficiencies(c echo.Context) (err error) {

	errorResponse := new(models.Errormessage)
	application := c.Param("application")
	if application == "" {
		errorResponse.Errorcode = "01"
		errorResponse.ErrorMessage = "Application not specified"
		log.WithField("microservice", "persian.black.authengine.service").Error("Calling application not specified")
		c.JSON(http.StatusBadRequest, errorResponse)
		return nil
	}
	fields := log.Fields{"microservice": "persian.black.authengine.service", "application": application}
	log.WithFields(fields).Info("Get proficiencies request received...")
	proficiencies, err := env.AuthDb.GetLanguageProficiencies(context.Background())
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Language Proficiencies not found"
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("proficiencies not found")
		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.WithFields(fields).Info(fmt.Sprintf("Successfully retrieved proficiency proficiencies: %v", proficiencies))
	proficienciesResponse := make([]string, len(proficiencies))
	for index, value := range proficiencies {
		proficienciesResponse[index] = value.Proficiency.String
	}
	proficiencyResponse := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
		ResponseDetails: proficienciesResponse,
	}
	c.JSON(http.StatusOK, proficiencyResponse)
	return err
}

// GetLanguageProficiency is used get proficiencies
func (env *Env) GetLanguageProficiency(c echo.Context) (err error) {

	errorResponse := new(models.Errormessage)
	application := c.Param("application")
	if application == "" {
		errorResponse.Errorcode = "01"
		errorResponse.ErrorMessage = "Application not specified"
		log.WithField("microservice", "persian.black.authengine.service").Error("Calling application not specified")
		c.JSON(http.StatusBadRequest, errorResponse)
		return nil
	}
	fields := log.Fields{"microservice": "persian.black.authengine.service", "application": application}
	log.WithFields(fields).Info("Get proficiencies request received...")
	proficiency := c.Param("proficiency")
	log.WithFields(fields).Info(fmt.Sprintf("LanguageProficiency: %s", proficiency))
	if proficiency == "" {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "LanguageProficiency not specified"
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("LanguageProficiency not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	proficiencies, err := env.AuthDb.GetLanguageProficiency(context.Background(), sql.NullString{String: strings.ToLower(proficiency), Valid: proficiency != ""})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "LanguageProficiency not found"
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("proficiencies not found")

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.WithFields(fields).Info(fmt.Sprintf("Successfully retrieved user proficiencies: %v", proficiencies))
	proficienciesResponse := proficiencies.Proficiency.String

	proficiencyResponse := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
		ResponseDetails: proficienciesResponse,
	}
	c.JSON(http.StatusOK, proficiencyResponse)
	return err
}

// AddLanguageProficiency is used add proficiencies
func (env *Env) AddLanguageProficiency(c echo.Context) (err error) {

	errorResponse := new(models.Errormessage)
	application := c.Param("application")
	if application == "" {
		errorResponse.Errorcode = "01"
		errorResponse.ErrorMessage = "Application not specified"
		log.WithField("microservice", "persian.black.authengine.service").Error("Calling application not specified")
		c.JSON(http.StatusBadRequest, errorResponse)
		return nil
	}
	fields := log.Fields{"microservice": "persian.black.authengine.service", "application": application}
	log.WithFields(fields).Info("Add proficiencies request received...")
	proficiency := c.Param("proficiency")

	log.WithFields(fields).Info(fmt.Sprintf("LanguageProficiency: %s", proficiency))
	if err != nil {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "LanguageProficiency not specified"
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("LanguageProficiency not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err

	}

	proficiencies, err := env.AuthDb.CreateLanguageProficiency(context.Background(), sql.NullString{String: strings.ToLower(proficiency), Valid: proficiency != ""})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not add proficiency. Duplicate found"
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured adding new proficiency")

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.WithFields(fields).Info(fmt.Sprintf("Successfully added proficiency: %v", proficiencies))

	proficiencyResponse := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	c.JSON(http.StatusOK, proficiencyResponse)
	return err
}

// UpdateLanguageProficiency is used add proficiencies
func (env *Env) UpdateLanguageProficiency(c echo.Context) (err error) {

	errorResponse := new(models.Errormessage)
	application := c.Param("application")
	if application == "" {
		errorResponse.Errorcode = "01"
		errorResponse.ErrorMessage = "Application not specified"
		log.WithField("microservice", "persian.black.authengine.service").Error("Calling application not specified")
		c.JSON(http.StatusBadRequest, errorResponse)
		return nil
	}
	fields := log.Fields{"microservice": "persian.black.authengine.service", "application": application}
	log.WithFields(fields).Info("Update proficiencies request received...")
	proficiency := c.Param("proficiency")

	log.WithFields(fields).Info(fmt.Sprintf("LanguageProficiency: %s", proficiency))
	if proficiency == "" {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "LanguageProficiency not specified"
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("LanguageProficiency not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	newLanguageProficiency := c.Param("newProficiency")

	log.Println(fmt.Sprintf("New LanguageProficiency: %s", strings.ToLower(proficiency)))
	if newLanguageProficiency == "" {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "New LanguageProficiency not specified"
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("New LanguageProficiency not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	proficiencies, err := env.AuthDb.UpdateLanguageProficiency(context.Background(), authdb.UpdateLanguageProficiencyParams{
		Proficiency:   sql.NullString{String: strings.ToLower(newLanguageProficiency), Valid: newLanguageProficiency != ""},
		Proficiency_2: sql.NullString{String: strings.ToLower(proficiency), Valid: proficiency != ""},
	})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not update proficiency. Duplicate found"
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured updating new proficiency")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	log.WithFields(fields).Info(fmt.Sprintf("Successfully updated proficiency: %v", proficiencies))
	proficienciesResponse := proficiencies.Proficiency.String

	proficiencyResponse := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
		ResponseDetails: proficienciesResponse,
	}
	c.JSON(http.StatusOK, proficiencyResponse)
	return err
}

// DeleteLanguageProficiency is used add proficiencies
func (env *Env) DeleteLanguageProficiency(c echo.Context) (err error) {

	errorResponse := new(models.Errormessage)
	application := c.Param("application")
	if application == "" {
		errorResponse.Errorcode = "01"
		errorResponse.ErrorMessage = "Application not specified"
		log.WithField("microservice", "persian.black.authengine.service").Error("Calling application not specified")
		c.JSON(http.StatusBadRequest, errorResponse)
		return nil
	}
	fields := log.Fields{"microservice": "persian.black.authengine.service", "application": application}
	log.WithFields(fields).Info("Delete proficiencies request received...")
	proficiency := c.Param("proficiency")

	log.WithFields(fields).Info(fmt.Sprintf("LanguageProficiency: %s", proficiency))
	if err != nil {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "LanguageProficiency not specified"
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("LanguageProficiency not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	err = env.AuthDb.DeleteLanguageProficiency(context.Background(), sql.NullString{String: strings.ToLower(proficiency), Valid: proficiency != ""})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not delete proficiency. LanguageProficiency not found"
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured deleting  proficiency")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	log.WithFields(fields).Info("Successfully deleted proficiency")

	proficiencyResponse := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	c.JSON(http.StatusOK, proficiencyResponse)
	return err
}
