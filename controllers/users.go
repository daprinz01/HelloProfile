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

// CheckAvailability is used to check user availablity
func (env *Env) CheckAvailability(w http.ResponseWriter, r *http.Request) {
	log.Println("Check availability Request received")

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

	user, err := env.AuthDb.GetUser(context.Background(), sql.NullString{String: strings.ToLower(username), Valid: true})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "User does not exist"
		log.Println(fmt.Sprintf("Error fetching user: %s", err))
		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)
		return
	}
	log.Println(fmt.Sprintf("User %s exists...", user.Username.String))

	resetResponse := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	responsebytes, err := json.MarshalIndent(resetResponse, "", "")
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responsebytes)
	return
}

// GetUser is used to fetch user details
func (env *Env) GetUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Get User Request received")

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

	user, err := env.AuthDb.GetUser(context.Background(), sql.NullString{String: strings.ToLower(username), Valid: true})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "User does not exist"
		log.Println(fmt.Sprintf("Error fetching user: %s", err))
		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)
		return
	}
	log.Println(fmt.Sprintf("User %s exists...", user.Username.String))

	resetResponse := &models.SuccessResponse{
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
			Phone:                     user.Phone.String,
		},
	}
	responsebytes, err := json.MarshalIndent(resetResponse, "", "")
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responsebytes)
	return
}

// UpdateUser is used to update User information. It can be used to update user details and timezone details as required. Only pass the details to be updated. Email or username is mandatory.
func (env *Env) UpdateUser(w http.ResponseWriter, r *http.Request) {

	log.Println("Update user request received...")
	var errorResponse models.Errormessage
	var err error
	var request models.UserDetail
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

	log.Println("Checking if user exist...")
	user, err := env.AuthDb.GetUser(context.Background(), sql.NullString{String: strings.ToLower(request.Email), Valid: true})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "User does not exist"
		log.Println(fmt.Sprintf("Error fetching user: %s", err))
		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)
		return
	}
	log.Println(fmt.Sprintf("User %s exists...", user.Username.String))
	go func() {
		log.Println("Updating user...")
		// Check all the user parameters passed and consolidate with existing user record

		_, err := env.AuthDb.UpdateUser(context.Background(), authdb.UpdateUserParams{
			Address:                   sql.NullString{String: getValue(request.Address, user.Address.String), Valid: true},
			City:                      sql.NullString{String: getValue(request.City, user.City.String), Valid: true},
			Country:                   sql.NullString{String: getValue(request.Country, user.Country.String), Valid: true},
			CreatedAt:                 user.CreatedAt,
			Email:                     user.Email,
			Firstname:                 sql.NullString{String: getValue(request.Firstname, user.Firstname.String), Valid: true},
			ImageUrl:                  sql.NullString{String: getValue(request.ProfilePicture, user.ProfilePicture.String), Valid: true},
			IsActive:                  user.IsActive,
			IsEmailConfirmed:          user.IsEmailConfirmed,
			IsLockedOut:               user.IsLockedOut,
			IsPasswordSystemGenerated: user.IsPasswordSystemGenerated,
			Lastname:                  sql.NullString{String: getValue(request.Lastname, user.Lastname.String), Valid: true},
			Password:                  user.Password,
			State:                     sql.NullString{String: getValue(request.State, user.State.String), Valid: true},
			Username:                  user.Username,
			Phone:                     sql.NullString{String: getValue(request.Phone, user.Phone.String), Valid: true},
			Username_2:                sql.NullString{String: user.Email, Valid: true},
		})

		if err != nil {
			log.Println(fmt.Sprintf("Error occured updating user information: %s", err))

		}
		log.Println("Successfully updated user details")
		// LanguageName              string    `json:"language_name"`
		// RoleName                  string    `json:"role_name"`
		// TimezoneName              string    `json:"timezone_name"`
		// Zone                      string    `json:"zone"`
		// ProviderName              string    `json:"provider_name"`
		// ClientID                  string    `json:"client_id"`
		// ClientSecret              string    `json:"client_secret"`
		// ProviderLogo              string    `json:"provider_logo"`
		log.Println("Updating User Timezone")
		if request.TimezoneName != "" {
			_, err := env.AuthDb.UpdateUserTimezone(context.Background(), authdb.UpdateUserTimezoneParams{
				Username:   sql.NullString{String: user.Email, Valid: true},
				Name:       request.TimezoneName,
				Username_2: sql.NullString{String: user.Email, Valid: true},
			})
			if err != nil {
				log.Println(fmt.Sprintf("Error occured updating timezone information: %s", err))

			}
			log.Println("Successfully updated user details")
		}
	}()
	resetResponse := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	responsebytes, err := json.MarshalIndent(resetResponse, "", "")
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responsebytes)
	return
}

// DeleteUser is used to disable a users account
func (env *Env) DeleteUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Delete user Request received")

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

	user, err := env.AuthDb.GetUser(context.Background(), sql.NullString{String: strings.ToLower(username), Valid: true})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "User does not exist"
		log.Println(fmt.Sprintf("Error fetching user: %s", err))
		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)
		return
	}
	log.Println(fmt.Sprintf("User %s exists...", user.Username.String))
	go func() {
		err = env.AuthDb.DeleteUser(context.Background(), user.Email)
		if err != nil {
			log.Println(fmt.Sprintf("Error occured while deleting user"))
		}
		log.Println("Successfully deactivated user")
	}()
	deleteResponse := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	responsebytes, err := json.MarshalIndent(deleteResponse, "", "")
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responsebytes)
	return
}

func getValue(request, user string) string {
	if request == "" {
		return user
	}
	return request
}
