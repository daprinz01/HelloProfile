package controllers

import (
	"authengine/models"
	"authengine/persistence/orm/authdb"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// GetRoles is used get roles
func (env *Env) GetRoles(c echo.Context) (err error) {
	log.Println("Get roles request received...")
	errorResponse := new(models.Errormessage)

	roles, err := env.AuthDb.GetRoles(context.Background())
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Roles not found"
		log.Println(fmt.Sprintf("Roles not found: %s", err))

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.Println("Successfully retrieved role...")
	roleResponse := make([]models.Role, len(roles))
	for index, value := range roles {
		role := models.Role{
			Role:        value.Name,
			Description: value.Description,
		}
		roleResponse[index] = role
	}
	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
		ResponseDetails: roleResponse,
	}
	c.JSON(http.StatusOK, response)
	return err
}

// GetRolesByApplication is used get roles
func (env *Env) GetRolesByApplication(c echo.Context) (err error) {
	log.Println("Get roles by application request received...")
	errorResponse := new(models.Errormessage)

	application := c.Param("application")

	if err != nil {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "Application not specified"
		log.Println("Application not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	log.Println(fmt.Sprintf("Application: %s", application))

	roles, err := env.AuthDb.GetRolesByApplication(context.Background(), application)
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Roles not found"
		log.Println(fmt.Sprintf("Roles not found: %s", err))
		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.Println("Successfully retrieved role...")
	roleResponse := make([]models.Role, len(roles))
	for index, value := range roles {
		role := models.Role{
			Role:        value.Name,
			Description: value.Description,
		}
		roleResponse[index] = role
	}
	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
		ResponseDetails: roleResponse,
	}
	c.JSON(http.StatusOK, response)
	return err
}

// AddApplicationRole Add Role to applications
func (env *Env) AddApplicationRole(c echo.Context) (err error) {
	log.Println("Add application to role request received...")

	errorResponse := new(models.Errormessage)

	role := c.Param("role")

	log.Println(fmt.Sprintf("Role: %s", role))
	if err != nil {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "Role not specified"
		log.Println("Role not specified")
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	application := c.Param("application")

	log.Println(fmt.Sprintf("Application: %s", application))
	if err != nil {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "Application not specified"
		log.Println("Application not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	dbRole, err := env.AuthDb.AddApplicationRole(context.Background(), authdb.AddApplicationRoleParams{
		Name:   strings.ToLower(application),
		Name_2: strings.ToLower(role),
	})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not add role to application. Not found"
		log.Println(fmt.Sprintf("Error occured adding  role to application: %s", err))

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.Println(fmt.Sprintf("Successfully added role to application: %v", dbRole))

	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	c.JSON(http.StatusOK, response)
	return err
}

// GetRole is used get role
func (env *Env) GetRole(c echo.Context) (err error) {
	log.Println("Get role request received...")

	errorResponse := new(models.Errormessage)

	role := c.Param("role")

	log.Println(fmt.Sprintf("Role: %s", role))
	if err != nil {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "Role not specified"
		log.Println("Role not specified")

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}

	dbRole, err := env.AuthDb.GetRole(context.Background(), strings.ToLower(role))
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Role not found"
		log.Println(fmt.Sprintf("Role not found: %s", err))

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.Println(fmt.Sprintf("Successfully retrieved role: %v", dbRole))
	roleResponse := models.Role{
		Role:        dbRole.Name,
		Description: dbRole.Description,
	}

	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
		ResponseDetails: roleResponse,
	}
	c.JSON(http.StatusOK, response)
	return err
}

// AddRole is used add role
func (env *Env) AddRole(c echo.Context) (err error) {
	log.Println("Add role request received...")

	errorResponse := new(models.Errormessage)

	request := new(models.Role)
	if err = c.Bind(request); err != nil {
		log.Println(fmt.Sprintf("Error occured while trying to marshal request: %s", err))
		errorResponse.Errorcode = "02"
		errorResponse.ErrorMessage = "Invalid request"
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	dbRole, err := env.AuthDb.CreateRole(context.Background(), authdb.CreateRoleParams{
		Name:        strings.ToLower(request.Role),
		Description: strings.ToLower(request.Description),
	})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not add role. Duplicate found"
		log.Println(fmt.Sprintf("Error occured adding new role: %s", err))

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.Println(fmt.Sprintf("Successfully added role: %v", dbRole))

	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	c.JSON(http.StatusOK, response)
	return err
}

// UpdateRole is used add role
func (env *Env) UpdateRole(c echo.Context) (err error) {
	log.Println("Update role request received...")

	errorResponse := new(models.Errormessage)

	role := c.Param("role")

	log.Println(fmt.Sprintf("Role: %s", role))
	if err != nil {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "Role not specified"
		log.Println("Role not specified")

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}

	request := new(models.Role)
	if err = c.Bind(request); err != nil {
		log.Println(fmt.Sprintf("Error occured while trying to marshal request: %s", err))
		errorResponse.Errorcode = "02"
		errorResponse.ErrorMessage = "Invalid request"
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	dbRole, err := env.AuthDb.UpdateRole(context.Background(), authdb.UpdateRoleParams{
		Name:        strings.ToLower(request.Role),
		Description: strings.ToLower(request.Description),
		Name_2:      strings.ToLower(role),
	})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not update role. Not found"
		log.Println(fmt.Sprintf("Error occured updating new role: %s", err))

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.Println(fmt.Sprintf("Successfully updated role: %v", dbRole))

	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	c.JSON(http.StatusOK, response)
	return err
}

// DeleteRole is used add role
func (env *Env) DeleteRole(c echo.Context) (err error) {
	log.Println("Delete role request received...")

	errorResponse := new(models.Errormessage)

	role := c.Param("role")

	log.Println(fmt.Sprintf("Role: %s", role))
	if err != nil {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "Role not specified"
		log.Println(fmt.Sprintf("Role not specified %s", err))

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}

	err = env.AuthDb.DeleteRoles(context.Background(), strings.ToLower(role))
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not delete role. Role not found"
		log.Println(fmt.Sprintf("Error occured deleting  role: %s", err))

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.Println("Successfully deleted role")

	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	c.JSON(http.StatusOK, response)
	return err
}
