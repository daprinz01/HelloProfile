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

// GetCountries is used get countries
func (env *Env) GetCountries(w http.ResponseWriter, r *http.Request) {
	log.Println("Get countries request received...")
	var errorResponse models.Errormessage
	var err error
	countries, err := env.AuthDb.GetCountries(context.Background())
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Countries not found"
		log.Println(fmt.Sprintf("Countries not found: %s", err))

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
	}
	log.Println("Successfully retrieved countries...")
	var countriesResponse []models.Country
	for index, value := range countries {
		country := models.Country{
			Country:     value.Name,
			FlagURL:     value.FlagImageUrl.String,
			CountryCode: value.CountryCode.String,
		}
		countriesResponse[index] = country
	}
	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
		ResponseDetails: countriesResponse,
	}
	responsebytes, err := json.MarshalIndent(response, "", "")
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responsebytes)
	return
}

// GetCountry is used get country
func (env *Env) GetCountry(w http.ResponseWriter, r *http.Request) {
	log.Println("Get country request received...")
	pathParams := mux.Vars(r)
	var errorResponse models.Errormessage
	var err error
	var country string
	if val2, ok := pathParams["country"]; ok {
		country = val2
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

	dbCountry, err := env.AuthDb.GetCountry(context.Background(), strings.ToLower(country))
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Country not found"
		log.Println(fmt.Sprintf("Country not found"))

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
	}
	log.Println(fmt.Sprintf("Successfully retrieved country: %v", dbCountry))
	countryResponse := models.Country{
		Country:     dbCountry.Name,
		FlagURL:     dbCountry.FlagImageUrl.String,
		CountryCode: dbCountry.CountryCode.String,
	}

	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
		ResponseDetails: countryResponse,
	}
	responsebytes, err := json.MarshalIndent(response, "", "")
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responsebytes)
	return
}

// AddCountry is used add country
func (env *Env) AddCountry(w http.ResponseWriter, r *http.Request) {
	log.Println("Add country request received...")

	var errorResponse models.Errormessage
	var err error
	var request models.Country
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

	dbCountry, err := env.AuthDb.CreateCountry(context.Background(), authdb.CreateCountryParams{
		Name:         strings.ToLower(request.Country),
		FlagImageUrl: sql.NullString{String: strings.ToLower(request.FlagURL), Valid: request.FlagURL != ""},
		CountryCode:  sql.NullString{String: strings.ToLower(request.CountryCode), Valid: request.CountryCode != ""},
	})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not add country. Duplicate found"
		log.Println(fmt.Sprintf("Error occured adding new country: %s", err))

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
	}
	log.Println(fmt.Sprintf("Successfully added country: %v", dbCountry))

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

// UpdateCountry is used add country
func (env *Env) UpdateCountry(w http.ResponseWriter, r *http.Request) {
	log.Println("Update country request received...")
	pathParams := mux.Vars(r)
	var errorResponse models.Errormessage
	var err error
	var country string
	if val, ok := pathParams["country"]; ok {
		country = val
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

	var request models.Country
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

	dbCountry, err := env.AuthDb.UpdateCountry(context.Background(), authdb.UpdateCountryParams{
		Name:         strings.ToLower(request.Country),
		FlagImageUrl: sql.NullString{String: strings.ToLower(request.FlagURL), Valid: true},
		CountryCode:  sql.NullString{String: strings.ToLower(request.CountryCode), Valid: true},
		Name_2:       strings.ToLower(country),
	})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not update country. Not found"
		log.Println(fmt.Sprintf("Error occured updating new country: %s", err))

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
	}
	log.Println(fmt.Sprintf("Successfully updated country: %v", dbCountry))

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

// DeleteCountry is used add country
func (env *Env) DeleteCountry(w http.ResponseWriter, r *http.Request) {
	log.Println("Delete country request received...")
	pathParams := mux.Vars(r)
	var errorResponse models.Errormessage
	var err error
	var country string
	if val, ok := pathParams["country"]; ok {
		country = val
		log.Println(fmt.Sprintf("Country: %s", country))
		if err != nil {
			errorResponse.Errorcode = "15"
			errorResponse.ErrorMessage = "Country not specified"
			log.Println(fmt.Sprintf("Country not specified %s", err))

			response, err := json.MarshalIndent(errorResponse, "", "")
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
			return
		}
	}

	err = env.AuthDb.DeleteCountry(context.Background(), strings.ToLower(country))
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not delete country. Country not found"
		log.Println(fmt.Sprintf("Error occured deleting  country: %s", err))

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
	}
	log.Println("Successfully deleted country")

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
