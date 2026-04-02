package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	handlers "github.com/Dev66-66/LB10/go-service/internal/handlers/http"
	"github.com/Dev66-66/LB10/go-service/internal/middleware"
	"github.com/Dev66-66/LB10/go-service/internal/store"
)

const testSecret = "test-secret"

func init() {
	gin.SetMode(gin.TestMode)
}

// newTestRouter builds a router that mirrors production: JWT on /workouts/*.
func newTestRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	r.POST("/auth/login", handlers.NewAuthHandler(testSecret).Login)

	protected := r.Group("/")
	protected.Use(middleware.JWT(testSecret))
	handlers.NewWorkoutHandler(store.NewWorkoutStore()).RegisterRoutes(protected)

	return r
}

// validBearerToken generates a signed token with testSecret valid for 1 hour.
func validBearerToken(t *testing.T) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": "admin",
		"exp":      time.Now().Add(time.Hour).Unix(),
	})
	signed, err := token.SignedString([]byte(testSecret))
	require.NoError(t, err)
	return "Bearer " + signed
}

// --- Workout handler tests ---

func TestPostWorkout_ValidBody_Returns201(t *testing.T) {
	r := newTestRouter()

	body := `{"name":"Push Day","type":"strength","duration":60,"difficulty":"hard","calories_burned":400}`
	req := httptest.NewRequest(http.MethodPost, "/workouts", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", validBearerToken(t))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "Push Day", resp["name"])
	assert.Equal(t, "strength", resp["type"])
}

func TestPostWorkout_EmptyName_Returns400(t *testing.T) {
	r := newTestRouter()

	body := `{"name":"","type":"cardio","duration":30,"difficulty":"easy","calories_burned":200}`
	req := httptest.NewRequest(http.MethodPost, "/workouts", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", validBearerToken(t))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Contains(t, resp["error"], "name")
}

func TestGetWorkouts_ReturnsArrayNotNull(t *testing.T) {
	r := newTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/workouts", nil)
	req.Header.Set("Authorization", validBearerToken(t))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.NotNil(t, resp)
	assert.Len(t, resp, 0)
}

func TestGetWorkoutByID_NotFound_Returns404(t *testing.T) {
	r := newTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/workouts/9999", nil)
	req.Header.Set("Authorization", validBearerToken(t))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestLoggerMiddleware_OutputsAllKeyValueFields(t *testing.T) {
	origStdout := os.Stdout
	pr, pw, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = pw

	r := gin.New()
	r.Use(middleware.Logger())
	r.GET("/ping", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	pw.Close()
	os.Stdout = origStdout

	var buf bytes.Buffer
	_, err = buf.ReadFrom(pr)
	require.NoError(t, err)
	pr.Close()

	line := buf.String()
	for _, field := range []string{"time=", "method=", "path=", "status=", "latency=", "client="} {
		assert.Contains(t, line, field, "expected field %q in logger output", field)
	}
}

// --- JWT middleware tests ---

func TestJWTMiddleware_ValidToken_PassesThrough(t *testing.T) {
	r := newTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/workouts", nil)
	req.Header.Set("Authorization", validBearerToken(t))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJWTMiddleware_InvalidToken_Returns401(t *testing.T) {
	r := newTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/workouts", nil)
	req.Header.Set("Authorization", "Bearer this.is.not.valid")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.NotEmpty(t, resp["error"])
}
