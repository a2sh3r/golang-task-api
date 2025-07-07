package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGzipMiddleware_CompressesResponse(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	mw := NewGzipMiddleware()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()

	mw(handler).ServeHTTP(w, req)

	assert.Equal(t, "gzip", w.Header().Get("Content-Encoding"))

	gr, err := gzip.NewReader(w.Body)
	assert.NoError(t, err)
	defer gr.Close()
	unzipped, err := io.ReadAll(gr)
	assert.NoError(t, err)
	assert.Equal(t, "hello world", string(unzipped))
}

func TestGzipMiddleware_PassesThroughIfNoAcceptGzip(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("plain"))
	})

	mw := NewGzipMiddleware()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	mw(handler).ServeHTTP(w, req)

	assert.NotEqual(t, "gzip", w.Header().Get("Content-Encoding"))
	assert.Equal(t, "plain", w.Body.String())
}

func TestGzipMiddleware_DecodeGzippedRequest(t *testing.T) {
	var received string
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		received = string(body)
		w.WriteHeader(http.StatusOK)
	})

	mw := NewGzipMiddleware()
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	zw.Write([]byte("gzipped data"))
	zw.Close()

	req := httptest.NewRequest(http.MethodPost, "/", &buf)
	req.Header.Set("Content-Encoding", "gzip")
	w := httptest.NewRecorder()

	mw(handler).ServeHTTP(w, req)
	assert.Equal(t, "gzipped data", received)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGzipMiddleware_InvalidGzipRequest(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called on invalid gzip")
	})

	mw := NewGzipMiddleware()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("not gzipped"))
	req.Header.Set("Content-Encoding", "gzip")
	w := httptest.NewRecorder()

	mw(handler).ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
