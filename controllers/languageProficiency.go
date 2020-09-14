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

// GetLanguageProficiencies is used get proficiencies
func (env *Env) GetLanguageProficiencies(w http.ResponseWriter, r *http.Request) {
	log.Println("Get proficiencies request received...")
	var errorResponse models.Errormessage
	var err error
	proficiencies, err := env.AuthDb.GetLanguageProficiencies(context.Background())
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Language Proficiencies not found"
		log.Println("proficiencies not found")

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
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
	responsebytes, err := json.MarshalIndent(proficiencyResponse, "", "")
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responsebytes)
	return
}

// GetLanguageProficiency is used get proficiencies
func (env *Env) GetLanguageProficiency(w http.ResponseWriter, r *http.Request) {
	log.Println("Get proficiencies request received...")
	pathParams := mux.Vars(r)
	var errorResponse models.Errormessage
	var err error
	var proficiency string
	if val2, ok := pathParams["proficiency"]; ok {
		proficiency = val2
		log.Println(fmt.Sprintf("LanguageProficiency: %s", proficiency))
		if err != nil {
			errorResponse.Errorcode = "15"
			errorResponse.ErrorMessage = "LanguageProficiency not specified"
			log.Println("LanguageProficiency not specified")

			response, err := json.MarshalIndent(errorResponse, "", "")
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
			return
		}
	}

	proficiencies, err := env.AuthDb.GetLanguageProficiency(context.Background(), sql.NullString{String: strings.ToLower(proficiency), Valid: proficiency != ""})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "LanguageProficiency not found"
		log.Println("proficiencies not found")

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
	}
	log.Println(fmt.Sprintf("Successfully retrieved user proficiencies: %v", proficiencies))
	proficienciesResponse := proficiencies.Proficiency.String

	proficiencyResponse := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
		ResponseDetails: proficienciesResponse,
	}
	responsebytes, err := json.MarshalIndent(proficiencyResponse, "", "")
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responsebytes)
	return
}

// AddLanguageProficiency is used add proficiencies
func (env *Env) AddLanguageProficiency(w http.ResponseWriter, r *http.Request) {
	log.Println("Add proficiencies request received...")
	pathParams := mux.Vars(r)
	var errorResponse models.Errormessage
	var err error
	var proficiency string
	if val, ok := pathParams["proficiency"]; ok {
		proficiency = val
		log.Println(fmt.Sprintf("LanguageProficiency: %s", proficiency))
		if err != nil {
			errorResponse.Errorcode = "15"
			errorResponse.ErrorMessage = "LanguageProficiency not specified"
			log.Println("LanguageProficiency not specified")

			response, err := json.MarshalIndent(errorResponse, "", "")
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
			return
		}
	}

	proficiencies, err := env.AuthDb.CreateLanguageProficiency(context.Background(), sql.NullString{String: strings.ToLower(proficiency), Valid: proficiency != ""})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not add proficiency. Duplicate found"
		log.Println(fmt.Sprintf("Error occured adding new proficiency: %s", err))

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
	}
	log.Println(fmt.Sprintf("Successfully added proficiency: %v", proficiencies))

	proficiencyResponse := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	responsebytes, err := json.MarshalIndent(proficiencyResponse, "", "")
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responsebytes)
	return
}

// UpdateLanguageProficiency is used add proficiencies
func (env *Env) UpdateLanguageProficiency(w http.ResponseWriter, r *http.Request) {
	log.Println("Update proficiencies request received...")
	pathParams := mux.Vars(r)
	var errorResponse models.Errormessage
	var err error
	var proficiency string
	if val, ok := pathParams["proficiency"]; ok {
		proficiency = val
		log.Println(fmt.Sprintf("LanguageProficiency: %s", proficiency))
		if err != nil {
			errorResponse.Errorcode = "15"
			errorResponse.ErrorMessage = "LanguageProficiency not specified"
			log.Println("LanguageProficiency not specified")

			response, err := json.MarshalIndent(errorResponse, "", "")
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
			return
		}
	}

	var newLanguageProficiency string
	if val, ok := pathParams["newProficiency"]; ok {
		newLanguageProficiency = val
		log.Println(fmt.Sprintf("New LanguageProficiency: %s", strings.ToLower(proficiency)))
		if err != nil {
			errorResponse.Errorcode = "15"
			errorResponse.ErrorMessage = "New LanguageProficiency not specified"
			log.Println("New LanguageProficiency not specified")

			response, err := json.MarshalIndent(errorResponse, "", "")
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
			return
		}
	}

	proficiencies, err := env.AuthDb.UpdateLanguageProficiency(context.Background(), authdb.UpdateLanguageProficiencyParams{
		Proficiency:   sql.NullString{String: strings.ToLower(newLanguageProficiency), Valid: newLanguageProficiency != ""},
		Proficiency_2: sql.NullString{String: strings.ToLower(proficiency), Valid: proficiency != ""},
	})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not update proficiency. Duplicate found"
		log.Println(fmt.Sprintf("Error occured updating new proficiency: %s", err))

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
	}
	log.Println(fmt.Sprintf("Successfully updated proficiency: %v", proficiencies))
	proficienciesResponse := proficiencies.Proficiency.String

	proficiencyResponse := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
		ResponseDetails: proficienciesResponse,
	}
	responsebytes, err := json.MarshalIndent(proficiencyResponse, "", "")
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responsebytes)
	return
}

// DeleteLanguageProficiency is used add proficiencies
func (env *Env) DeleteLanguageProficiency(w http.ResponseWriter, r *http.Request) {
	log.Println("Delete proficiencies request received...")
	pathParams := mux.Vars(r)
	var errorResponse models.Errormessage
	var err error
	var proficiency string
	if val, ok := pathParams["proficiency"]; ok {
		proficiency = val
		log.Println(fmt.Sprintf("LanguageProficiency: %s", proficiency))
		if err != nil {
			errorResponse.Errorcode = "15"
			errorResponse.ErrorMessage = "LanguageProficiency not specified"
			log.Println("LanguageProficiency not specified")

			response, err := json.MarshalIndent(errorResponse, "", "")
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
			return
		}
	}

	err = env.AuthDb.DeleteLanguageProficiency(context.Background(), sql.NullString{String: strings.ToLower(proficiency), Valid: proficiency != ""})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not delete proficiency. LanguageProficiency not found"
		log.Println(fmt.Sprintf("Error occured deleting  proficiency: %s", err))

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
	}
	log.Println("Successfully deleted proficiency")

	proficiencyResponse := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	responsebytes, err := json.MarshalIndent(proficiencyResponse, "", "")
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responsebytes)
	return
}
