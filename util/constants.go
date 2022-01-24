package util

const SUCCESS_RESPONSE_CODE = "00"
const SUCCESS_RESPONSE_MESSAGE = "Success"
const REGISTRATION_SUCCESS_RESPONSE_MESSAGE = "Registration successful..."
const OTP_SENT_RESPONSE_MESSAGE = "If you have an account with us, you should get an otp"

const APPLICATION_NOT_SPECIFIED_ERROR_CODE = "01"
const APPLICATION_NOT_SPECIFIED_ERROR_MESSAGE = "Application not specified"
const INVALID_APPLICATION_ERROR_MESSAGE = "Application is invalid"

const MODEL_VALIDATION_ERROR_CODE = "02"
const MODEL_VALIDATION_ERROR_MESSAGE = "Request is in an incorrect manner. Kindly check your request and try again."

const USER_NAME_OR_PASSWORD_INCORRECT_ERROR_CODE = "03"
const USER_NAME_OR_PASSWORD_INCORRECT_ERROR_MESSAGE = "Oops... something is wrong here... your email or password is incorrect..."

const USER_ALREADY_EXISTS_ERROR_CODE = "04"
const USER_ALREADY_EXISTS_ERROR_MESSAGE = "User already exists"

const OPERATION_FAILED_ERROR_CODE = "05"
const OPERATION_FAILED_ERROR_MESSAGE = "Operation failed. Please contact our support if error persists after 3 trials."

const NO_RECORD_FOUND_ERROR_CODE = "06"
const NO_RECORD_FOUND_ERROR_MESSAGE = "No record found"
const EMAIL_TOKEN_NOT_FOUND = "Sorry we could not verify your request. Please try registering again..."
const USER_NOT_FOUND_RESPONSE_MESSAGE = "User not found"

const DUPLICATE_RECORD_ERROR_CODE = "07"
const DUPLICATE_RECORD_ERROR_MESSAGE = "Record already exist"

const SESSION_EXPIRED_ERROR_CODE = "08"
const SESSION_EXPIRED_ERROR_MESSAGE = "Session expired. Kindly sign in again."
const RESET_PASSWORD_TOKEN_EXPIRED_MESSAGE = "Session expired. Kindly try generating one time password again"
const EMAIL_TOKEN_EXPIRED_MESSAGE = "Oops... something is wrong here... your email verification link has expired.. Kindly register again"
const OTP_EXPIRED_MESSAGE = "Oops... something is wrong here... your one time token has expired. Kindly request another one..."

const SESSION_STILL_ACTIVE_ERROR_CODE = "09"
const SESSION_STILL_ACTIVE_ERROR_MESSAGE = "Session is still valid..."

const INVALID_AUTHENTICATION_SCHEME_ERROR_CODE = "10"
const INVALID_AUTHENTICATION_SCHEME_ERROR_MESSAGE = "Unsupported authentication scheme type"

const UNAUTHORIZED_ERROR_CODE = "11"
const UNAUTHORIZED_ERROR_MESSAGE = "Sorry, you are not authorized to carry out this operation."
const UNAUTHORIZED_ERROR_MESSAGE_WRONG_JWT = "Sorry we cannot proceed with your login. Reason: Invalid token."

const ACCOUNT_LOCKOUT_ERROR_CODE = "12"
const ACCOUNT_LOCKOUT_ERROR_MESSAGE = "Sorry you exceeded the maximum login attempts... Kindly reset your password to continue..."

const SQL_ERROR_CODE = "98"
const SQL_ERROR_MESSAGE = "Error occured while running an operation. Please contact our support if error persists after 3 trials."

const GENERAL_ERROR_CODE = "99"
const GENERAL_ERROR_MESSAGE = "Unexpected error occured. Please contact our support if error persists after 3 trials."
