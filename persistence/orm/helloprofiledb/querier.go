// Code generated by sqlc. DO NOT EDIT.

package helloprofiledb

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Querier interface {
	AddAddress(ctx context.Context, arg AddAddressParams) (Address, error)
	AddAddressType(ctx context.Context, name string) (AddressType, error)
	AddContactCategory(ctx context.Context, name string) (ContactCategory, error)
	AddContacts(ctx context.Context, arg AddContactsParams) (Contact, error)
	AddProfile(ctx context.Context, arg AddProfileParams) (Profile, error)
	AddUserLanguage(ctx context.Context, arg AddUserLanguageParams) (UserLanguage, error)
	AddUserProvider(ctx context.Context, arg AddUserProviderParams) (UserProvider, error)
	AddUserRole(ctx context.Context, arg AddUserRoleParams) (UserRole, error)
	AddUserTimezone(ctx context.Context, arg AddUserTimezoneParams) (UserTimezone, error)
	CreateCountry(ctx context.Context, arg CreateCountryParams) (Country, error)
	CreateEmailVerification(ctx context.Context, arg CreateEmailVerificationParams) error
	CreateIdentityProvider(ctx context.Context, arg CreateIdentityProviderParams) (IdentityProvider, error)
	CreateLanguage(ctx context.Context, name string) (Language, error)
	CreateLanguageProficiency(ctx context.Context, proficiency sql.NullString) (LanguageProficiency, error)
	CreateOtp(ctx context.Context, arg CreateOtpParams) error
	CreateRefreshToken(ctx context.Context, arg CreateRefreshTokenParams) (RefreshToken, error)
	CreateRole(ctx context.Context, arg CreateRoleParams) (Role, error)
	CreateState(ctx context.Context, arg CreateStateParams) error
	CreateTimezone(ctx context.Context, arg CreateTimezoneParams) (Timezone, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	CreateUserLogin(ctx context.Context, arg CreateUserLoginParams) (UserLogin, error)
	DeleteAddress(ctx context.Context, id uuid.UUID) error
	DeleteAddressType(ctx context.Context, id uuid.UUID) error
	DeleteContact(ctx context.Context, arg DeleteContactParams) error
	DeleteContactCategory(ctx context.Context, id uuid.UUID) error
	DeleteCountry(ctx context.Context, name string) error
	DeleteEmailVerification(ctx context.Context, otp string) error
	DeleteIdentityProvider(ctx context.Context, name string) error
	DeleteLanguage(ctx context.Context, name string) error
	DeleteLanguageProficiency(ctx context.Context, proficiency sql.NullString) error
	DeleteOtp(ctx context.Context, arg DeleteOtpParams) error
	DeleteProfile(ctx context.Context, id uuid.UUID) error
	DeleteProviders(ctx context.Context, arg DeleteProvidersParams) error
	DeleteRefreshToken(ctx context.Context, token string) error
	DeleteRoles(ctx context.Context, name string) error
	DeleteState(ctx context.Context, name string) error
	DeleteTimezone(ctx context.Context, name string) error
	DeleteUser(ctx context.Context, email string) error
	DeleteUserLanguage(ctx context.Context, arg DeleteUserLanguageParams) error
	DeleteUserLogin(ctx context.Context, userID uuid.UUID) error
	GetAddress(ctx context.Context, id uuid.UUID) (Address, error)
	GetAddressType(ctx context.Context, id uuid.UUID) (AddressType, error)
	GetAllAddressTypes(ctx context.Context) ([]AddressType, error)
	GetAllAddresses(ctx context.Context) ([]Address, error)
	GetAllContactCategories(ctx context.Context) ([]ContactCategory, error)
	GetAllContacts(ctx context.Context) ([]Contact, error)
	GetAllOtp(ctx context.Context) ([]Otp, error)
	GetAllProfiles(ctx context.Context) ([]Profile, error)
	GetContactCategory(ctx context.Context, id uuid.UUID) (ContactCategory, error)
	GetContacts(ctx context.Context, userID uuid.UUID) ([]Contact, error)
	GetCountries(ctx context.Context) ([]Country, error)
	GetCountry(ctx context.Context, name string) (Country, error)
	GetEmailVerification(ctx context.Context, otp string) (EmailVerification, error)
	GetEmailVerifications(ctx context.Context) ([]EmailVerification, error)
	GetIdentityProvider(ctx context.Context, name string) (IdentityProvider, error)
	GetIdentityProviders(ctx context.Context) ([]IdentityProvider, error)
	GetLanguage(ctx context.Context, name string) (Language, error)
	GetLanguageProficiencies(ctx context.Context) ([]LanguageProficiency, error)
	GetLanguageProficiency(ctx context.Context, proficiency sql.NullString) (LanguageProficiency, error)
	GetLanguages(ctx context.Context) ([]Language, error)
	GetOtp(ctx context.Context, arg GetOtpParams) (Otp, error)
	GetProfile(ctx context.Context, id uuid.UUID) (Profile, error)
	GetProfiles(ctx context.Context, userID uuid.UUID) ([]Profile, error)
	GetRefreshToken(ctx context.Context, token string) (RefreshToken, error)
	GetRefreshTokens(ctx context.Context) ([]RefreshToken, error)
	GetRole(ctx context.Context, name string) (Role, error)
	GetRoles(ctx context.Context) ([]Role, error)
	GetState(ctx context.Context, name string) (State, error)
	GetStates(ctx context.Context) ([]State, error)
	GetStatesByCountry(ctx context.Context, name string) ([]State, error)
	GetTimezone(ctx context.Context, name string) (Timezone, error)
	GetTimezones(ctx context.Context) ([]Timezone, error)
	GetUnResoledLogins(ctx context.Context, userID uuid.UUID) ([]UserLogin, error)
	GetUser(ctx context.Context, username sql.NullString) (UserDetail, error)
	GetUserAddresses(ctx context.Context, userID uuid.NullUUID) ([]Address, error)
	GetUserLanguages(ctx context.Context, username sql.NullString) ([]GetUserLanguagesRow, error)
	GetUserLogin(ctx context.Context, userID uuid.UUID) ([]UserLogin, error)
	GetUserLogins(ctx context.Context) ([]UserLogin, error)
	GetUserProviders(ctx context.Context, username sql.NullString) ([]IdentityProvider, error)
	GetUserRoles(ctx context.Context, username sql.NullString) ([]string, error)
	GetUserTimezones(ctx context.Context, username sql.NullString) ([]Timezone, error)
	GetUsers(ctx context.Context) ([]UserDetail, error)
	UpdateAddress(ctx context.Context, arg UpdateAddressParams) error
	UpdateAddressType(ctx context.Context, arg UpdateAddressTypeParams) error
	UpdateContact(ctx context.Context, arg UpdateContactParams) error
	UpdateContactCategory(ctx context.Context, arg UpdateContactCategoryParams) error
	UpdateCountry(ctx context.Context, arg UpdateCountryParams) (Country, error)
	UpdateIdentityProvider(ctx context.Context, arg UpdateIdentityProviderParams) (IdentityProvider, error)
	UpdateLanguage(ctx context.Context, arg UpdateLanguageParams) (Language, error)
	UpdateLanguageProficiency(ctx context.Context, arg UpdateLanguageProficiencyParams) (LanguageProficiency, error)
	UpdateProfile(ctx context.Context, arg UpdateProfileParams) error
	UpdateRefreshToken(ctx context.Context, arg UpdateRefreshTokenParams) (RefreshToken, error)
	UpdateResolvedLogin(ctx context.Context, userID uuid.UUID) error
	UpdateRole(ctx context.Context, arg UpdateRoleParams) (Role, error)
	UpdateState(ctx context.Context, arg UpdateStateParams) error
	UpdateTimezone(ctx context.Context, arg UpdateTimezoneParams) (Timezone, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
	UpdateUserLanguage(ctx context.Context, arg UpdateUserLanguageParams) (UserLanguage, error)
	UpdateUserProvider(ctx context.Context, arg UpdateUserProviderParams) (UserProvider, error)
	UpdateUserRole(ctx context.Context, arg UpdateUserRoleParams) (UserRole, error)
	UpdateUserTimezone(ctx context.Context, arg UpdateUserTimezoneParams) (UserTimezone, error)
}

var _ Querier = (*Queries)(nil)
