package customjwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/q2rd/gRPC_sso_go/internal/domain/models"
	"time"
)

func NewToken(user models.UserDatabase, app models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uuid"] = user.Id
	claims["EMAIL"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["appId"] = app.Id

	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
