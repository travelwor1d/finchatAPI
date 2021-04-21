package tests

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogin(t *testing.T) {
	if userAuthToken == "" {
		t.Errorf("auth token is empty")
	}
}

func TestRegister(t *testing.T) {
	req := httptest.NewRequest("POST", "/auth/v1/register", strings.NewReader(`
	{
    "firstName": "Martin",
    "lastName": "Lukasik",
    "phoneNumber": "+48507968492",
    "countryCode": "PL",
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
		body, _ := io.ReadAll(resp.Body)
		t.Errorf("unsuccessful login: %s", body)
	}

	_, err = login("martilukas7@gmail.com", "password321")
	if err != nil {
		t.Error(err)
	}
}
