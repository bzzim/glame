package responses

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Base struct {
	Success bool        `json:"success"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data"`
}

func NewSuccessCodeResponse(c *gin.Context, code int, data interface{}) {
	c.JSON(code, Base{
		Success: true,
		Error:   "",
		Data:    data,
	})
}

func NewSuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Base{
		Success: true,
		Data:    data,
	})
}
