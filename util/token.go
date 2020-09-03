package util

import (
	"authengine/models"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// GenerateHash generates hash using bcrypt
func GenerateHash(rawString string) string {
	byteString := []byte(rawString)
	hash, err := bcrypt.GenerateFromPassword(byteString, bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

// VerifyHash verifies password using bcrypt
func VerifyHash(correctPassword string, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(correctPassword), []byte(plainPassword))
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

var jwtSecretKey = []byte("AnyString")

// GenerateJWT func will used to create the JWT while signing in and signing out
func GenerateJWT(email string, role string) (response string, refreshToken string, err error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &models.Claims{
		Email: email,
		Role:  role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "Persian Black",
		},
	}
	jwtKey := os.Getenv("JWT_SECRET_KEY")
	if jwtKey == "" {
		log.Println("JWT secret key cannot be empty")
		log.Println("JWT_SECRET_KEY cannot be empty, application intialization failed...")
	} else {
		log.Println(fmt.Sprintf("Setting JWT secret key..."))
		jwtSecretKey = []byte(jwtKey)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecretKey)
	if err == nil {
		log.Println("Auth token successfully generated...")
		hasher := md5.New()
		hasher.Write([]byte(tokenString))
		refreshToken := hex.EncodeToString(hasher.Sum(nil))
		log.Println(fmt.Sprintf("Refresh Token successfully generated..."))
		return tokenString, refreshToken, nil
	}
	log.Println(fmt.Sprintf("Error occured in generating token: %s", err))
	return "", "", err
}

// VerifyToken func will used to Verify the JWT Token while using APIS
func VerifyToken(tokenString string) (verifiedClaims models.VerifiedClaims, err error) {
	claims := &models.Claims{}

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
