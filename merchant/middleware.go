package merchant

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AccessKey is merchant write token
var AccessKey = "RakawmMerchantToken"

// CheckHeaders check header before continuing process
func CheckHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {

		/*
		   For now we just check if current request has valid
		   parameter from header.
		*/
		headers := c.Request.Header
		accessKey := headers.Get("Admin-Token")
		if accessKey != AccessKey {
			c.AbortWithStatus(http.StatusUnauthorized)
			c.JSON(http.StatusUnauthorized, gin.H{"status_code": http.StatusUnauthorized, "status_message": "Invalid access key or access key not found."})
		}

		c.Next()

	}
}
