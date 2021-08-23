# Authengine

Authentication engine for all Persian Black Applications

## Description

The authengine is used to authenticate applications across the Persian Black space. It is built to support OAuth 2.0 protocol.

### Response Codes

+ 00 | Success | The operation/registration was successful
+ 01 | Error   | Application not specified/Application is invalid
+ 02 | Error   | Request is in an incorrect manner. Kindly check your request and try again.
+ 03 | Error   | Oops... something is wrong here... your email or password is incorrect...
+ 04 | Error   | User already exists
+ 05 | Error   | Operation failed. Please contact our support if error persists after 3 trials.
+ 06 | Error   | No record found/Sorry we could not verify your request. Please try registering again.../User not found
+ 07 | Error   | Record already exist
+ 08 | Error   | Session expired. Kindly sign in again./Session expired. Kindly try generating one time password again/Oops... something is wrong here... your email verification link has expired.. Kindly register again"/Oops... something is wrong here... your one time token has expired. Kindly request another one..."
+ 09 | Error   | Session is still valid...
+ 10 | Error   | Unsupported authentication scheme type
+ 11 | Error   | Sorry, you are not authorized to carry out this operation.
+ 12 | Error   | Sorry you exceeded the maximum login attempts... Kindly reset your password to continue...
+ 98 | Error   | Error occured while running an operation. Please contact our support if error persists after 3 trials.
+ 99 | Error   | Unexpected error occured. Please contact our support if error persists after 3 trials.

