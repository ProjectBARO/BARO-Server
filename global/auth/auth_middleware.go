package auth

import (
	"errors"
	"gdsc/baro/global"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ContextKey string

const UserIDKey ContextKey = "user_id"

type authenticationMiddleware struct {
	secret string
}

func NewAuthentication(secret string) *authenticationMiddleware {
	return &authenticationMiddleware{secret: secret}
}

func (a *authenticationMiddleware) StripTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := getTokenFromRequest(c.Request)
		if err != nil {
			c.JSON(400, global.Response{
				Status:  400,
				Message: err.Error(),
				Data:    "failed",
			})
			return
		}

		claim, err := ValidateToken(token, a.secret)
		if err != nil {
			c.JSON(400, global.Response{
				Status:  400,
				Message: err.Error(),
				Data:    "failed",
			})
			return
		}

		c.Set(string(UserIDKey), claim["sub"])
		c.Next()
	}
}

func getTokenFromRequest(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	splitToken := strings.Split(authHeader, "Bearer")
	if len(splitToken) != 2 {
		return "", errors.New("invalid authorization header format")
	}

	return strings.TrimSpace(splitToken[1]), nil
}
