package tests

import (
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogin(t *testing.T) {
	if authToken == "" {
		t.Errorf("auth token is empty")
	}
}

func TestRegister(t *testing.T) {
	t.Run("register new user", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/auth/v1/register", strings.NewReader(`
	{
    "firstName": "Martin",
    "lastName": "Lukasik",
    "phone": "+48507968492",
    "email": "martilukas7@gmail.com",
    "password": "password321"
	}
	`))
		req.Header.Add("Content-Type", "application/json")
		resp, err := a.Test(req)
		if err != nil {
			t.Error(err)
		}
		if resp.StatusCode != 200 {
			body, _ := ioutil.ReadAll(resp.Body)
			t.Errorf("unsuccessful login: %s", body)
		}
	})

	t.Run("login a newly registed user", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/auth/v1/login", strings.NewReader(`
	{
		"email": "martilukas7@gmail.com",
		"password": "password321"
	}
	`))
		req.Header.Add("Content-Type", "application/json")
		resp, err := a.Test(req)
		if err != nil {
			t.Error(err)
		}
		if resp.StatusCode != 200 {
			body, _ := ioutil.ReadAll(resp.Body)
			t.Errorf("unsuccessful login: %s", body)
		}
	})
}
