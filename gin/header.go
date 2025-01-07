package stm_gin

import (
	"net/http"

	"github.com/amigus/go-stm"
	"github.com/gin-gonic/gin"
)

// Header adds a middleware to the gin engine that requires a valid token in the given header.
func HeaderChecker(r *gin.Engine, tm stm.TokenManager, headerName string) *gin.Engine {
	r.Use(func(c *gin.Context) {
		if tm.Check(c.GetHeader(headerName)) {
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		}
	})

	return r
}
