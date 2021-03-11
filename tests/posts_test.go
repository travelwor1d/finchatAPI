package tests

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPosts(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/v1/posts", strings.NewReader(`
	{
		"title": "First Post",
		"content": "..."
	}
	`))
	req.Header.Add("Authorization", "Bearer "+authToken)
	req.Header.Add("Content-Type", "application/json")
	resp, err := a.Test(req)
	if err != nil {
		log.Fatal(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		log.Fatalf("unsuccessful post creation: %s", body)
	}
	fmt.Println(string(body))
}
