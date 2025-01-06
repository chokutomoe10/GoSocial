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

	// t.Run("should hit the cache first and if not exists it sets the user on the cache", func(t *testing.T) {
	// 	mockCacheStore := app.cacheStorage.Users.(*cache.MockUserStore)

	// 	mockCacheStore.On("Get", int64(41)).Return(nil, nil)
	// 	mockCacheStore.On("Get", int64(1)).Return(nil, nil)
	// 	mockCacheStore.On("Set", mock.Anything, mock.Anything).Return(nil)

	// 	req, err := http.NewRequest(http.MethodGet, "/v1/users/41", nil)
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}

	// 	req.Header.Set("Authorization", "Bearer "+testoken)

	// 	rr := executeRequest(req, mux)

	// 	checkResponseCode(t, http.StatusOK, rr.Code)

	// 	mockCacheStore.AssertNumberOfCalls(t, "Get", 1)

	// 	mockCacheStore.Calls = nil
	// })
}
