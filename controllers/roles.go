package controllers

import (
	"authengine/models"
	"authengine/persistence/orm/authdb"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// GetRoles is used get roles
func (env *Env) GetRoles(w http.ResponseWriter, r *http.Request) {
	log.Println("Get roles request received...")
	var errorResponse models.Errormessage
	var err error
	roles, err := env.AuthDb.GetRoles(context.Background())
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Roles not found"
		log.Println(fmt.Sprintf("Roles not found: %s", err))

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
	}
	log.Println("Successfully retrieved role...")
	var roleResponse []models.Role
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
	responsebytes, err := json.MarshalIndent(response, "", "")
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responsebytes)
	return
}

// GetRolesByApplication is used get roles
func (env *Env) GetRolesByApplication(w http.ResponseWriter, r *http.Request) {
	log.Println("Get roles by application request received...")
	var errorResponse models.Errormessage
	var err error
	pathParams := mux.Vars(r)
	var application string
	if val2, ok := pathParams["application"]; ok {
		application = val2
		log.Println(fmt.Sprintf("Application: %s", application))
		if err != nil {
			errorResponse.Errorcode = "15"
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
	roles, err := env.AuthDb.GetRolesByApplication(context.Background(), application)
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Roles not found"
		log.Println(fmt.Sprintf("Roles not found: %s", err))

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
	}
	log.Println("Successfully retrieved role...")
	var roleResponse []models.Role
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
	responsebytes, err := json.MarshalIndent(response, "", "")
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responsebytes)
	return
}

// AddApplicationRole Add Role to applications
func (env *Env) AddApplicationRole(w http.ResponseWriter, r *http.Request) {
	log.Println("Add application to role request received...")
	pathParams := mux.Vars(r)
	var errorResponse models.Errormessage
	var err error
	var role string
	if val, ok := pathParams["role"]; ok {
		role = val
		log.Println(fmt.Sprintf("Role: %s", role))
		if err != nil {
			errorResponse.Errorcode = "15"
			errorResponse.ErrorMessage = "Role not specified"
			log.Println("Role not specified")

			response, err := json.MarshalIndent(errorResponse, "", "")
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
			return
		}
	}
	var application string
	if val, ok := pathParams["application"]; ok {
		application = val
		log.Println(fmt.Sprintf("Application: %s", application))
		if err != nil {
			errorResponse.Errorcode = "15"
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

	dbRole, err := env.AuthDb.AddApplicationRole(context.Background(), authdb.AddApplicationRoleParams{
		Name:   strings.ToLower(application),
		Name_2: strings.ToLower(role),
	})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not add role to application. Not found"
		log.Println(fmt.Sprintf("Error occured adding  role to application: %s", err))

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
	}
	log.Println(fmt.Sprintf("Successfully added role to application: %v", dbRole))

	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	responsebytes, err := json.MarshalIndent(response, "", "")
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responsebytes)
	return
}

// GetRole is used get role
func (env *Env) GetRole(w http.ResponseWriter, r *http.Request) {
	log.Println("Get role request received...")
	pathParams := mux.Vars(r)
	var errorResponse models.Errormessage
	var err error
	var role string
	if val2, ok := pathParams["role"]; ok {
		role = val2
		log.Println(fmt.Sprintf("Role: %s", role))
		if err != nil {
			errorResponse.Errorcode = "15"
			errorResponse.ErrorMessage = "Role not specified"
			log.Println("Role not specified")

			response, err := json.MarshalIndent(errorResponse, "", "")
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
			return
		}
	}

	dbRole, err := env.AuthDb.GetRole(context.Background(), strings.ToLower(role))
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Role not found"
		log.Println(fmt.Sprintf("Role not found: %s", err))

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
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
	responsebytes, err := json.MarshalIndent(response, "", "")
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responsebytes)
	return
}

// AddRole is used add role
func (env *Env) AddRole(w http.ResponseWriter, r *http.Request) {
	log.Println("Add role request received...")

	var errorResponse models.Errormessage
	var err error
	var request models.Role
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&request)
	defer r.Body.Close()
	if err != nil {
		errorResponse.Errorcode = "02"
		errorResponse.ErrorMessage = "Invalid request"
		log.Println(err)
		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)
		return
	}

	dbRole, err := env.AuthDb.CreateRole(context.Background(), authdb.CreateRoleParams{
		Name:        strings.ToLower(request.Role),
		Description: strings.ToLower(request.Description),
	})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not add role. Duplicate found"
		log.Println(fmt.Sprintf("Error occured adding new role: %s", err))

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
	}
	log.Println(fmt.Sprintf("Successfully added role: %v", dbRole))

	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	responsebytes, err := json.MarshalIndent(response, "", "")
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responsebytes)
	return
}

// UpdateRole is used add role
func (env *Env) UpdateRole(w http.ResponseWriter, r *http.Request) {
	log.Println("Update role request received...")
	pathParams := mux.Vars(r)
	var errorResponse models.Errormessage
	var err error
	var role string
	if val, ok := pathParams["role"]; ok {
		role = val
		log.Println(fmt.Sprintf("Role: %s", role))
		if err != nil {
			errorResponse.Errorcode = "15"
			errorResponse.ErrorMessage = "Role not specified"
			log.Println("Role not specified")

			response, err := json.MarshalIndent(errorResponse, "", "")
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
			return
		}
	}

	var request models.Role
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&request)
	defer r.Body.Close()
	if err != nil {
		errorResponse.Errorcode = "02"
		errorResponse.ErrorMessage = "Invalid request"
		log.Println(err)
		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)
		return
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

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
	}
	log.Println(fmt.Sprintf("Successfully updated role: %v", dbRole))

	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	responsebytes, err := json.MarshalIndent(response, "", "")
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responsebytes)
	return
}

// DeleteRole is used add role
func (env *Env) DeleteRole(w http.ResponseWriter, r *http.Request) {
	log.Println("Delete role request received...")
	pathParams := mux.Vars(r)
	var errorResponse models.Errormessage
	var err error
	var role string
	if val, ok := pathParams["role"]; ok {
		role = val
		log.Println(fmt.Sprintf("Role: %s", role))
		if err != nil {
			errorResponse.Errorcode = "15"
			errorResponse.ErrorMessage = "Role not specified"
			log.Println(fmt.Sprintf("Role not specified %s", err))

			response, err := json.MarshalIndent(errorResponse, "", "")
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
			return
		}
	}

	err = env.AuthDb.DeleteRoles(context.Background(), strings.ToLower(role))
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not delete role. Role not found"
		log.Println(fmt.Sprintf("Error occured deleting  role: %s", err))

		response, err := json.MarshalIndent(errorResponse, "", "")
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
	}
	log.Println("Successfully deleted role")

	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	responsebytes, err := json.MarshalIndent(response, "", "")
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responsebytes)
	return
}
