package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
	"persianblack.com/authengine/models"
	"persianblack.com/authengine/persistence/orm/authdb"
	"persianblack.com/authengine/util"

	"github.com/labstack/echo/v4"
)

// GetRoles is used get roles
func (env *Env) GetRoles(c echo.Context) (err error) {

	errorResponse := new(models.Errormessage)
	application := c.Param("application")
	if application == "" {
		errorResponse.Errorcode = util.APPLICATION_NOT_SPECIFIED_ERROR_CODE
		errorResponse.ErrorMessage = util.APPLICATION_NOT_SPECIFIED_ERROR_MESSAGE
		log.WithField("microservice", "persian.black.authengine.service").Error("Calling application not specified")
		c.JSON(http.StatusBadRequest, errorResponse)
		return nil
	}
	fields := log.Fields{"microservice": "persian.black.authengine.service", "application": application}
	log.WithFields(fields).Info("Get roles request received...")
	if c.QueryParam("email") != "" {
		log.WithFields(fields).Info(fmt.Sprintf("Getting roles for user %s", c.QueryParam("email")))
		roles, err := env.AuthDb.GetUserRoles(context.Background(), sql.NullString{String: strings.ToLower(c.QueryParam("email")), Valid: true})
		if err != nil {
			errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
			errorResponse.ErrorMessage = util.NO_RECORD_FOUND_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Roles not found for user")

			c.JSON(http.StatusNotFound, errorResponse)
			return err
		}
		log.WithFields(fields).Info("Successfully retrieved role...")
		roleResponse := make([]models.Role, len(roles))
		for index, value := range roles {
			role := models.Role{
				Role: value,
			}
			roleResponse[index] = role
		}
		response := &models.SuccessResponse{
			ResponseCode:    util.SUCCESS_RESPONSE_CODE,
			ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
			ResponseDetails: roleResponse,
		}
		c.JSON(http.StatusOK, response)

	} else {
		roles, err := env.AuthDb.GetRoles(context.Background())
		if err != nil {
			errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
			errorResponse.ErrorMessage = util.NO_RECORD_FOUND_ERROR_MESSAGE
			log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Roles not found")

			c.JSON(http.StatusNotFound, errorResponse)
			return err
		}
		log.WithFields(fields).Info("Successfully retrieved role...")
		roleResponse := make([]models.Role, len(roles))
		for index, value := range roles {
			role := models.Role{
				Role:        value.Name,
				Description: value.Description,
			}
			roleResponse[index] = role
		}
		response := &models.SuccessResponse{
			ResponseCode:    util.SUCCESS_RESPONSE_CODE,
			ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
			ResponseDetails: roleResponse,
		}
		c.JSON(http.StatusOK, response)
		return err
	}
	return err
}

