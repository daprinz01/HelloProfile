package controllers

import (
	"authengine/models"
	"authengine/persistence/orm/authdb"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// GetUserLanguages is used to retreive languages set by the user
func (env *Env) GetUserLanguages(w http.ResponseWriter, r *http.Request) {
	log.Println("Get user languages Request received")

	pathParams := mux.Vars(r)
	// file, fileHeader, err := r.FormFile("request.AttachmentName[i]")

	// file, err := os.Create(fmt.Sprintf("%s%s", attachmentPath, request.AttachmentName[i].FileName))
	// file.WriteString()

	var username string
	var errorResponse models.Errormessage
	var err error
	if val, ok := pathParams["username"]; ok {
		username = val
		log.Println(fmt.Sprintf("Username: %s", username))
		if err != nil {
			errorResponse.Errorcode = "15"
			errorResponse.ErrorMessage = "Username not specified"
			log.Println("Username not specified")

			response, err := json.MarshalIndent(errorResponse, "", "")
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
			return
		}
	}
	languages, err := env.AuthDb.GetUserLanguages(context.Background(), sql.NullString{String: username, Valid: true})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "User does not have any language yet"
		log.Println("User languages not found")

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)
		return
	}
	log.Println(fmt.Sprintf("Successfully retrieved user languages: %v", languages))
	var userLanguages []models.UserLanguage
	for index, value := range languages {
		userLanguages[index].Language = value.Name
		userLanguages[index].Proficiency = value.Proficiency
	}
	languageResponse := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
		ResponseDetails: userLanguages,
	}
	responsebytes, err := json.MarshalIndent(languageResponse, "", "")
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responsebytes)
	return
}

// AddUserLanguage is used to add languages to a users account
func (env *Env) AddUserLanguage(w http.ResponseWriter, r *http.Request) {
	log.Println("Add user languages Request received")
	pathParams := mux.Vars(r)
	// file, fileHeader, err := r.FormFile("request.AttachmentName[i]")

	// file, err := os.Create(fmt.Sprintf("%s%s", attachmentPath, request.AttachmentName[i].FileName))
	// file.WriteString()

	var username string
	var errorResponse models.Errormessage
	var err error
	if val, ok := pathParams["username"]; ok {
		username = val
		log.Println(fmt.Sprintf("Username: %s", username))
		if err != nil {
			errorResponse.Errorcode = "15"
			errorResponse.ErrorMessage = "Username not specified"
			log.Println("Username not specified")

			response, err := json.MarshalIndent(errorResponse, "", "")
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
			return
		}
	}
	var language string
	if val2, ok := pathParams["language"]; ok {
		language = val2
		log.Println(fmt.Sprintf("Language: %s", language))
		if err != nil {
			errorResponse.Errorcode = "15"
			errorResponse.ErrorMessage = "Language not specified"
			log.Println("Language not specified")

			response, err := json.MarshalIndent(errorResponse, "", "")
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
			return
		}
	}

	var proficiency string
	if val3, ok := pathParams["proficiency"]; ok {
		proficiency = val3
		log.Println(fmt.Sprintf("Proficiency: %s", proficiency))
		if err != nil {
			errorResponse.Errorcode = "15"
			errorResponse.ErrorMessage = "Proficiency not specified"
			log.Println("Proficiency not specified")

			response, err := json.MarshalIndent(errorResponse, "", "")
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
			return
		}
	}
	go func() {
		languages, err := env.AuthDb.AddUserLanguage(context.Background(), authdb.AddUserLanguageParams{
			Username:    sql.NullString{String: username, Valid: true},
			Name:        language,
			Proficiency: sql.NullString{String: proficiency, Valid: true},
		})
		if err != nil {

			log.Println(fmt.Sprintf("Error occured while adding user langauge: %s", err))

		}
		log.Println(fmt.Sprintf("Successfully added user languages: %v", languages))
	}()
	languageResponse := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	responsebytes, err := json.MarshalIndent(languageResponse, "", "")
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responsebytes)
	return
}

// DeleteUserLanguages is used to retreive languages set by the user
func (env *Env) DeleteUserLanguages(w http.ResponseWriter, r *http.Request) {
	log.Println("Get user languages Request received")

	pathParams := mux.Vars(r)
	// file, fileHeader, err := r.FormFile("request.AttachmentName[i]")

	// file, err := os.Create(fmt.Sprintf("%s%s", attachmentPath, request.AttachmentName[i].FileName))
	// file.WriteString()

	var username string
	var errorResponse models.Errormessage
	var err error
	if val, ok := pathParams["username"]; ok {
		username = val
		log.Println(fmt.Sprintf("Username: %s", username))
		if err != nil {
			errorResponse.Errorcode = "15"
			errorResponse.ErrorMessage = "Username not specified"
			log.Println("Username not specified")

			response, err := json.MarshalIndent(errorResponse, "", "")
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
			return
		}
	}
	var language string
	if val2, ok := pathParams["language"]; ok {
		language = val2
		log.Println(fmt.Sprintf("Language: %s", language))
		if err != nil {
			errorResponse.Errorcode = "15"
			errorResponse.ErrorMessage = "Language not specified"
			log.Println("Language not specified")

			response, err := json.MarshalIndent(errorResponse, "", "")
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
			return
		}
	}
	err = env.AuthDb.DeleteUserLanguage(context.Background(),
		authdb.DeleteUserLanguageParams{
			Username: sql.NullString{String: username, Valid: true},
			Name:     language,
		})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Cannot delete language"
		log.Println(fmt.Sprintf("Cannot delete language: %s", err))

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)
		return
	}
	log.Println(fmt.Sprintf("Successfully delete user languages"))

	languageResponse := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	responsebytes, err := json.MarshalIndent(languageResponse, "", "")
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responsebytes)
	return
}
