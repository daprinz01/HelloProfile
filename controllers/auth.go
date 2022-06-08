package controllers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"helloprofile.com/models"
	"helloprofile.com/persistence/orm/helloprofiledb"
	"helloprofile.com/util"

	"github.com/labstack/echo/v4"
)

//Login is used to sign users in
func (env *Env) Login(c echo.Context) (err error) {
	fields := log.Fields{"microservice": "helloprofile.service", "application": c.Param("application")}
	log.WithFields(fields).Info("Login Request received")

	errorResponse := new(models.Errormessage)

	request := new(models.LoginRequest)
	if err = c.Bind(request); err != nil {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to marshal request")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	user, err := env.HelloProfileDb.GetUser(context.Background(), sql.NullString{String: strings.ToLower(request.Username), Valid: true})
	if err != nil {
		errorResponse.Errorcode = util.USER_NAME_OR_PASSWORD_INCORRECT_ERROR_CODE
		errorResponse.ErrorMessage = util.USER_NAME_OR_PASSWORD_INCORRECT_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error fetching user")
		c.JSON(http.StatusUnauthorized, errorResponse)
		return err
	}
	if user.IsLockedOut {
		errorResponse.Errorcode = util.ACCOUNT_LOCKOUT_ERROR_CODE
		errorResponse.ErrorMessage = util.ACCOUNT_LOCKOUT_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Account was locked out....")
		c.JSON(http.StatusUnauthorized, errorResponse)
		return err
	}
	if util.VerifyHash(user.Password.String, request.Password) {
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
		profiles := make(chan []models.Profile)
		go env.getProfiles(user.ID, profiles, fields)
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
		//Add saved profiles to the user
		go env.saveProfile(user, fields)

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
				Country:                   user.City.String,
				City:                      user.City.String,
				Profiles:                  <-profiles,
			},
		}
		c.Response().Header().Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
		c.Response().Header().Set("Refresh-Token", refreshToken)
		c.Response().Header().Set("Role", strings.Join(userRoles, ":"))
		c.JSON(http.StatusOK, loginResponse)
		return err
	} else {
		errorResponse.Errorcode = util.USER_NAME_OR_PASSWORD_INCORRECT_ERROR_CODE
		errorResponse.ErrorMessage = util.USER_NAME_OR_PASSWORD_INCORRECT_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Password incorrect...")
		go func() {
			err = env.saveLogin(helloprofiledb.CreateUserLoginParams{
				UserID:              user.ID,
				ResponseCode:        sql.NullString{String: errorResponse.Errorcode, Valid: true},
				ResponseDescription: sql.NullString{String: errorResponse.ErrorMessage, Valid: true},
				LoginStatus:         false,
				Resolved:            false,
				IpAddress:           sql.NullString{String: c.RealIP(), Valid: true},
				Device:              sql.NullString{String: c.Request().UserAgent()},
			})
			if err != nil {
				log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Password incorrect...")
			}
		}()

		userLogins, err := env.HelloProfileDb.GetUnResoledLogins(context.Background(), user.ID)
		if err != nil {
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error ocurred fetching user unresolved logins")
		}
		var lockoutCount int
		lockOutCountKey := os.Getenv("LOCK_OUT_COUNT")
		if lockOutCountKey == "" {
			log.WithFields(fields).Error(`LOCK_OUT_COUNT cannot be empty`)
			log.WithFields(fields).Info(`LOCK_OUT_COUNT cannot be empty, setting default of 5...`)
		} else {
			log.WithFields(fields).Info(`Setting lock out count...`)
			lockoutCount, err = strconv.Atoi(lockOutCountKey)
			if err != nil {
				log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while converting lock out count")
			}
		}
		// Check if account has exceeded the lockout count
		if len(userLogins) >= lockoutCount {
			lockoutUpdate, err := env.HelloProfileDb.UpdateUser(context.Background(), helloprofiledb.UpdateUserParams{
				Username_2:                user.Username,
				IsLockedOut:               true,
				Email:                     user.Email,
				Firstname:                 user.Firstname,
				ImageUrl:                  user.ProfilePicture,
				IsActive:                  user.IsActive,
				IsEmailConfirmed:          user.IsEmailConfirmed,
				IsPasswordSystemGenerated: user.IsPasswordSystemGenerated,
				Lastname:                  user.Lastname,
				Password:                  user.Password,
				Username:                  user.Username,
				Phone:                     user.Phone,
			})
			if err != nil {
				log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to lockout account")
			}
			log.WithFields(fields).Info(fmt.Sprintf(`Account with username: %s has been locked out`, lockoutUpdate.Username.String))

			errorResponse.Errorcode = util.ACCOUNT_LOCKOUT_ERROR_CODE
			errorResponse.ErrorMessage = util.ACCOUNT_LOCKOUT_ERROR_MESSAGE
			c.JSON(http.StatusUnauthorized, errorResponse)
			return err
		}

		c.JSON(http.StatusUnauthorized, errorResponse)
		return err
	}
}

