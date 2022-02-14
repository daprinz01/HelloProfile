package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
	"helloprofile.com/models"
	"helloprofile.com/persistence/orm/helloprofiledb"
	"helloprofile.com/util"

	"github.com/labstack/echo/v4"
)

// CheckAvailability is used to check user availablity
func (env *Env) CheckAvailability(c echo.Context) (err error) {
	fields := log.Fields{"microservice": "helloprofile.service", "application": "backend"}
	log.WithFields(fields).Info("Check availability Request received")

	errorResponse := new(models.Errormessage)

	username := c.Param("username")
	if err != nil {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Username not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	log.WithFields(fields).Info(fmt.Sprintf("Username: %s", username))

	user, err := env.HelloProfileDb.GetUser(context.Background(), sql.NullString{String: strings.ToLower(username), Valid: true})
	if err != nil {
		errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
		errorResponse.ErrorMessage = util.NO_RECORD_FOUND_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error fetching user: %s")
		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.WithFields(fields).Info(fmt.Sprintf("User %s exists...", user.Username.String))

	response := &models.SuccessResponse{
		ResponseCode:    util.SUCCESS_RESPONSE_CODE,
		ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
	}
	c.JSON(http.StatusOK, response)
	return err
}

// GetUser is used to fetch user details
func (env *Env) GetUser(c echo.Context) (err error) {
	fields := log.Fields{"microservice": "helloprofile.service", "application": "backend"}
	log.WithFields(fields).Info("Get User Request received")

	email := c.Request().Header.Get("email")
	errorResponse := new(models.Errormessage)

	if err != nil {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Username not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	log.WithFields(fields).Info(fmt.Sprintf("Username: %s", email))

	user, err := env.HelloProfileDb.GetUser(context.Background(), sql.NullString{String: strings.ToLower(email), Valid: true})
	if err != nil {
		errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
		errorResponse.ErrorMessage = util.NO_RECORD_FOUND_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error fetching user")
		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.WithFields(fields).Info(fmt.Sprintf("User %s exists...", user.Username.String))
	profiles := make(chan []models.Profile)
	go env.getProfiles(user.ID, profiles, fields)
	response := &models.SuccessResponse{
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
			Profiles:                  <-profiles,
		},
	}
	c.JSON(http.StatusOK, response)
	return err
}

// GetUsers is used to fetch user details. This is an admin function. User must be an admin to access this function
func (env *Env) GetUsers(c echo.Context) (err error) {
	fields := log.Fields{"microservice": "helloprofile.service", "application": "backend"}
	log.WithFields(fields).Info("Get User Request received")

	errorResponse := new(models.Errormessage)

	users, err := env.HelloProfileDb.GetUsers(context.Background())
	if err != nil {
		errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
		errorResponse.ErrorMessage = util.NO_RECORD_FOUND_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error fetching user")
		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	userResponse := make([]models.UserDetail, len(users))

	for index, user := range users {
		profiles := make(chan []models.Profile)
		go env.getProfiles(user.ID, profiles, fields)
		tempUser := models.UserDetail{
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
			Profiles:                  <-profiles,
		}
		userResponse[index] = tempUser
	}
	log.WithFields(fields).Info(fmt.Sprintf("Returning %d users...", len(users)))

	resetResponse := &models.SuccessResponse{
		ResponseCode:    util.SUCCESS_RESPONSE_CODE,
		ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
		ResponseDetails: &userResponse,
	}
	c.JSON(http.StatusOK, resetResponse)
	return err
}

// UpdateUser is used to update User information. It can be used to update user details and timezone details as required. Only pass the details to be updated. Email or username is mandatory.
func (env *Env) UpdateUser(c echo.Context) (err error) {

	fields := log.Fields{"microservice": "helloprofile.service", "application": "backend"}
	log.WithFields(fields).Info("Update user request received...")
	errorResponse := new(models.Errormessage)

	request := new(models.UserDetail)
	if err = c.Bind(request); err != nil {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to marshal request")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	log.WithFields(fields).Info("Checking if user exist...")
	user, err := env.HelloProfileDb.GetUser(context.Background(), sql.NullString{String: strings.ToLower(request.Email), Valid: true})
	if err != nil {
		errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
		errorResponse.ErrorMessage = util.NO_RECORD_FOUND_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error fetching user")
		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.WithFields(fields).Info(fmt.Sprintf("User %s exists...", user.Username.String))
	go func() {
		log.WithFields(fields).Info("Updating user...")
		// Check all the user parameters passed and consolidate with existing user record

		_, err := env.HelloProfileDb.UpdateUser(context.Background(), helloprofiledb.UpdateUserParams{

			Email:                     user.Email,
			Firstname:                 sql.NullString{String: getValue(request.Firstname, user.Firstname.String), Valid: true},
			ImageUrl:                  sql.NullString{String: getValue(request.ProfilePicture, user.ProfilePicture.String), Valid: true},
			IsActive:                  user.IsActive,
			IsEmailConfirmed:          user.IsEmailConfirmed,
			IsLockedOut:               user.IsLockedOut,
			IsPasswordSystemGenerated: user.IsPasswordSystemGenerated,
			Lastname:                  sql.NullString{String: getValue(request.Lastname, user.Lastname.String), Valid: true},
			Password:                  user.Password,
			Username:                  sql.NullString{String: getValue(request.Username, user.Username.String), Valid: true},
			Phone:                     sql.NullString{String: getValue(request.Phone, user.Phone.String), Valid: true},
			Username_2:                sql.NullString{String: user.Email, Valid: true},
		})

		if err != nil {
			log.WithFields(fields).WithError(err).Error("Error occured updating user information")

		}
		log.WithFields(fields).Info("Successfully updated user details")

	}()
	response := &models.SuccessResponse{
		ResponseCode:    util.SUCCESS_RESPONSE_CODE,
		ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
	}
	c.JSON(http.StatusOK, response)
	return err
}

// UpdateUserRole Add Role to applications
func (env *Env) UpdateUserRole(c echo.Context) (err error) {
	fields := log.Fields{"microservice": "helloprofile.service", "application": "backend"}
	log.WithFields(fields).Info("Update user's role request received...")

	errorResponse := new(models.Errormessage)

	role := c.Param("newRole")

	log.WithFields(fields).Info(fmt.Sprintf("New Role: %s", role))
	if role == "" {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Role not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	oldRole := c.Param("oldRole")

	log.WithFields(fields).Info(fmt.Sprintf("Old Role: %s", role))
	if oldRole == "" {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Role not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	username := c.Param("username")
	if username == "" {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Username not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	log.WithFields(fields).Info(fmt.Sprintf("Username: %s", username))

	dbRole, err := env.HelloProfileDb.UpdateUserRole(context.Background(), helloprofiledb.UpdateUserRoleParams{
		Username:   sql.NullString{String: strings.ToLower(username), Valid: true},
		Username_2: sql.NullString{String: strings.ToLower(username), Valid: true},
		Name:       strings.ToLower(role),
		Name_2:     strings.ToLower(oldRole),
	})
	if err != nil {
		errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
		errorResponse.ErrorMessage = util.NO_RECORD_FOUND_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured updating  user role")

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.WithFields(fields).Info(fmt.Sprintf("Successfully updated user role: %v", dbRole))

	response := &models.SuccessResponse{
		ResponseCode:    util.SUCCESS_RESPONSE_CODE,
		ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
	}
	c.JSON(http.StatusOK, response)
	return err
}

// AddUserToRole Add user to a role. This increases the roles the user is being added to including the previous roles. At login a role must be selected else the the default role guest is selected for the user
func (env *Env) AddUserToRole(c echo.Context) (err error) {
	fields := log.Fields{"microservice": "helloprofile.service", "application": "backend"}
	log.WithFields(fields).Info("Add user to role request received...")

	errorResponse := new(models.Errormessage)

	role := c.Param("role")

	log.WithFields(fields).Info(fmt.Sprintf("Role: %s", role))
	if err != nil {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Role not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	username := c.Param("username")
	if err != nil {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Username not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	log.WithFields(fields).Info(fmt.Sprintf("Username: %s", username))

	dbRole, err := env.HelloProfileDb.AddUserRole(context.Background(), helloprofiledb.AddUserRoleParams{
		Username: sql.NullString{String: strings.ToLower(username), Valid: true},
		Name:     strings.ToLower(role),
	})
	if err != nil {
		errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
		errorResponse.ErrorMessage = util.NO_RECORD_FOUND_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error(fmt.Sprintf("Error occured adding  user to role: %s", err))

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.WithFields(fields).Info(fmt.Sprintf("Successfully added role to application: %v", dbRole))

	response := &models.SuccessResponse{
		ResponseCode:    util.SUCCESS_RESPONSE_CODE,
		ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
	}
	c.JSON(http.StatusOK, response)
	return err
}

// DeleteUser is used to disable a users account
func (env *Env) DeleteUser(c echo.Context) (err error) {
	fields := log.Fields{"microservice": "helloprofile.service", "application": "backend"}
	log.WithFields(fields).Info("Delete user Request received")

	username := c.Param("username")
	errorResponse := new(models.Errormessage)

	if err != nil {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Username not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	log.WithFields(fields).Info(fmt.Sprintf("Username: %s", username))

	user, err := env.HelloProfileDb.GetUser(context.Background(), sql.NullString{String: strings.ToLower(username), Valid: true})
	if err != nil {
		errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
		errorResponse.ErrorMessage = util.NO_RECORD_FOUND_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error(fmt.Sprintf("Error fetching user: %s", err))
		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.WithFields(fields).Info(fmt.Sprintf("User %s exists...", user.Username.String))
	go func() {
		err = env.HelloProfileDb.DeleteUser(context.Background(), user.Email)
		if err != nil {
			log.WithFields(fields).WithError(err).Error("Error occured while deleting user")
		}
		log.WithFields(fields).Info("Successfully deactivated user")
	}()
	response := &models.SuccessResponse{
		ResponseCode:    util.SUCCESS_RESPONSE_CODE,
		ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
	}
	c.JSON(http.StatusOK, response)
	return err
}

func getValue(request, user string) string {
	if request == "" {
		return user
	}
	return request
}
