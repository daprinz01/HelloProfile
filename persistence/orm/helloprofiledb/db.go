// Code generated by sqlc. DO NOT EDIT.

package helloprofiledb

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.addBasicBlockStmt, err = db.PrepareContext(ctx, addBasicBlock); err != nil {
		return nil, fmt.Errorf("error preparing query AddBasicBlock: %w", err)
	}
	if q.addContactBlockStmt, err = db.PrepareContext(ctx, addContactBlock); err != nil {
		return nil, fmt.Errorf("error preparing query AddContactBlock: %w", err)
	}
	if q.addContactCategoryStmt, err = db.PrepareContext(ctx, addContactCategory); err != nil {
		return nil, fmt.Errorf("error preparing query AddContactCategory: %w", err)
	}
	if q.addContactsStmt, err = db.PrepareContext(ctx, addContacts); err != nil {
		return nil, fmt.Errorf("error preparing query AddContacts: %w", err)
	}
	if q.addProfileStmt, err = db.PrepareContext(ctx, addProfile); err != nil {
		return nil, fmt.Errorf("error preparing query AddProfile: %w", err)
	}
	if q.addProfileContentStmt, err = db.PrepareContext(ctx, addProfileContent); err != nil {
		return nil, fmt.Errorf("error preparing query AddProfileContent: %w", err)
	}
	if q.addProfileSocialStmt, err = db.PrepareContext(ctx, addProfileSocial); err != nil {
		return nil, fmt.Errorf("error preparing query AddProfileSocial: %w", err)
	}
	if q.addSocialStmt, err = db.PrepareContext(ctx, addSocial); err != nil {
		return nil, fmt.Errorf("error preparing query AddSocial: %w", err)
	}
	if q.addUserRoleStmt, err = db.PrepareContext(ctx, addUserRole); err != nil {
		return nil, fmt.Errorf("error preparing query AddUserRole: %w", err)
	}
	if q.createEmailVerificationStmt, err = db.PrepareContext(ctx, createEmailVerification); err != nil {
		return nil, fmt.Errorf("error preparing query CreateEmailVerification: %w", err)
	}
	if q.createOtpStmt, err = db.PrepareContext(ctx, createOtp); err != nil {
		return nil, fmt.Errorf("error preparing query CreateOtp: %w", err)
	}
	if q.createRefreshTokenStmt, err = db.PrepareContext(ctx, createRefreshToken); err != nil {
		return nil, fmt.Errorf("error preparing query CreateRefreshToken: %w", err)
	}
	if q.createRoleStmt, err = db.PrepareContext(ctx, createRole); err != nil {
		return nil, fmt.Errorf("error preparing query CreateRole: %w", err)
	}
	if q.createSavedProfileStmt, err = db.PrepareContext(ctx, createSavedProfile); err != nil {
		return nil, fmt.Errorf("error preparing query CreateSavedProfile: %w", err)
	}
	if q.createUserStmt, err = db.PrepareContext(ctx, createUser); err != nil {
		return nil, fmt.Errorf("error preparing query CreateUser: %w", err)
	}
	if q.createUserLoginStmt, err = db.PrepareContext(ctx, createUserLogin); err != nil {
		return nil, fmt.Errorf("error preparing query CreateUserLogin: %w", err)
	}
	if q.deleteBasicBlockStmt, err = db.PrepareContext(ctx, deleteBasicBlock); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteBasicBlock: %w", err)
	}
	if q.deleteContactStmt, err = db.PrepareContext(ctx, deleteContact); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteContact: %w", err)
	}
	if q.deleteContactBlockStmt, err = db.PrepareContext(ctx, deleteContactBlock); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteContactBlock: %w", err)
	}
	if q.deleteContactCategoryStmt, err = db.PrepareContext(ctx, deleteContactCategory); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteContactCategory: %w", err)
	}
	if q.deleteEmailVerificationStmt, err = db.PrepareContext(ctx, deleteEmailVerification); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteEmailVerification: %w", err)
	}
	if q.deleteOtpStmt, err = db.PrepareContext(ctx, deleteOtp); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteOtp: %w", err)
	}
	if q.deleteProfileStmt, err = db.PrepareContext(ctx, deleteProfile); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteProfile: %w", err)
	}
	if q.deleteProfileContentStmt, err = db.PrepareContext(ctx, deleteProfileContent); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteProfileContent: %w", err)
	}
	if q.deleteProfileSocialStmt, err = db.PrepareContext(ctx, deleteProfileSocial); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteProfileSocial: %w", err)
	}
	if q.deleteRefreshTokenStmt, err = db.PrepareContext(ctx, deleteRefreshToken); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteRefreshToken: %w", err)
	}
	if q.deleteRolesStmt, err = db.PrepareContext(ctx, deleteRoles); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteRoles: %w", err)
	}
	if q.deleteSavedProfileStmt, err = db.PrepareContext(ctx, deleteSavedProfile); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteSavedProfile: %w", err)
	}
	if q.deleteSocialStmt, err = db.PrepareContext(ctx, deleteSocial); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteSocial: %w", err)
	}
	if q.deleteUserStmt, err = db.PrepareContext(ctx, deleteUser); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteUser: %w", err)
	}
	if q.deleteUserLoginStmt, err = db.PrepareContext(ctx, deleteUserLogin); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteUserLogin: %w", err)
	}
	if q.getAllContactCategoriesStmt, err = db.PrepareContext(ctx, getAllContactCategories); err != nil {
		return nil, fmt.Errorf("error preparing query GetAllContactCategories: %w", err)
	}
	if q.getAllContactsStmt, err = db.PrepareContext(ctx, getAllContacts); err != nil {
		return nil, fmt.Errorf("error preparing query GetAllContacts: %w", err)
	}
	if q.getAllContentTypesStmt, err = db.PrepareContext(ctx, getAllContentTypes); err != nil {
		return nil, fmt.Errorf("error preparing query GetAllContentTypes: %w", err)
	}
	if q.getAllOtpStmt, err = db.PrepareContext(ctx, getAllOtp); err != nil {
		return nil, fmt.Errorf("error preparing query GetAllOtp: %w", err)
	}
	if q.getAllProfilesStmt, err = db.PrepareContext(ctx, getAllProfiles); err != nil {
		return nil, fmt.Errorf("error preparing query GetAllProfiles: %w", err)
	}
	if q.getBasicBlockStmt, err = db.PrepareContext(ctx, getBasicBlock); err != nil {
		return nil, fmt.Errorf("error preparing query GetBasicBlock: %w", err)
	}
	if q.getCallToActionStmt, err = db.PrepareContext(ctx, getCallToAction); err != nil {
		return nil, fmt.Errorf("error preparing query GetCallToAction: %w", err)
	}
	if q.getCallToActionsStmt, err = db.PrepareContext(ctx, getCallToActions); err != nil {
		return nil, fmt.Errorf("error preparing query GetCallToActions: %w", err)
	}
	if q.getContactBlockStmt, err = db.PrepareContext(ctx, getContactBlock); err != nil {
		return nil, fmt.Errorf("error preparing query GetContactBlock: %w", err)
	}
	if q.getContactCategoryStmt, err = db.PrepareContext(ctx, getContactCategory); err != nil {
		return nil, fmt.Errorf("error preparing query GetContactCategory: %w", err)
	}
	if q.getContactsStmt, err = db.PrepareContext(ctx, getContacts); err != nil {
		return nil, fmt.Errorf("error preparing query GetContacts: %w", err)
	}
	if q.getEmailVerificationStmt, err = db.PrepareContext(ctx, getEmailVerification); err != nil {
		return nil, fmt.Errorf("error preparing query GetEmailVerification: %w", err)
	}
	if q.getEmailVerificationsStmt, err = db.PrepareContext(ctx, getEmailVerifications); err != nil {
		return nil, fmt.Errorf("error preparing query GetEmailVerifications: %w", err)
	}
	if q.getOtpStmt, err = db.PrepareContext(ctx, getOtp); err != nil {
		return nil, fmt.Errorf("error preparing query GetOtp: %w", err)
	}
	if q.getProfileStmt, err = db.PrepareContext(ctx, getProfile); err != nil {
		return nil, fmt.Errorf("error preparing query GetProfile: %w", err)
	}
	if q.getProfileContentStmt, err = db.PrepareContext(ctx, getProfileContent); err != nil {
		return nil, fmt.Errorf("error preparing query GetProfileContent: %w", err)
	}
	if q.getProfileContentsStmt, err = db.PrepareContext(ctx, getProfileContents); err != nil {
		return nil, fmt.Errorf("error preparing query GetProfileContents: %w", err)
	}
	if q.getProfileSocialStmt, err = db.PrepareContext(ctx, getProfileSocial); err != nil {
		return nil, fmt.Errorf("error preparing query GetProfileSocial: %w", err)
	}
	if q.getProfileSocialsStmt, err = db.PrepareContext(ctx, getProfileSocials); err != nil {
		return nil, fmt.Errorf("error preparing query GetProfileSocials: %w", err)
	}
	if q.getProfilesStmt, err = db.PrepareContext(ctx, getProfiles); err != nil {
		return nil, fmt.Errorf("error preparing query GetProfiles: %w", err)
	}
	if q.getRefreshTokenStmt, err = db.PrepareContext(ctx, getRefreshToken); err != nil {
		return nil, fmt.Errorf("error preparing query GetRefreshToken: %w", err)
	}
	if q.getRefreshTokensStmt, err = db.PrepareContext(ctx, getRefreshTokens); err != nil {
		return nil, fmt.Errorf("error preparing query GetRefreshTokens: %w", err)
	}
	if q.getRoleStmt, err = db.PrepareContext(ctx, getRole); err != nil {
		return nil, fmt.Errorf("error preparing query GetRole: %w", err)
	}
	if q.getRolesStmt, err = db.PrepareContext(ctx, getRoles); err != nil {
		return nil, fmt.Errorf("error preparing query GetRoles: %w", err)
	}
	if q.getSavedProfileStmt, err = db.PrepareContext(ctx, getSavedProfile); err != nil {
		return nil, fmt.Errorf("error preparing query GetSavedProfile: %w", err)
	}
	if q.getSavedProfilesStmt, err = db.PrepareContext(ctx, getSavedProfiles); err != nil {
		return nil, fmt.Errorf("error preparing query GetSavedProfiles: %w", err)
	}
	if q.getSavedProfilesByEmailStmt, err = db.PrepareContext(ctx, getSavedProfilesByEmail); err != nil {
		return nil, fmt.Errorf("error preparing query GetSavedProfilesByEmail: %w", err)
	}
	if q.getSavedProfilesByProfileIdStmt, err = db.PrepareContext(ctx, getSavedProfilesByProfileId); err != nil {
		return nil, fmt.Errorf("error preparing query GetSavedProfilesByProfileId: %w", err)
	}
	if q.getSocialStmt, err = db.PrepareContext(ctx, getSocial); err != nil {
		return nil, fmt.Errorf("error preparing query GetSocial: %w", err)
	}
	if q.getSocialsStmt, err = db.PrepareContext(ctx, getSocials); err != nil {
		return nil, fmt.Errorf("error preparing query GetSocials: %w", err)
	}
	if q.getUnResoledLoginsStmt, err = db.PrepareContext(ctx, getUnResoledLogins); err != nil {
		return nil, fmt.Errorf("error preparing query GetUnResoledLogins: %w", err)
	}
	if q.getUserStmt, err = db.PrepareContext(ctx, getUser); err != nil {
		return nil, fmt.Errorf("error preparing query GetUser: %w", err)
	}
	if q.getUserLoginStmt, err = db.PrepareContext(ctx, getUserLogin); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserLogin: %w", err)
	}
	if q.getUserLoginsStmt, err = db.PrepareContext(ctx, getUserLogins); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserLogins: %w", err)
	}
	if q.getUserRolesStmt, err = db.PrepareContext(ctx, getUserRoles); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserRoles: %w", err)
	}
	if q.getUsersStmt, err = db.PrepareContext(ctx, getUsers); err != nil {
		return nil, fmt.Errorf("error preparing query GetUsers: %w", err)
	}
	if q.isProfileExistStmt, err = db.PrepareContext(ctx, isProfileExist); err != nil {
		return nil, fmt.Errorf("error preparing query IsProfileExist: %w", err)
	}
	if q.isUrlExistsStmt, err = db.PrepareContext(ctx, isUrlExists); err != nil {
		return nil, fmt.Errorf("error preparing query IsUrlExists: %w", err)
	}
	if q.updateBasicBlockStmt, err = db.PrepareContext(ctx, updateBasicBlock); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateBasicBlock: %w", err)
	}
	if q.updateContactBlockStmt, err = db.PrepareContext(ctx, updateContactBlock); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateContactBlock: %w", err)
	}
	if q.updateContactCategoryStmt, err = db.PrepareContext(ctx, updateContactCategory); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateContactCategory: %w", err)
	}
	if q.updateProfileStmt, err = db.PrepareContext(ctx, updateProfile); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateProfile: %w", err)
	}
	if q.updateProfileContentStmt, err = db.PrepareContext(ctx, updateProfileContent); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateProfileContent: %w", err)
	}
	if q.updateProfileSocialStmt, err = db.PrepareContext(ctx, updateProfileSocial); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateProfileSocial: %w", err)
	}
	if q.updateProfileUrlStmt, err = db.PrepareContext(ctx, updateProfileUrl); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateProfileUrl: %w", err)
	}
	if q.updateProfileWithBasicBlockIdStmt, err = db.PrepareContext(ctx, updateProfileWithBasicBlockId); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateProfileWithBasicBlockId: %w", err)
	}
	if q.updateProfileWithContactBlockIdStmt, err = db.PrepareContext(ctx, updateProfileWithContactBlockId); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateProfileWithContactBlockId: %w", err)
	}
	if q.updateRefreshTokenStmt, err = db.PrepareContext(ctx, updateRefreshToken); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateRefreshToken: %w", err)
	}
	if q.updateResolvedLoginStmt, err = db.PrepareContext(ctx, updateResolvedLogin); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateResolvedLogin: %w", err)
	}
	if q.updateRoleStmt, err = db.PrepareContext(ctx, updateRole); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateRole: %w", err)
	}
	if q.updateSavedProfileStmt, err = db.PrepareContext(ctx, updateSavedProfile); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateSavedProfile: %w", err)
	}
	if q.updateSocialStmt, err = db.PrepareContext(ctx, updateSocial); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateSocial: %w", err)
	}
	if q.updateUserStmt, err = db.PrepareContext(ctx, updateUser); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateUser: %w", err)
	}
	if q.updateUserRoleStmt, err = db.PrepareContext(ctx, updateUserRole); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateUserRole: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.addBasicBlockStmt != nil {
		if cerr := q.addBasicBlockStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addBasicBlockStmt: %w", cerr)
		}
	}
	if q.addContactBlockStmt != nil {
		if cerr := q.addContactBlockStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addContactBlockStmt: %w", cerr)
		}
	}
	if q.addContactCategoryStmt != nil {
		if cerr := q.addContactCategoryStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addContactCategoryStmt: %w", cerr)
		}
	}
	if q.addContactsStmt != nil {
		if cerr := q.addContactsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addContactsStmt: %w", cerr)
		}
	}
	if q.addProfileStmt != nil {
		if cerr := q.addProfileStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addProfileStmt: %w", cerr)
		}
	}
	if q.addProfileContentStmt != nil {
		if cerr := q.addProfileContentStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addProfileContentStmt: %w", cerr)
		}
	}
	if q.addProfileSocialStmt != nil {
		if cerr := q.addProfileSocialStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addProfileSocialStmt: %w", cerr)
		}
	}
	if q.addSocialStmt != nil {
		if cerr := q.addSocialStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addSocialStmt: %w", cerr)
		}
	}
	if q.addUserRoleStmt != nil {
		if cerr := q.addUserRoleStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addUserRoleStmt: %w", cerr)
		}
	}
	if q.createEmailVerificationStmt != nil {
		if cerr := q.createEmailVerificationStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createEmailVerificationStmt: %w", cerr)
		}
	}
	if q.createOtpStmt != nil {
		if cerr := q.createOtpStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createOtpStmt: %w", cerr)
		}
	}
	if q.createRefreshTokenStmt != nil {
		if cerr := q.createRefreshTokenStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createRefreshTokenStmt: %w", cerr)
		}
	}
	if q.createRoleStmt != nil {
		if cerr := q.createRoleStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createRoleStmt: %w", cerr)
		}
	}
	if q.createSavedProfileStmt != nil {
		if cerr := q.createSavedProfileStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createSavedProfileStmt: %w", cerr)
		}
	}
	if q.createUserStmt != nil {
		if cerr := q.createUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createUserStmt: %w", cerr)
		}
	}
	if q.createUserLoginStmt != nil {
		if cerr := q.createUserLoginStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createUserLoginStmt: %w", cerr)
		}
	}
	if q.deleteBasicBlockStmt != nil {
		if cerr := q.deleteBasicBlockStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteBasicBlockStmt: %w", cerr)
		}
	}
	if q.deleteContactStmt != nil {
		if cerr := q.deleteContactStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteContactStmt: %w", cerr)
		}
	}
	if q.deleteContactBlockStmt != nil {
		if cerr := q.deleteContactBlockStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteContactBlockStmt: %w", cerr)
		}
	}
	if q.deleteContactCategoryStmt != nil {
		if cerr := q.deleteContactCategoryStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteContactCategoryStmt: %w", cerr)
		}
	}
	if q.deleteEmailVerificationStmt != nil {
		if cerr := q.deleteEmailVerificationStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteEmailVerificationStmt: %w", cerr)
		}
	}
	if q.deleteOtpStmt != nil {
		if cerr := q.deleteOtpStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteOtpStmt: %w", cerr)
		}
	}
	if q.deleteProfileStmt != nil {
		if cerr := q.deleteProfileStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteProfileStmt: %w", cerr)
		}
	}
	if q.deleteProfileContentStmt != nil {
		if cerr := q.deleteProfileContentStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteProfileContentStmt: %w", cerr)
		}
	}
	if q.deleteProfileSocialStmt != nil {
		if cerr := q.deleteProfileSocialStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteProfileSocialStmt: %w", cerr)
		}
	}
	if q.deleteRefreshTokenStmt != nil {
		if cerr := q.deleteRefreshTokenStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteRefreshTokenStmt: %w", cerr)
		}
	}
	if q.deleteRolesStmt != nil {
		if cerr := q.deleteRolesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteRolesStmt: %w", cerr)
		}
	}
	if q.deleteSavedProfileStmt != nil {
		if cerr := q.deleteSavedProfileStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteSavedProfileStmt: %w", cerr)
		}
	}
	if q.deleteSocialStmt != nil {
		if cerr := q.deleteSocialStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteSocialStmt: %w", cerr)
		}
	}
	if q.deleteUserStmt != nil {
		if cerr := q.deleteUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteUserStmt: %w", cerr)
		}
	}
	if q.deleteUserLoginStmt != nil {
		if cerr := q.deleteUserLoginStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteUserLoginStmt: %w", cerr)
		}
	}
	if q.getAllContactCategoriesStmt != nil {
		if cerr := q.getAllContactCategoriesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAllContactCategoriesStmt: %w", cerr)
		}
	}
	if q.getAllContactsStmt != nil {
		if cerr := q.getAllContactsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAllContactsStmt: %w", cerr)
		}
	}
	if q.getAllContentTypesStmt != nil {
		if cerr := q.getAllContentTypesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAllContentTypesStmt: %w", cerr)
		}
	}
	if q.getAllOtpStmt != nil {
		if cerr := q.getAllOtpStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAllOtpStmt: %w", cerr)
		}
	}
	if q.getAllProfilesStmt != nil {
		if cerr := q.getAllProfilesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAllProfilesStmt: %w", cerr)
		}
	}
	if q.getBasicBlockStmt != nil {
		if cerr := q.getBasicBlockStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getBasicBlockStmt: %w", cerr)
		}
	}
	if q.getCallToActionStmt != nil {
		if cerr := q.getCallToActionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getCallToActionStmt: %w", cerr)
		}
	}
	if q.getCallToActionsStmt != nil {
		if cerr := q.getCallToActionsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getCallToActionsStmt: %w", cerr)
		}
	}
	if q.getContactBlockStmt != nil {
		if cerr := q.getContactBlockStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getContactBlockStmt: %w", cerr)
		}
	}
	if q.getContactCategoryStmt != nil {
		if cerr := q.getContactCategoryStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getContactCategoryStmt: %w", cerr)
		}
	}
	if q.getContactsStmt != nil {
		if cerr := q.getContactsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getContactsStmt: %w", cerr)
		}
	}
	if q.getEmailVerificationStmt != nil {
		if cerr := q.getEmailVerificationStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getEmailVerificationStmt: %w", cerr)
		}
	}
	if q.getEmailVerificationsStmt != nil {
		if cerr := q.getEmailVerificationsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getEmailVerificationsStmt: %w", cerr)
		}
	}
	if q.getOtpStmt != nil {
		if cerr := q.getOtpStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getOtpStmt: %w", cerr)
		}
	}
	if q.getProfileStmt != nil {
		if cerr := q.getProfileStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getProfileStmt: %w", cerr)
		}
	}
	if q.getProfileContentStmt != nil {
		if cerr := q.getProfileContentStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getProfileContentStmt: %w", cerr)
		}
	}
	if q.getProfileContentsStmt != nil {
		if cerr := q.getProfileContentsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getProfileContentsStmt: %w", cerr)
		}
	}
	if q.getProfileSocialStmt != nil {
		if cerr := q.getProfileSocialStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getProfileSocialStmt: %w", cerr)
		}
	}
	if q.getProfileSocialsStmt != nil {
		if cerr := q.getProfileSocialsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getProfileSocialsStmt: %w", cerr)
		}
	}
	if q.getProfilesStmt != nil {
		if cerr := q.getProfilesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getProfilesStmt: %w", cerr)
		}
	}
	if q.getRefreshTokenStmt != nil {
		if cerr := q.getRefreshTokenStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getRefreshTokenStmt: %w", cerr)
		}
	}
	if q.getRefreshTokensStmt != nil {
		if cerr := q.getRefreshTokensStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getRefreshTokensStmt: %w", cerr)
		}
	}
	if q.getRoleStmt != nil {
		if cerr := q.getRoleStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getRoleStmt: %w", cerr)
		}
	}
	if q.getRolesStmt != nil {
		if cerr := q.getRolesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getRolesStmt: %w", cerr)
		}
	}
	if q.getSavedProfileStmt != nil {
		if cerr := q.getSavedProfileStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getSavedProfileStmt: %w", cerr)
		}
	}
	if q.getSavedProfilesStmt != nil {
		if cerr := q.getSavedProfilesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getSavedProfilesStmt: %w", cerr)
		}
	}
	if q.getSavedProfilesByEmailStmt != nil {
		if cerr := q.getSavedProfilesByEmailStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getSavedProfilesByEmailStmt: %w", cerr)
		}
	}
	if q.getSavedProfilesByProfileIdStmt != nil {
		if cerr := q.getSavedProfilesByProfileIdStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getSavedProfilesByProfileIdStmt: %w", cerr)
		}
	}
	if q.getSocialStmt != nil {
		if cerr := q.getSocialStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getSocialStmt: %w", cerr)
		}
	}
	if q.getSocialsStmt != nil {
		if cerr := q.getSocialsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getSocialsStmt: %w", cerr)
		}
	}
	if q.getUnResoledLoginsStmt != nil {
		if cerr := q.getUnResoledLoginsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUnResoledLoginsStmt: %w", cerr)
		}
	}
	if q.getUserStmt != nil {
		if cerr := q.getUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserStmt: %w", cerr)
		}
	}
	if q.getUserLoginStmt != nil {
		if cerr := q.getUserLoginStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserLoginStmt: %w", cerr)
		}
	}
	if q.getUserLoginsStmt != nil {
		if cerr := q.getUserLoginsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserLoginsStmt: %w", cerr)
		}
	}
	if q.getUserRolesStmt != nil {
		if cerr := q.getUserRolesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserRolesStmt: %w", cerr)
		}
	}
	if q.getUsersStmt != nil {
		if cerr := q.getUsersStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUsersStmt: %w", cerr)
		}
	}
	if q.isProfileExistStmt != nil {
		if cerr := q.isProfileExistStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing isProfileExistStmt: %w", cerr)
		}
	}
	if q.isUrlExistsStmt != nil {
		if cerr := q.isUrlExistsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing isUrlExistsStmt: %w", cerr)
		}
	}
	if q.updateBasicBlockStmt != nil {
		if cerr := q.updateBasicBlockStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateBasicBlockStmt: %w", cerr)
		}
	}
	if q.updateContactBlockStmt != nil {
		if cerr := q.updateContactBlockStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateContactBlockStmt: %w", cerr)
		}
	}
	if q.updateContactCategoryStmt != nil {
		if cerr := q.updateContactCategoryStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateContactCategoryStmt: %w", cerr)
		}
	}
	if q.updateProfileStmt != nil {
		if cerr := q.updateProfileStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateProfileStmt: %w", cerr)
		}
	}
	if q.updateProfileContentStmt != nil {
		if cerr := q.updateProfileContentStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateProfileContentStmt: %w", cerr)
		}
	}
	if q.updateProfileSocialStmt != nil {
		if cerr := q.updateProfileSocialStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateProfileSocialStmt: %w", cerr)
		}
	}
	if q.updateProfileUrlStmt != nil {
		if cerr := q.updateProfileUrlStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateProfileUrlStmt: %w", cerr)
		}
	}
	if q.updateProfileWithBasicBlockIdStmt != nil {
		if cerr := q.updateProfileWithBasicBlockIdStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateProfileWithBasicBlockIdStmt: %w", cerr)
		}
	}
	if q.updateProfileWithContactBlockIdStmt != nil {
		if cerr := q.updateProfileWithContactBlockIdStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateProfileWithContactBlockIdStmt: %w", cerr)
		}
	}
	if q.updateRefreshTokenStmt != nil {
		if cerr := q.updateRefreshTokenStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateRefreshTokenStmt: %w", cerr)
		}
	}
	if q.updateResolvedLoginStmt != nil {
		if cerr := q.updateResolvedLoginStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateResolvedLoginStmt: %w", cerr)
		}
	}
	if q.updateRoleStmt != nil {
		if cerr := q.updateRoleStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateRoleStmt: %w", cerr)
		}
	}
	if q.updateSavedProfileStmt != nil {
		if cerr := q.updateSavedProfileStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateSavedProfileStmt: %w", cerr)
		}
	}
	if q.updateSocialStmt != nil {
		if cerr := q.updateSocialStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateSocialStmt: %w", cerr)
		}
	}
	if q.updateUserStmt != nil {
		if cerr := q.updateUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateUserStmt: %w", cerr)
		}
	}
	if q.updateUserRoleStmt != nil {
		if cerr := q.updateUserRoleStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateUserRoleStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                                  DBTX
	tx                                  *sql.Tx
	addBasicBlockStmt                   *sql.Stmt
	addContactBlockStmt                 *sql.Stmt
	addContactCategoryStmt              *sql.Stmt
	addContactsStmt                     *sql.Stmt
	addProfileStmt                      *sql.Stmt
	addProfileContentStmt               *sql.Stmt
	addProfileSocialStmt                *sql.Stmt
	addSocialStmt                       *sql.Stmt
	addUserRoleStmt                     *sql.Stmt
	createEmailVerificationStmt         *sql.Stmt
	createOtpStmt                       *sql.Stmt
	createRefreshTokenStmt              *sql.Stmt
	createRoleStmt                      *sql.Stmt
	createSavedProfileStmt              *sql.Stmt
	createUserStmt                      *sql.Stmt
	createUserLoginStmt                 *sql.Stmt
	deleteBasicBlockStmt                *sql.Stmt
	deleteContactStmt                   *sql.Stmt
	deleteContactBlockStmt              *sql.Stmt
	deleteContactCategoryStmt           *sql.Stmt
	deleteEmailVerificationStmt         *sql.Stmt
	deleteOtpStmt                       *sql.Stmt
	deleteProfileStmt                   *sql.Stmt
	deleteProfileContentStmt            *sql.Stmt
	deleteProfileSocialStmt             *sql.Stmt
	deleteRefreshTokenStmt              *sql.Stmt
	deleteRolesStmt                     *sql.Stmt
	deleteSavedProfileStmt              *sql.Stmt
	deleteSocialStmt                    *sql.Stmt
	deleteUserStmt                      *sql.Stmt
	deleteUserLoginStmt                 *sql.Stmt
	getAllContactCategoriesStmt         *sql.Stmt
	getAllContactsStmt                  *sql.Stmt
	getAllContentTypesStmt              *sql.Stmt
	getAllOtpStmt                       *sql.Stmt
	getAllProfilesStmt                  *sql.Stmt
	getBasicBlockStmt                   *sql.Stmt
	getCallToActionStmt                 *sql.Stmt
	getCallToActionsStmt                *sql.Stmt
	getContactBlockStmt                 *sql.Stmt
	getContactCategoryStmt              *sql.Stmt
	getContactsStmt                     *sql.Stmt
	getEmailVerificationStmt            *sql.Stmt
	getEmailVerificationsStmt           *sql.Stmt
	getOtpStmt                          *sql.Stmt
	getProfileStmt                      *sql.Stmt
	getProfileContentStmt               *sql.Stmt
	getProfileContentsStmt              *sql.Stmt
	getProfileSocialStmt                *sql.Stmt
	getProfileSocialsStmt               *sql.Stmt
	getProfilesStmt                     *sql.Stmt
	getRefreshTokenStmt                 *sql.Stmt
	getRefreshTokensStmt                *sql.Stmt
	getRoleStmt                         *sql.Stmt
	getRolesStmt                        *sql.Stmt
	getSavedProfileStmt                 *sql.Stmt
	getSavedProfilesStmt                *sql.Stmt
	getSavedProfilesByEmailStmt         *sql.Stmt
	getSavedProfilesByProfileIdStmt     *sql.Stmt
	getSocialStmt                       *sql.Stmt
	getSocialsStmt                      *sql.Stmt
	getUnResoledLoginsStmt              *sql.Stmt
	getUserStmt                         *sql.Stmt
	getUserLoginStmt                    *sql.Stmt
	getUserLoginsStmt                   *sql.Stmt
	getUserRolesStmt                    *sql.Stmt
	getUsersStmt                        *sql.Stmt
	isProfileExistStmt                  *sql.Stmt
	isUrlExistsStmt                     *sql.Stmt
	updateBasicBlockStmt                *sql.Stmt
	updateContactBlockStmt              *sql.Stmt
	updateContactCategoryStmt           *sql.Stmt
	updateProfileStmt                   *sql.Stmt
	updateProfileContentStmt            *sql.Stmt
	updateProfileSocialStmt             *sql.Stmt
	updateProfileUrlStmt                *sql.Stmt
	updateProfileWithBasicBlockIdStmt   *sql.Stmt
	updateProfileWithContactBlockIdStmt *sql.Stmt
	updateRefreshTokenStmt              *sql.Stmt
	updateResolvedLoginStmt             *sql.Stmt
	updateRoleStmt                      *sql.Stmt
	updateSavedProfileStmt              *sql.Stmt
	updateSocialStmt                    *sql.Stmt
	updateUserStmt                      *sql.Stmt
	updateUserRoleStmt                  *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                                  tx,
		tx:                                  tx,
		addBasicBlockStmt:                   q.addBasicBlockStmt,
		addContactBlockStmt:                 q.addContactBlockStmt,
		addContactCategoryStmt:              q.addContactCategoryStmt,
		addContactsStmt:                     q.addContactsStmt,
		addProfileStmt:                      q.addProfileStmt,
		addProfileContentStmt:               q.addProfileContentStmt,
		addProfileSocialStmt:                q.addProfileSocialStmt,
		addSocialStmt:                       q.addSocialStmt,
		addUserRoleStmt:                     q.addUserRoleStmt,
		createEmailVerificationStmt:         q.createEmailVerificationStmt,
		createOtpStmt:                       q.createOtpStmt,
		createRefreshTokenStmt:              q.createRefreshTokenStmt,
		createRoleStmt:                      q.createRoleStmt,
		createSavedProfileStmt:              q.createSavedProfileStmt,
		createUserStmt:                      q.createUserStmt,
		createUserLoginStmt:                 q.createUserLoginStmt,
		deleteBasicBlockStmt:                q.deleteBasicBlockStmt,
		deleteContactStmt:                   q.deleteContactStmt,
		deleteContactBlockStmt:              q.deleteContactBlockStmt,
		deleteContactCategoryStmt:           q.deleteContactCategoryStmt,
		deleteEmailVerificationStmt:         q.deleteEmailVerificationStmt,
		deleteOtpStmt:                       q.deleteOtpStmt,
		deleteProfileStmt:                   q.deleteProfileStmt,
		deleteProfileContentStmt:            q.deleteProfileContentStmt,
		deleteProfileSocialStmt:             q.deleteProfileSocialStmt,
		deleteRefreshTokenStmt:              q.deleteRefreshTokenStmt,
		deleteRolesStmt:                     q.deleteRolesStmt,
		deleteSavedProfileStmt:              q.deleteSavedProfileStmt,
		deleteSocialStmt:                    q.deleteSocialStmt,
		deleteUserStmt:                      q.deleteUserStmt,
		deleteUserLoginStmt:                 q.deleteUserLoginStmt,
		getAllContactCategoriesStmt:         q.getAllContactCategoriesStmt,
		getAllContactsStmt:                  q.getAllContactsStmt,
		getAllContentTypesStmt:              q.getAllContentTypesStmt,
		getAllOtpStmt:                       q.getAllOtpStmt,
		getAllProfilesStmt:                  q.getAllProfilesStmt,
		getBasicBlockStmt:                   q.getBasicBlockStmt,
		getCallToActionStmt:                 q.getCallToActionStmt,
		getCallToActionsStmt:                q.getCallToActionsStmt,
		getContactBlockStmt:                 q.getContactBlockStmt,
		getContactCategoryStmt:              q.getContactCategoryStmt,
		getContactsStmt:                     q.getContactsStmt,
		getEmailVerificationStmt:            q.getEmailVerificationStmt,
		getEmailVerificationsStmt:           q.getEmailVerificationsStmt,
		getOtpStmt:                          q.getOtpStmt,
		getProfileStmt:                      q.getProfileStmt,
		getProfileContentStmt:               q.getProfileContentStmt,
		getProfileContentsStmt:              q.getProfileContentsStmt,
		getProfileSocialStmt:                q.getProfileSocialStmt,
		getProfileSocialsStmt:               q.getProfileSocialsStmt,
		getProfilesStmt:                     q.getProfilesStmt,
		getRefreshTokenStmt:                 q.getRefreshTokenStmt,
		getRefreshTokensStmt:                q.getRefreshTokensStmt,
		getRoleStmt:                         q.getRoleStmt,
		getRolesStmt:                        q.getRolesStmt,
		getSavedProfileStmt:                 q.getSavedProfileStmt,
		getSavedProfilesStmt:                q.getSavedProfilesStmt,
		getSavedProfilesByEmailStmt:         q.getSavedProfilesByEmailStmt,
		getSavedProfilesByProfileIdStmt:     q.getSavedProfilesByProfileIdStmt,
		getSocialStmt:                       q.getSocialStmt,
		getSocialsStmt:                      q.getSocialsStmt,
		getUnResoledLoginsStmt:              q.getUnResoledLoginsStmt,
		getUserStmt:                         q.getUserStmt,
		getUserLoginStmt:                    q.getUserLoginStmt,
		getUserLoginsStmt:                   q.getUserLoginsStmt,
		getUserRolesStmt:                    q.getUserRolesStmt,
		getUsersStmt:                        q.getUsersStmt,
		isProfileExistStmt:                  q.isProfileExistStmt,
		isUrlExistsStmt:                     q.isUrlExistsStmt,
		updateBasicBlockStmt:                q.updateBasicBlockStmt,
		updateContactBlockStmt:              q.updateContactBlockStmt,
		updateContactCategoryStmt:           q.updateContactCategoryStmt,
		updateProfileStmt:                   q.updateProfileStmt,
		updateProfileContentStmt:            q.updateProfileContentStmt,
		updateProfileSocialStmt:             q.updateProfileSocialStmt,
		updateProfileUrlStmt:                q.updateProfileUrlStmt,
		updateProfileWithBasicBlockIdStmt:   q.updateProfileWithBasicBlockIdStmt,
		updateProfileWithContactBlockIdStmt: q.updateProfileWithContactBlockIdStmt,
		updateRefreshTokenStmt:              q.updateRefreshTokenStmt,
		updateResolvedLoginStmt:             q.updateResolvedLoginStmt,
		updateRoleStmt:                      q.updateRoleStmt,
		updateSavedProfileStmt:              q.updateSavedProfileStmt,
		updateSocialStmt:                    q.updateSocialStmt,
		updateUserStmt:                      q.updateUserStmt,
		updateUserRoleStmt:                  q.updateUserRoleStmt,
	}
}
