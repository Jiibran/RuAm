package service

import (
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"RuAm/internal/pkg/config"
	"RuAm/internal/pkg/db"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey []byte

func InitService(cfg config.Config) {
	jwtKey = []byte(cfg.JWT.SecretKey)
}

func Register(email, password string) error {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return err
	}

	query := `INSERT INTO users (email, password) VALUES ($1, $2)`
	_, err = db.DB.Exec(query, email, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}

func Login(email, password string) (string, error) {
	var storedPassword string
	query := `SELECT password FROM users WHERE email = $1`
	err := db.DB.QueryRow(query, email).Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("user not found")
		}
		return "", err
	}

	if !checkPasswordHash(password, storedPassword) {
		return "", errors.New("invalid password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateRefreshToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func RefreshToken(refreshToken string) (string, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		email := claims["email"].(string)
		newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"email": email,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		})

		tokenString, err := newToken.SignedString(jwtKey)
		if err != nil {
			return "", err
		}

		return tokenString, nil
	} else {
		return "", errors.New("invalid refresh token")
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
