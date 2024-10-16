package echo

import (
	"os"
	"strconv"
	"time"

	"github.com/charmbracelet/log"
	"github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	Uuid string `json:"name"`
	jwt.RegisteredClaims
}

func CreateClaim(uuid string) (string, error) {
	tokenAgeEnv := os.Getenv("TOKEN_AGE")
	seconds, err := strconv.Atoi(tokenAgeEnv)
	if err != nil {
		log.Warn("Invalid TokenAge env, revert to 3600")
		seconds = 3600
	}

	tokenAge := time.Now().Add(time.Duration(seconds) * time.Second)

	claims := &JwtClaims{
		Uuid: uuid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(tokenAge),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("ACCESS_TOKEN_KEY")
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return t, nil
}