// GetRolesByApplication is used get roles
func (env *Env) GetRolesByApplication(c echo.Context) (err error) {
	log.Println("Get roles by application request received...")
	errorResponse := new(models.Errormessage)

	application := c.Param("referenceApplication")
	if application == "" {
		errorResponse.Errorcode = util.APPLICATION_NOT_SPECIFIED_ERROR_CODE
		errorResponse.ErrorMessage = util.APPLICATION_NOT_SPECIFIED_ERROR_MESSAGE
		log.WithField("microservice", "persian.black.authengine.service").Error("Calling application not specified")
		c.JSON(http.StatusBadRequest, errorResponse)
		return nil
	}
	fields := log.Fields{"microservice": "persian.black.authengine.service", "application": application}
	log.WithFields(fields).Info(fmt.Sprintf("Application: %s", application))

	roles, err := env.AuthDb.GetRolesByApplication(context.Background(), application)
	if err != nil {
		errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
		errorResponse.ErrorMessage = util.NO_RECORD_FOUND_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Roles not found")
		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.WithFields(fields).Info("Successfully retrieved role...")
	roleResponse := make([]models.Role, len(roles))
	for index, value := range roles {
		role := models.Role{
			Role:        value.Name,
			Description: value.Description,
		}
		roleResponse[index] = role
	}
	response := &models.SuccessResponse{
		ResponseCode:    util.SUCCESS_RESPONSE_CODE,
		ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
		ResponseDetails: roleResponse,
	}
	c.JSON(http.StatusOK, response)
	return err
}

// AddApplicationRole Add Role to applications
func (env *Env) AddApplicationRole(c echo.Context) (err error) {

	errorResponse := new(models.Errormessage)
	application := c.Param("referenceApplication")
	if application == "" {
		errorResponse.Errorcode = util.APPLICATION_NOT_SPECIFIED_ERROR_CODE
		errorResponse.ErrorMessage = util.APPLICATION_NOT_SPECIFIED_ERROR_MESSAGE
		log.WithField("microservice", "persian.black.authengine.service").Error("Calling application not specified")
		c.JSON(http.StatusBadRequest, errorResponse)
		return nil
	}
	fields := log.Fields{"microservice": "persian.black.authengine.service", "application": application}
	log.WithFields(fields).Info("Add application to role request received...")
	role := c.Param("role")

	log.WithFields(fields).Info(fmt.Sprintf("Role: %s", role))
	if err != nil {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Role not specified")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	log.WithFields(fields).Info(fmt.Sprintf("Application: %s", application))
	if err != nil {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.APPLICATION_NOT_SPECIFIED_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error(util.APPLICATION_NOT_SPECIFIED_ERROR_MESSAGE)

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	dbRole, err := env.AuthDb.AddApplicationRole(context.Background(), authdb.AddApplicationRoleParams{
		Name:   strings.ToLower(application),
		Name_2: strings.ToLower(role),
	})
	if err != nil {
		errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
		errorResponse.ErrorMessage = util.NO_RECORD_FOUND_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured adding  role to application")

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

// GetRole is used get role
func (env *Env) GetRole(c echo.Context) (err error) {

	errorResponse := new(models.Errormessage)
	application := c.Param("application")
	if application == "" {
		errorResponse.Errorcode = util.APPLICATION_NOT_SPECIFIED_ERROR_CODE
		errorResponse.ErrorMessage = util.APPLICATION_NOT_SPECIFIED_ERROR_MESSAGE
		log.WithField("microservice", "persian.black.authengine.service").Error("Calling application not specified")
		c.JSON(http.StatusBadRequest, errorResponse)
		return nil
	}
	fields := log.Fields{"microservice": "persian.black.authengine.service", "application": application}
	log.WithFields(fields).Info("Get role request received...")
	role := c.Param("role")

	log.WithFields(fields).Info(fmt.Sprintf("Role: %s", role))
	if err != nil {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Role not specified")

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}

	dbRole, err := env.AuthDb.GetRole(context.Background(), strings.ToLower(role))
	if err != nil {
		errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
		errorResponse.ErrorMessage = util.NO_RECORD_FOUND_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Role not found")

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.WithFields(fields).Info(fmt.Sprintf("Successfully retrieved role: %v", dbRole))
	roleResponse := models.Role{
		Role:        dbRole.Name,
		Description: dbRole.Description,
	}

	response := &models.SuccessResponse{
		ResponseCode:    util.SUCCESS_RESPONSE_CODE,
		ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
		ResponseDetails: roleResponse,
	}
	c.JSON(http.StatusOK, response)
	return err
}

// AddRole is used add role
func (env *Env) AddRole(c echo.Context) (err error) {

	errorResponse := new(models.Errormessage)
	application := c.Param("application")
	if application == "" {
		errorResponse.Errorcode = util.APPLICATION_NOT_SPECIFIED_ERROR_CODE
		errorResponse.ErrorMessage = util.APPLICATION_NOT_SPECIFIED_ERROR_MESSAGE
		log.WithField("microservice", "persian.black.authengine.service").Error("Calling application not specified")
		c.JSON(http.StatusBadRequest, errorResponse)
		return nil
	}
	fields := log.Fields{"microservice": "persian.black.authengine.service", "application": application}
	log.WithFields(fields).Info("Add role request received...")
	request := new(models.Role)
	if err = c.Bind(request); err != nil {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to marshal request")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	dbRole, err := env.AuthDb.CreateRole(context.Background(), authdb.CreateRoleParams{
		Name:        strings.ToLower(request.Role),
		Description: strings.ToLower(request.Description),
	})
	if err != nil {
		errorResponse.Errorcode = util.DUPLICATE_RECORD_ERROR_CODE
		errorResponse.ErrorMessage = util.DUPLICATE_RECORD_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured adding new role")

		c.JSON(http.StatusNotModified, errorResponse)
		return err
	}
	log.WithFields(fields).Info(fmt.Sprintf("Successfully added role: %v", dbRole))

	response := &models.SuccessResponse{
		ResponseCode:    util.SUCCESS_RESPONSE_CODE,
		ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
	}
	c.JSON(http.StatusOK, response)
	return err
}

// UpdateRole is used add role
func (env *Env) UpdateRole(c echo.Context) (err error) {

	errorResponse := new(models.Errormessage)
	application := c.Param("application")
	if application == "" {
		errorResponse.Errorcode = util.APPLICATION_NOT_SPECIFIED_ERROR_CODE
		errorResponse.ErrorMessage = util.APPLICATION_NOT_SPECIFIED_ERROR_MESSAGE
		log.WithField("microservice", "persian.black.authengine.service").Error("Calling application not specified")
		c.JSON(http.StatusBadRequest, errorResponse)
		return nil
	}
	fields := log.Fields{"microservice": "persian.black.authengine.service", "application": application}
	log.WithFields(fields).Info("Update role request received...")
	role := c.Param("role")

	log.WithFields(fields).Info(fmt.Sprintf("Role: %s", role))
	if err != nil {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Role not specified")

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}

	request := new(models.Role)
	if err = c.Bind(request); err != nil {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured while trying to marshal request")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	dbRole, err := env.AuthDb.UpdateRole(context.Background(), authdb.UpdateRoleParams{
		Name:        strings.ToLower(request.Role),
		Description: strings.ToLower(request.Description),
		Name_2:      strings.ToLower(role),
	})
	if err != nil {
		errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
		errorResponse.ErrorMessage = util.NO_RECORD_FOUND_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured updating new role")

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.WithFields(fields).Info(fmt.Sprintf("Successfully updated role: %v", dbRole))

	response := &models.SuccessResponse{
		ResponseCode:    util.SUCCESS_RESPONSE_CODE,
		ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
	}
	c.JSON(http.StatusOK, response)
	return err
}

// DeleteRole is used add role
func (env *Env) DeleteRole(c echo.Context) (err error) {

	errorResponse := new(models.Errormessage)
	application := c.Param("application")
	if application == "" {
		errorResponse.Errorcode = util.APPLICATION_NOT_SPECIFIED_ERROR_CODE
		errorResponse.ErrorMessage = util.APPLICATION_NOT_SPECIFIED_ERROR_MESSAGE
		log.WithField("microservice", "persian.black.authengine.service").Error("Calling application not specified")
		c.JSON(http.StatusBadRequest, errorResponse)
		return nil
	}
	fields := log.Fields{"microservice": "persian.black.authengine.service", "application": application}

	log.WithFields(fields).Info("Delete role request received...")
	role := c.Param("role")

	log.WithFields(fields).Info(fmt.Sprintf("Role: %s", role))
	if err != nil {
		errorResponse.Errorcode = util.MODEL_VALIDATION_ERROR_CODE
		errorResponse.ErrorMessage = util.MODEL_VALIDATION_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Role not specified")

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}

	err = env.AuthDb.DeleteRoles(context.Background(), strings.ToLower(role))
	if err != nil {
		errorResponse.Errorcode = util.NO_RECORD_FOUND_ERROR_CODE
		errorResponse.ErrorMessage = util.NO_RECORD_FOUND_ERROR_MESSAGE
		log.WithFields(fields).WithError(err).WithFields(log.Fields{"responseCode": errorResponse.Errorcode, "responseDescription": errorResponse.ErrorMessage}).Error("Error occured deleting  role")

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.WithFields(fields).Info("Successfully deleted role")

	response := &models.SuccessResponse{
		ResponseCode:    util.SUCCESS_RESPONSE_CODE,
		ResponseMessage: util.SUCCESS_RESPONSE_MESSAGE,
	}
	c.JSON(http.StatusOK, response)
	return err
}
