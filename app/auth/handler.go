package auth

import (
	"net/http"

	"github.com/AnggaPutraa/talk-backend/exceptions"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service AuthService
}

func NewAuthHandler(service AuthService) *AuthHandler {
	handler := &AuthHandler{
		service: service,
	}
	return handler
}

func (h *AuthHandler) Register(c *gin.Context) {
	var request RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, exceptions.ErrorResponse(err))
		return
	}
	token, err := h.service.CreateUser(c, &request)
	if err != nil {
		c.JSON(http.StatusBadRequest, exceptions.ErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, token)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, exceptions.ErrorResponse(err))
		return
	}
	token, err := h.service.LoginUser(c, &request)
	if err != nil {
		c.JSON(http.StatusBadRequest, exceptions.ErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, token)
}
