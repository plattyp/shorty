package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Index response to generic index endpoint
func Index(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		gin.H{
			"status":  true,
			"message": "All systems go.",
		},
	)
}
