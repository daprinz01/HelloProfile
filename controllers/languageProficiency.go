package controllers

import (
	"authengine/models"
	"authengine/persistence/orm/authdb"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// GetLanguageProficiencies is used get proficiencies
func (env *Env) GetLanguageProficiencies(c echo.Context) (err error) {
	log.Println("Get proficiencies request received...")
	errorResponse := new(models.Errormessage)

	proficiencies, err := env.AuthDb.GetLanguageProficiencies(context.Background())
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Language Proficiencies not found"
		log.Println("proficiencies not found")
		log.Println(err)
		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.Println(fmt.Sprintf("Successfully retrieved proficiency proficiencies: %v", proficiencies))
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
	log.Println("Get proficiencies request received...")

	errorResponse := new(models.Errormessage)

	proficiency := c.Param("proficiency")

	log.Println(fmt.Sprintf("LanguageProficiency: %s", proficiency))
	if proficiency == "" {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "LanguageProficiency not specified"
		log.Println("LanguageProficiency not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	proficiencies, err := env.AuthDb.GetLanguageProficiency(context.Background(), sql.NullString{String: strings.ToLower(proficiency), Valid: proficiency != ""})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "LanguageProficiency not found"
		log.Println("proficiencies not found")

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.Println(fmt.Sprintf("Successfully retrieved user proficiencies: %v", proficiencies))
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
	log.Println("Add proficiencies request received...")

	errorResponse := new(models.Errormessage)

	proficiency := c.Param("proficiency")

	log.Println(fmt.Sprintf("LanguageProficiency: %s", proficiency))
	if err != nil {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "LanguageProficiency not specified"
		log.Println("LanguageProficiency not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err

	}

	proficiencies, err := env.AuthDb.CreateLanguageProficiency(context.Background(), sql.NullString{String: strings.ToLower(proficiency), Valid: proficiency != ""})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not add proficiency. Duplicate found"
		log.Println(fmt.Sprintf("Error occured adding new proficiency: %s", err))

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.Println(fmt.Sprintf("Successfully added proficiency: %v", proficiencies))

	proficiencyResponse := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	c.JSON(http.StatusOK, proficiencyResponse)
	return err
}

// UpdateLanguageProficiency is used add proficiencies
func (env *Env) UpdateLanguageProficiency(c echo.Context) (err error) {
	log.Println("Update proficiencies request received...")

	errorResponse := new(models.Errormessage)

	proficiency := c.Param("proficiency")

	log.Println(fmt.Sprintf("LanguageProficiency: %s", proficiency))
	if proficiency == "" {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "LanguageProficiency not specified"
		log.Println("LanguageProficiency not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	newLanguageProficiency := c.Param("newProficiency")

	log.Println(fmt.Sprintf("New LanguageProficiency: %s", strings.ToLower(proficiency)))
	if newLanguageProficiency == "" {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "New LanguageProficiency not specified"
		log.Println("New LanguageProficiency not specified")

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
		log.Println(fmt.Sprintf("Error occured updating new proficiency: %s", err))

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	log.Println(fmt.Sprintf("Successfully updated proficiency: %v", proficiencies))
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
	log.Println("Delete proficiencies request received...")

	errorResponse := new(models.Errormessage)

	proficiency := c.Param("proficiency")

	log.Println(fmt.Sprintf("LanguageProficiency: %s", proficiency))
	if err != nil {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "LanguageProficiency not specified"
		log.Println("LanguageProficiency not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	err = env.AuthDb.DeleteLanguageProficiency(context.Background(), sql.NullString{String: strings.ToLower(proficiency), Valid: proficiency != ""})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not delete proficiency. LanguageProficiency not found"
		log.Println(fmt.Sprintf("Error occured deleting  proficiency: %s", err))

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	log.Println("Successfully deleted proficiency")

	proficiencyResponse := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	c.JSON(http.StatusOK, proficiencyResponse)
	return err
}
