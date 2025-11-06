package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Martins-Iroka/MyGallery-Backend/config"
	"github.com/Martins-Iroka/MyGallery-Backend/internal"
	"github.com/Martins-Iroka/MyGallery-Backend/internal/auth"
	"go.uber.org/zap"
)

func newTestApplication(t *testing.T, config config.Configuration) *application {
	t.Helper()

	logger := zap.NewNop().Sugar()

	store := internal.NewMockStorate()

	mockVerification := auth.MockOTPVerification{}
	mockAuth := auth.TestAuthenticator{}

	return &application{
		logger:          logger,
		config:          config,
		store:           store,
		otpVerification: mockVerification,
		auth:            mockAuth,
	}
}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d", expected, actual)
	}
}
