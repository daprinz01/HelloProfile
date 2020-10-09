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

// GetApplications is used get countries
func (env *Env) GetApplications(c echo.Context) (err error) {
	log.Println("Get applications request received...")
	errorResponse := new(models.Errormessage)

	applications, err := env.AuthDb.GetApplications(context.Background())
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Applications not found"
		log.Println(fmt.Sprintf("Application not found: %s", err))

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	log.Println("Successfully retrieved application...")
	applicationResponse := make([]models.Application, len(applications))
	for index, value := range applications {
		application := models.Application{
			Application: value.Name,
			Description: value.Description,
		}
		applicationResponse[index] = application
	}
	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
		ResponseDetails: applicationResponse,
	}
	c.JSON(http.StatusOK, response)
	return err
}

// GetApplication is used get application
func (env *Env) GetApplication(c echo.Context) (err error) {
	log.Println("Get application request received...")

	errorResponse := new(models.Errormessage)

	application := c.Param("application")
	if application == "" {
		errorResponse.Errorcode = "01"
		errorResponse.ErrorMessage = "Application not specified"
		log.Println("Application not specified")
		c.JSON(http.StatusBadRequest, errorResponse)
		return nil
	}

	dbApplication, err := env.AuthDb.GetApplication(context.Background(), strings.ToLower(application))
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Application not found"
		log.Println(fmt.Sprintf("Application not found: %s", err))

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	log.Println(fmt.Sprintf("Successfully retrieved application: %v", dbApplication))
	applicationResponse := models.Application{
		Application: dbApplication.Name,
		Description: dbApplication.Description,
	}

	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
		ResponseDetails: applicationResponse,
	}
	c.JSON(http.StatusOK, response)
	return err
}

// AddApplication is used add application
func (env *Env) AddApplication(c echo.Context) (err error) {
	log.Println("Add application request received...")

	errorResponse := new(models.Errormessage)

	request := new(models.Application)
	if err = c.Bind(request); err != nil {
		log.Println(fmt.Sprintf("Error occured while trying to marshal request: %s", err))
		errorResponse.Errorcode = "02"
		errorResponse.ErrorMessage = "Invalid request"
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	dbApplication, err := env.AuthDb.CreateApplication(context.Background(), authdb.CreateApplicationParams{
		Name:        strings.ToLower(request.Application),
		Description: strings.ToLower(request.Description),
	})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not add application. Duplicate found"
		log.Println(fmt.Sprintf("Error occured adding new application: %s", err))

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	log.Println(fmt.Sprintf("Successfully added application: %v", dbApplication))

	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	c.JSON(http.StatusOK, response)
	return err
}

// UpdateApplication is used add application
func (env *Env) UpdateApplication(c echo.Context) (err error) {
	log.Println("Update application request received...")

	errorResponse := new(models.Errormessage)

	application := c.Param("application")
	if application == "" {
		errorResponse.Errorcode = "01"
		errorResponse.ErrorMessage = "Application not specified"
		log.Println("Application not specified")
		c.JSON(http.StatusBadRequest, errorResponse)
		return nil
	}

	request := new(models.Application)
	if err = c.Bind(request); err != nil {
		log.Println(fmt.Sprintf("Error occured while trying to marshal request: %s", err))
		errorResponse.Errorcode = "02"
		errorResponse.ErrorMessage = "Invalid request"
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	dbApplication, err := env.AuthDb.UpdateApplication(context.Background(), authdb.UpdateApplicationParams{
		Name:        strings.ToLower(request.Application),
		Description: strings.ToLower(request.Description),
		Name_2:      strings.ToLower(application),
	})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not update application. Not found"
		log.Println(fmt.Sprintf("Error occured updating new application: %s", err))

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	log.Println(fmt.Sprintf("Successfully updated application: %v", dbApplication))

	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	c.JSON(http.StatusOK, response)
	return err
}

// DeleteApplication is used add application
func (env *Env) DeleteApplication(c echo.Context) (err error) {
	log.Println("Delete application request received...")

	errorResponse := new(models.Errormessage)

	application := c.Param("application")
	if application == "" {
		errorResponse.Errorcode = "01"
		errorResponse.ErrorMessage = "Application not specified"
		log.Println("Application not specified")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	err = env.AuthDb.DeleteApplication(context.Background(), strings.ToLower(application))
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not delete application. Application not found"
		log.Println(fmt.Sprintf("Error occured deleting  application: %s", err))

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.Println("Successfully deleted application")

	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	c.JSON(http.StatusOK, response)
	return err
}
