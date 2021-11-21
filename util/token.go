package util

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"helloprofile.com/models"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// GenerateHash generates hash using bcrypt
func GenerateHash(rawString string) string {
	fields := log.Fields{"microservice": "helloprofile.service", "function": "GenerateHash"}
	byteString := []byte(rawString)
	hash, err := bcrypt.GenerateFromPassword(byteString, bcrypt.DefaultCost)
	if err != nil {
		log.WithFields(fields).Error("Error occured generating hash")
	}
	return string(hash)
}

// VerifyHash verifies password using bcrypt
func VerifyHash(correctPassword string, plainPassword string) bool {
	fields := log.Fields{"microservice": "helloprofile.service", "function": "VerifyHash"}
	err := bcrypt.CompareHashAndPassword([]byte(correctPassword), []byte(plainPassword))
	if err != nil {
		log.WithFields(fields).Error("Error occured trying to verify hash")
		return false
	}

	return true
}

var jwtSecretKey = []byte("AnyString")

// GenerateJWT func will used to create the JWT while signing in and signing out
func GenerateJWT(email string, role []string) (response string, refreshToken string, err error) {
	fields := log.Fields{"microservice": "helloprofile.service", "function": "VerifyHash"}
	authExpirationTime, _ := time.ParseDuration(os.Getenv("TOKEN_LIFESPAN"))
	expirationTime := time.Now().Add(authExpirationTime)
	claims := &models.Claims{
		Email: email,
		Role:  strings.Join(role, ":"),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "Persian Black",
		},
	}
	jwtKey := os.Getenv("JWT_SECRET_KEY")
	if jwtKey == "" {
		log.WithFields(fields).Error("JWT secret key cannot be empty")
		log.WithFields(fields).Error("JWT_SECRET_KEY cannot be empty, application intialization failed...")
	} else {
		log.WithFields(fields).Error("Setting JWT secret key...")
		jwtSecretKey = []byte(jwtKey)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecretKey)
	if err == nil {
		log.WithFields(fields).WithError(err).Error("Auth token successfully generated...")
		hasher := md5.New()
		hasher.Write([]byte(tokenString))
		refreshToken := hex.EncodeToString(hasher.Sum(nil))
		log.WithFields(fields).WithError(err).Error("Refresh Token successfully generated...")
		return tokenString, refreshToken, nil
	}
	log.WithFields(fields).WithError(err).Error("Error occured in generating token")
	return "", "", err
}

// VerifyToken func will used to Verify the JWT Token while using APIS
func VerifyToken(tokenString string) (verifiedClaims models.VerifiedClaims, err error) {
	fields := log.Fields{"microservice": "helloprofile.service", "function": "VerifyHash"}
	claims := &models.Claims{}
	jwtKey := os.Getenv("JWT_SECRET_KEY")
	if jwtKey == "" {
		log.WithFields(fields).Error("JWT secret key cannot be empty")
		log.WithFields(fields).Error("JWT_SECRET_KEY cannot be empty, application intialization failed...")
	} else {
		log.WithFields(fields).Error("Setting JWT secret key...")
		jwtSecretKey = []byte(jwtKey)
	}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})

	verifiedClaim := models.VerifiedClaims{}

	if token != nil && token.Valid {
		verifiedClaim.Email = claims.Email
		verifiedClaim.Role = claims.Role
		return verifiedClaim, nil
	}
	if token != nil && !token.Valid {
		verifiedClaim.Email = claims.Email
		verifiedClaim.Role = claims.Role
		return verifiedClaim, err
	}
	return verifiedClaim, err

}

const otpChars = "1234567890"

// GenerateOTP is used to generate OTP of fixed length
func GenerateOTP(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	otpCharsLength := len(otpChars)
	for i := 0; i < length; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}

	return string(buffer), nil
}
