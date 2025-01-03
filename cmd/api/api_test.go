package main

import (
	"net/http"
	"net/http/httptest"
	"p1/internal/ratelimiter"
	"testing"
	"time"
)

func TestRateLimiterMiddleware(t *testing.T) {
	config := config{
		rateLimiter: ratelimiter.Config{
			RequestPerTimeFrame: 20,
			TimeFrame:           time.Second * 5,
			Enabled:             true,
		},
		addr: "8080",
	}

	app := newTestApplication(t, config)
	ts := httptest.NewServer(app.mount())
	defer ts.Close()

	client := &http.Client{}
	mockIP := "192.168.1.1"
	marginOfError := 2

	for i := 0; i < config.rateLimiter.RequestPerTimeFrame+marginOfError; i++ {
		req, err := http.NewRequest(http.MethodGet, ts.URL+"/v1/health", nil)
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}

		req.Header.Set("X-Forwarded-For", mockIP)

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("could not send request: %v", err)
		}

		defer resp.Body.Close()

		if i < config.rateLimiter.RequestPerTimeFrame {
			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected status OK; got %v", resp.Status)
			}
		} else {
			if resp.StatusCode != http.StatusTooManyRequests {
				t.Errorf("expected status too many request; got %v", resp.Status)
			}
		}
	}
}
