package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
	userRole            = "userRole"
)

func (h *Handler) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)

			return
		}

		c.Next()
	}
}

func (h *Handler) loggingMiddleware(c *gin.Context) {
	logrus.Infof("[%s] - %s", c.Request.Method, c.Request.RequestURI)
}

func (h *Handler) userIdentify(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")

		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")

		return
	}

	if headerParts[0] != "Bearer" {
		newErrorResponse(c, http.StatusUnauthorized, "invalid bearer")

		return
	}

	if headerParts[1] == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty token")

		return
	}

	user, err := h.userService.GetMe(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.Set(userCtx, user.Id)
	c.Set(userRole, user.Role)
}

func (h *Handler) userIsAdmin(c *gin.Context) {
	role, exist := c.Get(userRole)
	if !exist {
		newErrorResponse(c, http.StatusUnauthorized, "you are unauthorized")

		return
	}

	roleStr, ok := role.(string)
	if !ok {
		newErrorResponse(c, http.StatusUnauthorized, "type error")

		return
	}

	if roleStr != "ADMIN" {
		newErrorResponse(c, http.StatusUnauthorized, "you don't have admin role")

		return
	}
}
