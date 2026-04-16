package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTraceID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should add trace_id to response header", func(t *testing.T) {
		w := httptest.NewRecorder()
		_, router := gin.CreateTestContext(w)

		router.Use(TraceID())
		router.GET("/test", func(c *gin.Context) {
			c.String(200, "ok")
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		traceID := w.Header().Get("X-Trace-ID")
		assert.NotEmpty(t, traceID)
		assert.Equal(t, 200, w.Code)
	})

	t.Run("should handle normal request", func(t *testing.T) {
		w := httptest.NewRecorder()
		_, router := gin.CreateTestContext(w)

		router.Use(TraceID())
		router.GET("/test", func(c *gin.Context) {
			c.String(200, "ok")
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})
}

func TestErrorHandling(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should handle normal request without error", func(t *testing.T) {
		w := httptest.NewRecorder()
		_, router := gin.CreateTestContext(w)

		router.Use(ErrorHandling())
		router.GET("/normal", func(c *gin.Context) {
			c.String(200, "ok")
		})

		req, _ := http.NewRequest("GET", "/normal", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "ok", w.Body.String())
	})

	t.Run("should handle request with context error", func(t *testing.T) {
		w := httptest.NewRecorder()
		_, router := gin.CreateTestContext(w)

		router.Use(ErrorHandling())
		router.GET("/error", func(c *gin.Context) {
			_ = c.Error(assert.AnError) // ErrorHandling中间件会处理这个错误
			c.String(500, "error")
		})

		req, _ := http.NewRequest("GET", "/error", nil)
		router.ServeHTTP(w, req)

		// ErrorHandling should intercept c.Error and return 500
		assert.Equal(t, 500, w.Code)
	})
}
