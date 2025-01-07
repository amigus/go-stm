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

func TestPublisher(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	tm := stm.UUIDTokenManager(1, 5, time.Minute)

	r = TokenPublisher(r, tm, "/")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Expected status OK")
	assert.NotEmpty(t, w.Body.String(), "Expected a non-empty token in response")
	if err := uuid.Validate(w.Body.String()); err != nil {
		t.Error("Expected a valid UUID in response", err)
	}
}
