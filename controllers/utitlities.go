package controllers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"helloprofile.com/models"
	"helloprofile.com/persistence/orm/helloprofiledb"
)

func (env *Env) getProfiles(userID uuid.UUID, profiles chan []models.Profile, fields log.Fields) {

	socials := make(chan []models.Socials)
	linkContents := make(chan []models.Content)
	articleContents := make(chan []models.Content)
	videoContents := make(chan []models.Content)
	audioContents := make(chan []models.Content)
	meetingContents := make(chan []models.Content)
	eventsContents := make(chan []models.Content)
	formContents := make(chan []models.Content)
	basicBlock := make(chan models.Basic)
	contactBlock := make(chan models.ContactBlock)

	go func() {
		log.WithFields(fields).Info(`Fetching profiles for user...`)
		var profileList []models.Profile
		dbProfiles, err := env.HelloProfileDb.GetProfiles(context.Background(), userID)
		if err != nil {
			log.WithFields(fields).WithError(err).Error(`Error occured while fetching profiles for user...`)
			profiles <- profileList
		}
		for _, value := range dbProfiles {
			go env.getBasicBlock(value.BasicBlockID.UUID, basicBlock, fields)
			go env.getContactBlock(value.ContactBlockID.UUID, contactBlock, fields)
			go env.getSocails(value.ID, socials, fields)
			go env.getContents(value.ID, linkContents, articleContents, videoContents, audioContents, meetingContents, eventsContents, formContents, fields)
			profileList = append(profileList, models.Profile{
				ID:           value.ID,
				Status:       value.Status,
				ProfileName:  value.ProfileName,
				PageColor:    value.PageColor,
				Font:         value.Font,
				IsDefault:    value.IsDefault,
				Basic:        <-basicBlock,
				ContactBlock: <-contactBlock,
				Socials:      <-socials,
				Links:        <-linkContents,
				Articles:     <-articleContents,
				Videos:       <-videoContents,
				Audios:       <-audioContents,
				Forms:        <-formContents,
				Meetings:     <-meetingContents,
				Events:       <-eventsContents,
			})
		}
		profiles <- profileList
	}()
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
		return
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

func (env *Env) getBasicBlock(id uuid.UUID, basicBlock chan models.Basic, fields log.Fields) {
	log.WithFields(fields).Info(`Getting the basic block for the user profile %v`, id)
	basic := new(models.Basic)
	dbBasic, err := env.HelloProfileDb.GetBasicBlock(context.Background(), id)
	if err != nil {
		log.WithFields(fields).WithError(err).Error(`Error occured fetching basic block for user profile %v`, id)
		basicBlock <- *basic
		return
	}
	basic.Bio = dbBasic.Bio
	basic.Fullname = dbBasic.Fullname
	basic.Title = dbBasic.Title
	basic.CoverColour = dbBasic.CoverColour.String
	basic.CoverPhotoUrl = dbBasic.CoverPhotoUrl.String
	basic.ProfilePhotoUrl = dbBasic.ProfilePhotoUrl.String
	basic.ID = dbBasic.ID
	basicBlock <- *basic
}

func (env *Env) getContactBlock(id uuid.UUID, contactBlock chan models.ContactBlock, fields log.Fields) {
	log.WithFields(fields).Info(`Getting the contact block for the user profile %v`, id)
	contact := new(models.ContactBlock)
	dbContact, err := env.HelloProfileDb.GetContactBlock(context.Background(), id)
	if err != nil {
		log.WithFields(fields).WithError(err).Error(`Error occured fetching contact block for user profile %v`, id)
		contactBlock <- *contact
		return
	}
	contact.Address = dbContact.Address
	contact.Email = dbContact.Email
	contact.ID = dbContact.ID
	contact.Phone = dbContact.Phone
	contact.Website = dbContact.Website
	contactBlock <- *contact
}
