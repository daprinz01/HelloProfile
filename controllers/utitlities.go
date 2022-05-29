package controllers

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"helloprofile.com/models"
	"helloprofile.com/persistence/orm/helloprofiledb"
)

func (env *Env) getProfiles(userID uuid.UUID, profiles chan []models.Profile, fields log.Fields) {

	socials := make(chan []models.Socials)
	contents := make(chan []models.Content)
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
			log.WithFields(fields).Info(`dbprofile to get %v`, value)
			go env.getBasicBlock(value.BasicBlockID.UUID, basicBlock, fields)
			go env.getContactBlock(value.ContactBlockID.UUID, contactBlock, fields)
			go env.getSocails(value.ID, socials, fields)
			go env.getContents(value.ID, contents, fields)
			profileList = append(profileList, models.Profile{
				ID:           value.ID,
				Status:       value.Status,
				ProfileName:  value.ProfileName,
				PageColor:    value.PageColor,
				Font:         value.Font,
				IsDefault:    value.IsDefault,
				Url:          env.GetValue(fmt.Sprintf("%s/%s", os.Getenv("HELLOPROFILE_HOME"), value.Url.String), fmt.Sprintf("%s/%s", os.Getenv("HELLOPROFILE_HOME"), value.ID)),
				Basic:        <-basicBlock,
				ContactBlock: <-contactBlock,
				Socials:      <-socials,
				Contents:     <-contents,
			})
		}
		profiles <- profileList
	}()
}

func (env *Env) getProfile(profileID uuid.UUID, profile chan models.Profile, fields log.Fields) {

	socials := make(chan []models.Socials)
	contents := make(chan []models.Content)
	basicBlock := make(chan models.Basic)
	contactBlock := make(chan models.ContactBlock)

	go func() {
		log.WithFields(fields).Info(`Fetching profiles for user...`)
		var profileResult models.Profile
		dbProfile, err := env.HelloProfileDb.GetProfile(context.Background(), profileID)
		if err != nil {
			log.WithFields(fields).WithError(err).Error(`Error occured while fetching profiles for user...`)
			profile <- profileResult
		}

		go env.getBasicBlock(dbProfile.BasicBlockID.UUID, basicBlock, fields)
		go env.getContactBlock(dbProfile.ContactBlockID.UUID, contactBlock, fields)
		go env.getSocails(dbProfile.ID, socials, fields)
		go env.getContents(dbProfile.ID, contents, fields)
		profileResult = models.Profile{
			ID:           dbProfile.ID,
			Status:       dbProfile.Status,
			ProfileName:  dbProfile.ProfileName,
			PageColor:    dbProfile.PageColor,
			Font:         dbProfile.Font,
			IsDefault:    dbProfile.IsDefault,
			Url:          env.GetValue(dbProfile.Url.String, fmt.Sprintf("%s/%s", os.Getenv("HELLOPROFILE_HOME"), dbProfile.ID)),
			Basic:        <-basicBlock,
			ContactBlock: <-contactBlock,
			Socials:      <-socials,
			Contents:     <-contents,
		}

		profile <- profileResult
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
			SocialsID:   value.SocialsID,
			ProfileID:   value.ProfileID,
			ID:          value.ID,
		}
		socialsList = append(socialsList, social)
	}
	socials <- socialsList
}

func (env *Env) getContents(profileID uuid.UUID, contents chan []models.Content, fields log.Fields) {
	log.WithFields(fields).Info(`Getting the socials for the user profile`)
	var contentList []models.Content

	dbContents, err := env.HelloProfileDb.GetProfileContents(context.Background(), profileID)
	if err != nil {
		log.WithFields(fields).WithError(err).Error(`Error occured fetching contents for user profile`)
		contents <- contentList

		return
	}
	dbCallToActions, err := env.HelloProfileDb.GetCallToActions(context.Background())
	if err != nil {
		log.WithFields(fields).WithError(err).Error(`Error occured fetching call to action for profile contents`)
		contents <- contentList

		return
	}
	dbContentTypes, err := env.HelloProfileDb.GetAllContentTypes(context.Background())
	if err != nil {
		log.WithFields(fields).WithError(err).Error(`Error occured fetching call to action for profile contents`)
		contents <- contentList

		return
	}
	//Sort contents
	//Get links
	go env.enrichContents(dbContents, dbCallToActions, contents, dbContentTypes)
	// for _, value := range dbContentTypes {
	// 	if value.Type == "links" {
	// 		go env.sortContents(dbContents, dbCallToActions, linkContents, value.Type)
	// 	} else if value.Type == "articles" {
	// 		go env.sortContents(dbContents, dbCallToActions, articleContents, value.ID)
	// 	} else if value.Type == "embedded videos" {
	// 		go env.sortContents(dbContents, dbCallToActions, videoContents, value.ID)
	// 	} else if value.Type == "embedded audios" {
	// 		go env.sortContents(dbContents, dbCallToActions, audioContents, value.ID)
	// 	} else if value.Type == "forms" {
	// 		go env.sortContents(dbContents, dbCallToActions, formContents, value.ID)
	// 	} else if value.Type == "meetings" {
	// 		go env.sortContents(dbContents, dbCallToActions, meetingContents, value.ID)
	// 	} else if value.Type == "events" {
	// 		go env.sortContents(dbContents, dbCallToActions, eventsContents, value.ID)
	// 	}
	// }

}

func (env *Env) enrichContents(contentResult []helloprofiledb.ProfileContent, callToActions []helloprofiledb.CallToAction, contentChannel chan []models.Content, contentTypes []helloprofiledb.Content) {
	var contentList []models.Content
	var contentType helloprofiledb.Content
	for _, value := range contentResult {

		callToAction := new(helloprofiledb.CallToAction)
		for _, action := range callToActions {
			if value.CallToActionID == action.ID {
				callToAction.DisplayName = action.DisplayName
				callToAction.ID = action.ID
				callToAction.Type = action.Type
			}
		}
		for _, contentValue := range contentTypes {
			if contentValue.ID == value.ContentID {
				contentType = contentValue
			}
		}
		content := models.Content{
			ID:          value.ID,
			Title:       value.DisplayTitle,
			Description: value.Description,
			Url:         value.Url,
			Order:       value.Order,
			Type:        contentType.Type,
			CallToAction: models.CallToAction{
				ID:          callToAction.ID,
				Type:        callToAction.Type,
				DisplayName: callToAction.DisplayName,
			},
			CallToActionID: callToAction.ID,
			ContentID:      contentType.ID,
		}
		contentList = append(contentList, content)

	}
	contentChannel <- contentList
}

func (env *Env) getBasicBlock(id uuid.UUID, basicBlock chan models.Basic, fields log.Fields) {
	log.WithFields(fields).Info(`Getting the basic block for the user profile `, id)
	basic := new(models.Basic)
	dbBasic, err := env.HelloProfileDb.GetBasicBlock(context.Background(), id)
	if err != nil {
		log.WithFields(fields).WithError(err).Error(`Error occured fetching basic block for user profile `, id)
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
	log.WithFields(fields).Info(`Getting the contact block for the user profile `, id)
	contact := new(models.ContactBlock)
	dbContact, err := env.HelloProfileDb.GetContactBlock(context.Background(), id)
	if err != nil {
		log.WithFields(fields).WithError(err).Error(`Error occured fetching contact block for user profile `, id)
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
