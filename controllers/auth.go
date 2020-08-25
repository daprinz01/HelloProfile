package controllers

import (
	"authengine/models"
	"authengine/persistence/orm/authdb"
	"authengine/util"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

//Login is used to sign users in
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
	if util.VerifyHash(user.Password.String, request.Password) {
		response, err := json.MarshalIndent(user, "", "")
		if err != nil {
			log.Println(err)
		}
		log.Println("Generating authentication token...")
		authToken, refreshToken, err := util.GenerateJWT(user.Email, "Guest")
		if err != nil {
			errorResponse.Errorcode = "05"
			errorResponse.ErrorMessage = fmt.Sprintf("Error occured generating auth token: %s", err)
			response, err := json.MarshalIndent(errorResponse, "", "")
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
			return
		}
		log.Println("Storing refresh token in separate thread...")
		// store refresh token add later
		go func() {
			dbRefreshToken, err := env.AuthDb.CreateRefreshToken(context.Background(), authdb.CreateRefreshTokenParams{
				UserID: user.ID,
				Token:  refreshToken,
			})
			if err != nil {
				log.Println(err)
			}

			log.Println(fmt.Sprintf("Refresh Token Id: %d", dbRefreshToken.ID))
		}()
		w.Header().Set("Authorization", authToken)
		w.Header().Set("Refresh-Token", refreshToken)
		// log.Println(fmt.Sprintf("Auth token: %s, Refresh token: %s, Return object: %v", authToken, refreshToken, registerResponse))

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

// Register is used to register new users
func (env *Env) Register(w http.ResponseWriter, r *http.Request) {
	log.Println("Register Request received")
	pathParams := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	var applicationName string
	var errorResponse models.Errormessage
	var err error
	if val, ok := pathParams["application"]; ok {
		applicationName = val
		log.Println(fmt.Sprintf("Application: %s", applicationName))
		if err != nil {
			errorResponse.Errorcode = "01"
			errorResponse.ErrorMessage = "Application not specified"
			log.Println(err)
			response, err := json.MarshalIndent(errorResponse, "", "")
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
			return
		}
	}
	application, err := env.AuthDb.GetApplication(context.Background(), applicationName)
	if err != nil {
		errorResponse.Errorcode = "02"
		errorResponse.ErrorMessage = "Application is invalid"
		log.Println(err)
		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)
		return
	}

	log.Println(fmt.Sprintf("Applicaiton ID: %d", application.ID))

	var request models.UserDetail
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&request)
	// defer r.Body.Close()
	var hashedPassword string
	if err != nil {
		errorResponse.Errorcode = "03"
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
	// log.Println(fmt.Sprintf("Password: %s", request.Password))
	if request.Password != "" {
		hashedPassword = util.GenerateHash(request.Password)

	}
	log.Println("Successfully hashed password")
	log.Println("Inserting user...")
	user, err := env.AuthDb.CreateUser(context.Background(), authdb.CreateUserParams{
		// Address: sql.NullString{String: request.Address, Valid: true},
		// City:    sql.NullString{String: request.City, Valid: true},
		// Country: sql.NullString{String: request.Country, Valid: false},
		CreatedAt: time.Now(),
		Email:     request.Email,
		// Firstname:                 sql.NullString{String: request.Firstname, Valid: false},
		// ImageUrl:                  sql.NullString{String: request.ProfilePicture, Valid: false},
		IsActive: true,
		// IsEmailConfirmed:          request.IsEmailConfirmed,
		// IsLockedOut:               request.IsLockedOut,
		// IsPasswordSystemGenerated: request.IsPasswordSystemGenerated,
		// Lastname:                  sql.NullString{String: request.Lastname, Valid: false},
		Password: sql.NullString{String: hashedPassword, Valid: true},
		// State:                     sql.NullString{String: request.State, Valid: false},
		Username: sql.NullString{String: request.Username, Valid: true},
	})

	if err != nil {
		errorResponse.Errorcode = "04"
		errorResponse.ErrorMessage = fmt.Sprintf("Error occured creating user: %s", err)
		log.Println(fmt.Sprintf("Error occured creating user: %s", err))
		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)
		return
	}
	registerResponse := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
		ResponseDetails: &models.UserDetail{
			Address:                   user.Address.String,
			City:                      user.City.String,
			Country:                   user.Country.String,
			CreatedAt:                 user.CreatedAt,
			Email:                     user.Email,
			Firstname:                 user.Firstname.String,
			ProfilePicture:            user.ImageUrl.String,
			IsActive:                  user.IsActive,
			IsEmailConfirmed:          user.IsEmailConfirmed,
			IsLockedOut:               user.IsLockedOut,
			IsPasswordSystemGenerated: user.IsPasswordSystemGenerated,
			Lastname:                  user.Lastname.String,
			Password:                  "",
			State:                     user.State.String,
			Username:                  user.Username.String,
		},
	}
	responseString, err := json.MarshalIndent(registerResponse, "", "")
	if err != nil {
		log.Println(err)
	}
	// log.Println(fmt.Sprintf("Got to response string: %s", responseString))
	log.Println("Generating authentication token...")
	authToken, refreshToken, err := util.GenerateJWT(user.Email, "Guest")
	if err != nil {
		errorResponse.Errorcode = "05"
		errorResponse.ErrorMessage = fmt.Sprintf("Error occured generating auth token: %s", err)
		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)
		return
	}
	log.Println("Storing refresh token in separate thread...")
	// store refresh token add later
	go func() {
		dbRefreshToken, err := env.AuthDb.CreateRefreshToken(context.Background(), authdb.CreateRefreshTokenParams{
			UserID: user.ID,
			Token:  refreshToken,
		})
		if err != nil {
			log.Println(err)
		}

		log.Println(fmt.Sprintf("Refresh Token Id: %d", dbRefreshToken.ID))
	}()
	w.Header().Set("Authorization", authToken)
	w.Header().Set("Refresh-Token", refreshToken)
	// log.Println(fmt.Sprintf("Auth token: %s, Refresh token: %s, Return object: %v", authToken, refreshToken, registerResponse))
	w.WriteHeader(http.StatusOK)
	w.Write(responseString)
	log.Println("Successfully processed registration request")
	return
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
