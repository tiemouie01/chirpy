package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type RequestError struct {
	StatusCode int
	Err        error
}

func (r *RequestError) Error() string {
	return fmt.Sprintf("status %d: err %v", r.StatusCode, r.Err)
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", errors.New(err.Error())
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	currentTime := time.Now().UTC()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    "chirpy",
		IssuedAt:  currentTime.Unix(),
		ExpiresAt: currentTime.UnixMicro() + expiresIn.Microseconds(),
		Subject:   userID.String(),
	})

	signedToken, err := token.SignedString(tokenSecret)
	if err != nil {
		return "", errors.New(err.Error())
	}
	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		formattedError := RequestError{
			StatusCode: 401,
		}
		return uuid.Nil, &formattedError
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			return uuid.Nil, errors.New("invalid user ID in token")
		}
		return userID, nil
	}

	return uuid.Nil, errors.New("invalid token")
}

func GetBearerToken(headers http.Header) (string, error) {
	reqToken := headers.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]

	if reqToken != "" {
		return "", errors.New("token does not exist")
	}
	return reqToken, nil
}
