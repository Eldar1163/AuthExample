package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	TOKENHEADER = "X-Token"
)

func (h *Handler) clearAudit(c *gin.Context) {
	token := c.GetHeader(TOKENHEADER)
	if len(token) == 0 {
		newErrorResponse(c, http.StatusBadRequest, "Token must be in request")
		return
	}
	userId, err := h.services.CheckToken(token)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Cannot check token, try again later")
		return
	}
	err = h.services.ClearAudit(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Cannot clear audit, try again later")
		return
	}
	c.JSONP(http.StatusOK, map[string]interface{}{
		"status": "audit is clean",
	})
}

func (h *Handler) getAudit(c *gin.Context) {
	token := c.GetHeader(TOKENHEADER)
	if len(token) == 0 {
		newErrorResponse(c, http.StatusBadRequest, "Token must be in request")
		return
	}
	events, err := h.services.GetUserEvents(token)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Try again later")
		return
	}

	if len(events) != 0 {
		c.JSONP(http.StatusOK, events)
		return
	}
	c.JSONP(http.StatusOK, map[string]interface{}{})
}
