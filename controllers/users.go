package controllers

import (
	"authengine/models"
	"authengine/persistence/orm/authdb"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// CheckAvailability is used to check user availablity
func (env *Env) CheckAvailability(c echo.Context) (err error) {
	log.Println("Check availability Request received")

	errorResponse := new(models.Errormessage)

	username := c.Param("username")
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Username not specified"
		log.Println("Username not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	log.Println(fmt.Sprintf("Username: %s", username))

	user, err := env.AuthDb.GetUser(context.Background(), sql.NullString{String: strings.ToLower(username), Valid: true})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "User does not exist"
		log.Println(fmt.Sprintf("Error fetching user: %s", err))
		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.Println(fmt.Sprintf("User %s exists...", user.Username.String))

	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	c.JSON(http.StatusOK, response)
	return err
}

// GetUser is used to fetch user details
func (env *Env) GetUser(c echo.Context) (err error) {
	log.Println("Get User Request received")

	// file, fileHeader, err := r.FormFile("request.AttachmentName[i]")

	// file, err := os.Create(fmt.Sprintf("%s%s", attachmentPath, request.AttachmentName[i].FileName))
	// file.WriteString()

	username := c.Param("username")
	errorResponse := new(models.Errormessage)

	if err != nil {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "Username not specified"
		log.Println("Username not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	log.Println(fmt.Sprintf("Username: %s", username))

	user, err := env.AuthDb.GetUser(context.Background(), sql.NullString{String: strings.ToLower(username), Valid: true})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "User does not exist"
		log.Println(fmt.Sprintf("Error fetching user: %s", err))
		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.Println(fmt.Sprintf("User %s exists...", user.Username.String))

	response := &models.SuccessResponse{
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
	c.JSON(http.StatusOK, response)
	return err
}

// GetUsers is used to fetch user details. This is an admin function. User must be an admin to access this function
func (env *Env) GetUsers(c echo.Context) (err error) {
	log.Println("Get User Request received")

	errorResponse := new(models.Errormessage)

	users, err := env.AuthDb.GetUsers(context.Background())
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Users does not exist"
		log.Println(fmt.Sprintf("Error fetching user: %s", err))
		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	userResponse := make([]models.UserDetail, len(users))
	for index, user := range users {
		tempUser := models.UserDetail{
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
		}
		userResponse[index] = tempUser
	}
	log.Println(fmt.Sprintf("Returning %d users...", len(users)))

	resetResponse := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
		ResponseDetails: &userResponse,
	}
	c.JSON(http.StatusOK, resetResponse)
	return err
}

// UpdateUser is used to update User information. It can be used to update user details and timezone details as required. Only pass the details to be updated. Email or username is mandatory.
func (env *Env) UpdateUser(c echo.Context) (err error) {

	log.Println("Update user request received...")
	errorResponse := new(models.Errormessage)

	request := new(models.UserDetail)
	if err = c.Bind(request); err != nil {
		log.Println(fmt.Sprintf("Error occured while trying to marshal request: %s", err))
		errorResponse.Errorcode = "02"
		errorResponse.ErrorMessage = "Invalid request"
		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	log.Println("Checking if user exist...")
	user, err := env.AuthDb.GetUser(context.Background(), sql.NullString{String: strings.ToLower(request.Email), Valid: true})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "User does not exist"
		log.Println(fmt.Sprintf("Error fetching user: %s", err))
		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.Println(fmt.Sprintf("User %s exists...", user.Username.String))
	go func() {
		log.Println("Updating user...")
		// Check all the user parameters passed and consolidate with existing user record

		_, err := env.AuthDb.UpdateUser(context.Background(), authdb.UpdateUserParams{
			Address:                   sql.NullString{String: getValue(request.Address, user.Address.String), Valid: true},
			City:                      sql.NullString{String: getValue(request.City, user.City.String), Valid: true},
			Country:                   sql.NullString{String: getValue(request.Country, user.Country.String), Valid: true},
			CreatedAt:                 user.CreatedAt,
			Email:                     user.Email,
			Firstname:                 sql.NullString{String: getValue(request.Firstname, user.Firstname.String), Valid: true},
			ImageUrl:                  sql.NullString{String: getValue(request.ProfilePicture, user.ProfilePicture.String), Valid: true},
			IsActive:                  user.IsActive,
			IsEmailConfirmed:          user.IsEmailConfirmed,
			IsLockedOut:               user.IsLockedOut,
			IsPasswordSystemGenerated: user.IsPasswordSystemGenerated,
			Lastname:                  sql.NullString{String: getValue(request.Lastname, user.Lastname.String), Valid: true},
			Password:                  user.Password,
			State:                     sql.NullString{String: getValue(request.State, user.State.String), Valid: true},
			Username:                  sql.NullString{String: getValue(request.Username, user.Username.String), Valid: true},
			Phone:                     sql.NullString{String: getValue(request.Phone, user.Phone.String), Valid: true},
			Username_2:                sql.NullString{String: user.Email, Valid: true},
		})

		if err != nil {
			log.Println(fmt.Sprintf("Error occured updating user information: %s", err))

		}
		log.Println("Successfully updated user details")
		// LanguageName              string    `json:"language_name"`
		// RoleName                  string    `json:"role_name"`
		// TimezoneName              string    `json:"timezone_name"`
		// Zone                      string    `json:"zone"`
		// ProviderName              string    `json:"provider_name"`
		// ClientID                  string    `json:"client_id"`
		// ClientSecret              string    `json:"client_secret"`
		// ProviderLogo              string    `json:"provider_logo"`
		log.Println("Updating User Timezone")
		if request.TimezoneName != "" {
			_, err := env.AuthDb.UpdateUserTimezone(context.Background(), authdb.UpdateUserTimezoneParams{
				Username:   sql.NullString{String: user.Email, Valid: true},
				Name:       request.TimezoneName,
				Username_2: sql.NullString{String: user.Email, Valid: true},
			})
			if err != nil {
				log.Println(fmt.Sprintf("Error occured updating timezone information: %s", err))

			}
			log.Println("Successfully updated user details")
		}
	}()
	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	c.JSON(http.StatusOK, response)
	return err
}

// UpdateUserRole Add Role to applications
func (env *Env) UpdateUserRole(c echo.Context) (err error) {
	log.Println("Update user's role request received...")

	errorResponse := new(models.Errormessage)

	role := c.Param("newRole")

	log.Println(fmt.Sprintf("New Role: %s", role))
	if role == "" {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "New Role not specified"
		log.Println("Role not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	oldRole := c.Param("oldRole")

	log.Println(fmt.Sprintf("Old Role: %s", role))
	if oldRole == "" {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Old Role not specified"
		log.Println("Role not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	username := c.Param("username")
	if username == "" {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Username not specified"
		log.Println("Username not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	log.Println(fmt.Sprintf("Username: %s", username))

	dbRole, err := env.AuthDb.UpdateUserRole(context.Background(), authdb.UpdateUserRoleParams{
		Username:   sql.NullString{String: strings.ToLower(username), Valid: true},
		Username_2: sql.NullString{String: strings.ToLower(username), Valid: true},
		Name:       strings.ToLower(role),
		Name_2:     strings.ToLower(oldRole),
	})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not update user to role. Not found"
		log.Println(fmt.Sprintf("Error occured updating  user role: %s", err))

		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.Println(fmt.Sprintf("Successfully updated user role: %v", dbRole))

	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
	}
	c.JSON(http.StatusOK, response)
	return err
}

// AddUserToRole Add user to a role. This increases the roles the user is being added to including the previous roles. At login a role must be selected else the the default role guest is selected for the user
func (env *Env) AddUserToRole(c echo.Context) (err error) {
	log.Println("Add user to role request received...")

	errorResponse := new(models.Errormessage)

	role := c.Param("role")

	log.Println(fmt.Sprintf("Role: %s", role))
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Role not specified"
		log.Println("Role not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}

	username := c.Param("username")
	if err != nil {
		errorResponse.Errorcode = "15"
		errorResponse.ErrorMessage = "Username not specified"
		log.Println("Username not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	log.Println(fmt.Sprintf("Username: %s", username))

	dbRole, err := env.AuthDb.AddUserRole(context.Background(), authdb.AddUserRoleParams{
		Username: sql.NullString{String: strings.ToLower(username), Valid: true},
		Name:     strings.ToLower(role),
	})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Could not add user to role. Not found"
		log.Println(fmt.Sprintf("Error occured adding  user to role: %s", err))

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

// DeleteUser is used to disable a users account
func (env *Env) DeleteUser(c echo.Context) (err error) {
	log.Println("Delete user Request received")

	// file, fileHeader, err := r.FormFile("request.AttachmentName[i]")

	// file, err := os.Create(fmt.Sprintf("%s%s", attachmentPath, request.AttachmentName[i].FileName))
	// file.WriteString()

	username := c.Param("username")
	errorResponse := new(models.Errormessage)

	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "Username not specified"
		log.Println("Username not specified")

		c.JSON(http.StatusBadRequest, errorResponse)
		return err
	}
	log.Println(fmt.Sprintf("Username: %s", username))

	user, err := env.AuthDb.GetUser(context.Background(), sql.NullString{String: strings.ToLower(username), Valid: true})
	if err != nil {
		errorResponse.Errorcode = "03"
		errorResponse.ErrorMessage = "User does not exist"
		log.Println(fmt.Sprintf("Error fetching user: %s", err))
		c.JSON(http.StatusNotFound, errorResponse)
		return err
	}
	log.Println(fmt.Sprintf("User %s exists...", user.Username.String))
	go func() {
		err = env.AuthDb.DeleteUser(context.Background(), user.Email)
		if err != nil {
			log.Println(fmt.Sprintf("Error occured while deleting user"))
		}
		log.Println("Successfully deactivated user")
	}()
	response := &models.SuccessResponse{
		ResponseCode:    "00",
		ResponseMessage: "Success",
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
