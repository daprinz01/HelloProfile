package controllers

import (
	"authengine/models"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//Login is used to log user in
func (env *Env) Login(w http.ResponseWriter, r *http.Request) {
	log.Println("Login Request received")
	pathParams := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	var application string
	var errorResponse models.Errormessage
	var err error
	if val, ok := pathParams["application"]; ok {
		application = val
		log.Println(fmt.Sprintf("Application: %s", application))
		if err != nil {
			errorResponse.Errorcode = "01"
			errorResponse.ErrorMessage = "Application not specified"
			response, err := json.MarshalIndent(errorResponse, "", "")
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
			return
		}
	}
	var request models.LoginRequest
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&request)
	defer r.Body.Close()
	if err != nil {
		errorResponse.Errorcode = "02"
		errorResponse.ErrorMessage = "Invalid request"
		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)
		return
	}

	var username sql.NullString
	username.String = request.Username
	username.Valid = true
	user, err := env.AuthDb.GetUser(context.Background(), username)
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Error fetching user"
		log.Println(fmt.Sprintf("Error fetching user: %s", err))
		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)
		return
	}
	if user.Password.String == request.Password {
		response, err := json.MarshalIndent(user, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(response)

	} else {
		errorResponse.Errorcode = "04"
		errorResponse.ErrorMessage = "Error fetching user"
		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)

	}
	return
	// commentID := -1
	// if val, ok := pathParams["commentID"]; ok {
	// 	commentID, err = strconv.Atoi(val)
	// 	if err != nil {
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		w.Write([]byte(`{"message": "need a number"}`))
	// 		return

	// 		query := r.URL.Query()
	// 		name := query.Get("name")
	// 		if name == "" {
	// 			name = "Guest"
	// 		}
	// 		log.Printf("Received request for %s\n", name)
	// 		w.Write([]byte(fmt.Sprintf("Hello, %s\n", name)))
	// 	}
	// }
}

// func (env *orm.Env) login(w http.ResponseWriter, r *http.Request) {
// 	log.Println("Login Request received")
// 	pathParams := mux.Vars(r)
// 	w.Header().Set("Content-Type", "application/json")
// 	var application string
// 	var errorResponse models.Errormessage
// 	var err error
// 	if val, ok := pathParams["application"]; ok {
// 		application = val
// 		log.Println(fmt.Sprintf("Application: %s", application))
// 		if err != nil {
// 			errorResponse.Errorcode = "01"
// 			errorResponse.ErrorMessage = "Application not specified"
// 			response, err := json.MarshalIndent(errorResponse, "", "")
// 			if err != nil {
// 				log.Println(err)
// 			}
// 			w.WriteHeader(http.StatusBadRequest)
// 			w.Write(response)
// 			return
// 		}
// 	}
// 	var request models.LoginRequest
// 	decoder := json.NewDecoder(r.Body)
// 	err = decoder.Decode(&request)
// 	if err != nil {
// 		errorResponse.Errorcode = "02"
// 		errorResponse.ErrorMessage = "Invalid request"
// 		response, err := json.MarshalIndent(errorResponse, "", "")
// 		if err != nil {
// 			log.Println(err)
// 		}
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write(response)
// 		return
// 	}
// 	ctx := context.Background()
// 	var username sql.NullString
// 	username.String = request.Username
// 	username.Valid = true
// 	user, err := env.AuthDb.GetUser(ctx, username)
// 	if err != nil {
// 		errorResponse.Errorcode = "03"
// 		errorResponse.ErrorMessage = "Error fetching user"
// 		response, err := json.MarshalIndent(errorResponse, "", "")
// 		if err != nil {
// 			log.Println(err)
// 		}
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write(response)
// 		return
// 	}
// 	// commentID := -1
// 	// if val, ok := pathParams["commentID"]; ok {
// 	// 	commentID, err = strconv.Atoi(val)
// 	// 	if err != nil {
// 	// 		w.WriteHeader(http.StatusInternalServerError)
// 	// 		w.Write([]byte(`{"message": "need a number"}`))
// 	// 		return

// 	// 		query := r.URL.Query()
// 	// 		name := query.Get("name")
// 	// 		if name == "" {
// 	// 			name = "Guest"
// 	// 		}
// 	// 		log.Printf("Received request for %s\n", name)
// 	// 		w.Write([]byte(fmt.Sprintf("Hello, %s\n", name)))
// 	// 	}
// 	// }
// }
