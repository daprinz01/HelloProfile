package controllers

import (
	"authengine/models"
	"authengine/persistence/orm/authdb"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// GetTimezones is used get languages
func (env *Env) GetTimezones(c echo.Context) (err error) {
	log.Println("Get timezones request received...")
	errorResponse := new(models.Errormessage)

	timezones, err := env.AuthDb.GetTimezones(context.Background())
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Timezones not found"
		log.Println(fmt.Sprintf("timezones not found %s", err))

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.Println("Successfully retrieved timezones...")
	timezonesResponse := make([]models.Timezone, len(timezones))
	for index, value := range timezones {
		timezone := models.Timezone{
			Timezone: value.Name,
			Zone:     value.Zone,
		}
		timezonesResponse[index] = timezone
	}
	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
		ResponseDetails: timezonesResponse,
	}
	c.JSON(http.StatusOK, response)
	return err
}

// GetTimezone is used get timezone
func (env *Env) GetTimezone(c echo.Context) (err error) {
	log.Println("Get timezone request received...")

	errorResponse := new(models.Errormessage)

	timezone := c.Param("timezone")

	log.Println(fmt.Sprintf("Timezone: %s", timezone))
	if err != nil {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "Timezone not specified"
		log.Println("Timezone not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	dbTimezone, err := env.AuthDb.GetTimezone(context.Background(), strings.ToLower(timezone))
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Timezone not found"
		log.Println(fmt.Sprintf("Timezone not found"))

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.Println(fmt.Sprintf("Successfully retrieved user languages: %v", dbTimezone))
	timezoneResponse := models.Timezone{
		Timezone: dbTimezone.Name,
		Zone:     dbTimezone.Zone,
	}

	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
		ResponseDetails: timezoneResponse,
	}
	c.JSON(http.StatusOK, response)
	return err
}

// AddTimezone is used add timezones
func (env *Env) AddTimezone(c echo.Context) (err error) {
	log.Println("Add timezone request received...")

	errorResponse := new(models.Errormessage)

	request := new(models.Timezone)
	if err = c.Bind(request); err != nil {
		log.Println(fmt.Sprintf("Error occured while trying to marshal request: %s", err))
		errorResponse.Errorcode = "02"
		errorResponse.ErrorMessage = "Invalid request"
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	dbTimezone, err := env.AuthDb.CreateTimezone(context.Background(), authdb.CreateTimezoneParams{
		Name: strings.ToLower(request.Timezone),
		Zone: strings.ToLower(request.Zone),
	})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not add timezone. Duplicate found"
		log.Println(fmt.Sprintf("Error occured adding new timezone: %s", err))

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.Println(fmt.Sprintf("Successfully added timezone: %v", dbTimezone))

	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	c.JSON(http.StatusOK, response)
	return err
}

// UpdateTimezone is used add timezone
func (env *Env) UpdateTimezone(c echo.Context) (err error) {
	log.Println("Update timezone request received...")

	errorResponse := new(models.Errormessage)

	timezone := c.Param("timezone")

	log.Println(fmt.Sprintf("Timezone: %s", timezone))
	if err != nil {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "Timezone not specified"
		log.Println("Timezone not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	request := new(models.Timezone)
	if err = c.Bind(request); err != nil {
		log.Println(fmt.Sprintf("Error occured while trying to marshal request: %s", err))
		errorResponse.Errorcode = "02"
		errorResponse.ErrorMessage = "Invalid request"
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	dbTimezone, err := env.AuthDb.UpdateTimezone(context.Background(), authdb.UpdateTimezoneParams{
		Name:   strings.ToLower(request.Timezone),
		Zone:   strings.ToLower(request.Zone),
		Name_2: strings.ToLower(timezone),
	})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not update timezone. Not found"
		log.Println(fmt.Sprintf("Error occured updating new timezone: %s", err))

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.Println(fmt.Sprintf("Successfully updated timezone: %v", dbTimezone))

	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	c.JSON(http.StatusOK, response)
	return err
}

// DeleteTimezone is used add languages
func (env *Env) DeleteTimezone(c echo.Context) (err error) {
	log.Println("Delete timezone request received...")

	errorResponse := new(models.Errormessage)

	timezone := c.Param("timezone")

	log.Println(fmt.Sprintf("Timezone: %s", timezone))
	if err != nil {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "Timezone not specified"
		log.Println("Timezone not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	err = env.AuthDb.DeleteTimezone(context.Background(), strings.ToLower(timezone))
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not delete timezone. Timezone not found"
		log.Println(fmt.Sprintf("Error occured deleting  timezone: %s", err))
		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.Println("Successfully deleted timezone")

	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	c.JSON(http.StatusOK, response)
	return err
}
