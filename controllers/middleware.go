package controllers

import (
	"authengine/models"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// CheckApplication checks if the application passed is valid
func (env *Env) CheckApplication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Println("Checking application")

		errorResponse := new(models.Errormessage)

		application := c.Param("application")
		if application == "" {
			errorResponse.Errorcode = "01"
			errorResponse.ErrorMessage = "Application not specified"
			log.Println("Application not specified")
			c.JSON(http.StatusBadRequest, errorResponse)
			return nil
		}
		log.Println(fmt.Sprintf("Application: %s", application))
		applicationObject, err := env.AuthDb.GetApplication(context.Background(), strings.ToLower(application))
		if err != nil {
			errorResponse.Errorcode = "06"
			errorResponse.ErrorMessage = "Application is invalid"
			log.Println(err)
			c.JSON(http.StatusBadRequest, errorResponse)
			return err
		}
		log.Println(fmt.Sprintf("Applicaiton ID: %d", applicationObject.ID))
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		if err := next(c); err != nil {
			c.Error(err)
		}
		return nil
	}
}

// TrackResponseTime is used to track the response time of api calls
func TrackResponseTime(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Measure response time
		start := time.Now()
		if err := next(c); err != nil {
			c.Error(err)
		}
		responseTime := time.Since(start)

		// Write it to the log
		log.Println(fmt.Sprintf("Request executed in %v", responseTime))
		return nil
	}

}

// // Authorize is used to check if requests are authorized
// func Authorize(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		log.Println("Checking authorization...")
// 		errorResponse := new(models.Errormessage)
//
// 		var authCode string
// 		authArray := strings.Split(r.Header.Get("Authorization"), " ")
// 		if len(authArray) != 2 {
// 			errorResponse.Errorcode = "11"
// 			errorResponse.ErrorMessage = "Unsupported authentication scheme type"
// 			log.Println("Unsupported authentication scheme type")
// 			response, err := json.MarshalIndent(errorResponse, "", "")
// 			if err != nil {
// 				log.Println(err)
// 			}
// 			w.WriteHeader(http.StatusBadRequest)
// 			w.Write(response)
// 			return
// 		}
// 		authCode = authArray[1]

// 		verifiedClaims, err := util.VerifyToken(authCode)

// 		if err != nil || verifiedClaims.Email == "" {
// 			errorResponse.Errorcode = "09"
// 			errorResponse.ErrorMessage = "Session expired. Kindly login again"
// 			log.Println("Token has expired...")
// 			response, err := json.MarshalIndent(errorResponse, "", "")
// 			if err != nil {
// 				log.Println(err)
// 			}
// 			w.WriteHeader(http.StatusUnauthorized)
// 			w.Write(response)
// 			return
// 		}
// 		// Call the next handler, which can be another middleware in the chain, or the final handler.
// 		next.ServeHTTP(w, r)
// 	})
// }

// // AuthorizeAdmin is used to check if requests are authorized
// func AuthorizeAdmin(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		log.Println("Checking admin authorization...")
// 		errorResponse := new(models.Errormessage)
//
// 		var authCode string
// 		authArray := strings.Split(r.Header.Get("Authorization"), " ")
// 		if len(authArray) != 2 {
// 			errorResponse.Errorcode = "11"
// 			errorResponse.ErrorMessage = "Unsupported authentication scheme type"
// 			log.Println("Unsupported authentication scheme type")
// 			response, err := json.MarshalIndent(errorResponse, "", "")
// 			if err != nil {
// 				log.Println(err)
// 			}
// 			w.WriteHeader(http.StatusBadRequest)
// 			w.Write(response)
// 			return
// 		}
// 		authCode = authArray[1]

// 		verifiedClaims, err := util.VerifyToken(authCode)

// 		if err != nil || verifiedClaims.Email == "" {
// 			errorResponse.Errorcode = "09"
// 			errorResponse.ErrorMessage = "Session expired. Kindly login again..."
// 			log.Println("Token has expired...")
// 			response, err := json.MarshalIndent(errorResponse, "", "")
// 			if err != nil {
// 				log.Println(err)
// 			}
// 			w.WriteHeader(http.StatusUnauthorized)
// 			w.Write(response)
// 			return
// 		}
// 		if !(strings.ToLower(verifiedClaims.Role) == "admin" || strings.ToLower(verifiedClaims.Role) == "superadmin") {
// 			errorResponse.Errorcode = "09"
// 			errorResponse.ErrorMessage = "Sorry, you are not authorized to carry out this operation."
// 			log.Println(fmt.Sprintf("User is not authorised to perform this operation with role %s...", verifiedClaims.Role))
// 			response, err := json.MarshalIndent(errorResponse, "", "")
// 			if err != nil {
// 				log.Println(err)
// 			}
// 			w.WriteHeader(http.StatusUnauthorized)
// 			w.Write(response)
// 			return
// 		}
// 		// Call the next handler, which can be another middleware in the chain, or the final handler.
// 		next.ServeHTTP(w, r)
// 	})
// }
