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

// GetApplications is used get countries
func (env *Env) GetApplications(w http.ResponseWriter, r *http.Request) {
	log.Println("Get applications request received...")
	var errorResponse models.Errormessage
	var err error
	applications, err := env.AuthDb.GetApplications(context.Background())
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Applications not found"
		log.Println(fmt.Sprintf("Applications not found: %s", err))

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
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
	responsebytes, err := json.MarshalIndent(response, "", "")
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responsebytes)
	return
}

// GetApplication is used get application
func (env *Env) GetApplication(w http.ResponseWriter, r *http.Request) {
	log.Println("Get application request received...")
	pathParams := mux.Vars(r)
	var errorResponse models.Errormessage
	var err error
	var application string
	if val2, ok := pathParams["application"]; ok {
		application = val2
		log.Println(fmt.Sprintf("Application: %s", application))
		if err != nil {
			errorResponse.Errorcode = "15"
			errorResponse.ErrorMessage = "Application not specified"
			log.Println("Application not specified")

			response, err := json.MarshalIndent(errorResponse, "", "")
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
			return
		}
	}

	dbApplication, err := env.AuthDb.GetApplication(context.Background(), strings.ToLower(application))
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Application not found"
		log.Println(fmt.Sprintf("Application not found: %s", err))

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
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
	responsebytes, err := json.MarshalIndent(response, "", "")
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responsebytes)
	return
}

// AddApplication is used add application
func (env *Env) AddApplication(w http.ResponseWriter, r *http.Request) {
	log.Println("Add application request received...")

	var errorResponse models.Errormessage
	var err error
	var request models.Application
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

	dbApplication, err := env.AuthDb.CreateApplication(context.Background(), authdb.CreateApplicationParams{
		Name:        strings.ToLower(request.Application),
		Description: strings.ToLower(request.Description),
	})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not add application. Duplicate found"
		log.Println(fmt.Sprintf("Error occured adding new application: %s", err))

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
	}
	log.Println(fmt.Sprintf("Successfully added application: %v", dbApplication))

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

// UpdateApplication is used add application
func (env *Env) UpdateApplication(w http.ResponseWriter, r *http.Request) {
	log.Println("Update application request received...")
	pathParams := mux.Vars(r)
	var errorResponse models.Errormessage
	var err error
	var application string
	if val, ok := pathParams["application"]; ok {
		application = val
		log.Println(fmt.Sprintf("Application: %s", application))
		if err != nil {
			errorResponse.Errorcode = "15"
			errorResponse.ErrorMessage = "Application not specified"
			log.Println("Application not specified")

			response, err := json.MarshalIndent(errorResponse, "", "")
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
			return
		}
	}

	var request models.Application
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

	dbApplication, err := env.AuthDb.UpdateApplication(context.Background(), authdb.UpdateApplicationParams{
		Name:        strings.ToLower(request.Application),
		Description: strings.ToLower(request.Description),
		Name_2:      strings.ToLower(application),
	})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not update application. Not found"
		log.Println(fmt.Sprintf("Error occured updating new application: %s", err))

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
	}
	log.Println(fmt.Sprintf("Successfully updated application: %v", dbApplication))

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

// DeleteApplication is used add application
func (env *Env) DeleteApplication(w http.ResponseWriter, r *http.Request) {
	log.Println("Delete application request received...")
	pathParams := mux.Vars(r)
	var errorResponse models.Errormessage
	var err error
	var application string
	if val, ok := pathParams["application"]; ok {
		application = val
		log.Println(fmt.Sprintf("Application: %s", application))
		if err != nil {
			errorResponse.Errorcode = "15"
			errorResponse.ErrorMessage = "Application not specified"
			log.Println(fmt.Sprintf("Application not specified %s", err))

			response, err := json.MarshalIndent(errorResponse, "", "")
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
			return
		}
	}

	err = env.AuthDb.DeleteApplication(context.Background(), strings.ToLower(application))
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not delete application. Application not found"
		log.Println(fmt.Sprintf("Error occured deleting  application: %s", err))

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
	}
	log.Println("Successfully deleted application")

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
