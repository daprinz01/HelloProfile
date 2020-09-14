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
	"strings"

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
	languages, err := env.AuthDb.GetUserLanguages(context.Background(), sql.NullString{String: strings.ToLower(username), Valid: true})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "User does not have any language yet"
		log.Println(fmt.Sprintf("User languages not found: %s", err))

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)
		return
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

	languages, err := env.AuthDb.AddUserLanguage(context.Background(), authdb.AddUserLanguageParams{
		Username:    sql.NullString{String: strings.ToLower(username), Valid: true},
		Name:        strings.ToLower(language),
		Proficiency: sql.NullString{String: strings.ToLower(proficiency), Valid: true},
	})
	if err != nil {

		log.Println(fmt.Sprintf("Error occured while adding user langauge: %s", err))
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Language does not exist"

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)
		return
	}
	log.Println(fmt.Sprintf("Successfully added user languages: %v", languages))

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
			Username: sql.NullString{String: strings.ToLower(username), Valid: true},
			Name:     strings.ToLower(language),
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

// GetLanguages is used get languages
func (env *Env) GetLanguages(w http.ResponseWriter, r *http.Request) {
	log.Println("Get languages request received...")
	var errorResponse models.Errormessage
	var err error
	languages, err := env.AuthDb.GetLanguages(context.Background())
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Languages not found"
		log.Println("languages not found")

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
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
	responsebytes, err := json.MarshalIndent(languageResponse, "", "")
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responsebytes)
	return
}

// GetLanguage is used get languages
func (env *Env) GetLanguage(w http.ResponseWriter, r *http.Request) {
	log.Println("Get languages request received...")
	pathParams := mux.Vars(r)
	var errorResponse models.Errormessage
	var err error
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

	languages, err := env.AuthDb.GetLanguage(context.Background(), strings.ToLower(language))
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Language not found"
		log.Println("languages not found")

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
	}
	log.Println(fmt.Sprintf("Successfully retrieved user languages: %v", languages))
	languagesResponse := languages.Name

	languageResponse := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
		ResponseDetails: languagesResponse,
	}
	responsebytes, err := json.MarshalIndent(languageResponse, "", "")
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responsebytes)
	return
}

// AddLanguage is used add languages
func (env *Env) AddLanguage(w http.ResponseWriter, r *http.Request) {
	log.Println("Add languages request received...")
	pathParams := mux.Vars(r)
	var errorResponse models.Errormessage
	var err error
	var language string
	if val, ok := pathParams["language"]; ok {
		language = val
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

	languages, err := env.AuthDb.CreateLanguage(context.Background(), strings.ToLower(language))
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not add language. Duplicate found"
		log.Println(fmt.Sprintf("Error occured adding new language: %s", err))

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
	}
	log.Println(fmt.Sprintf("Successfully added language: %v", languages))

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

// UpdateLanguage is used add languages
func (env *Env) UpdateLanguage(w http.ResponseWriter, r *http.Request) {
	log.Println("Update languages request received...")
	pathParams := mux.Vars(r)
	var errorResponse models.Errormessage
	var err error
	var language string
	if val, ok := pathParams["language"]; ok {
		language = val
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

	var newLanguage string
	if val, ok := pathParams["newLanguage"]; ok {
		newLanguage = val
		log.Println(fmt.Sprintf("New Language: %s", strings.ToLower(language)))
		if err != nil {
			errorResponse.Errorcode = "15"
			errorResponse.ErrorMessage = "New Language not specified"
			log.Println("New Language not specified")

			response, err := json.MarshalIndent(errorResponse, "", "")
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
			return
		}
	}

	languages, err := env.AuthDb.UpdateLanguage(context.Background(), authdb.UpdateLanguageParams{
		Name:   strings.ToLower(newLanguage),
		Name_2: strings.ToLower(language),
	})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not update language. Duplicate found"
		log.Println(fmt.Sprintf("Error occured updating new language: %s", err))

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
	}
	log.Println(fmt.Sprintf("Successfully updated language: %v", languages))
	languagesResponse := languages.Name

	languageResponse := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
		ResponseDetails: languagesResponse,
	}
	responsebytes, err := json.MarshalIndent(languageResponse, "", "")
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responsebytes)
	return
}

// DeleteLanguage is used add languages
func (env *Env) DeleteLanguage(w http.ResponseWriter, r *http.Request) {
	log.Println("Delete languages request received...")
	pathParams := mux.Vars(r)
	var errorResponse models.Errormessage
	var err error
	var language string
	if val, ok := pathParams["language"]; ok {
		language = val
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

	err = env.AuthDb.DeleteLanguage(context.Background(), strings.ToLower(language))
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not delete language. Language not found"
		log.Println(fmt.Sprintf("Error occured deleting  language: %s", err))

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
	}
	log.Println("Successfully deleted language")

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
