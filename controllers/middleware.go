package controllers

import (
	"authengine/models"
	"authengine/util"
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

// Authorize is used to check if requests are authorized
func Authorize(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Println("Checking authorization...")
		errorResponse := new(models.Errormessage)

		var authCode string
		authArray := strings.Split(c.Request().Header.Get("Authorization"), " ")
		if len(authArray) != 2 {
			errorResponse.Errorcode = "11"
			errorResponse.ErrorMessage = "Unsupported authentication scheme type"
			log.Println("Unsupported authentication scheme type")
			c.JSON(http.StatusUnauthorized, errorResponse)
			return nil
		}
		authCode = authArray[1]

		verifiedClaims, err := util.VerifyToken(authCode)

		if err != nil || verifiedClaims.Email == "" {
			errorResponse.Errorcode = "09"
			errorResponse.ErrorMessage = "Session expired. Kindly login again"
			log.Println("Token has expired...")
			log.Println(err)
			c.JSON(http.StatusUnauthorized, errorResponse)
			return nil
		}
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		if err := next(c); err != nil {
			c.Error(err)
		}
		return nil
	}
}

// AuthorizeAdmin is used to check if requests are authorized
func AuthorizeAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Println("Checking admin authorization...")
		errorResponse := new(models.Errormessage)

		var authCode string
		authArray := strings.Split(c.Request().Header.Get("Authorization"), " ")
		if len(authArray) != 2 {
			errorResponse.Errorcode = "11"
			errorResponse.ErrorMessage = "Unsupported authentication scheme type"
			log.Println("Unsupported authentication scheme type")

			c.JSON(http.StatusUnauthorized, errorResponse)
			return nil
		}
		authCode = authArray[1]

		verifiedClaims, err := util.VerifyToken(authCode)

		if err != nil || verifiedClaims.Email == "" {
			errorResponse.Errorcode = "09"
			errorResponse.ErrorMessage = "Session expired. Kindly login again..."
			log.Println("Token has expired...")
			log.Println(err)
			c.JSON(http.StatusUnauthorized, errorResponse)
			return err
		}
		if !(strings.ToLower(verifiedClaims.Role) == "admin" || strings.ToLower(verifiedClaims.Role) == "superadmin") {
			errorResponse.Errorcode = "09"
			errorResponse.ErrorMessage = "Sorry, you are not authorized to carry out this operation."
			log.Println(fmt.Sprintf("User is not authorised to perform this operation with role %s...", verifiedClaims.Role))
			c.JSON(http.StatusUnauthorized, errorResponse)
			return nil
		}
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		if err := next(c); err != nil {
			c.Error(err)
		}
		return nil
	}
}
