package controllers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"helloprofile.com/models"
	"helloprofile.com/persistence/orm/helloprofiledb"
)

func (env *Env) getPrimaryAddress(userID uuid.UUID, primaryAddress chan models.Address, fields log.Fields) {
	address, err := env.HelloProfileDb.GetPrimaryAddress(context.Background(), uuid.NullUUID{UUID: userID, Valid: true})
	if err != nil {
		log.WithFields(fields).Error(`No address found for user`)
		primaryAddress <- models.Address{}
	} else {
		primaryAddress <- models.Address{
			ID:      address.ID,
			Street:  address.Street,
			City:    address.City,
			State:   address.State.String,
			Country: address.Country.String,
		}
	}
}
func (env *Env) getUserAddresses(userID uuid.UUID, addresses chan []models.Address, fields log.Fields) {
	log.WithFields(fields).Info(`Getting the primary address for the user`)
	var addressList []models.Address
	dbAddresses, err := env.HelloProfileDb.GetUserAddresses(context.Background(), uuid.NullUUID{UUID: userID, Valid: true})
	if err != nil {
		log.WithFields(fields).WithError(err).Error(`Error occured fetching address for user`)
		addresses <- addressList
	}

	for _, value := range dbAddresses {
		address := models.Address{
			ID:      value.ID,
			Street:  value.Street,
			City:    value.City,
			State:   value.State.String,
			Country: value.Country.String,
		}
		addressList = append(addressList, address)
	}
	addresses <- addressList
}

func (env *Env) getProfiles(userID uuid.UUID, profiles chan []models.Profile, fields log.Fields) {
	addresses := make(chan []models.Address)
	go env.getUserAddresses(userID, addresses, fields)
	userAddress := <-addresses

	go func(userAddresses []models.Address) {
		log.WithFields(fields).Info(`Fetching profiles for user...`)
		var profileList []models.Profile
		dbProfiles, err := env.HelloProfileDb.GetProfiles(context.Background(), userID)
		if err != nil {
			log.WithFields(fields).WithError(err).Error(`Error occured while fetching profiles for user...`)
			profiles <- profileList
		}
		for _, value := range dbProfiles {
			address := new(models.Address)
			for _, addressValue := range userAddresses {
				if addressValue.ID == value.AddressID.UUID {
					address = &addressValue
				}
			}
			profileList = append(profileList, models.Profile{
				ID:             value.ID,
				Status:         value.Status,
				ProfileName:    value.ProfileName,
				Fullname:       value.Fullname,
				Title:          value.Title,
				Bio:            value.Bio,
				Company:        value.Company,
				CompanyAddress: value.CompanyAddress,
				ImageUrl:       value.ImageUrl.String,
				Phone:          value.Phone,
				Email:          value.Email,
				Address:        *address,
				Website:        value.Website.String,
				IsDefault:      value.IsDefault,
				Color:          value.Color.Int32,
			})
		}
		profiles <- profileList
	}(userAddress)
}

// saveLogin is used to log a login request that failed with incorrect password or was successful
func (env *Env) saveLogin(createParams helloprofiledb.CreateUserLoginParams) error {
	fields := log.Fields{"microservice": "helloprofile.service", "function": "saveLogin"}

	userLogin, err := env.HelloProfileDb.CreateUserLogin(context.Background(), createParams)
	if err != nil {
		log.WithFields(fields).WithError(err).Error("Error occured saving user login")
		return err
	}
	log.WithFields(fields).Info(fmt.Sprintf("Successfully saved user login, user login id: %s", userLogin.ID))
	return err
}

// func (env *Env) getSocails(profileID uuid.UUID, socials chan []models)