func (env *Env) saveProfile(user helloprofiledb.UserDetail, fields log.Fields) {
	savedProfiles, err := env.HelloProfileDb.GetSavedProfilesByEmail(context.Background(), helloprofiledb.GetSavedProfilesByEmailParams{
		Email:   strings.ToLower(user.Email),
		IsAdded: false,
	})
	if err != nil {
		log.WithFields(fields).WithError(err).Error("Error occured fetching saved profiles")
		return
	}
	var addedContactCount int
	for _, savedProfile := range savedProfiles {
		_, err := env.HelloProfileDb.AddContacts(context.Background(), helloprofiledb.AddContactsParams{
			UserID:    user.ID,
			ProfileID: savedProfile.ProfileID,
		})
		if err != nil {
			log.WithFields(fields).WithError(err).Error("Error occured while adding saved profile with ID: ", savedProfile.ID)
			continue
		}
		_, err = env.HelloProfileDb.UpdateSavedProfile(context.Background(), helloprofiledb.UpdateSavedProfileParams{
			FirstName: savedProfile.FirstName,
			LastName:  savedProfile.LastName,
			IsAdded:   true,
		})
		if err != nil {
			log.WithFields(fields).WithError(err).Error("Error occured while updating saved profile with ID: ", savedProfile.ID)
			_ = env.HelloProfileDb.DeleteContact(context.Background(), helloprofiledb.DeleteContactParams{
				UserID:    user.ID,
				ProfileID: savedProfile.ProfileID,
			})
		}
		addedContactCount++
	}
	log.WithFields(fields).WithError(err).Error("Added %d saved contacts to user: %s", addedContactCount, user.Email)
}

