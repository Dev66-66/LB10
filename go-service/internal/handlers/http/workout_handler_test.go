package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	handlers "github.com/Dev66-66/LB10/go-service/internal/handlers/http"
	"github.com/Dev66-66/LB10/go-service/internal/middleware"
	"github.com/Dev66-66/LB10/go-service/internal/store"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newTestRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	h := handlers.NewWorkoutHandler(store.NewWorkoutStore())
	h.RegisterRoutes(r)
	return r
}

func TestPostWorkout_ValidBody_Returns201(t *testing.T) {
	r := newTestRouter()

	body := `{"name":"Push Day","type":"strength","duration":60,"difficulty":"hard","calories_burned":400}`
	req := httptest.NewRequest(http.MethodPost, "/workouts", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
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
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestLoggerMiddleware_OutputsAllKeyValueFields(t *testing.T) {
	// Redirect os.Stdout to a pipe to capture fmt.Printf output from the logger.
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

	// Close the write end and restore stdout before reading.
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
