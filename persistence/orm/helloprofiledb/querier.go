// Code generated by sqlc. DO NOT EDIT.

package helloprofiledb

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Querier interface {
	AddBasicBlock(ctx context.Context, arg AddBasicBlockParams) (BasicBlock, error)
	AddContactBlock(ctx context.Context, arg AddContactBlockParams) (ContactBlock, error)
	AddContactCategory(ctx context.Context, name string) (ContactCategory, error)
	AddContacts(ctx context.Context, arg AddContactsParams) (Contact, error)
	AddProfile(ctx context.Context, arg AddProfileParams) (Profile, error)
	AddProfileContent(ctx context.Context, arg AddProfileContentParams) (ProfileContent, error)
	AddProfileSocial(ctx context.Context, arg AddProfileSocialParams) (ProfileSocial, error)
	AddSocial(ctx context.Context, arg AddSocialParams) (Social, error)
	AddUserRole(ctx context.Context, arg AddUserRoleParams) (UserRole, error)
	CreateEmailVerification(ctx context.Context, arg CreateEmailVerificationParams) error
	CreateOtp(ctx context.Context, arg CreateOtpParams) error
	CreateRefreshToken(ctx context.Context, arg CreateRefreshTokenParams) (RefreshToken, error)
	CreateRole(ctx context.Context, arg CreateRoleParams) (Role, error)
	CreateSavedProfile(ctx context.Context, arg CreateSavedProfileParams) (SavedProfile, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	CreateUserLogin(ctx context.Context, arg CreateUserLoginParams) (UserLogin, error)
	DeleteBasicBlock(ctx context.Context, id uuid.UUID) error
	DeleteContact(ctx context.Context, arg DeleteContactParams) error
	DeleteContactBlock(ctx context.Context, id uuid.UUID) error
	DeleteContactCategory(ctx context.Context, name string) error
	DeleteEmailVerification(ctx context.Context, otp string) error
	DeleteOtp(ctx context.Context, arg DeleteOtpParams) error
	DeleteProfile(ctx context.Context, id uuid.UUID) error
	DeleteProfileContent(ctx context.Context, id uuid.UUID) error
	DeleteProfileSocial(ctx context.Context, id uuid.UUID) error
	DeleteRefreshToken(ctx context.Context, token string) error
	DeleteRoles(ctx context.Context, name string) error
	DeleteSavedProfile(ctx context.Context, id uuid.UUID) error
	DeleteSocial(ctx context.Context, id uuid.UUID) error
	DeleteUser(ctx context.Context, email string) error
	DeleteUserLogin(ctx context.Context, userID uuid.UUID) error
	GetAllContactCategories(ctx context.Context) ([]ContactCategory, error)
	GetAllContacts(ctx context.Context) ([]Contact, error)
	GetAllContentTypes(ctx context.Context) ([]Content, error)
	GetAllOtp(ctx context.Context) ([]Otp, error)
	GetAllProfiles(ctx context.Context) ([]Profile, error)
	GetBasicBlock(ctx context.Context, id uuid.UUID) (BasicBlock, error)
	GetCallToAction(ctx context.Context, id uuid.UUID) (CallToAction, error)
	GetCallToActions(ctx context.Context) ([]CallToAction, error)
	GetContactBlock(ctx context.Context, id uuid.UUID) (ContactBlock, error)
	GetContactCategory(ctx context.Context, name string) (ContactCategory, error)
	GetContacts(ctx context.Context, userID uuid.UUID) ([]Contact, error)
	GetEmailVerification(ctx context.Context, otp string) (EmailVerification, error)
	GetEmailVerifications(ctx context.Context) ([]EmailVerification, error)
	GetOtp(ctx context.Context, arg GetOtpParams) (Otp, error)
	GetProfile(ctx context.Context, id uuid.UUID) (Profile, error)
	GetProfileContent(ctx context.Context, id uuid.UUID) (ProfileContent, error)
	GetProfileContents(ctx context.Context, profileID uuid.UUID) ([]ProfileContent, error)
	GetProfileSocial(ctx context.Context, id uuid.UUID) (ProfileSocial, error)
	GetProfileSocials(ctx context.Context, profileID uuid.UUID) ([]GetProfileSocialsRow, error)
	GetProfiles(ctx context.Context, userID uuid.UUID) ([]Profile, error)
	GetRefreshToken(ctx context.Context, token string) (RefreshToken, error)
	GetRefreshTokens(ctx context.Context) ([]RefreshToken, error)
	GetRole(ctx context.Context, name string) (Role, error)
	GetRoles(ctx context.Context) ([]Role, error)
	GetSavedProfile(ctx context.Context, id uuid.UUID) (SavedProfile, error)
	GetSavedProfiles(ctx context.Context) ([]SavedProfile, error)
	GetSavedProfilesByEmail(ctx context.Context, arg GetSavedProfilesByEmailParams) ([]SavedProfile, error)
	GetSavedProfilesByProfileId(ctx context.Context, arg GetSavedProfilesByProfileIdParams) ([]SavedProfile, error)
	GetSocial(ctx context.Context, id uuid.UUID) (Social, error)
	GetSocials(ctx context.Context) ([]Social, error)
	GetUnResoledLogins(ctx context.Context, userID uuid.UUID) ([]UserLogin, error)
	GetUser(ctx context.Context, username sql.NullString) (UserDetail, error)
	GetUserLogin(ctx context.Context, userID uuid.UUID) ([]UserLogin, error)
	GetUserLogins(ctx context.Context) ([]UserLogin, error)
	GetUserRoles(ctx context.Context, username sql.NullString) ([]string, error)
	GetUsers(ctx context.Context) ([]UserDetail, error)
	IsProfileExist(ctx context.Context, id uuid.UUID) (bool, error)
	UpdateBasicBlock(ctx context.Context, arg UpdateBasicBlockParams) error
	UpdateContactBlock(ctx context.Context, arg UpdateContactBlockParams) error
	UpdateContactCategory(ctx context.Context, arg UpdateContactCategoryParams) error
	UpdateProfile(ctx context.Context, arg UpdateProfileParams) error
	UpdateProfileContent(ctx context.Context, arg UpdateProfileContentParams) error
	UpdateProfileSocial(ctx context.Context, arg UpdateProfileSocialParams) error
	UpdateProfileWithBasicBlockId(ctx context.Context, arg UpdateProfileWithBasicBlockIdParams) error
	UpdateProfileWithContactBlockId(ctx context.Context, arg UpdateProfileWithContactBlockIdParams) error
	UpdateRefreshToken(ctx context.Context, arg UpdateRefreshTokenParams) (RefreshToken, error)
	UpdateResolvedLogin(ctx context.Context, userID uuid.UUID) error
	UpdateRole(ctx context.Context, arg UpdateRoleParams) (Role, error)
	UpdateSavedProfile(ctx context.Context, arg UpdateSavedProfileParams) (SavedProfile, error)
	UpdateSocial(ctx context.Context, arg UpdateSocialParams) error
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
	UpdateUserRole(ctx context.Context, arg UpdateUserRoleParams) (UserRole, error)
}

var _ Querier = (*Queries)(nil)
