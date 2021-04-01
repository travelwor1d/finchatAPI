package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPosts(t *testing.T) {
	title := "First Post"
	content := "..."
	reqBody := fmt.Sprintf(`{
		"title": "%s",
		"content": "%s"
	}`, title, content)
	req := httptest.NewRequest("POST", "/api/v1/posts", strings.NewReader(reqBody))
	req.Header.Add("Authorization", "Bearer "+goatAuthToken)
	req.Header.Add("Content-Type", "application/json")
	resp, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("unsuccessful post creation: %s", body)
	}
	var body struct {
		ID      int    `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.Title != title {
		t.Errorf("got %q, want %q", body.Title, title)
	}
	if body.Content != content {
		t.Errorf("got %q, want %q", body.Content, content)
	}
}
