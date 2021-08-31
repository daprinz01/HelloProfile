package controllers

import (
	"authengine/models"
	"authengine/util"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
)

// CheckApplication checks if the application passed is valid
func (env *Env) CheckApplication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		log.WithField("microservice", "persian.black.authengine.service").Info("Checking application")

		errorResponse := new(models.Errormessage)

		application := c.Param("application")
		if application == "" {
			errorResponse.Errorcode = util.APPLICATION_NOT_SPECIFIED_ERROR_CODE
			errorResponse.ErrorMessage = util.APPLICATION_NOT_SPECIFIED_ERROR_MESSAGE
			log.WithField("microservice", "persian.black.authengine.service").WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error(util.APPLICATION_NOT_SPECIFIED_ERROR_MESSAGE)
			c.JSON(http.StatusBadRequest, errorResponse)
			return nil
		}
		fields := log.Fields{"microservice": "persian.black.authengine.service", "application": application}
		log.WithFields(fields).Info(fmt.Sprintf("Application: %s", application))
		applicationObject, err := env.AuthDb.GetApplication(context.Background(), strings.ToLower(application))
		if err != nil {
			errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
			errorResponse.ErrorMessage = util.INVALID_APPLICATION_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error(util.INVALID_APPLICATION_ERROR_MESSAGE)
			c.JSON(http.StatusNotFound, errorResponse)
			return err
		}
		log.WithFields(fields).Info(fmt.Sprintf("Applicaiton ID: %s", applicationObject.ID))
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
		log.WithField("microservice", "persian.black.authengine.service").Info(fmt.Sprintf("Request executed in %v", responseTime))
		return nil
	}

}

// Authorize is used to check if requests are authorized
func Authorize(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.WithField("microservice", "persian.black.authengine.service").Info("Checking authorization...")
		errorResponse := new(models.Errormessage)

		var authCode string
		authArray := strings.Split(c.Request().Header.Get("Authorization"), " ")
		if len(authArray) != 2 {
			errorResponse.Errorcode = util.INVALID_AUTHENTICATION_SCHEME_ERROR_CODE
			errorResponse.ErrorMessage = util.INVALID_AUTHENTICATION_SCHEME_ERROR_MESSAGE
			log.WithField("microservice", "persian.black.authengine.service").WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error(util.INVALID_AUTHENTICATION_SCHEME_ERROR_MESSAGE)
			c.JSON(http.StatusUnauthorized, errorResponse)
			return nil
		}
		authCode = authArray[1]

		verifiedClaims, err := util.VerifyToken(authCode)

		if err != nil || verifiedClaims.Email == "" {
			errorResponse.Errorcode = util.SESSION_EXPIRED_ERROR_CODE
			errorResponse.ErrorMessage = util.SESSION_EXPIRED_ERROR_MESSAGE
			log.WithField("microservice", "persian.black.authengine.service").WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).WithError(err).Error("Token has expired...")
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
		log.WithField("microservice", "persian.black.authengine.service").Info("Checking admin authorization...")
		errorResponse := new(models.Errormessage)

		var authCode string
		authArray := strings.Split(c.Request().Header.Get("Authorization"), " ")
		if len(authArray) != 2 {
			errorResponse.Errorcode = util.INVALID_AUTHENTICATION_SCHEME_ERROR_CODE
			errorResponse.ErrorMessage = util.INVALID_AUTHENTICATION_SCHEME_ERROR_MESSAGE
			log.WithField("microservice", "persian.black.authengine.service").WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error(util.INVALID_AUTHENTICATION_SCHEME_ERROR_MESSAGE)

			c.JSON(http.StatusUnauthorized, errorResponse)
			return nil
		}
		authCode = authArray[1]

		verifiedClaims, err := util.VerifyToken(authCode)

		if err != nil || verifiedClaims.Email == "" {
			errorResponse.Errorcode = util.SESSION_EXPIRED_ERROR_CODE
			errorResponse.ErrorMessage = util.SESSION_EXPIRED_ERROR_MESSAGE
			log.WithField("microservice", "persian.black.authengine.service").WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Token has expired...")
			c.JSON(http.StatusUnauthorized, errorResponse)
			return err
		}
		if !(strings.Contains(strings.ToLower(verifiedClaims.Role), "admin") || strings.Contains(strings.ToLower(verifiedClaims.Role), "superadmin")) {
			errorResponse.Errorcode = util.UNAUTHORIZED_ERROR_CODE
			errorResponse.ErrorMessage = util.UNAUTHORIZED_ERROR_MESSAGE
			log.WithField("microservice", "persian.black.authengine.service").WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error(fmt.Sprintf("User is not authorised to perform this operation with role(s) %s...", verifiedClaims.Role))
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
