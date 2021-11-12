package controllertests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestSignIn(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()
	if err != nil {
		fmt.Printf("the error %v\n", err)
	}

	samples := []struct {
		email        string
		password     string
		errorMessage string
	}{
		{
			email:        user.Email,
			password:     "password", // Not the hashed from database
			errorMessage: "",
		},
		{
			email:        user.Email,
			password:     "Wrong Password",
			errorMessage: "crypto/bcrypt: hashedPassword is not the hash of the given password",
		},
		{
			email:        "Wrong Email",
			password:     "password",
			errorMessage: "record not found",
		},
	}

	for _, v := range samples {
		token, err := server.SignIn(v.email, v.password)
		if err != nil {
			assert.Equal(t, err, errors.New(v.errorMessage))
		} else {
			assert.NotEqual(t, token, "")
		}
	}
}

func TestLogin(t *testing.T) {
	refreshUserTable()

	_, err := seedOneUser()
	if err != nil {
		fmt.Printf("the error %v\n", err)
	}

	samples := []struct {
		inputJSON    string
		statusCode   int
		email        string
		password     string
		errorMessage string
	}{
		{
			inputJSON:    `{"email": "berg@gmail.com", "password": "password"}`,
			statusCode:   200,
			errorMessage: "",
		},
		{
			inputJSON:    `{"email": "berg@gmail.com", "password": "wrong password"}`,
			statusCode:   422,
			errorMessage: "incorrect password",
		},
		{
			inputJSON:    `{"email": "frank@gmail.com", "password": "password"}`,
			statusCode:   422,
			errorMessage: "incorrect details",
		},
		{
			inputJSON:    `{"email": "kangmail.com", "password": "password"}`,
			statusCode:   422,
			errorMessage: "invalid email",
		},
		{
			inputJSON:    `{"email": "", "password": "password"}`,
			statusCode:   422,
			errorMessage: "required email",
		},
		{
			inputJSON:    `{"email": "kan@gmail.com", "password": ""}`,
			statusCode:   422,
			errorMessage: "required password",
		},
		{
			inputJSON:    `{"email": "", "password": "password"}`,
			statusCode:   422,
			errorMessage: "required email",
		},
	}

	for i, v := range samples {
		req, err := http.NewRequest("POST", "/login", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("the error: %v", err)
		}

		rr := httptest.NewRecorder()
		fmt.Printf("%d. body: %s - rr: %d\n", i, v.inputJSON, rr.Code) // debug
		handler := http.HandlerFunc(server.Login)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.NotEqual(t, rr.Body.String(), "")
		}

		if v.statusCode == 422 && v.errorMessage != "" {
			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.Bytes()), &responseMap)
			if err != nil {
				t.Errorf("cannot convert to json: %v", err)
			}
			fmt.Println("133 ", (i), responseMap["error"]) // debug
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
