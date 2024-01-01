package auth

import (
	"errors"
	"gdsc/baro/models"
	"gdsc/baro/types"
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
			c.JSON(400, types.Response{
				Status:  400,
				Message: err.Error(),
				Data:    "failed",
			})
			return
		}

		claim, err := ValidateToken(token, a.secret)
		if err != nil {
			c.JSON(400, types.Response{
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

func FindCurrentUser(c *gin.Context) *models.User {
	userID, exists := c.Get(string(UserIDKey))
	if !exists {
		c.JSON(400, types.Response{
			Status:  400,
			Message: "Not Found UserID in Context!",
			Data:    "failed",
		})
		return nil
	}

	var user models.User
	if err := models.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(400, types.Response{
			Status:  400,
			Message: "Not Found User!",
			Data:    "failed",
		})
		return nil
	}

	return &user
}
