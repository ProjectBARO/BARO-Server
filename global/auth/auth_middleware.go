package auth

import (
	"context"
	"errors"
	"gdsc/baro/global"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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

func UnaryAuthInterceptor(c context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if info.FullMethod == "/video.VideoService/GetVideos" ||
		info.FullMethod == "/video.VideoService/GetVideosByCategory" ||
		info.FullMethod == "/user.UserService/Login" {
		return handler(c, req)
	}

	md, ok := metadata.FromIncomingContext(c)
	if !ok {
		return nil, errors.New("missing metadata")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, errors.New("authorization header is required")
	}

	token := strings.TrimPrefix(authHeader[0], "Bearer ")
	claim, err := ValidateToken(token, os.Getenv("JWT_SECRET"))
	if err != nil {
		return nil, errors.New("invalid token")
	}

	c = context.WithValue(c, UserIDKey, claim["sub"])

	return handler(c, req)
}
