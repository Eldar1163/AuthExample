package handler

import (
	"TestTask/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signUp(c *gin.Context) {
	var input model.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	_, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSONP(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}

func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.services.CheckUserCredentials(input.Username, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "Incorrect credentials")
		return
	}
	badAuths, err := h.services.GetBadAuthAttemptsCnt(user.Id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Internal error")
		return
	}
	if badAuths > 5 {
		if errBlock := h.services.BlockUser(input.Username); errBlock != nil {
			newErrorResponse(c, http.StatusInternalServerError, "Internal error")
			return
		} else {
			newErrorResponse(c, http.StatusUnauthorized, "You have been blocked by system")
			return
		}
	}

	token, err := h.services.GenerateToken(user.Id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Internal error")
	}

	c.JSONP(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
