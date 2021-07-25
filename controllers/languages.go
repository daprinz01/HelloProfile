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

// GetUserLanguages is used to retreive languages set by the user
func (env *Env) GetUserLanguages(c echo.Context) (err error) {
	log.Println("Get user languages Request received")

	// file, fileHeader, err := r.FormFile("request.AttachmentName[i]")

	// file, err := os.Create(fmt.Sprintf("%s%s", attachmentPath, request.AttachmentName[i].FileName))
	// file.WriteString()

	errorResponse := new(models.Errormessage)

	username := c.Param("username")
	log.Println(fmt.Sprintf("Username: %s", username))
	if username == "" {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "Username not specified"
		log.Println("Username not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	languages, err := env.AuthDb.GetUserLanguages(context.Background(), sql.NullString{String: strings.ToLower(username), Valid: true})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "User does not have any language yet"
		log.Println(fmt.Sprintf("User languages not found: %s", err))

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	log.Println(fmt.Sprintf("Successfully retrieved user languages: %v", languages))
	userLanguages := make([]models.UserLanguage, len(languages))
	for index, value := range languages {
		userLanguages[index].Language = value.Name
		userLanguages[index].Proficiency = value.Proficiency.String
	}
	languageResponse := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
		ResponseDetails: userLanguages,
	}
	c.JSON(http.StatusOK, languageResponse)
	return err
}

// AddUserLanguage is used to add languages to a users account
func (env *Env) AddUserLanguage(c echo.Context) (err error) {
	log.Println("Add user languages Request received")

	// file, fileHeader, err := r.FormFile("request.AttachmentName[i]")

	// file, err := os.Create(fmt.Sprintf("%s%s", attachmentPath, request.AttachmentName[i].FileName))
	// file.WriteString()

	username := c.Param("username")
	errorResponse := new(models.Errormessage)

	log.Println(fmt.Sprintf("Username: %s", username))
	if err != nil {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "Username not specified"
		log.Println("Username not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	language := c.Param("language")

	log.Println(fmt.Sprintf("Language: %s", language))
	if err != nil {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "Language not specified"
		log.Println("Language not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	proficiency := c.Param("proficiency")

	log.Println(fmt.Sprintf("Proficiency: %s", proficiency))
	if err != nil {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "Proficiency not specified"
		log.Println("Proficiency not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	languages, err := env.AuthDb.AddUserLanguage(context.Background(), authdb.AddUserLanguageParams{
		Username:    sql.NullString{String: strings.ToLower(username), Valid: true},
		Name:        strings.ToLower(language),
		Proficiency: sql.NullString{String: strings.ToLower(proficiency), Valid: true},
	})
	if err != nil {

		log.Println(fmt.Sprintf("Error occured while adding user langauge: %s", err))
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Language does not exist"

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	log.Println(fmt.Sprintf("Successfully added user languages: %v", languages))

	languageResponse := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	c.JSON(http.StatusOK, languageResponse)
	return err
}

// DeleteUserLanguages is used to retreive languages set by the user
func (env *Env) DeleteUserLanguages(c echo.Context) (err error) {
	log.Println("Get user languages Request received")

	// file, fileHeader, err := r.FormFile("request.AttachmentName[i]")

	// file, err := os.Create(fmt.Sprintf("%s%s", attachmentPath, request.AttachmentName[i].FileName))
	// file.WriteString()

	username := c.Param("username")
	errorResponse := new(models.Errormessage)

	log.Println(fmt.Sprintf("Username: %s", username))
	if err != nil {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "Username not specified"
		log.Println("Username not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err

	}
	language := c.Param("language")

	if err != nil {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "Language not specified"
		log.Println("Language not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err

	}
	log.Println(fmt.Sprintf("Language: %s", language))
	err = env.AuthDb.DeleteUserLanguage(context.Background(),
		authdb.DeleteUserLanguageParams{
			Username: sql.NullString{String: strings.ToLower(username), Valid: true},
			Name:     strings.ToLower(language),
		})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Cannot delete language"
		log.Println(fmt.Sprintf("Cannot delete language: %s", err))

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.Println("Successfully delete user languages")

	languageResponse := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	c.JSON(http.StatusOK, languageResponse)
	return err
}

// GetLanguages is used get languages
func (env *Env) GetLanguages(c echo.Context) (err error) {
	log.Println("Get languages request received...")
	errorResponse := new(models.Errormessage)

	languages, err := env.AuthDb.GetLanguages(context.Background())
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Languages not found"
		log.Println("languages not found")

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.Println(fmt.Sprintf("Successfully retrieved user languages: %v", languages))
	languagesResponse := make([]string, len(languages))
	for index, value := range languages {
		languagesResponse[index] = value.Name
	}
	languageResponse := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
		ResponseDetails: languagesResponse,
	}
	c.JSON(http.StatusOK, languageResponse)
	return err
}

// GetLanguage is used get languages
func (env *Env) GetLanguage(c echo.Context) (err error) {
	log.Println("Get languages request received...")

	errorResponse := new(models.Errormessage)
	language := c.Param("language")

	if err != nil {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "Language not specified"
		log.Println("Language not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err

	}
	log.Println(fmt.Sprintf("Language: %s", language))
	languages, err := env.AuthDb.GetLanguage(context.Background(), strings.ToLower(language))
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Language not found"
		log.Println("languages not found")

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.Println(fmt.Sprintf("Successfully retrieved user languages: %v", languages))
	languagesResponse := languages.Name

	languageResponse := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
		ResponseDetails: languagesResponse,
	}
	c.JSON(http.StatusOK, languageResponse)
	return err
}

// AddLanguage is used add languages
func (env *Env) AddLanguage(c echo.Context) (err error) {
	log.Println("Add languages request received...")

	errorResponse := new(models.Errormessage)

	language := c.Param("language")

	if err != nil {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "Language not specified"
		log.Println("Language not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	log.Println(fmt.Sprintf("Language: %s", language))
	languages, err := env.AuthDb.CreateLanguage(context.Background(), strings.ToLower(language))
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not add language. Duplicate found"
		log.Println(fmt.Sprintf("Error occured adding new language: %s", err))

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.Println(fmt.Sprintf("Successfully added language: %v", languages))

	languageResponse := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	c.JSON(http.StatusOK, languageResponse)
	return err
}

// UpdateLanguage is used add languages
func (env *Env) UpdateLanguage(c echo.Context) (err error) {
	log.Println("Update languages request received...")

	errorResponse := new(models.Errormessage)

	language := c.Param("language")

	if err != nil {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "Language not specified"
		log.Println("Language not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	log.Println(fmt.Sprintf("Language: %s", language))

	newLanguage := c.Param("newLanguage")

	log.Println(fmt.Sprintf("New Language: %s", strings.ToLower(language)))
	if err != nil {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "New Language not specified"
		log.Println("New Language not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	languages, err := env.AuthDb.UpdateLanguage(context.Background(), authdb.UpdateLanguageParams{
		Name:   strings.ToLower(newLanguage),
		Name_2: strings.ToLower(language),
	})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not update language. Duplicate found"
		log.Println(fmt.Sprintf("Error occured updating new language: %s", err))

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.Println(fmt.Sprintf("Successfully updated language: %v", languages))
	languagesResponse := languages.Name

	languageResponse := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
		ResponseDetails: languagesResponse,
	}
	c.JSON(http.StatusOK, languageResponse)
	return err
}

// DeleteLanguage is used add languages
func (env *Env) DeleteLanguage(c echo.Context) (err error) {
	log.Println("Delete languages request received...")

	errorResponse := new(models.Errormessage)

	language := c.Param("language")

	if err != nil {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "Language not specified"
		log.Println("Language not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	log.Println(fmt.Sprintf("Language: %s", language))
	err = env.AuthDb.DeleteLanguage(context.Background(), strings.ToLower(language))
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not delete language. Language not found"
		log.Println(fmt.Sprintf("Error occured deleting  language: %s", err))

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	log.Println("Successfully deleted language")

	languageResponse := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	c.JSON(http.StatusOK, languageResponse)
	return err
}
