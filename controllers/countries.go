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

// GetCountries is used get countries
func (env *Env) GetCountries(c echo.Context) (err error) {
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
	log.WithFields(fields).Info("Get countries request received...")
	countries, err := env.AuthDb.GetCountries(context.Background())
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Countries not found"
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Countries not found")

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.WithFields(fields).Info("Successfully retrieved countries...")
	countriesResponse := make([]models.Country, len(countries))
	for index, value := range countries {
		country := models.Country{
			Country:     value.Name,
			FlagURL:     value.FlagImageUrl.String,
			CountryCode: value.CountryCode.String,
		}
		countriesResponse[index] = country
	}
	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
		ResponseDetails: countriesResponse,
	}
	c.JSON(http.StatusOK, response)
	return err
}

// GetCountry is used get country
func (env *Env) GetCountry(c echo.Context) (err error) {

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

	log.WithFields(fields).Info("Get country request received...")
	country := c.Param("country")
	if country == "" {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "Country not specified"
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Country not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	log.WithFields(fields).Info(fmt.Sprintf("Country: %s", country))

	dbCountry, err := env.AuthDb.GetCountry(context.Background(), strings.ToLower(country))
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Country not found"
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Country not found")

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.WithFields(fields).Info(fmt.Sprintf("Successfully retrieved country: %v", dbCountry))
	countryResponse := models.Country{
		Country:     dbCountry.Name,
		FlagURL:     dbCountry.FlagImageUrl.String,
		CountryCode: dbCountry.CountryCode.String,
	}

	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
		ResponseDetails: countryResponse,
	}
	c.JSON(http.StatusOK, response)
	return err
}

// AddCountry is used add country
func (env *Env) AddCountry(c echo.Context) (err error) {

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
	log.WithFields(fields).Info("Add country request received...")
	request := new(models.Country)
	if err = c.Bind(request); err != nil {
		errorResponse.Errorcode = "02"
		errorResponse.ErrorMessage = "Invalid request"
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to marshal request")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	dbCountry, err := env.AuthDb.CreateCountry(context.Background(), authdb.CreateCountryParams{
		Name:         strings.ToLower(request.Country),
		FlagImageUrl: sql.NullString{String: strings.ToLower(request.FlagURL), Valid: request.FlagURL != ""},
		CountryCode:  sql.NullString{String: strings.ToLower(request.CountryCode), Valid: request.CountryCode != ""},
	})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not add country. Duplicate found"
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured adding new country")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	log.WithFields(fields).Info(fmt.Sprintf("Successfully added country: %v", dbCountry))

	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	c.JSON(http.StatusOK, response)
	return err
}

// UpdateCountry is used add country
func (env *Env) UpdateCountry(c echo.Context) (err error) {

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
	log.WithFields(fields).Info("Update country request received...")

	country := c.Param("country")
	if country == "" {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "Country not specified"
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Country not specified")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	log.WithFields(fields).Info(fmt.Sprintf("Country: %s", country))

	request := new(models.Country)
	if err = c.Bind(request); err != nil {
		errorResponse.Errorcode = "02"
		errorResponse.ErrorMessage = "Invalid request"
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to marshal request")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	dbCountry, err := env.AuthDb.UpdateCountry(context.Background(), authdb.UpdateCountryParams{
		Name:         strings.ToLower(request.Country),
		FlagImageUrl: sql.NullString{String: strings.ToLower(request.FlagURL), Valid: true},
		CountryCode:  sql.NullString{String: strings.ToLower(request.CountryCode), Valid: true},
		Name_2:       strings.ToLower(country),
	})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not update country. Country not found"
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured updating new country")

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.WithFields(fields).Info(fmt.Sprintf("Successfully updated country: %v", dbCountry))

	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	c.JSON(http.StatusOK, response)
	return err
}

// DeleteCountry is used add country
func (env *Env) DeleteCountry(c echo.Context) (err error) {

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
	log.WithFields(fields).Info("Delete country request received...")
	country := c.Param("country")
	if country == "" {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "Country not specified"
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Country not specified")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	log.WithFields(fields).Info(fmt.Sprintf("Country: %s", country))

	err = env.AuthDb.DeleteCountry(context.Background(), strings.ToLower(country))
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not delete country. Country not found"
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured deleting  country")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	log.WithFields(fields).Info("Successfully deleted country")

	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	c.JSON(http.StatusOK, response)
	return err
}
