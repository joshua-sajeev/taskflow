package ratelimiter

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

func TestNewIPRateLimiter(t *testing.T) {
	limiter := NewIPRateLimiter(rate.Limit(5), 10)

	assert.NotNil(t, limiter)
	assert.NotNil(t, limiter.ips)
	assert.NotNil(t, limiter.mu)
	assert.Equal(t, rate.Limit(5), limiter.r)
	assert.Equal(t, 10, limiter.b)
}

func TestAddIP(t *testing.T) {
	limiter := NewIPRateLimiter(rate.Limit(5), 10)
	ip := "192.168.1.1"

	rateLimiter := limiter.AddIP(ip)

	assert.NotNil(t, rateLimiter)
	assert.Equal(t, 1, len(limiter.ips))
}

func TestGetLimiter(t *testing.T) {
	limiter := NewIPRateLimiter(rate.Limit(5), 10)
	ip := "192.168.1.1"

	// First call should create a new limiter
	rateLimiter1 := limiter.GetLimiter(ip)
	assert.NotNil(t, rateLimiter1)
	assert.Equal(t, 1, len(limiter.ips))

	// Second call should return the same limiter
	rateLimiter2 := limiter.GetLimiter(ip)
	assert.NotNil(t, rateLimiter2)
	assert.Equal(t, 1, len(limiter.ips))
	assert.Equal(t, rateLimiter1, rateLimiter2)
}

func TestMiddleware_AllowsRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a rate limiter with generous limits
	limiter := NewIPRateLimiter(rate.Limit(100), 100)

	r := gin.New()
	r.Use(limiter.Middleware())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Make a request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["message"])
}

func TestMiddleware_RateLimitsRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a rate limiter with very strict limits (1 request per second, burst of 2)
	limiter := NewIPRateLimiter(rate.Limit(1), 2)

	r := gin.New()
	r.Use(limiter.Middleware())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	successCount := 0
	rateLimitedCount := 0

	// Make 5 rapid requests
	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:1234" // Same IP for all requests
		r.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			successCount++
		} else if w.Code == http.StatusTooManyRequests {
			rateLimitedCount++
		}
	}

	// First 2 requests should succeed (burst capacity)
	// Remaining 3 should be rate limited
	assert.Equal(t, 2, successCount, "Expected 2 successful requests")
	assert.Equal(t, 3, rateLimitedCount, "Expected 3 rate limited requests")
}

func TestMiddleware_DifferentIPs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a rate limiter with strict limits
	limiter := NewIPRateLimiter(rate.Limit(1), 2)

	r := gin.New()
	r.Use(limiter.Middleware())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Test with different IPs
	ips := []string{"192.168.1.1:1234", "192.168.1.2:1234", "192.168.1.3:1234"}

	for _, ip := range ips {
		// Each IP should get 2 successful requests (burst)
		for i := 0; i < 2; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			req.RemoteAddr = ip
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code, "Request should succeed for different IPs")
		}

		// Third request from same IP should be rate limited
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = ip
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTooManyRequests, w.Code, "Third request should be rate limited")
	}
}

func TestMiddleware_RateLimitResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a rate limiter with very strict limits
	limiter := NewIPRateLimiter(rate.Limit(1), 1)

	r := gin.New()
	r.Use(limiter.Middleware())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// First request should succeed
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:1234"
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Second request should be rate limited
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:1234"
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTooManyRequests, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "rate limit exceeded, please try again later", response["error"])
}

func TestCleanupOldEntries(t *testing.T) {
	limiter := NewIPRateLimiter(rate.Limit(5), 10)

	// Add some IPs
	limiter.GetLimiter("192.168.1.1")
	limiter.GetLimiter("192.168.1.2")
	limiter.GetLimiter("192.168.1.3")

	assert.Equal(t, 3, len(limiter.ips))

	// Clean up entries
	limiter.CleanupOldEntries()

	assert.Equal(t, 0, len(limiter.ips))
}

func TestStartCleanupRoutine(t *testing.T) {
	limiter := NewIPRateLimiter(rate.Limit(5), 10)

	// Add some IPs
	limiter.GetLimiter("192.168.1.1")
	limiter.GetLimiter("192.168.1.2")
	assert.Equal(t, 2, len(limiter.ips))

	// Start cleanup routine with very short interval
	limiter.StartCleanupRoutine(100 * time.Millisecond)

	// Wait for cleanup to run
	time.Sleep(200 * time.Millisecond)

	// IPs should be cleaned up
	assert.Equal(t, 0, len(limiter.ips))
}
