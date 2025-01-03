package main

import (
	"net/http"
	"net/http/httptest"
	"p1/internal/auth"
	"p1/internal/ratelimiter"
	"p1/internal/store"
	"p1/internal/store/cache"
	"testing"

	"go.uber.org/zap"
)

func newTestApplication(t *testing.T, config config) *application {
	t.Helper()

	// logger := zap.NewNop().Sugar()
	logger := zap.Must(zap.NewProduction()).Sugar()
	mockStore := store.NewMockStore()
	mockCacheStore := cache.NewMockStore()

	mockAuthenticator := &auth.MockAuthenticator{}

	rateLimiter := ratelimiter.NewFixedWindowLimiter(
		config.rateLimiter.RequestPerTimeFrame,
		config.rateLimiter.TimeFrame,
	)

	return &application{
		logger:        logger,
		store:         mockStore,
		config:        config,
		cacheStorage:  mockCacheStore,
		authenticator: mockAuthenticator,
		rateLimiter:   rateLimiter,
	}
}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		// if actual != expected {
		t.Errorf("expected the response code to be %d and we got %d", expected, actual)
	}
}
