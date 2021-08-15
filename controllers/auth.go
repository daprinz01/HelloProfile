package controllers

import (
	"authengine/models"
	"authengine/persistence/orm/authdb"
	"authengine/util"
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

	"github.com/labstack/echo/v4"
)

//Login is used to sign users in
func (env *Env) Login(c echo.Context) (err error) {
	fields := log.Fields{"microservice": "persian.black.authengine.service", "application": c.Param("application")}
	log.WithFields(fields).Info("Login Request received")

	errorResponse := new(models.Errormessage)

	application := c.Param("application")
	if application == "" {
		errorResponse.Errorcode = "01"
		errorResponse.ErrorMessage = "Application not specified"
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Application not specified")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	applicationObject, err := env.AuthDb.GetApplication(context.Background(), strings.ToLower(application))
	if err != nil {
		errorResponse.Errorcode = "06"
		errorResponse.ErrorMessage = "Application is invalid"
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Application is invalid")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	log.WithFields(fields).Info(fmt.Sprintf("Applicaiton ID: %d", applicationObject.ID))
	request := new(models.LoginRequest)
	if err = c.Bind(request); err != nil {
		errorResponse.Errorcode = "02"
		errorResponse.ErrorMessage = "Invalid request"
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to marshal request")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	var username sql.NullString
	username.String = strings.ToLower(request.Username)
	username.Valid = true
	user, err := env.AuthDb.GetUser(context.Background(), username)
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Oops... something is wrong here... your email or password is incorrect..."
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error fetching user")
		c.JSON(http.StatusUnauthorized, errorResponse)
		return err
	}
	if user.IsLockedOut {
		errorResponse.Errorcode = "13"
		errorResponse.ErrorMessage = "Sorry you exceeded the maximum login attempts... Kindly reset your password to continue..."
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Account was locked out....")
		c.JSON(http.StatusUnauthorized, errorResponse)
		return err
	}
	if util.VerifyHash(user.Password.String, request.Password) {
		log.WithFields(fields).Info("Verifying that user is in the role access is being requested...")

		userRoles, err := env.AuthDb.GetUserRoles(context.Background(), sql.NullString{String: user.Email, Valid: true})
		if err != nil {
			log.WithFields(fields).WithError(err).Error(`Invalid role entered... Changing to default role of "Guest"`)
			userRoles[0] = "guest"
		}

		log.WithFields(fields).Info(fmt.Sprintf("Generating authentication token for user: %s role: %v...", user.Email, userRoles))
		authToken, refreshToken, err := util.GenerateJWT(user.Email, userRoles)
		if err != nil {
			errorResponse.Errorcode = "05"
			errorResponse.ErrorMessage = fmt.Sprintf("Error occured generating auth token: %s", err)
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured generating quth token")
			c.JSON(http.StatusConflict, errorResponse)
			return err
		}
		log.WithFields(fields).Info("Storing refresh token in separate thread...")
		// store refresh token add later
		go func() {
			dbRefreshToken, err := env.AuthDb.CreateRefreshToken(context.Background(), authdb.CreateRefreshTokenParams{
				UserID: user.ID,
				Token:  refreshToken,
			})
			if err != nil {
				log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while saving refresh token")
			}

			log.WithFields(fields).Info(fmt.Sprintf("Refresh Token Id: %d", dbRefreshToken.ID))
		}()
		go func() {
			err = env.saveLogin(authdb.CreateUserLoginParams{
				ApplicationID:       applicationObject.ID,
				UserID:              user.ID,
				ResponseCode:        sql.NullString{String: "00", Valid: true},
				ResponseDescription: sql.NullString{String: "Success", Valid: true},
				LoginStatus:         true,
			})
			if err != nil {
				log.WithFields(fields).Info("Successful login...")
			}
			err := env.AuthDb.UpdateResolvedLogin(context.Background(), user.ID)
			if err != nil {
				log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured clearing failed user logins after successful login...")
			}
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
				Phone:                     user.Phone.String,
			},
		}
		c.Response().Header().Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
		c.Response().Header().Set("Refresh-Token", refreshToken)
		c.Response().Header().Set("Role", strings.Join(userRoles, ":"))
		c.JSON(http.StatusOK, loginResponse)

	} else {
		errorResponse.Errorcode = "04"
		errorResponse.ErrorMessage = "Oops... something is wrong here... your email or password is incorrect..."
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Password incorrect...")
		go func() {
			err = env.saveLogin(authdb.CreateUserLoginParams{
				ApplicationID:       applicationObject.ID,
				UserID:              user.ID,
				ResponseCode:        sql.NullString{String: errorResponse.Errorcode, Valid: true},
				ResponseDescription: sql.NullString{String: errorResponse.ErrorMessage, Valid: true},
				LoginStatus:         false,
				Resolved:            false,
			})
			if err != nil {
				log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Password incorrect...")
			}
		}()

		userLogins, err := env.AuthDb.GetUnResoledLogins(context.Background(), user.ID)
		if err != nil {
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error ocurred fetching user logins")
		}
		var lockoutCount int
		lockOutCountKey := os.Getenv("LOCK_OUT_COUNT")
		if lockOutCountKey == "" {
			log.WithFields(fields).Error("LOCK_OUT_COUNT cannot be empty")
			log.WithFields(fields).Info("LOCK_OUT_COUNT cannot be empty, setting default of 5...")
		} else {
			log.WithFields(fields).Info("Setting lock out count...")
			lockoutCount, err = strconv.Atoi(lockOutCountKey)
			if err != nil {
				log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while converting lock out count")
			}
		}
		// Check if account has exceeded the lockout count
		if len(userLogins) >= lockoutCount {
			lockoutUpdate, err := env.AuthDb.UpdateUser(context.Background(), authdb.UpdateUserParams{
				Username_2:                user.Username,
				IsLockedOut:               true,
				Address:                   user.Address,
				City:                      user.City,
				Country:                   user.Country,
				CreatedAt:                 user.CreatedAt,
				Email:                     user.Email,
				Firstname:                 user.Firstname,
				ImageUrl:                  user.ProfilePicture,
				IsActive:                  user.IsActive,
				IsEmailConfirmed:          user.IsEmailConfirmed,
				IsPasswordSystemGenerated: user.IsPasswordSystemGenerated,
				Lastname:                  user.Lastname,
				Password:                  user.Password,
				State:                     user.State,
				Username:                  user.Username,
				Phone:                     user.Phone,
			})
			if err != nil {
				log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to lockout account")
			}
			log.WithFields(fields).Info(fmt.Sprintf("Account with username: %s has been locked out", lockoutUpdate.Username.String))

			errorResponse.Errorcode = "12"
			errorResponse.ErrorMessage = "Account locked out, kindly reset your password to continue..."
			c.JSON(http.StatusUnauthorized, errorResponse)
			return err
		}

		c.JSON(http.StatusBadRequest, errorResponse)

	}

	return err

}

// saveLogin is used to log a login request that failed with incorrect password or was successful
func (env *Env) saveLogin(createParams authdb.CreateUserLoginParams) error {
	fields := log.Fields{"microservice": "persian.black.authengine.service", "function": "saveLogin"}

	userLogin, err := env.AuthDb.CreateUserLogin(context.Background(), createParams)
	if err != nil {
		log.WithFields(fields).WithError(err).Error("Error occured saving user login")
		return err
	}
	log.WithFields(fields).Info(fmt.Sprintf("Successfully saved user login, user login id: %d", userLogin.ID))
	return err
}

// Register is used to register new users
func (env *Env) Register(c echo.Context) (err error) {
	fields := log.Fields{"microservice": "persian.black.authengine.service", "application": c.Param("application")}
	log.WithFields(fields).Info("Register Request received")

	errorResponse := new(models.Errormessage)

	applicationName := c.Param("application")
	if applicationName == "" {
		errorResponse.Errorcode = "01"
		errorResponse.ErrorMessage = "Application not specified"
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Application not specified")
		c.JSON(http.StatusBadRequest, errorResponse)
		return nil
	}
	application, err := env.AuthDb.GetApplication(context.Background(), strings.ToLower(applicationName))
	if err != nil {
		errorResponse.Errorcode = "06"
		errorResponse.ErrorMessage = "Application is invalid"
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Application is invalid")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	log.WithFields(fields).Info(fmt.Sprintf("Applicaiton ID: %d", application.ID))

	request := new(models.UserDetail)
	if err = c.Bind(request); err != nil {
		errorResponse.Errorcode = "02"
		errorResponse.ErrorMessage = "Invalid request"
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
	user, err := env.AuthDb.CreateUser(context.Background(), authdb.CreateUserParams{
		Address:                   sql.NullString{String: request.Address, Valid: request.Address != ""},
		City:                      sql.NullString{String: request.City, Valid: request.City != ""},
		Country:                   sql.NullString{String: request.Country, Valid: request.Country != ""},
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
		State:                     sql.NullString{String: request.State, Valid: request.State != ""},
		Username:                  sql.NullString{String: strings.ToLower(request.Username), Valid: request.Username != ""},
		Phone:                     sql.NullString{String: request.Phone, Valid: request.Phone != ""},
	})

	if err != nil {
		errorResponse.Errorcode = "04"
		errorResponse.ErrorMessage = "User already exist"
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured creating user")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	go func() {
		err = env.saveLogin(authdb.CreateUserLoginParams{
			ApplicationID:       application.ID,
			UserID:              user.ID,
			ResponseCode:        sql.NullString{String: "00", Valid: true},
			ResponseDescription: sql.NullString{String: "Registration successful...", Valid: true},
			LoginStatus:         true,
		})
		if err != nil {
			log.WithFields(fields).WithError(err).Error("Error occured while saving login...")
		}
	}()
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
			Phone:                     user.Phone.String,
		},
	}

	// log.Println(fmt.Sprintf("Got to response string: %s", responseString))
	log.WithFields(fields).Info("Generating authentication token...")
	role := c.Request().Header.Get("Role")
	dbRole, err := env.AuthDb.GetRole(context.Background(), strings.ToLower(role))
	if err != nil {
		log.WithFields(fields).WithError(err).Error(`Invalid role entered... Changing to default role of "Guest"`)
		role = "Guest"
	} else {
		log.WithFields(fields).Info(fmt.Sprintf("Creating token for user: %s | role: %s", user.Email, dbRole.Name))

	}
	go func() {
		log.WithFields(fields).Info("Verifying that role exist for the application")
		applicationRole, err := env.AuthDb.GetApplicationRole(context.Background(), authdb.GetApplicationRoleParams{
			ApplicationsID: application.ID,
			RolesID:        dbRole.ID,
		})
		if err != nil {
			log.WithFields(fields).WithError(err).Error("Error occured fetching applicationRole")
		}
		log.WithFields(fields).Info(fmt.Sprintf("Role is valid for application. Application Role Id: %d", applicationRole.ID))
		log.WithFields(fields).Info("Adding user to role...")
		userRole, err := env.AuthDb.AddUserRole(context.Background(), authdb.AddUserRoleParams{
			Name:     strings.ToLower(role),
			Username: user.Username,
		})
		if err != nil {
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error(fmt.Sprintf("Error occured adding user: %s to role: %s", user.Username.String, role))
		}
		log.WithFields(fields).Info(fmt.Sprintf("Successfully added user to role.. User Role Id: %d", userRole.ID))
	}()
	authToken, refreshToken, err := util.GenerateJWT(user.Email, strings.Split(role, ":"))
	if err != nil {
		errorResponse.Errorcode = "05"
		errorResponse.ErrorMessage = fmt.Sprintf("Error occured generating auth token: %s", err)
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured generating auth token")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	log.WithFields(fields).Info("Storing refresh token in separate thread...")
	// store refresh token add later
	go func() {
		dbRefreshToken, err := env.AuthDb.CreateRefreshToken(context.Background(), authdb.CreateRefreshTokenParams{
			UserID: user.ID,
			Token:  refreshToken,
		})
		if err != nil {
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to save refresh token")
		}

		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error(fmt.Sprintf("Refresh Token Id: %d", dbRefreshToken.ID))
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
	fields := log.Fields{"microservice": "persian.black.authengine.service", "application": c.Param("application")}
	log.WithFields(fields).Info("Refresh token Request received")

	errorResponse := new(models.Errormessage)

	var authCode string
	authArray := strings.Split(c.Request().Header.Get("Authorization"), " ")
	if len(authArray) != 2 {
		errorResponse.Errorcode = "11"
		errorResponse.ErrorMessage = "Unsupported authentication scheme type"
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Unsupported authentication scheme type")
		c.JSON(http.StatusUnauthorized, errorResponse)
		return err
	}
	authCode = authArray[1]
	refreshToken := c.Request().Header.Get("Refresh-Token")

	verifiedClaims, err := util.VerifyToken(authCode)
	if err == nil {
		errorResponse.Errorcode = "10"
		errorResponse.ErrorMessage = "Session is still valid..."
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Token is still valid...")
		c.JSON(http.StatusTooEarly, errorResponse)
		return err
	}
	if err != nil && verifiedClaims.Email == "" {
		errorResponse.Errorcode = "02"
		errorResponse.ErrorMessage = "Invalid request"
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Invalid request")
		c.JSON(http.StatusUnauthorized, errorResponse)
		return err
	}

	dbRefreshToken, err := env.AuthDb.GetRefreshToken(context.Background(), refreshToken)
	if err != nil {
		errorResponse.Errorcode = "08"
		errorResponse.ErrorMessage = "Cannot refresh session. Kindly login again"
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured refreshing token")
		c.JSON(http.StatusUnauthorized, errorResponse)
		return err
	}
	defer func() {
		err = env.AuthDb.DeleteRefreshToken(context.Background(), refreshToken)
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
			errorResponse.Errorcode = "05"
			errorResponse.ErrorMessage = fmt.Sprintf("Error occured generating auth token: %s", err)
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured generating auth token")
			c.JSON(http.StatusBadRequest, errorResponse)
			return err
		}
		log.WithFields(fields).Info("Fetching user...")
		user, err := env.AuthDb.GetUser(context.Background(), sql.NullString{String: strings.ToLower(verifiedClaims.Email), Valid: true})
		if err != nil {
			errorResponse.Errorcode = "03"
			errorResponse.ErrorMessage = "Oops... something is wrong here... your email or password is incorrect..."
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error fetching user")
			c.JSON(http.StatusUnauthorized, errorResponse)
			return err
		}
		log.WithFields(fields).Info("Storing refresh token in separate thread...")
		// store refresh token add later
		go func() {
			dbRefreshToken, err := env.AuthDb.CreateRefreshToken(context.Background(), authdb.CreateRefreshTokenParams{
				UserID: user.ID,
				Token:  refreshToken,
			})
			if err != nil {
				log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to create refresh token")
			}

			log.WithFields(fields).Info(fmt.Sprintf("Refresh Token Id: %d", dbRefreshToken.ID))
		}()
		resetResponse := &models.RefreshResponse{
			ResponseCode:    "00",
			ResponseMessage: "Success",
		}

		c.Response().Header().Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
		c.Response().Header().Set("Refresh-Token", refreshToken)
		c.Response().Header().Set("Role", verifiedClaims.Role)
		c.JSON(http.StatusOK, resetResponse)

	} else {
		errorResponse.Errorcode = "09"
		errorResponse.ErrorMessage = "Session expired. Kindly login again to continue"
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Token has expired...")
		c.JSON(http.StatusUnauthorized, errorResponse)
		return err

	}
	return err

}

// SendOtp is used to send OTP request after validating user exist
func (env *Env) SendOtp(c echo.Context) (err error) {
	fields := log.Fields{"microservice": "persian.black.authengine.service", "application": c.Param("application")}
	log.WithFields(fields).Info("Send OTP Request received")
	errorResponse := new(models.Errormessage)

	request := new(models.SendOtpRequest)
	if err = c.Bind(request); err != nil {
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to marshal request")
		errorResponse.Errorcode = "02"
		errorResponse.ErrorMessage = "Invalid request"
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	user, err := env.AuthDb.GetUser(context.Background(), sql.NullString{String: strings.ToLower(request.Email), Valid: true})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "If you have an account with us, you should get an otp"
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
		errorResponse.Errorcode = "14"
		errorResponse.ErrorMessage = "Error occured generating otp"
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured generating otp")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	// Save otp to db in another thread
	go func() {
		err = env.AuthDb.CreateOtp(context.Background(), authdb.CreateOtpParams{
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
				From:    models.EmailAddress{Email: os.Getenv("SMTP_USER"), Name: "Persian Black"},
				To:      []models.EmailAddress{{Email: user.Email, Name: fmt.Sprintf("%s %s", user.Firstname.String, user.Lastname.String)}},
				Subject: fmt.Sprintf("%s OTP", request.Purpose),
				Message: fmt.Sprintf("<h5>Hey %s,</h5><p>Kindly use the otp below to complete your request:</p><h4>%s</h4><p>Your account security is paramount to us. Don't share your otp with anyone.</p><h5>Micheal from Persian Black.</h5>", user.Firstname.String, otp),
			}
			emailRequestBytes, _ := json.Marshal(emailRequest)
			emailRequestReader := bytes.NewReader(emailRequestBytes)
			log.WithFields(fields).Info("Sending otp email...")

			client := &http.Client{}
			req, _ := http.NewRequest("POST", fmt.Sprintf("%s%s", communicationEndpoint, emailPath), emailRequestReader)
			req.Header.Add("Authorization", "Bearer persianblack")
			req.Header.Add("Accept", "application/json")
			req.Header.Add("Content-Type", "application/json")
			emailResponse, err := client.Do(req)

			// emailResponse, err := http.Post(fmt.Sprintf("%s%s", communicationEndpoint, emailPath), "application/json", bytes.NewBuffer(emailRequestBytes))
			if err != nil {
				log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured sending otp")
			}
			if emailResponse.StatusCode == 200 {
				log.WithFields(fields).Info("OTP sent successfully")
			} else {
				log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured sending OTP")
			}
			emailBody, _ := ioutil.ReadAll(emailResponse.Body)
			log.WithFields(fields).Info(fmt.Sprintf("Response body from email request: %s", emailBody))
		} else {
			if user.Phone.String == "" {
				log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Phonenumber not available")
			} else {
				smsPath := os.Getenv("SMS_PATH")
				smsRequest := models.SendSmsRequest{
					Phone:   user.Phone.String,
					Message: fmt.Sprintf("Your Persian Black %s code is:\n%s", request.Purpose, otp),
				}
				smsRequestBytes, _ := json.Marshal(smsRequest)
				smsRequestReader := bytes.NewReader(smsRequestBytes)
				log.WithFields(fields).Info("Sending otp sms...")

				client := &http.Client{}
				req, _ := http.NewRequest("POST", fmt.Sprintf("%s%s", communicationEndpoint, smsPath), smsRequestReader)
				req.Header.Add("Authorization", "Bearer persianblack")
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
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	c.JSON(http.StatusOK, resetResponse)
	return err
}

// SendOtp is used to send OTP request after validating user exist
func (env *Env) DoEmailVerification(c echo.Context) (err error) {
	fields := log.Fields{"microservice": "persian.black.authengine.service", "application": c.Param("application")}
	log.WithFields(fields).Info("Generate OTP Request received")
	errorResponse := new(models.Errormessage)

	request := new(models.SendOtpRequest)
	if err = c.Bind(request); err != nil {
		errorResponse.Errorcode = "02"
		errorResponse.ErrorMessage = "Invalid request"
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
		errorResponse.Errorcode = "14"
		errorResponse.ErrorMessage = "Error occured generating otp"
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured generating otp")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	log.WithFields(fields).Info("Successfully generated otp")
	//Save otp to db in another thread
	go func() {
		err = env.AuthDb.CreateEmailVerification(context.Background(), authdb.CreateEmailVerificationParams{
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
			From:    models.EmailAddress{Email: os.Getenv("SMTP_USER"), Name: "Persian Black"},
			To:      []models.EmailAddress{{Email: request.Email}},
			Subject: fmt.Sprintf("%s Email Verification", strings.ToTitle(request.Application)),
			Message: fmt.Sprintf("<h5>Hey,</h5><p>Kindly click the link below to confirm your email address</p><a href=\"%s/%s/%s\">click here</a><h5>Micheal from Persian Black.</h5>", request.VerifyPath, request.Email, otp),
		}
		emailRequestBytes, _ := json.Marshal(emailRequest)
		emailRequestReader := bytes.NewReader(emailRequestBytes)
		log.WithFields(fields).Info("Sending email verification email...")

		client := &http.Client{}
		req, _ := http.NewRequest("POST", fmt.Sprintf("%s%s", communicationEndpoint, emailPath), emailRequestReader)
		req.Header.Add("Authorization", "Bearer persianblack")
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
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	c.JSON(http.StatusOK, resetResponse)
	return err
}

func (env *Env) VerifyEmailToken(c echo.Context) (err error) {
	fields := log.Fields{"microservice": "persian.black.authengine.service", "application": c.Param("application")}
	log.WithFields(fields).Info("Verify Email token request received")

	errorResponse := new(models.Errormessage)

	request := new(models.VerifyOtpRequest)
	if err = c.Bind(request); err != nil {
		errorResponse.Errorcode = "02"
		errorResponse.ErrorMessage = "Invalid request"
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to marshal request")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	dbOtp, err := env.AuthDb.GetEmailVerification(context.Background(), request.OTP)
	if err != nil {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "Sorry we could not verify your request. Please try registering again..."
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Invalid request")
		c.JSON(http.StatusBadRequest, errorResponse)
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
			ResponseCode:    "00",
			ResponseMessage: "Success",
		}
		c.JSON(http.StatusOK, verifyResponse)
	} else {
		errorResponse.Errorcode = "16"
		errorResponse.ErrorMessage = "Oops... something is wrong here... your email verification link has expired.. Kindly register again"
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Email otp has expired or is invalid...")
		c.JSON(http.StatusForbidden, errorResponse)

	}
	log.WithFields(fields).Info("Finished processing Verify otp request...")
	return err
}

// VerifyOtp is used to verify and otp against a user. Authentication token is generated that is used in subsequent requests.
func (env *Env) VerifyOtp(c echo.Context) (err error) {
	fields := log.Fields{"microservice": "persian.black.authengine.service", "application": c.Param("application")}
	log.WithFields(fields).Info("Verify otp request received")

	// file, fileHeader, err := r.FormFile("request.AttachmentName[i]")

	// file, err := os.Create(fmt.Sprintf("%s%s", attachmentPath, request.AttachmentName[i].FileName))
	// file.WriteString()

	errorResponse := new(models.Errormessage)

	application := c.Param("application")
	if application == "" {
		errorResponse.Errorcode = "01"
		errorResponse.ErrorMessage = "Application not specified"
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Application not specified")
		c.JSON(http.StatusBadRequest, errorResponse)
		return nil
	}
	applicationObject, err := env.AuthDb.GetApplication(context.Background(), strings.ToLower(application))
	if err != nil {
		errorResponse.Errorcode = "06"
		errorResponse.ErrorMessage = "Application is invalid"
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Application is invalid")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	log.WithFields(fields).Info(fmt.Sprintf("Applicaiton ID: %d", applicationObject.ID))
	// errorResponse := new(models.Errormessage)
	//
	request := new(models.VerifyOtpRequest)
	if err = c.Bind(request); err != nil {
		errorResponse.Errorcode = "02"
		errorResponse.ErrorMessage = "Invalid request"
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to marshal request")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	dbOtp, err := env.AuthDb.GetOtp(context.Background(), authdb.GetOtpParams{
		OtpToken: sql.NullString{String: request.OTP, Valid: true},
		Username: sql.NullString{String: strings.ToLower(request.Email), Valid: true},
	})
	if err != nil {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "Incorrect OTP. Please try again..."
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Invalid request")
		c.JSON(http.StatusBadRequest, errorResponse)
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

		userRoles, err := env.AuthDb.GetUserRoles(context.Background(), sql.NullString{String: request.Email, Valid: true})
		if err != nil {
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error(`Invalid role entered... Changing to default role of "Guest"`)
			userRoles[0] = "Guest"
		}

		log.WithFields(fields).Info(fmt.Sprintf("Generating authentication token for user: %s role: %s...", request.Email, strings.Join(userRoles, ":")))
		authToken, refreshToken, err := util.GenerateJWT(request.Email, userRoles)
		if err != nil {
			errorResponse.Errorcode = "05"
			errorResponse.ErrorMessage = "Error occured generating auth token"
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured generating wuth token")
			c.JSON(http.StatusBadRequest, errorResponse)
			return err
		}
		log.WithFields(fields).Info("Storing refresh token in separate thread...")
		// store refresh token add later
		go func() {
			dbRefreshToken, err := env.AuthDb.CreateRefreshToken(context.Background(), authdb.CreateRefreshTokenParams{
				UserID: dbOtp.UserID,
				Token:  refreshToken,
			})
			if err != nil {
				log.WithFields(fields).WithError(err).Error("Error occured saving refresh token")
			}

			log.WithFields(fields).Info(fmt.Sprintf("Refresh Token Id: %d", dbRefreshToken.ID))
		}()
		go func() {
			err = env.saveLogin(authdb.CreateUserLoginParams{
				ApplicationID:       applicationObject.ID,
				UserID:              dbOtp.UserID,
				ResponseCode:        sql.NullString{String: "00", Valid: true},
				ResponseDescription: sql.NullString{String: "Success", Valid: true},
				LoginStatus:         true,
			})
			if err != nil {
				log.WithFields(fields).Info("Successful login...")
			}
			err := env.AuthDb.UpdateResolvedLogin(context.Background(), dbOtp.UserID)
			if err != nil {
				log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured clearing failed user logins after successful login...")
			}
		}()
		loginResponse := &models.SuccessResponse{
			ResponseCode:    "00",
			ResponseMessage: "Success",
		}

		c.Response().Header().Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
		c.Response().Header().Set("Refresh-Token", refreshToken)
		c.Response().Header().Set("Role", strings.Join(userRoles, ":"))
		c.JSON(http.StatusOK, loginResponse)

	} else {
		errorResponse.Errorcode = "16"
		errorResponse.ErrorMessage = "Oops... something is wrong here... your one time token has expired. Kindly request another one..."
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Otp has expired...")
		c.JSON(http.StatusBadRequest, errorResponse)

	}
	log.WithFields(fields).Info("Finished processing Verify otp request...")
	return err
}

// ResetPassword password is used to reset account password
func (env *Env) ResetPassword(c echo.Context) (err error) {
	fields := log.Fields{"microservice": "persian.black.authengine.service", "application": c.Param("application")}
	log.WithFields(fields).Info("Password Reset Request received")
	errorResponse := new(models.Errormessage)

	var authCode string
	authArray := strings.Split(c.Request().Header.Get("Authorization"), " ")
	if len(authArray) != 2 {
		errorResponse.Errorcode = "11"
		errorResponse.ErrorMessage = "Unsupported authentication scheme type"
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Unsupported authentication scheme type")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	authCode = authArray[1]

	verifiedClaims, err := util.VerifyToken(authCode)

	if err != nil || verifiedClaims.Email == "" {
		errorResponse.Errorcode = "09"
		errorResponse.ErrorMessage = "Session expired. Kindly try generating one time password again"
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Token has expired...")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	request := new(models.ResetPasswordRequest)
	if err = c.Bind(request); err != nil {
		errorResponse.Errorcode = "02"
		errorResponse.ErrorMessage = "Invalid request"
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to marshal request")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	// Try to update password
	user, err := env.AuthDb.GetUser(context.Background(), sql.NullString{String: strings.ToLower(verifiedClaims.Email), Valid: true})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Oops... something is wrong here... your email is incorrect..."
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error fetching user")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	go func() {
		var hashedPassword string
		if request.NewPassword != "" {
			hashedPassword = util.GenerateHash(request.NewPassword)

		}
		_, err := env.AuthDb.UpdateUser(context.Background(), authdb.UpdateUserParams{
			Username_2:                user.Username,
			IsLockedOut:               false,
			Address:                   user.Address,
			City:                      user.City,
			Country:                   user.Country,
			CreatedAt:                 user.CreatedAt,
			Email:                     user.Email,
			Firstname:                 user.Firstname,
			ImageUrl:                  user.ProfilePicture,
			IsActive:                  user.IsActive,
			IsEmailConfirmed:          user.IsEmailConfirmed,
			IsPasswordSystemGenerated: user.IsPasswordSystemGenerated,
			Lastname:                  user.Lastname,
			Password:                  sql.NullString{String: hashedPassword, Valid: true},
			State:                     user.State,
			Username:                  user.Username,
			Phone:                     user.Phone,
		})
		if err != nil {
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to update account")
		}
		log.WithFields(fields).Info("Successfully changed password...")

	}()
	go func() {
		err := env.AuthDb.UpdateResolvedLogin(context.Background(), user.ID)
		if err != nil {
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured clearing failed user logins after successful login...")
		}
		log.WithFields(fields).Info("Successsfully updated failed logins ")

	}()
	resetResponse := &models.RefreshResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	c.JSON(http.StatusOK, resetResponse)
	return err
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
