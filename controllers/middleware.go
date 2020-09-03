package controllers

import (
	"authengine/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

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
