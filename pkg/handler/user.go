package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sber/types"
)

type Response struct {
	Data Data `json:"data"`
}
type Data struct {
	User *types.User `json:"user"`
}

// @Summary GetById
// @Tags api
// @Description Get employee tree by Id
// @ID getById
// @Accept json
// @Produce json
// @Param id path string true "employee id"
// @Success 200 {object} Response
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/{id} [get]
func (h Handler) GetById(c *gin.Context) {
	userId := c.Param("id")

	user, err := h.service.User.GetById(c.Request.Context(), userId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, Response{
		Data: Data{
			User: user,
		},
	})
}
