package main

import (
	"net/http"
	"testing"
)

func TestGetUser(t *testing.T) {
	app := newTestApplication(t, config{})
	mux := app.mount()

	testoken, err := app.authenticator.GenerateToken(nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("should not allow unauthenticated request", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/users/41", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should allow authenticated request", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/users/41", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer "+testoken)

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusOK, rr.Code)
	})
}
