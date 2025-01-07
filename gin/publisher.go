package stm_gin

import (
	"net/http"

	"github.com/amigus/go-stm"
	"github.com/gin-gonic/gin"
)

// TokenPublisher adds a route to the gin engine for GET on the given path that returns a valid token.
func TokenPublisher(r *gin.Engine, tm stm.TokenManager, path string) *gin.Engine {
	r.GET(path, func(c *gin.Context) { c.String(http.StatusOK, tm.Get()) })

	return r
}
