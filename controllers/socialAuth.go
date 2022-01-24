package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"helloprofile.com/models"
	"helloprofile.com/persistence/orm/helloprofiledb"
	"helloprofile.com/util"
)

func (env Env) GoogleLoginHandler(c echo.Context) (err error) {

	fields := log.Fields{"microservice": "helloprofile.service", "application": c.Param("application")}
	log.WithFields(fields).Info("Login Request received")

	errorResponse := new(models.Errormessage)
	request := new(models.GoogleJWT)
	if err = c.Bind(request); err != nil {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to marshal request")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	// Validate the JWT is valid
	claims, err := util.ValidateGoogleJWT(request.GoogleJWT, strings.ToLower(c.Request().Header.Get("os")))
	if err != nil {

		errorResponse.Errorcode = util.UNAUTHORIZED_ERROR_CODE
		errorResponse.ErrorMessage = util.UNAUTHORIZED_ERROR_MESSAGE_WRONG_JWT
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to marshal request")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	user, err := env.HelloProfileDb.GetUser(context.Background(), sql.NullString{String: strings.ToLower(claims.Email), Valid: true})
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			//Create user
			user, err := env.HelloProfileDb.CreateUser(context.Background(), helloprofiledb.CreateUserParams{
				Firstname:        sql.NullString{String: claims.FirstName, Valid: true},
				Lastname:         sql.NullString{String: claims.LastName, Valid: true},
				Username:         sql.NullString{String: claims.Email, Valid: true},
				Email:            claims.Email,
				IsEmailConfirmed: claims.EmailVerified,
				IsActive:         true,
			})
			if err != nil {
				errorResponse.Errorcode = util.USER_ALREADY_EXISTS_ERROR_CODE
				errorResponse.ErrorMessage = util.USER_ALREADY_EXISTS_ERROR_MESSAGE
				log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured creating user")
				c.JSON(http.StatusAlreadyReported, errorResponse)
				return err
			}
			go func() {
				err = env.saveLogin(helloprofiledb.CreateUserLoginParams{
					UserID:              user.ID,
					ResponseCode:        sql.NullString{String: util.SUCCESS_RESPONSE_CODE, Valid: true},
					ResponseDescription: sql.NullString{String: util.REGISTRATION_SUCCESS_RESPONSE_MESSAGE, Valid: true},
					LoginStatus:         true,
					IpAddress:           sql.NullString{String: c.RealIP(), Valid: true},
					Device:              sql.NullString{String: c.Request().UserAgent()},
				})
				if err != nil {
					log.WithFields(fields).WithError(err).Error("Error occured while saving login...")
				}
			}()
			registerResponse := &models.SuccessResponse{
				ResponseCode:    util.SUCCESS_RESPONSE_CODE,
				ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
				ResponseDetails: &models.UserDetail{
					CreatedAt:                 user.CreatedAt,
					Email:                     user.Email,
					Firstname:                 user.Firstname.String,
					ProfilePicture:            user.ImageUrl.String,
					IsActive:                  user.IsActive,
					IsEmailConfirmed:          user.IsEmailConfirmed,
					IsLockedOut:               user.IsLockedOut,
					IsPasswordSystemGenerated: user.IsPasswordSystemGenerated,
					Lastname:                  user.Lastname.String,
					Username:                  user.Username.String,
					Phone:                     user.Phone.String,
				},
			}

			// log.Println(fmt.Sprintf("Got to response string: %s", responseString))
			log.WithFields(fields).Info("Generating authentication token...")
			role := c.Request().Header.Get("Role")
			dbRole, err := env.HelloProfileDb.GetRole(context.Background(), strings.ToLower(role))
			if err != nil {
				log.WithFields(fields).WithError(err).Error(`Invalid role entered... Changing to default role of "Guest"`)
				role = "Guest"
			} else {
				log.WithFields(fields).Info(fmt.Sprintf("Creating token for user: %s | role: %s", user.Email, dbRole.Name))
			}
			go func() {
				log.WithFields(fields).Info("Adding user to role...")
				userRole, err := env.HelloProfileDb.AddUserRole(context.Background(), helloprofiledb.AddUserRoleParams{
					Name:     strings.ToLower(role),
					Username: user.Username,
				})
				if err != nil {
					log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error(fmt.Sprintf("Error occured adding user: %s to role: %s", user.Username.String, role))
				}
				log.WithFields(fields).Info(fmt.Sprintf("Successfully added user to role.. User Role Id: %s", userRole.ID))
			}()
			authToken, refreshToken, err := util.GenerateJWT(user.Email, strings.Split(role, ":"))
			if err != nil {
				errorResponse.Errorcode = util.OPERATION_FAILED_ERROR_CODE
				errorResponse.ErrorMessage = util.OPERATION_FAILED_ERROR_MESSAGE
				log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured generating auth token")
				c.JSON(http.StatusBadRequest, errorResponse)
				return err
			}
			log.WithFields(fields).Info("Storing refresh token in separate thread...")
			// store refresh token add later
			go func() {
				dbRefreshToken, err := env.HelloProfileDb.CreateRefreshToken(context.Background(), helloprofiledb.CreateRefreshTokenParams{
					UserID: user.ID,
					Token:  refreshToken,
				})
				if err != nil {
					log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to save refresh token")
				}

				log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error(fmt.Sprintf("Refresh Token Id: %s", dbRefreshToken.ID))
			}()
			c.Response().Header().Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
			c.Response().Header().Set("Refresh-Token", refreshToken)
			c.Response().Header().Set("Role", role)
			// log.Println(fmt.Sprintf("Auth token: %s, Refresh token: %s, Return object: %v", authToken, refreshToken, registerResponse))
			c.JSON(http.StatusOK, registerResponse)

			log.WithFields(fields).Info("Successfully processed registration request")
			return err
		} else {
			errorResponse.Errorcode = util.USER_NAME_OR_PASSWORD_INCORRECT_ERROR_CODE
			errorResponse.ErrorMessage = util.USER_NAME_OR_PASSWORD_INCORRECT_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error fetching user")
			c.JSON(http.StatusUnauthorized, errorResponse)
			return err
		}
	}
	if user.IsLockedOut {
		errorResponse.Errorcode = util.ACCOUNT_LOCKOUT_ERROR_CODE
		errorResponse.ErrorMessage = util.ACCOUNT_LOCKOUT_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Account was locked out....")
		c.JSON(http.StatusUnauthorized, errorResponse)
		return err
	}
	log.WithFields(fields).Info("Verifying that user is in the role access is being requested...")

	userRoles, err := env.HelloProfileDb.GetUserRoles(context.Background(), sql.NullString{String: user.Email, Valid: true})
	if err != nil {
		log.WithFields(fields).WithError(err).Error(`Invalid role entered... Changing to default role of "Guest"`)
		userRoles[0] = "guest"
	}

	log.WithFields(fields).Info(fmt.Sprintf("Generating authentication token for user: %s role: %v...", user.Email, userRoles))
	authToken, refreshToken, err := util.GenerateJWT(user.Email, userRoles)
	if err != nil {
		errorResponse.Errorcode = util.OPERATION_FAILED_ERROR_CODE
		errorResponse.ErrorMessage = util.OPERATION_FAILED_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured generating auth token")
		c.JSON(http.StatusConflict, errorResponse)
		return err
	}
	log.WithFields(fields).Info("Storing refresh token in separate thread...")
	// store refresh token add later
	go func() {
		dbRefreshToken, err := env.HelloProfileDb.CreateRefreshToken(context.Background(), helloprofiledb.CreateRefreshTokenParams{
			UserID: user.ID,
			Token:  refreshToken,
		})
		if err != nil {
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while saving refresh token")
		}

		log.WithFields(fields).Info(fmt.Sprintf("Refresh Token Id: %s", dbRefreshToken.ID))
	}()
	go func() {
		err = env.saveLogin(helloprofiledb.CreateUserLoginParams{
			UserID:              user.ID,
			ResponseCode:        sql.NullString{String: util.SUCCESS_RESPONSE_CODE, Valid: true},
			ResponseDescription: sql.NullString{String: util.SUCCESS_RESPONSE_MESSAGE, Valid: true},
			LoginStatus:         true,
			IpAddress:           sql.NullString{String: c.RealIP(), Valid: true},
			Device:              sql.NullString{String: c.Request().UserAgent()},
		})
		if err != nil {
			log.WithFields(fields).Info("Successful login...")
		}
		err := env.HelloProfileDb.UpdateResolvedLogin(context.Background(), user.ID)
		if err != nil {
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured clearing failed user logins after successful login...")
		}
	}()
	// address := make(chan models.Address)
	go func() {
		log.WithFields(fields).Info(`Getting the primary address for the user`)

	}()
	loginResponse := &models.SuccessResponse{
		ResponseCode:    util.SUCCESS_RESPONSE_CODE,
		ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
		ResponseDetails: &models.UserDetail{
			CreatedAt:                 user.CreatedAt,
			Email:                     user.Email,
			Firstname:                 user.Firstname.String,
			ProfilePicture:            user.ProfilePicture.String,
			IsActive:                  user.IsActive,
			IsEmailConfirmed:          user.IsEmailConfirmed,
			IsLockedOut:               user.IsLockedOut,
			IsPasswordSystemGenerated: user.IsPasswordSystemGenerated,
			Lastname:                  user.Lastname.String,
			Username:                  user.Username.String,
			Phone:                     user.Phone.String,
		},
	}
	c.Response().Header().Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
	c.Response().Header().Set("Refresh-Token", refreshToken)
	c.Response().Header().Set("Role", strings.Join(userRoles, ":"))
	c.JSON(http.StatusOK, loginResponse)
	return err
}
