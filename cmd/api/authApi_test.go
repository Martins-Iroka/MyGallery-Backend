package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Martins-Iroka/MyGallery-Backend/config"
)

const path = "/v1/authentication"

func TestRegisterUserHandler(t *testing.T) {
	app := newTestApplication(t, config.Configuration{})
	mux := app.Mount()
	register := path + "/register"
	email := "test@example.com"
	t.Run("should return 400 as a result of wrong email format", func(t *testing.T) {
		payload := RegisterUserRequestPayload{
			Username: "testuser",
			Email:    "test",
			Password: "password123",
		}
		body, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, register, bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should return 400 as a result of unknown field added", func(t *testing.T) {
		payload := struct {
			Username string `json:"username" validate:"required,max=100"`
			Email    string `json:"email" validate:"required,email,max=255"`
			Password string `json:"password" validate:"required,min=5,max=72"`
			Unknown  string `json:"unknown"`
		}{
			Username: "testuser",
			Email:    "test",
			Password: "password123",
			Unknown:  "unknown",
		}
		body, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, register, bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should return 400 as a result of password being less than min", func(t *testing.T) {
		payload := RegisterUserRequestPayload{
			Username: "testuser",
			Email:    email,
			Password: "123",
		}
		body, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, register, bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should return 201", func(t *testing.T) {
		payload := RegisterUserRequestPayload{
			Username: "testuser",
			Email:    email,
			Password: "123456",
		}
		body, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, path+register, bytes.NewReader(body))
		if err != nil {
			t.Fatalf("error from new request %s", err)
		}

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusCreated, rr.Code)
	})
}

func TestVerifyUserHandler(t *testing.T) {
	app := newTestApplication(t, config.Configuration{})
	mux := app.Mount()

	verifyPath := path + "/verify"
	email := "test@g.com"

	t.Run("should return 400 because of an unknown field", func(t *testing.T) {
		payload := struct {
			Code    string `json:"code" validate:"required,len=6"`
			Email   string `json:"email" validate:"required,email,max=255"`
			Token   string `json:"token" validate:"required"`
			Unknown string `json:"unknown"`
		}{
			Code:    "123458",
			Email:   email,
			Token:   "token",
			Unknown: "unknown",
		}
		body, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, verifyPath, bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should return 400 because code length is less than 6", func(t *testing.T) {
		payload := struct {
			Code  string `json:"code" validate:"required,len=6"`
			Email string `json:"email" validate:"required,email,max=255"`
			Token string `json:"token" validate:"required"`
		}{
			Code:  "1234",
			Email: email,
			Token: "token",
		}
		body, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, verifyPath, bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should return 400 because code length is greater than 6", func(t *testing.T) {
		payload := VerifyUserRequestPayload{
			Code:  "12345678",
			Email: email,
			Token: "token",
		}
		body, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, verifyPath, bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should return 400 because of wrong email format", func(t *testing.T) {
		payload := VerifyUserRequestPayload{
			Code:  "123456",
			Email: "test",
			Token: "token",
		}
		body, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, verifyPath, bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should return 400 if token wasn't sent", func(t *testing.T) {
		payload := VerifyUserRequestPayload{
			Code:  "123456",
			Email: "test@example.com",
		}

		body, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, verifyPath, bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should return 200", func(t *testing.T) {
		payload := VerifyUserRequestPayload{
			Code:  "123456",
			Email: "t@example.com",
			Token: "token",
		}

		body, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, verifyPath, bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusOK, rr.Code)
	})
}

func TestLoginUserHandler(t *testing.T) {
	app := newTestApplication(t, config.Configuration{})
	mux := app.Mount()

	loginPath := path + "/login"

	t.Run("should return 400 because of an unknown field", func(t *testing.T) {
		payload := struct {
			Email    string `json:"email" validate:"required,email,max=255"`
			Password string `json:"password" validate:"required,min=5,max=72"`
			Unknown  string `json:"unknown"`
		}{
			Email:    "test@g.com",
			Password: "12345",
			Unknown:  "unknown",
		}
		body, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, loginPath, bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should return 400 because of wrong email format", func(t *testing.T) {
		payload := LoginUserRequestPayload{
			Email:    "test",
			Password: "12345",
		}
		body, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, loginPath, bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should return 400 because of password is less than 5", func(t *testing.T) {
		payload := LoginUserRequestPayload{
			Email:    "test@a.com",
			Password: "1234",
		}
		body, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, loginPath, bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should return 200", func(t *testing.T) {
		payload := LoginUserRequestPayload{
			Email:    "test@a.com",
			Password: "12345",
		}
		body, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, loginPath, bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusOK, rr.Code)
	})
}
