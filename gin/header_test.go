package stm_gin

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/amigus/go-stm"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_TokenHeaderChecker(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	tm := stm.UUIDTokenManager(1, 1, time.Hour)

	validToken := tm.Get()
	r = HeaderChecker(r, tm, "X-Token")

	r.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Token", validToken)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Expected status OK for valid token")
}

func Test_TokenHeaderCheckerWithInvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	tm := stm.UUIDTokenManager(1, 1, time.Hour)

	invalidToken := uuid.NewString()
	r = HeaderChecker(r, tm, "X-Token")

	r.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Token", invalidToken)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code, "Expected status OK for valid token")
}
