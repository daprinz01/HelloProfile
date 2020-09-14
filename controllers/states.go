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

// GetStates is used get states
func (env *Env) GetStates(w http.ResponseWriter, r *http.Request) {
	log.Println("Get states request received...")
	var errorResponse models.Errormessage
	var err error
	states, err := env.AuthDb.GetStates(context.Background())
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "States not found"
		log.Println(fmt.Sprintf("States not found %s", err))

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
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
	responsebytes, err := json.MarshalIndent(response, "", "")
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responsebytes)
	return
}

// GetStatesByCountry is used get states
func (env *Env) GetStatesByCountry(w http.ResponseWriter, r *http.Request) {
	log.Println("Get states by country request received...")
	var errorResponse models.Errormessage
	var err error
	var country string
	pathParams := mux.Vars(r)
	if val1, ok := pathParams["country"]; ok {
		country = val1
		log.Println(fmt.Sprintf("Country: %s", country))
		if err != nil {
			errorResponse.Errorcode = "15"
			errorResponse.ErrorMessage = "Country not specified"
			log.Println("Country not specified")

			response, err := json.MarshalIndent(errorResponse, "", "")
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
			return
		}
	}
	states, err := env.AuthDb.GetStatesByCountry(context.Background(), strings.ToLower(country))
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "States not found"
		log.Println(fmt.Sprintf("States not found %s", err))

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
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
	responsebytes, err := json.MarshalIndent(response, "", "")
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responsebytes)
	return
}

// GetState is used get states
func (env *Env) GetState(w http.ResponseWriter, r *http.Request) {
	log.Println("Get state request received...")
	pathParams := mux.Vars(r)
	var errorResponse models.Errormessage
	var err error
	var state string
	if val2, ok := pathParams["state"]; ok {
		state = val2
		log.Println(fmt.Sprintf("State: %s", state))
		if err != nil {
			errorResponse.Errorcode = "15"
			errorResponse.ErrorMessage = "State not specified"
			log.Println("State not specified")

			response, err := json.MarshalIndent(errorResponse, "", "")
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
			return
		}
	}

	dbState, err := env.AuthDb.GetState(context.Background(), strings.ToLower(state))
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "State not found"
		log.Println(fmt.Sprintf("State not found"))

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
	}
	log.Println(fmt.Sprintf("Successfully retrieved states: %v", dbState))
	stateResponse := dbState.Name

	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
		ResponseDetails: stateResponse,
	}
	responsebytes, err := json.MarshalIndent(response, "", "")
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responsebytes)
	return
}

// AddState is used add states
func (env *Env) AddState(w http.ResponseWriter, r *http.Request) {
	log.Println("Add state request received...")
	pathParams := mux.Vars(r)
	var errorResponse models.Errormessage
	var err error
	var state string
	if val2, ok := pathParams["state"]; ok {
		state = val2
		log.Println(fmt.Sprintf("State: %s", state))
		if err != nil {
			errorResponse.Errorcode = "15"
			errorResponse.ErrorMessage = "State not specified"
			log.Println("State not specified")

			response, err := json.MarshalIndent(errorResponse, "", "")
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
			return
		}
	}

	var country string
	if val1, ok := pathParams["country"]; ok {
		country = val1
		log.Println(fmt.Sprintf("Country: %s", country))
		if err != nil {
			errorResponse.Errorcode = "15"
			errorResponse.ErrorMessage = "Country not specified"
			log.Println("Country not specified")

			response, err := json.MarshalIndent(errorResponse, "", "")
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
			return
		}
	}

	err = env.AuthDb.CreateState(context.Background(), authdb.CreateStateParams{
		Name:   strings.ToLower(state),
		Name_2: strings.ToLower(country),
	})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not add state. Duplicate found"
		log.Println(fmt.Sprintf("Error occured adding new state: %s", err))

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
	}
	log.Println("Successfully added timezone ")

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

// UpdateState is used add state
func (env *Env) UpdateState(w http.ResponseWriter, r *http.Request) {
	log.Println("Update state request received...")
	pathParams := mux.Vars(r)
	var errorResponse models.Errormessage
	var err error
	var state string
	if val, ok := pathParams["state"]; ok {
		state = val
		log.Println(fmt.Sprintf("State: %s", state))
		if err != nil {
			errorResponse.Errorcode = "15"
			errorResponse.ErrorMessage = "State not specified"
			log.Println("State not specified")

			response, err := json.MarshalIndent(errorResponse, "", "")
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
			return
		}
	}

	var newState string
	if val2, ok := pathParams["newState"]; ok {
		newState = val2
		log.Println(fmt.Sprintf("New State: %s", newState))
		if err != nil {
			errorResponse.Errorcode = "15"
			errorResponse.ErrorMessage = "New State not specified"
			log.Println("New State not specified")

			response, err := json.MarshalIndent(errorResponse, "", "")
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
			return
		}
	}

	err = env.AuthDb.UpdateState(context.Background(), authdb.UpdateStateParams{
		Name:   strings.ToLower(newState),
		Name_2: strings.ToLower(state),
	})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not update state. Not found"
		log.Println(fmt.Sprintf("Error occured updating new state: %s", err))

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
	}
	log.Println("Successfully updated state")

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

// DeleteState is used add state
func (env *Env) DeleteState(w http.ResponseWriter, r *http.Request) {
	log.Println("Delete state request received...")
	pathParams := mux.Vars(r)
	var errorResponse models.Errormessage
	var err error
	var state string
	if val, ok := pathParams["state"]; ok {
		state = val
		log.Println(fmt.Sprintf("State: %s", state))
		if err != nil {
			errorResponse.Errorcode = "15"
			errorResponse.ErrorMessage = "State not specified"
			log.Println("State not specified")

			response, err := json.MarshalIndent(errorResponse, "", "")
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
			return
		}
	}

	err = env.AuthDb.DeleteState(context.Background(), strings.ToLower(state))
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not delete state. State not found"
		log.Println(fmt.Sprintf("Error occured deleting  state: %s", err))

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
	}
	log.Println("Successfully deleted state")

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
