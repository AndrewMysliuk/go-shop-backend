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

	user, err := h.userService.GetMe(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
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

	roleStr := role.(string)
	if roleStr != "ADMIN" {
		newErrorResponse(c, http.StatusUnauthorized, "you don't have admin role")
		return
	}
}
