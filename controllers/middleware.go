package controllers

import (
	"authengine/models"
	"authengine/util"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// CheckApplication checks if the application passed is valid
func (env *Env) CheckApplication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Checking application")
		pathParams := mux.Vars(r)
		w.Header().Set("Content-Type", "application/json")
		// file, fileHeader, err := r.FormFile("request.AttachmentName[i]")

		// file, err := os.Create(fmt.Sprintf("%s%s", attachmentPath, request.AttachmentName[i].FileName))
		// file.WriteString()

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
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

// Authorize is used to check if requests are authorized
func Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Checking authorization...")
		var errorResponse models.Errormessage
		var err error
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

		verifiedClaims, err := util.VerifyToken(authCode)

		if err != nil || verifiedClaims.Email == "" {
			errorResponse.Errorcode = "09"
			errorResponse.ErrorMessage = "Session expired. Kindly login again"
			log.Println("Token has expired...")
			response, err := json.MarshalIndent(errorResponse, "", "")
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(response)
			return
		}
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

// AuthorizeAdmin is used to check if requests are authorized
func AuthorizeAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Checking admin authorization...")
		var errorResponse models.Errormessage
		var err error
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

		verifiedClaims, err := util.VerifyToken(authCode)

		if err != nil || verifiedClaims.Email == "" {
			errorResponse.Errorcode = "09"
			errorResponse.ErrorMessage = "Session expired. Kindly login again..."
			log.Println("Token has expired...")
			response, err := json.MarshalIndent(errorResponse, "", "")
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(response)
			return
		}
		if !(strings.ToLower(verifiedClaims.Role) == "admin" || strings.ToLower(verifiedClaims.Role) == "superadmin") {
			errorResponse.Errorcode = "09"
			errorResponse.ErrorMessage = "Sorry, you are not authorized to carry out this operation."
			log.Println(fmt.Sprintf("User is not authorised to perform this operation with role %s...", verifiedClaims.Role))
			response, err := json.MarshalIndent(errorResponse, "", "")
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(response)
			return
		}
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

// TrackResponseTime is used to track the response time of api calls
func TrackResponseTime(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Measure response time
		start := time.Now()
		next.ServeHTTP(w, r)
		responseTime := time.Since(start)

		// Write it to the log
		log.Println(fmt.Sprintf("Request executed in %v", responseTime))

		// Make sure to pass the error back!

	})
}
