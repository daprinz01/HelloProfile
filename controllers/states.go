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

// GetStates is used get states
func (env *Env) GetStates(c echo.Context) (err error) {
	log.Println("Get states request received...")
	errorResponse := new(models.Errormessage)

	states, err := env.AuthDb.GetStates(context.Background())
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "States not found"
		log.Println(fmt.Sprintf("States not found %s", err))

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.Println("Successfully retrieved states...")
	statesResponse := make([]string, len(states))
	for index, value := range states {
		statesResponse[index] = value.Name
	}
	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
		ResponseDetails: statesResponse,
	}
	c.JSON(http.StatusOK, response)
	return err
}

// GetStatesByCountry is used get states
func (env *Env) GetStatesByCountry(c echo.Context) (err error) {
	log.Println("Get states by country request received...")
	errorResponse := new(models.Errormessage)

	country := c.Param("country")

	log.Println(fmt.Sprintf("Country: %s", country))
	if err != nil {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "Country not specified"
		log.Println("Country not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	states, err := env.AuthDb.GetStatesByCountry(context.Background(), strings.ToLower(country))
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "States not found"
		log.Println(fmt.Sprintf("States not found %s", err))
		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.Println("Successfully retrieved states...")
	statesResponse := make([]string, len(states))
	for index, value := range states {
		statesResponse[index] = value.Name
	}
	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
		ResponseDetails: statesResponse,
	}
	c.JSON(http.StatusOK, response)
	return err
}

// GetState is used get states
func (env *Env) GetState(c echo.Context) (err error) {
	log.Println("Get state request received...")

	errorResponse := new(models.Errormessage)

	state := c.Param("state")

	log.Println(fmt.Sprintf("State: %s", state))
	if err != nil {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "State not specified"
		log.Println("State not specified")

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}

	dbState, err := env.AuthDb.GetState(context.Background(), strings.ToLower(state))
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "State not found"
		log.Println(fmt.Sprintf("State not found"))

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.Println(fmt.Sprintf("Successfully retrieved states: %v", dbState))
	stateResponse := dbState.Name

	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
		ResponseDetails: stateResponse,
	}
	c.JSON(http.StatusOK, response)
	return err
}

// AddState is used add states
func (env *Env) AddState(c echo.Context) (err error) {
	log.Println("Add state request received...")

	errorResponse := new(models.Errormessage)

	state := c.Param("state")

	log.Println(fmt.Sprintf("State: %s", state))
	if err != nil {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "State not specified"
		log.Println("State not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	country := c.Param("country")

	log.Println(fmt.Sprintf("Country: %s", country))
	if err != nil {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "Country not specified"
		log.Println("Country not specified")

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}

	err = env.AuthDb.CreateState(context.Background(), authdb.CreateStateParams{
		Name:   strings.ToLower(state),
		Name_2: strings.ToLower(country),
	})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not add state. Duplicate found"
		log.Println(fmt.Sprintf("Error occured adding new state: %s", err))

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.Println("Successfully added timezone ")

	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	c.JSON(http.StatusOK, response)
	return err
}

// UpdateState is used add state
func (env *Env) UpdateState(c echo.Context) (err error) {
	log.Println("Update state request received...")

	errorResponse := new(models.Errormessage)

	state := c.Param("state")

	log.Println(fmt.Sprintf("State: %s", state))
	if err != nil {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "State not specified"
		log.Println("State not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	newState := c.Param("newState")

	log.Println(fmt.Sprintf("New State: %s", newState))
	if err != nil {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "New State not specified"
		log.Println("New State not specified")

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}

	err = env.AuthDb.UpdateState(context.Background(), authdb.UpdateStateParams{
		Name:   strings.ToLower(newState),
		Name_2: strings.ToLower(state),
	})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not update state. Not found"
		log.Println(fmt.Sprintf("Error occured updating new state: %s", err))

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.Println("Successfully updated state")

	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	c.JSON(http.StatusOK, response)
	return err
}

// DeleteState is used add state
func (env *Env) DeleteState(c echo.Context) (err error) {
	log.Println("Delete state request received...")

	errorResponse := new(models.Errormessage)

	state := c.Param("state")

	log.Println(fmt.Sprintf("State: %s", state))
	if err != nil {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "State not specified"
		log.Println("State not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	err = env.AuthDb.DeleteState(context.Background(), strings.ToLower(state))
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not delete state. State not found"
		log.Println(fmt.Sprintf("Error occured deleting  state: %s", err))

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.Println("Successfully deleted state")

	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	c.JSON(http.StatusOK, response)
	return err
}