// Register is used to register new users
func (env *Env) Register(c echo.Context) (err error) {
	fields := log.Fields{"microservice": "helloprofile.service", "application": c.Param("application")}
	log.WithFields(fields).Info("Register Request received")

	errorResponse := new(models.Errormessage)

	request := new(models.UserDetail)
	if err = c.Bind(request); err != nil {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to marshal request")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	var hashedPassword string

	if request.Password != "" {
		hashedPassword = util.GenerateHash(request.Password)

	}
	log.WithFields(fields).Info("Successfully hashed password")
	log.WithFields(fields).Info("Inserting user...")
	user, err := env.HelloProfileDb.CreateUser(context.Background(), helloprofiledb.CreateUserParams{
		CreatedAt:                 time.Now(),
		Email:                     strings.ToLower(request.Email),
		Firstname:                 sql.NullString{String: request.Firstname, Valid: request.Firstname != ""},
		ImageUrl:                  sql.NullString{String: request.ProfilePicture, Valid: request.ProfilePicture != ""},
		IsActive:                  true,
		IsEmailConfirmed:          request.IsEmailConfirmed,
		IsLockedOut:               request.IsLockedOut,
		IsPasswordSystemGenerated: request.IsPasswordSystemGenerated,
		Lastname:                  sql.NullString{String: request.Lastname, Valid: request.Lastname != ""},
		Password:                  sql.NullString{String: hashedPassword, Valid: hashedPassword != ""},
		Username:                  sql.NullString{String: strings.ToLower(request.Username), Valid: request.Username != ""},
		Phone:                     sql.NullString{String: request.Phone, Valid: request.Phone != ""},
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
	//Add saved profiles to the user
	go env.saveProfile(helloprofiledb.UserDetail{
		ID:    user.ID,
		Email: user.Email,
	}, fields)
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
	return
}

// RefreshToken is used to register expired token
func (env *Env) RefreshToken(c echo.Context) (err error) {
	fields := log.Fields{"microservice": "helloprofile.service", "application": c.Param("application")}
	log.WithFields(fields).Info("Refresh token Request received")

	errorResponse := new(models.Errormessage)

	var authCode string
	authArray := strings.Split(c.Request().Header.Get("Authorization"), " ")
	if len(authArray) != 2 {
		errorResponse.Errorcode = util.INVALID_AUTHENTICATION_SCHEME_ERROR_CODE
		errorResponse.ErrorMessage = util.INVALID_AUTHENTICATION_SCHEME_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error(util.INVALID_AUTHENTICATION_SCHEME_ERROR_MESSAGE)
		c.JSON(http.StatusUnauthorized, errorResponse)
		return err
	}
	authCode = authArray[1]
	refreshToken := c.Request().Header.Get("Refresh-Token")

	verifiedClaims, err := util.VerifyToken(authCode)
	if err == nil {
		errorResponse.Errorcode = util.SESSION_STILL_ACTIVE_ERROR_CODE
		errorResponse.ErrorMessage = util.SESSION_STILL_ACTIVE_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Token is still valid...")
		c.JSON(http.StatusTooEarly, errorResponse)
		return err
	}
	if err != nil && verifiedClaims.Email == "" {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error(util.MODEL_VALIDATION_ERROR_MESSAGE)
		c.JSON(http.StatusUnauthorized, errorResponse)
		return err
	}

	dbRefreshToken, err := env.HelloProfileDb.GetRefreshToken(context.Background(), refreshToken)
	if err != nil {
		errorResponse.Errorcode = util.SESSION_EXPIRED_ERROR_CODE
		errorResponse.ErrorMessage = util.SESSION_EXPIRED_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured refreshing token")
		c.JSON(http.StatusUnauthorized, errorResponse)
		return err
	}
	defer func() {
		err = env.HelloProfileDb.DeleteRefreshToken(context.Background(), refreshToken)
		if err != nil {
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to delete refresh token")
		}
		log.WithFields(fields).Info("Successfully deleted old refresh token...")
	}()
	var refreshTokenDuration time.Duration
	refreshTokenLifespan := os.Getenv("SESSION_LIFESPAN")
	if refreshTokenLifespan == "" {
		log.WithFields(fields).Error("Session lifespan cannot be empty")
		log.WithFields(fields).Info("SESSION_LIFESPAN cannot be empty, setting duration to default of 15 mins ...")
		refreshTokenDuration, err = time.ParseDuration("15m")
	} else {
		log.WithFields(fields).Info("Setting Refresh token lifespan...")
		refreshTokenDuration, err = time.ParseDuration(refreshTokenLifespan)
		if err != nil {
			log.WithFields(fields).WithError(err).Error("Error converting refresh token duration to number")
		}
	}
	if !dbRefreshToken.CreatedAt.Add(refreshTokenDuration).Before(time.Now()) {
		log.WithFields(fields).Info("Generating authentication token...")
		authToken, refreshToken, err := util.GenerateJWT(verifiedClaims.Email, strings.Split(verifiedClaims.Role, ":"))
		if err != nil {
			errorResponse.Errorcode = util.OPERATION_FAILED_ERROR_CODE
			errorResponse.ErrorMessage = util.OPERATION_FAILED_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured generating auth token")
			c.JSON(http.StatusBadRequest, errorResponse)
			return err
		}
		log.WithFields(fields).Info("Fetching user...")
		user, err := env.HelloProfileDb.GetUser(context.Background(), sql.NullString{String: strings.ToLower(verifiedClaims.Email), Valid: true})
		if err != nil {
			errorResponse.Errorcode = util.USER_NAME_OR_PASSWORD_INCORRECT_ERROR_CODE
			errorResponse.ErrorMessage = util.USER_NAME_OR_PASSWORD_INCORRECT_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error fetching user")
			c.JSON(http.StatusUnauthorized, errorResponse)
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
				log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to create refresh token")
			}

			log.WithFields(fields).Info(fmt.Sprintf("Refresh Token Id: %s", dbRefreshToken.ID))
		}()
		resetResponse := &models.RefreshResponse{
			ResponseCode:    util.SUCCESS_RESPONSE_CODE,
			ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
		}

		c.Response().Header().Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
		c.Response().Header().Set("Refresh-Token", refreshToken)
		c.Response().Header().Set("Role", verifiedClaims.Role)
		c.JSON(http.StatusOK, resetResponse)

	} else {
		errorResponse.Errorcode = util.SESSION_EXPIRED_ERROR_CODE
		errorResponse.ErrorMessage = util.SESSION_EXPIRED_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Token has expired...")
		c.JSON(http.StatusUnauthorized, errorResponse)
		return err

	}
	return err

}

// SendOtp is used to send OTP request after validating user exist
func (env *Env) SendOtp(c echo.Context) (err error) {
	fields := log.Fields{"microservice": "helloprofile.service", "application": c.Param("application")}
	log.WithFields(fields).Info("Send OTP Request received")
	errorResponse := new(models.Errormessage)

	request := new(models.SendOtpRequest)
	if err = c.Bind(request); err != nil {
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to marshal request")
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	user, err := env.HelloProfileDb.GetUser(context.Background(), sql.NullString{String: strings.ToLower(request.Email), Valid: true})
	if err != nil {
		errorResponse.Errorcode = util.SUCCESS_RESPONSE_CODE
		errorResponse.ErrorMessage = util.OTP_SENT_RESPONSE_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error fetching user")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	log.WithFields(fields).Info(fmt.Sprintf("User %s exists... Generating OTP code for user", user.Username.String))
	var otpLength int
	otpLengthKey := os.Getenv("OTP_LENGTH")
	if otpLengthKey != "" {
		otpLength, err = strconv.Atoi(otpLengthKey)
		if err != nil {
			log.WithFields(fields).WithError(err).Error("Error occure converting otp lenght. Setting a default of 6")
		}
	}
	otp, err := util.GenerateOTP(otpLength)
	if err != nil {
		errorResponse.Errorcode = util.OPERATION_FAILED_ERROR_CODE
		errorResponse.ErrorMessage = util.OPERATION_FAILED_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured generating otp")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	// Save otp to db in another thread
	go func() {
		err = env.HelloProfileDb.CreateOtp(context.Background(), helloprofiledb.CreateOtpParams{
			OtpToken:         sql.NullString{String: otp, Valid: true},
			UserID:           user.ID,
			IsEmailPreferred: request.IsEmailPrefered,
			IsSmsPreferred:   !request.IsEmailPrefered,
			Purpose:          sql.NullString{String: request.Purpose, Valid: true},
		})
		if err != nil {
			log.WithFields(fields).WithError(err).Error("Error occured saving otp")
		}
		log.WithFields(fields).Info("Successfully saved OTP...")
		log.WithFields(fields).Info("Sending OTP through preferred channel...")
		communicationEndpoint := os.Getenv("COMMUNICATION_SERVICE_ENDPOINT")
		if request.IsEmailPrefered {
			emailPath := os.Getenv("EMAIL_PATH")
			emailRequest := models.SendEmailRequest{
				From:    models.EmailAddress{Email: os.Getenv("SMTP_USER"), Name: "HelloProfile"},
				To:      []models.EmailAddress{{Email: user.Email, Name: fmt.Sprintf("%s %s", user.Firstname.String, user.Lastname.String)}},
				Subject: fmt.Sprintf(util.SEND_OTP_EMAIL_SUBJECT, request.Purpose),
				Message: fmt.Sprintf(util.SEND_OTP_EMAIL_MESSAGE_BODY, user.Firstname.String, otp),
			}
			emailRequestBytes, _ := json.Marshal(emailRequest)
			emailRequestReader := bytes.NewReader(emailRequestBytes)
			log.WithFields(fields).Info("Sending otp email...")

			client := &http.Client{}
			req, _ := http.NewRequest("POST", fmt.Sprintf("%s%s", communicationEndpoint, emailPath), emailRequestReader)
			req.Header.Add("Client_Id", "persianblack")
			req.Header.Add("Accept", "application/json")
			req.Header.Add("Content-Type", "application/json")
			emailResponse, err := client.Do(req)

			// emailResponse, err := http.Post(fmt.Sprintf("%s%s", communicationEndpoint, emailPath), "application/json", bytes.NewBuffer(emailRequestBytes))
			if err != nil {
				log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured sending otp")
			} else {
				if emailResponse.StatusCode == 200 {
					log.WithFields(fields).Info("OTP sent successfully")
				} else {
					log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured sending OTP")
				}
				emailBody, _ := ioutil.ReadAll(emailResponse.Body)
				log.WithFields(fields).Info(fmt.Sprintf("Response body from email request: %s", emailBody))
			}

		} else {
			if user.Phone.String == "" {
				log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Phonenumber not available")
			} else {
				smsPath := os.Getenv("SMS_PATH")
				smsRequest := models.SendSmsRequest{
					Phone:   user.Phone.String,
					Message: fmt.Sprintf(util.SENT_OTP_SMS_MESSAGE, request.Purpose, otp),
				}
				smsRequestBytes, _ := json.Marshal(smsRequest)
				smsRequestReader := bytes.NewReader(smsRequestBytes)
				log.WithFields(fields).Info("Sending otp sms...")

				client := &http.Client{}
				req, _ := http.NewRequest("POST", fmt.Sprintf("%s%s", communicationEndpoint, smsPath), smsRequestReader)
				req.Header.Add("Client_Id", "persianblack")
				req.Header.Add("Accept", "application/json")
				req.Header.Add("Content-Type", "application/json")
				smsResponse, err := client.Do(req)

				if err != nil {
					log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured sending otp")
				}
				if smsResponse.StatusCode == 200 {
					log.WithFields(fields).Info("OTP sent successfully")
				} else {
					log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": smsResponse.StatusCode, "responseDescription": smsResponse.Status}).Error("Error occured sending OTP")
				}
				smsBody, _ := ioutil.ReadAll(smsResponse.Body)
				log.WithFields(fields).Info(fmt.Sprintf("Response body from sms request: %s", smsBody))
			}
		}
	}()
	log.WithFields(fields).Info("Successfully generated otp")
	resetResponse := &models.SuccessResponse{
		ResponseCode:    util.SUCCESS_RESPONSE_CODE,
		ResponseMessage: util.OTP_SENT_RESPONSE_MESSAGE,
	}
	c.JSON(http.StatusOK, resetResponse)
	return err
}

// SendOtp is used to send OTP request after validating user exist
func (env *Env) DoEmailVerification(c echo.Context) (err error) {
	fields := log.Fields{"microservice": "helloprofile.service", "application": c.Param("application")}
	log.WithFields(fields).Info("Generate OTP Request received")
	errorResponse := new(models.Errormessage)

	request := new(models.SendOtpRequest)
	if err = c.Bind(request); err != nil {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to marshal request")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	var otpLength int
	otpLengthKey := os.Getenv("OTP_LENGTH")
	if otpLengthKey != "" {
		otpLength, err = strconv.Atoi(otpLengthKey)
		if err != nil {
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured converting otp lenght. Setting a default of 6: %s")
		}
	}
	otp, err := util.GenerateOTP(otpLength)
	if err != nil {
		errorResponse.Errorcode = util.OPERATION_FAILED_ERROR_CODE
		errorResponse.ErrorMessage = util.OPERATION_FAILED_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured generating otp")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	log.WithFields(fields).Info("Successfully generated otp")
	//Save otp to db in another thread
	go func() {
		err = env.HelloProfileDb.CreateEmailVerification(context.Background(), helloprofiledb.CreateEmailVerificationParams{
			Otp:   otp,
			Email: sql.NullString{String: strings.ToLower(request.Email), Valid: true},
		})
		if err != nil {
			log.WithFields(fields).WithError(err).Error("Error occured saving otp")
		}
		log.WithFields(fields).Info("Successfully saved OTP...")
		log.WithFields(fields).Info("Sending OTP through preferred channel...")
		communicationEndpoint := os.Getenv("COMMUNICATION_SERVICE_ENDPOINT")

		emailPath := os.Getenv("EMAIL_PATH")
		emailRequest := models.SendEmailRequest{
			From:    models.EmailAddress{Email: os.Getenv("SMTP_USER"), Name: "HelloProfile"},
			To:      []models.EmailAddress{{Email: request.Email}},
			Subject: fmt.Sprintf(util.EMAIL_VERIFICATION_SUBJECT, strings.ToTitle(request.Application)),
			Message: fmt.Sprintf(util.EMAIL_VERIFICATION_MESSAGE, otp),
		}
		emailRequestBytes, _ := json.Marshal(emailRequest)
		emailRequestReader := bytes.NewReader(emailRequestBytes)
		log.WithFields(fields).Info("Sending email verification email...")

		client := &http.Client{}
		req, _ := http.NewRequest("POST", fmt.Sprintf("%s%s", communicationEndpoint, emailPath), emailRequestReader)
		req.Header.Add("Client_Id", "persianblack")
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
		emailResponse, err := client.Do(req)

		// emailResponse, err := http.Post(fmt.Sprintf("%s%s", communicationEndpoint, emailPath), "application/json", bytes.NewBuffer(emailRequestBytes))
		if err != nil {
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured sending otp")
		} else {
			if emailResponse.StatusCode == 200 {
				log.WithFields(fields).Info("OTP send successfully")
			} else {
				log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured sending OTP")
			}
			emailBody, _ := ioutil.ReadAll(emailResponse.Body)
			log.WithFields(fields).Info(fmt.Sprintf("Response body from email request: %s", emailBody))
		}
	}()
	resetResponse := &models.SuccessResponse{
		ResponseCode:    util.SUCCESS_RESPONSE_CODE,
		ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
	}
	c.JSON(http.StatusOK, resetResponse)
	return err
}

func (env *Env) VerifyEmailToken(c echo.Context) (err error) {
	fields := log.Fields{"microservice": "helloprofile.service", "application": c.Param("application")}
	log.WithFields(fields).Info("Verify Email token request received")

	errorResponse := new(models.Errormessage)

	request := new(models.VerifyOtpRequest)
	if err = c.Bind(request); err != nil {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to marshal request")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	dbOtp, err := env.HelloProfileDb.GetEmailVerification(context.Background(), request.OTP)
	if err != nil {
		errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
		errorResponse.ErrorMessage = util.NO_RECORD_FOUND_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error(util.MODEL_VALIDATION_ERROR_MESSAGE)
		c.JSON(http.StatusUnauthorized, errorResponse)
		return err
	}
	var otpDuration int
	otpDurationKey := os.Getenv("OTP_VALIDITY_PERIOD")
	if otpDurationKey != "" {
		otpDuration, err = strconv.Atoi(otpDurationKey)
		if err != nil {
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("OTP_VALIDITY_PERIOD is not a valid number")
			log.WithFields(fields).Info("Setting default of 30mins")
			otpDuration = 30
		}
	}
	if !dbOtp.CreatedAt.Add(time.Duration(otpDuration)*time.Minute).Before(time.Now()) && strings.EqualFold(strings.ToLower(request.Email), strings.ToLower(dbOtp.Email.String)) {
		log.WithFields(fields).Info("Email token verification successful")
		verifyResponse := &models.SuccessResponse{
			ResponseCode:    util.SUCCESS_RESPONSE_CODE,
			ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
		}
		c.JSON(http.StatusOK, verifyResponse)
	} else {
		errorResponse.Errorcode = util.SESSION_EXPIRED_ERROR_CODE
		errorResponse.ErrorMessage = util.EMAIL_TOKEN_EXPIRED_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Email otp has expired or is invalid...")
		c.JSON(http.StatusForbidden, errorResponse)

	}
	log.WithFields(fields).Info("Finished processing Verify otp request...")
	return err
}

// VerifyOtp is used to verify and otp against a user. Authentication token is generated that is used in subsequent requests.
func (env *Env) VerifyOtp(c echo.Context) (err error) {
	fields := log.Fields{"microservice": "helloprofile.service", "application": c.Param("application")}
	log.WithFields(fields).Info("Verify otp request received")

	// file, fileHeader, err := r.FormFile("request.AttachmentName[i]")

	// file, err := os.Create(fmt.Sprintf("%s%s", attachmentPath, request.AttachmentName[i].FileName))
	// file.WriteString()

	errorResponse := new(models.Errormessage)

	request := new(models.VerifyOtpRequest)
	if err = c.Bind(request); err != nil {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to marshal request")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	dbOtp, err := env.HelloProfileDb.GetOtp(context.Background(), helloprofiledb.GetOtpParams{
		OtpToken: sql.NullString{String: request.OTP, Valid: true},
		Username: sql.NullString{String: strings.ToLower(request.Email), Valid: true},
	})
	if err != nil {
		errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
		errorResponse.ErrorMessage = util.NO_RECORD_FOUND_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error(util.MODEL_VALIDATION_ERROR_MESSAGE)
		c.JSON(http.StatusUnauthorized, errorResponse)
		return err
	}
	var otpDuration int
	otpDurationKey := os.Getenv("OTP_VALIDITY_PERIOD")
	if otpDurationKey != "" {
		otpDuration, err = strconv.Atoi(otpDurationKey)
		if err != nil {
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("OTP_VALIDITY_PERIOD is not a valid number")
			log.WithFields(fields).Info("Setting default of 5mins")
			otpDuration = 5
		}
	}
	if !dbOtp.CreatedAt.Add(time.Duration(otpDuration) * time.Minute).Before(time.Now()) {
		log.WithFields(fields).Info("Verifying that user is in the role access is being requested...")

		userRoles, err := env.HelloProfileDb.GetUserRoles(context.Background(), sql.NullString{String: request.Email, Valid: true})
		if err != nil {
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error(`Invalid role entered... Changing to default role of "Guest"`)
			userRoles[0] = "Guest"
		}

		log.WithFields(fields).Info(fmt.Sprintf("Generating authentication token for user: %s role: %s...", request.Email, strings.Join(userRoles, ":")))
		authToken, refreshToken, err := util.GenerateJWT(request.Email, userRoles)
		if err != nil {
			errorResponse.Errorcode = util.OPERATION_FAILED_ERROR_CODE
			errorResponse.ErrorMessage = util.OPERATION_FAILED_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured generating wuth token")
			c.JSON(http.StatusBadRequest, errorResponse)
			return err
		}
		log.WithFields(fields).Info("Storing refresh token in separate thread...")
		// store refresh token add later
		go func() {
			dbRefreshToken, err := env.HelloProfileDb.CreateRefreshToken(context.Background(), helloprofiledb.CreateRefreshTokenParams{
				UserID: dbOtp.UserID,
				Token:  refreshToken,
			})
			if err != nil {
				log.WithFields(fields).WithError(err).Error("Error occured saving refresh token")
			}

			log.WithFields(fields).Info(fmt.Sprintf("Refresh Token Id: %s", dbRefreshToken.ID))
		}()
		go func() {
			err = env.saveLogin(helloprofiledb.CreateUserLoginParams{
				UserID:              dbOtp.UserID,
				ResponseCode:        sql.NullString{String: util.SUCCESS_RESPONSE_CODE, Valid: true},
				ResponseDescription: sql.NullString{String: util.SUCCESS_RESPONSE_MESSAGE, Valid: true},
				LoginStatus:         true,
				IpAddress:           sql.NullString{String: c.RealIP(), Valid: true},
				Device:              sql.NullString{String: c.Request().UserAgent()},
			})
			if err != nil {
				log.WithFields(fields).Info("Successful login...")
			}
			err := env.HelloProfileDb.UpdateResolvedLogin(context.Background(), dbOtp.UserID)
			if err != nil {
				log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured clearing failed user logins after successful login...")
			}
		}()
		loginResponse := &models.SuccessResponse{
			ResponseCode:    util.SUCCESS_RESPONSE_CODE,
			ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
		}

		c.Response().Header().Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
		c.Response().Header().Set("Refresh-Token", refreshToken)
		c.Response().Header().Set("Role", strings.Join(userRoles, ":"))
		c.JSON(http.StatusOK, loginResponse)

	} else {
		errorResponse.Errorcode = util.SESSION_EXPIRED_ERROR_CODE
		errorResponse.ErrorMessage = util.OTP_EXPIRED_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Otp has expired...")
		c.JSON(http.StatusUnauthorized, errorResponse)

	}
	log.WithFields(fields).Info("Finished processing Verify otp request...")
	return err
}

// ResetPassword password is used to reset account password
func (env *Env) ResetPassword(c echo.Context) (err error) {
	fields := log.Fields{"microservice": "helloprofile.service", "application": c.Param("application"), "function": "ResetPassword"}
	log.WithFields(fields).Info("Password Reset Request received")
	errorResponse := new(models.Errormessage)

	var authCode string
	authArray := strings.Split(c.Request().Header.Get("Authorization"), " ")
	if len(authArray) != 2 {
		errorResponse.Errorcode = util.INVALID_AUTHENTICATION_SCHEME_ERROR_CODE
		errorResponse.ErrorMessage = util.INVALID_AUTHENTICATION_SCHEME_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error(util.INVALID_AUTHENTICATION_SCHEME_ERROR_MESSAGE)
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	authCode = authArray[1]

	verifiedClaims, err := util.VerifyToken(authCode)

	if err != nil || verifiedClaims.Email == "" {
		errorResponse.Errorcode = util.SESSION_EXPIRED_ERROR_CODE
		errorResponse.ErrorMessage = util.RESET_PASSWORD_TOKEN_EXPIRED_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Token has expired...")
		c.JSON(http.StatusUnauthorized, errorResponse)
		return err
	}

	request := new(models.ResetPasswordRequest)
	if err = c.Bind(request); err != nil {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to marshal request")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	// Try to update password
	user, err := env.HelloProfileDb.GetUser(context.Background(), sql.NullString{String: strings.ToLower(verifiedClaims.Email), Valid: true})
	if err != nil {
		errorResponse.Errorcode = util.USER_NAME_OR_PASSWORD_INCORRECT_ERROR_CODE
		errorResponse.ErrorMessage = util.USER_NAME_OR_PASSWORD_INCORRECT_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error fetching user")
		c.JSON(http.StatusUnauthorized, errorResponse)
		return err
	}
	go func() {
		var hashedPassword string
		if request.NewPassword != "" {
			hashedPassword = util.GenerateHash(request.NewPassword)
		}
		_, err := env.HelloProfileDb.UpdateUser(context.Background(), helloprofiledb.UpdateUserParams{
			Username_2:                user.Username,
			IsLockedOut:               false,
			Email:                     user.Email,
			Firstname:                 user.Firstname,
			ImageUrl:                  user.ProfilePicture,
			IsActive:                  user.IsActive,
			IsEmailConfirmed:          user.IsEmailConfirmed,
			IsPasswordSystemGenerated: user.IsPasswordSystemGenerated,
			Lastname:                  user.Lastname,
			Password:                  sql.NullString{String: hashedPassword, Valid: true},
			Username:                  user.Username,
			Phone:                     user.Phone,
		})
		if err != nil {
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to update account")
		}
		log.WithFields(fields).Info("Successfully changed password...")

	}()
	go func() {
		err := env.HelloProfileDb.UpdateResolvedLogin(context.Background(), user.ID)
		if err != nil {
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured clearing failed user logins after successful login...")
		}
		log.WithFields(fields).Info("Successsfully updated failed logins ")

	}()
	resetResponse := &models.SuccessResponse{
		ResponseCode:    util.SUCCESS_RESPONSE_CODE,
		ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
	}
	c.JSON(http.StatusOK, resetResponse)
	return err
}

func (env *Env) Feedback(c echo.Context) (err error) {
	fields := log.Fields{"microservice": "helloprofile.service", "application": c.Param("application"), "function": "Feedback"}
	log.WithFields(fields).Info("Password Reset Request received")
	errorResponse := new(models.Errormessage)
	request := new(models.FeedbackRequest)
	if err = c.Bind(request); err != nil {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to marshal request")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	log.WithFields(fields).Info("Sending feedback message to the support team...")

	communicationEndpoint := os.Getenv("COMMUNICATION_SERVICE_ENDPOINT")

	emailPath := os.Getenv("EMAIL_PATH")
	emailRequest := models.SendEmailRequest{
		From:    models.EmailAddress{Email: os.Getenv("SMTP_USER"), Name: "HelloProfile"},
		To:      []models.EmailAddress{{Email: os.Getenv("SMTP_USER"), Name: "Support"}},
		CC:      []models.EmailAddress{{Email: "daprinz.op@gmail.com", Name: "Prince Okechukwu"}, {Email: "calveen.chikezie@gmail.com", Name: "Kelvin Chikezie"}, {Email: "amehugochukwu@gmail.com", Name: "Julius Ameh"}},
		Subject: util.FEEDBACK_SUBJECT,
		Message: fmt.Sprintf(util.FEEDBACK_MESSAGE, request.Sender, request.Message, strings.Join(request.AttachmentUrl, "</br>")),
	}
	emailRequestBytes, _ := json.Marshal(emailRequest)
	emailRequestReader := bytes.NewReader(emailRequestBytes)
	log.WithFields(fields).Info("Sending support email to the support team...")

	client := &http.Client{}
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s%s", communicationEndpoint, emailPath), emailRequestReader)
	req.Header.Add("Client_Id", "persianblack")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	emailResponse, err := client.Do(req)

	// emailResponse, err := http.Post(fmt.Sprintf("%s%s", communicationEndpoint, emailPath), "application/json", bytes.NewBuffer(emailRequestBytes))
	if err != nil {
		errorResponse.Errorcode = util.SUPPORT_EMAIL_SENDING_FAILURE_CODE
		errorResponse.ErrorMessage = util.SUPPORT_EMAIL_SENDING_FAILURE_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured sending otp")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	} else {
		if emailResponse.StatusCode == 200 {
			log.WithFields(fields).Info("Email sent successfully")
		} else {
			errorResponse.Errorcode = util.SUPPORT_EMAIL_SENDING_FAILURE_CODE
			errorResponse.ErrorMessage = util.SUPPORT_EMAIL_SENDING_FAILURE_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured sending Email")
			c.JSON(http.StatusBadRequest, errorResponse)
			return err
		}
		emailBody, _ := ioutil.ReadAll(emailResponse.Body)
		log.WithFields(fields).Info(fmt.Sprintf("Response body from email request: %s", emailBody))
		c.JSON(http.StatusOK, models.SuccessResponse{
			ResponseCode:    util.SUCCESS_RESPONSE_CODE,
			ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
		})
		return err
	}

}
