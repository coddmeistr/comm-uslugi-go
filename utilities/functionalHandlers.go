package utilities

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handling OPTION request coming to all cors requests to prevent block
func HandleOptions(c *gin.Context) {
	if c.Request.Method == http.MethodOptions {
		c.JSON(http.StatusOK, gin.H{})
	}
}
