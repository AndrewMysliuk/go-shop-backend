package handler

import (
	"github.com/AndrewMislyuk/go-shop-backend/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type getCreationId struct {
	Id string `json:"id"`
}

type getUserToken struct {
	Token string `json:"access_token"`
}

// @Summary SignUp
// @Tags Auth
// @Description user sign-up
// @ID sign-up
// @Accept  json
// @Produce  json
// @Param input body domain.UserSignUp true "User signup"
// @Success 200 {object} getCreationId
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var input domain.UserSignUp

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.userService.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getCreationId{
		Id: id,
	})
}

// @Summary SignIn
// @Tags Auth
// @Description user sign-in
// @ID sign-in
// @Accept  json
// @Produce  json
// @Param input body domain.UserSignIn true "User login"
// @Success 200 {object} getUserToken
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var input domain.UserSignIn

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.userService.GenerateToken(input.Email, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getUserToken{
		Token: token,
	})
}
