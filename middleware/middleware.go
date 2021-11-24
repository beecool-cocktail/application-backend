package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"time"
)

type Handler struct {
	redisAdp *redis.Client
}

func NewMiddlewareHandler(redis *redis.Client) *Handler {
	return &Handler{redisAdp:redis}
}

func (h *Handler) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Content-OrderType", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Headers"," Content-Type,Access-Control-Allow-Origin,Access-Control-Allow-Headers,x-token")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

//func (h *Handler) AuthMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		session := sessions.Default(c)
//		id := session.Get("id").(int64)
//		key := "client:client_id:" + strconv.FormatInt(id, 10)
//		token := c.Request.Header.Get("x-token")
//		client, err := h.redisAdp.HMGet(key, "access_token", "token_expire").Result()
//		if err != nil {
//			c.AbortWithStatusJSON(http.StatusUnauthorized, viewmodels.ResponseData{ErrorCode:"", ErrorMessage:""})
//		}
//
//		accessToken, ok := client["access_token"].(string)
//		if !ok {
//			c.AbortWithStatusJSON(http.StatusUnauthorized, viewmodels.ResponseData{ErrorCode:"", ErrorMessage:""})
//		}
//
//		tokenExpire, ok := client["token_expire"].(string)
//		if !ok {
//			c.AbortWithStatusJSON(http.StatusUnauthorized, viewmodels.ResponseData{ErrorCode:"", ErrorMessage:""})
//		}
//
//		if !isTokenValid(token, accessToken, tokenExpire) {
//			c.AbortWithStatusJSON(http.StatusUnauthorized, viewmodels.ResponseData{ErrorCode:"", ErrorMessage:""})
//		}
//
//		c.Next()
//	}
//}

func isTokenValid(requestToken string, serverToken string, expire string) bool {
	expireTime := parseTime(expire)
	if requestToken != serverToken || expireTime.Before(time.Now()) {
		return false
	}

	return true
}

func parseTime(myTime string) time.Time {
	parseTime, _ := time.Parse("2006-01-02 15:04:05", myTime)
	return parseTime
}