package middleware

import (
	"errors"
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/service"
	"github.com/beecool-cocktail/application-backend/viewmodels"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type PayloadData struct {
	UserID  int64  `json:"user_id"`
	Account string `json:"account"`
	Name    string `json:"name"`
}

type MyClaims struct {
	PayloadData
	jwt.StandardClaims
}

var secret = []byte("secret")

type Handler struct {
	service *service.Service
}

func NewMiddlewareHandler(redis *service.Service) *Handler {
	return &Handler{service: redis}
}

func (h *Handler) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Content-OrderType", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Headers", " Content-Type,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func GenToken(data PayloadData) (string, error) {
	c := MyClaims{
		PayloadData: data,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(2 * time.Hour).Unix(),
			Issuer:    "GiftForm69King",
			Audience:  data.Account,
		},
	}
	// Choose specific algorithm
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// Choose specific Signature
	return token.SignedString(secret)
}

func parseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	// Valid token
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// JWTAuthMiddleware Middleware of JWT
func (h *Handler) JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// Get token from Header.Authorization field.
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, viewmodels.ResponseData{ErrorCode: domain.GetErrorCode(domain.ErrTokenExpired), ErrorMessage: domain.ErrTokenExpired.Error()})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, viewmodels.ResponseData{ErrorCode: domain.GetErrorCode(domain.ErrTokenExpired), ErrorMessage: domain.ErrTokenExpired.Error()})
			return
		}
		// parts[0] is Bearer, parts is token.
		token := parts[1]
		mc, err := parseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, viewmodels.ResponseData{ErrorCode: domain.GetErrorCode(domain.ErrTokenExpired), ErrorMessage: domain.ErrTokenExpired.Error()})
			return
		}

		//key := "admin:admin_id:" + strconv.FormatInt(mc.UserID, 10)
		//user, err := h.service.Redis.HMGet(key, "access_token").Result()
		//if err != nil {
		//	c.AbortWithStatusJSON(http.StatusUnauthorized, viewmodels.ResponseData{ErrorCode: domain.GetErrorCode(domain.ErrTokenExpired), ErrorMessage: domain.ErrTokenExpired.Error()})
		//	return
		//}
		//
		//serverAccessToken, ok := user[0].(string)
		//if !ok {
		//	c.AbortWithStatusJSON(domain.GetStatusCode(domain.ErrTokenExpired), viewmodels.ResponseData{ErrorCode: domain.GetErrorCode(domain.ErrTokenExpired), ErrorMessage: domain.ErrTokenExpired.Error()})
		//}
		//
		//if !isTokenValid(token, serverAccessToken) {
		//	c.AbortWithStatusJSON(domain.GetStatusCode(domain.ErrTokenExpired), viewmodels.ResponseData{ErrorCode: domain.GetErrorCode(domain.ErrTokenExpired), ErrorMessage: domain.ErrTokenExpired.Error()})
		//}

		// Store info into Context
		c.Set("account", mc.Account)
		c.Set("user_id", mc.UserID)
		c.Next()
	}
}

// JWTAuthMiddlewareIfExist Middleware of JWT if exist
// some api needs to authenticate specific resource status if they are login
func (h *Handler) JWTAuthMiddlewareIfExist() func(c *gin.Context) {
	return func(c *gin.Context) {
		// Get token from Header.Authorization field.
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.Set("auth", false)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, viewmodels.ResponseData{ErrorCode: domain.GetErrorCode(domain.ErrTokenExpired), ErrorMessage: domain.ErrTokenExpired.Error()})
			return
		}
		// parts[0] is Bearer, parts is token.
		token := parts[1]
		mc, err := parseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, viewmodels.ResponseData{ErrorCode: domain.GetErrorCode(domain.ErrTokenExpired), ErrorMessage: domain.ErrTokenExpired.Error()})
			return
		}

		//key := "admin:admin_id:" + strconv.FormatInt(mc.UserID, 10)
		//user, err := h.service.Redis.HMGet(key, "access_token").Result()
		//if err != nil {
		//	c.AbortWithStatusJSON(http.StatusUnauthorized, viewmodels.ResponseData{ErrorCode: domain.GetErrorCode(domain.ErrTokenExpired), ErrorMessage: domain.ErrTokenExpired.Error()})
		//	return
		//}
		//
		//serverAccessToken, ok := user[0].(string)
		//if !ok {
		//	c.AbortWithStatusJSON(domain.GetStatusCode(domain.ErrTokenExpired), viewmodels.ResponseData{ErrorCode: domain.GetErrorCode(domain.ErrTokenExpired), ErrorMessage: domain.ErrTokenExpired.Error()})
		//}
		//
		//if !isTokenValid(token, serverAccessToken) {
		//	c.AbortWithStatusJSON(domain.GetStatusCode(domain.ErrTokenExpired), viewmodels.ResponseData{ErrorCode: domain.GetErrorCode(domain.ErrTokenExpired), ErrorMessage: domain.ErrTokenExpired.Error()})
		//}

		// Store info into Context
		c.Set("auth", true)
		c.Set("account", mc.Account)
		c.Set("user_id", mc.UserID)
		c.Next()
	}
}

func isTokenValid(requestToken string, serverToken string) bool {
	if requestToken != serverToken {
		return false
	}

	return true
}

func parseTime(myTime string) time.Time {
	parseTime, _ := time.Parse("2006-01-02 15:04:05", myTime)
	return parseTime
}
