package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggingMiddleware_Basic(t *testing.T) {
	called := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusTeapot)
		w.Write([]byte("hello"))
	})

	mw := NewLoggingMiddleware()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	mw(handler).ServeHTTP(w, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusTeapot, w.Code)
	assert.Equal(t, "hello", w.Body.String())
}

func TestLoggingMiddleware_StatusAndSize(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("abc"))
	})

	mw := NewLoggingMiddleware()
	req := httptest.NewRequest(http.MethodGet, "/size", nil)
	w := httptest.NewRecorder()

	mw(handler).ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "abc", w.Body.String())
}
