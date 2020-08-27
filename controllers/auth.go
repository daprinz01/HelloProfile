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
	"os"
	"strconv"
	"strings"
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
	applicationObject, err := env.AuthDb.GetApplication(context.Background(), strings.ToLower(application))
	if err != nil {
		errorResponse.Errorcode = "06"
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

	log.Println(fmt.Sprintf("Applicaiton ID: %d", applicationObject.ID))
	var request models.LoginRequest
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&request)
	defer r.Body.Close()
	if err != nil {
		errorResponse.Errorcode = "02"
		errorResponse.ErrorMessage = "Invalid request"
		log.Println(fmt.Sprintf("Invalid request: %s", err))
		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)
		return
	}

	var username sql.NullString
	username.String = strings.ToLower(request.Username)
	username.Valid = true
	user, err := env.AuthDb.GetUser(context.Background(), username)
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Oops... something is wrong here... your email or password is incorrect..."
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
		log.Println("Verifying that user is in the role access is being requested...")
		role := r.Header.Get("Role")
		userRoles, err := env.AuthDb.GetUserRoles(context.Background(), sql.NullInt64{Int64: user.ID, Valid: true})
		if err != nil {
			log.Println(`Invalid role entered... Changing to default role of "Guest"`)
			role = "Guest"
		}
		testRole := "Guest"
		log.Println("Searching roles...")
		// log.Println(userRoles)
		for i := 0; i < len(userRoles); i++ {

			if userRoles[i] == role {
				testRole = userRoles[i]
				break
			}
			log.Println(fmt.Sprintf("role: %s doesn't match role: %s ", userRoles[i], role))
		}
		role = testRole

		log.Println(fmt.Sprintf("Generating authentication token for user: %s role: %s...", user.Email, role))
		authToken, refreshToken, err := util.GenerateJWT(user.Email, role)
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
		loginResponse := &models.SuccessResponse{
			ResponseCode:    "00",
			ResponseMessage: "Success",
			ResponseDetails: &models.UserDetail{
				Address:                   user.Address.String,
				City:                      user.City.String,
				Country:                   user.Country.String,
				CreatedAt:                 user.CreatedAt,
				Email:                     user.Email,
				Firstname:                 user.Firstname.String,
				ProfilePicture:            user.ProfilePicture.String,
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
		responsebytes, err := json.MarshalIndent(loginResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
		w.Header().Set("Refresh-Token", refreshToken)
		w.Header().Set("Role", role)
		w.WriteHeader(http.StatusOK)
		w.Write(responsebytes)

	} else {
		errorResponse.Errorcode = "04"
		errorResponse.ErrorMessage = "Oops... something is wrong here... your email or password is incorrect..."
		log.Println("Password incorrect...")
		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)

	}
	return

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
	application, err := env.AuthDb.GetApplication(context.Background(), strings.ToLower(applicationName))
	if err != nil {
		errorResponse.Errorcode = "06"
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
	defer r.Body.Close()
	var hashedPassword string
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
		Email:     strings.ToLower(request.Email),
		// Firstname:                 sql.NullString{String: request.Firstname, Valid: false},
		// ImageUrl:                  sql.NullString{String: request.ProfilePicture, Valid: false},
		IsActive: true,
		// IsEmailConfirmed:          request.IsEmailConfirmed,
		// IsLockedOut:               request.IsLockedOut,
		// IsPasswordSystemGenerated: request.IsPasswordSystemGenerated,
		// Lastname:                  sql.NullString{String: request.Lastname, Valid: false},
		Password: sql.NullString{String: hashedPassword, Valid: true},
		// State:                     sql.NullString{String: request.State, Valid: false},
		Username: sql.NullString{String: strings.ToLower(request.Username), Valid: true},
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
	role := r.Header.Get("Role")
	dbRole, err := env.AuthDb.GetRole(context.Background(), strings.ToLower(role))
	if err != nil {
		log.Println(`Invalid role entered... Changing to default role of "Guest"`)
		role = "Guest"
	} else {
		log.Println(fmt.Sprintf("Creating token for user: %s | role: %s", user.Email, dbRole.Name))

	}
	go func() {
		log.Println("Verifying that role exist for the application")
		applicationRole, err := env.AuthDb.GetApplicationRole(context.Background(), authdb.GetApplicationRoleParams{
			ApplicationsID: sql.NullInt64{Int64: application.ID, Valid: true},
			RolesID:        sql.NullInt64{Int64: dbRole.ID, Valid: true},
		})
		if err != nil {
			log.Println(fmt.Sprintf("Error occured fetching applicationRole: %s", err))
		}
		log.Println(fmt.Sprintf("Role is valid for application. Application Role Id: %d", applicationRole.ID))
		log.Println("Adding user to role...")
		userRole, err := env.AuthDb.AddUserRole(context.Background(), authdb.AddUserRoleParams{
			Name:     strings.ToLower(role),
			Username: user.Username,
		})
		if err != nil {
			log.Println(fmt.Sprintf("Error occured adding user: %s to role: %s", user.Username.String, role))
		}
		log.Println(fmt.Sprintf("Successfully added user to role.. User Role Id: %d", userRole.ID))
	}()
	authToken, refreshToken, err := util.GenerateJWT(user.Email, role)
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
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
	w.Header().Set("Refresh-Token", refreshToken)
	w.Header().Set("Role", role)
	// log.Println(fmt.Sprintf("Auth token: %s, Refresh token: %s, Return object: %v", authToken, refreshToken, registerResponse))
	w.WriteHeader(http.StatusOK)
	w.Write(responseString)
	log.Println("Successfully processed registration request")
	return
}

// RefreshToken is used to register expired token
func (env *Env) RefreshToken(w http.ResponseWriter, r *http.Request) {
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
	application, err := env.AuthDb.GetApplication(context.Background(), strings.ToLower(applicationName))
	if err != nil {
		errorResponse.Errorcode = "06"
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
	var authCode string
	authArray := strings.Split(r.Header.Get("Authorization"), " ")
	if len(authArray) != 2 {
		errorResponse.Errorcode = "11"
		errorResponse.ErrorMessage = "Unsupported authentication scheme type"
		log.Println("Unsupported authentication scheme type")
		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)
		return
	}
	authCode = authArray[1]
	refreshToken := r.Header.Get("Refresh-Token")

	verifiedClaims, err := util.VerifyToken(authCode)
	if err == nil {
		errorResponse.Errorcode = "10"
		errorResponse.ErrorMessage = "Session is still valid..."
		log.Println(fmt.Sprintf("Token is still valid..."))
		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusTooEarly)
		w.Write(response)
		return
	}
	if err != nil && verifiedClaims.Email == "" {
		errorResponse.Errorcode = "02"
		errorResponse.ErrorMessage = "Invalid request"
		log.Println(fmt.Sprintf("Invalid request: %s", err))
		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)
		return
	}

	dbRefreshToken, err := env.AuthDb.GetRefreshToken(context.Background(), refreshToken)
	if err != nil {
		errorResponse.Errorcode = "08"
		errorResponse.ErrorMessage = "Cannot refresh session. Kindly login again"
		log.Println(fmt.Sprintf("Error occured refreshing token: %s", err))
		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)
		return
	}
	defer func() {
		err = env.AuthDb.DeleteRefreshToken(context.Background(), refreshToken)
		if err != nil {
			log.Println(err)
		}
	}()
	var refreshTokenDuration int
	refreshTokenLifespan := os.Getenv("SESSION_LIFESPAN")
	if refreshTokenLifespan == "" {
		log.Println("Session lifespan cannot be empty")
		log.Println("SESSION_LIFESPAN cannot be empty, setting duration to default of 15 mins ...")
		refreshTokenDuration = 15
	} else {
		log.Println(fmt.Sprintf("Setting Refresh token lifespan..."))
		refreshTokenDuration, err = strconv.Atoi(refreshTokenLifespan)
		if err != nil {
			log.Println(fmt.Sprintf("Error converting refresh token duration to number: %s", err))
		}
	}
	if !dbRefreshToken.CreatedAt.Add(time.Duration(refreshTokenDuration) * time.Minute).Before(time.Now()) {
		log.Println("Generating authentication token...")
		authToken, refreshToken, err := util.GenerateJWT(verifiedClaims.Email, verifiedClaims.Role)
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
		log.Println("Fetching user...")
		user, err := env.AuthDb.GetUser(context.Background(), sql.NullString{String: verifiedClaims.Email, Valid: true})
		if err != nil {
			errorResponse.Errorcode = "03"
			errorResponse.ErrorMessage = "Oops... something is wrong here... your email or password is incorrect..."
			log.Println(fmt.Sprintf("Error fetching user: %s", err))
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
		resetResponse := &models.RefreshResponse{
			ResponseCode:    "00",
			ResponseMessage: "Success",
		}
		responsebytes, err := json.MarshalIndent(resetResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
		w.Header().Set("Refresh-Token", refreshToken)
		w.Header().Set("Role", verifiedClaims.Role)
		w.WriteHeader(http.StatusOK)
		w.Write(responsebytes)
	} else {
		errorResponse.Errorcode = "09"
		errorResponse.ErrorMessage = "Session expired. Kindly login again to continue"
		log.Println("Token has expired...")
		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)

	}
	return

}

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
