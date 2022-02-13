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
	socials := make(chan []models.Socials)
	linkContents := make(chan []models.Content)
	articleContents := make(chan []models.Content)
	videoContents := make(chan []models.Content)
	audioContents := make(chan []models.Content)
	meetingContents := make(chan []models.Content)
	eventsContents := make(chan []models.Content)
	formContents := make(chan []models.Content)
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
			go env.getSocails(value.ID, socials, fields)
			go env.getContents(value.ID, linkContents, articleContents, videoContents, audioContents, meetingContents, eventsContents, formContents, fields)
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
				Website:        value.Website.String,
				IsDefault:      value.IsDefault,
				Color:          value.Color.Int32,
				Address:        *address,
				Socials:        <-socials,
				Links:          <-linkContents,
				Articles:       <-articleContents,
				Videos:         <-videoContents,
				Audios:         <-audioContents,
				Forms:          <-formContents,
				Meetings:       <-meetingContents,
				Events:         <-eventsContents,
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

func (env *Env) getSocails(profileID uuid.UUID, socials chan []models.Socials, fields log.Fields) {
	log.WithFields(fields).Info(`Getting the socials for the user profile`)
	var socialsList []models.Socials
	dbSocials, err := env.HelloProfileDb.GetProfileSocials(context.Background(), profileID)
	if err != nil {
		log.WithFields(fields).WithError(err).Error(`Error occured fetching socials for user profile`)
		socials <- socialsList
	}

	for _, value := range dbSocials {
		social := models.Socials{
			Username:    value.Username,
			Name:        value.Name.String,
			Placeholder: value.Placeholder.String,
			ImageUrl:    value.ImageUrl.String,
			Order:       value.Order,
		}
		socialsList = append(socialsList, social)
	}
	socials <- socialsList
}

func (env *Env) getContents(profileID uuid.UUID, linkContents chan []models.Content, articleContents chan []models.Content, videoContents chan []models.Content, audioContents chan []models.Content, meetingContents chan []models.Content, eventsContents chan []models.Content, formContents chan []models.Content, fields log.Fields) {
	log.WithFields(fields).Info(`Getting the socials for the user profile`)
	var linkContentList []models.Content
	var articleContentList []models.Content
	var videoContentList []models.Content
	var audioContentList []models.Content
	var meetingContentList []models.Content
	var eventsContentList []models.Content
	var formContentList []models.Content
	dbContents, err := env.HelloProfileDb.GetProfileContents(context.Background(), profileID)
	if err != nil {
		log.WithFields(fields).WithError(err).Error(`Error occured fetching contents for user profile`)
		linkContents <- linkContentList
		articleContents <- articleContentList
		videoContents <- videoContentList
		audioContents <- audioContentList
		meetingContents <- meetingContentList
		eventsContents <- eventsContentList
		formContents <- formContentList
		return
	}
	dbCallToActions, err := env.HelloProfileDb.GetCallToActions(context.Background())
	if err != nil {
		log.WithFields(fields).WithError(err).Error(`Error occured fetching call to action for profile contents`)
		linkContents <- linkContentList
		articleContents <- articleContentList
		videoContents <- videoContentList
		audioContents <- audioContentList
		meetingContents <- meetingContentList
		eventsContents <- eventsContentList
		formContents <- formContentList
		return
	}
	dbContentTypes, err := env.HelloProfileDb.GetAllContentTypes(context.Background())
	if err != nil {
		log.WithFields(fields).WithError(err).Error(`Error occured fetching call to action for profile contents`)
		linkContents <- linkContentList
		articleContents <- articleContentList
		videoContents <- videoContentList
		audioContents <- audioContentList
		meetingContents <- meetingContentList
		eventsContents <- eventsContentList
		formContents <- formContentList
		return
	}
	//Sort contents
	//Get links
	for _, value := range dbContentTypes {
		if value.Type == "links" {
			go env.sortContents(dbContents, dbCallToActions, linkContents, value.ID)
		} else if value.Type == "articles" {
			go env.sortContents(dbContents, dbCallToActions, articleContents, value.ID)
		} else if value.Type == "embedded videos" {
			go env.sortContents(dbContents, dbCallToActions, videoContents, value.ID)
		} else if value.Type == "embedded audios" {
			go env.sortContents(dbContents, dbCallToActions, audioContents, value.ID)
		} else if value.Type == "forms" {
			go env.sortContents(dbContents, dbCallToActions, formContents, value.ID)
		} else if value.Type == "meetings" {
			go env.sortContents(dbContents, dbCallToActions, meetingContents, value.ID)
		} else if value.Type == "events" {
			go env.sortContents(dbContents, dbCallToActions, eventsContents, value.ID)
		}
	}

}

func (env *Env) sortContents(contentResult []helloprofiledb.ProfileContent, callToActions []helloprofiledb.CallToAction, contentChannel chan []models.Content, contentTypeID uuid.UUID) {
	var contentList []models.Content
	for _, value := range contentResult {
		if value.CallToActionID == contentTypeID {
			callToAction := new(helloprofiledb.CallToAction)
			for _, action := range callToActions {
				if value.CallToActionID == action.ID {
					callToAction.DisplayName = action.DisplayName
					callToAction.ID = action.ID
					callToAction.Type = action.Type
				}
			}
			content := models.Content{
				ID:          value.ID,
				Title:       value.DisplayTitle,
				Description: value.Description,
				Url:         value.Url,
				Order:       value.Order,
				CallToAction: models.CallToAction{
					ID:          callToAction.ID,
					Type:        callToAction.Type,
					DisplayName: callToAction.DisplayName,
				},
			}
			contentList = append(contentList, content)
		}
	}
	contentChannel <- contentList
}
