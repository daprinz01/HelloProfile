package controllers

import (
	"authengine/models"
	"authengine/persistence/orm/authdb"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// GetTimezones is used get languages
func (env *Env) GetTimezones(w http.ResponseWriter, r *http.Request) {
	log.Println("Get timezones request received...")
	var errorResponse models.Errormessage
	var err error
	timezones, err := env.AuthDb.GetTimezones(context.Background())
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Timezones not found"
		log.Println(fmt.Sprintf("timezones not found %s", err))

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
	}
	log.Println("Successfully retrieved timezones...")
	var timezonesResponse []models.Timezone
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
	responsebytes, err := json.MarshalIndent(response, "", "")
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responsebytes)
	return
}

// GetTimezone is used get timezone
func (env *Env) GetTimezone(w http.ResponseWriter, r *http.Request) {
	log.Println("Get timezone request received...")
	pathParams := mux.Vars(r)
	var errorResponse models.Errormessage
	var err error
	var timezone string
	if val2, ok := pathParams["timezone"]; ok {
		timezone = val2
		log.Println(fmt.Sprintf("Timezone: %s", timezone))
		if err != nil {
			errorResponse.Errorcode = "15"
			errorResponse.ErrorMessage = "Timezone not specified"
			log.Println("Timezone not specified")

			response, err := json.MarshalIndent(errorResponse, "", "")
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
			return
		}
	}

	dbTimezone, err := env.AuthDb.GetTimezone(context.Background(), strings.ToLower(timezone))
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Timezone not found"
		log.Println(fmt.Sprintf("Timezone not found"))

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
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
	responsebytes, err := json.MarshalIndent(response, "", "")
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responsebytes)
	return
}

// AddTimezone is used add timezones
func (env *Env) AddTimezone(w http.ResponseWriter, r *http.Request) {
	log.Println("Add timezone request received...")

	var errorResponse models.Errormessage
	var err error
	var request models.Timezone
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&request)
	defer r.Body.Close()
	if err != nil {
		errorResponse.Errorcode = "02"
		errorResponse.ErrorMessage = "Invalid request"
		log.Println(err)
		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)
		return
	}

	dbTimezone, err := env.AuthDb.CreateTimezone(context.Background(), authdb.CreateTimezoneParams{
		Name: strings.ToLower(request.Timezone),
		Zone: strings.ToLower(request.Zone),
	})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not add timezone. Duplicate found"
		log.Println(fmt.Sprintf("Error occured adding new timezone: %s", err))

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
	}
	log.Println(fmt.Sprintf("Successfully added timezone: %v", dbTimezone))

	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	responsebytes, err := json.MarshalIndent(response, "", "")
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responsebytes)
	return
}

// UpdateTimezone is used add timezone
func (env *Env) UpdateTimezone(w http.ResponseWriter, r *http.Request) {
	log.Println("Update timezone request received...")
	pathParams := mux.Vars(r)
	var errorResponse models.Errormessage
	var err error
	var timezone string
	if val, ok := pathParams["timezone"]; ok {
		timezone = val
		log.Println(fmt.Sprintf("Timezone: %s", timezone))
		if err != nil {
			errorResponse.Errorcode = "15"
			errorResponse.ErrorMessage = "Timezone not specified"
			log.Println("Timezone not specified")

			response, err := json.MarshalIndent(errorResponse, "", "")
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
			return
		}
	}

	var request models.Timezone
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&request)
	defer r.Body.Close()
	if err != nil {
		errorResponse.Errorcode = "02"
		errorResponse.ErrorMessage = "Invalid request"
		log.Println(err)
		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)
		return
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

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
	}
	log.Println(fmt.Sprintf("Successfully updated timezone: %v", dbTimezone))

	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	responsebytes, err := json.MarshalIndent(response, "", "")
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responsebytes)
	return
}

// DeleteTimezone is used add languages
func (env *Env) DeleteTimezone(w http.ResponseWriter, r *http.Request) {
	log.Println("Delete timezone request received...")
	pathParams := mux.Vars(r)
	var errorResponse models.Errormessage
	var err error
	var timezone string
	if val, ok := pathParams["timezone"]; ok {
		timezone = val
		log.Println(fmt.Sprintf("Timezone: %s", timezone))
		if err != nil {
			errorResponse.Errorcode = "15"
			errorResponse.ErrorMessage = "Timezone not specified"
			log.Println("Timezone not specified")

			response, err := json.MarshalIndent(errorResponse, "", "")
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
			return
		}
	}

	err = env.AuthDb.DeleteTimezone(context.Background(), strings.ToLower(timezone))
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not delete timezone. Timezone not found"
		log.Println(fmt.Sprintf("Error occured deleting  timezone: %s", err))

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
	}
	log.Println("Successfully deleted timezone")

	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	responsebytes, err := json.MarshalIndent(response, "", "")
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responsebytes)
	return
}
