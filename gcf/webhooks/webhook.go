package webhooks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// AuthEvent is the payload of a Firestore Auth event.
type AuthEvent struct {
	Email    string `json:"email"`
	Metadata struct {
		CreatedAt time.Time `json:"createdAt"`
	} `json:"metadata"`
	UID string `json:"uid"`
}

var (
	token    = os.Getenv("WEBHOOK_TOKEN")
	endpoint = os.Getenv("WEBHOOK_ENDPOINT")
)

func CreateUserWebhook(ctx context.Context, e AuthEvent) error {
	payload, err := json.Marshal(map[string]interface{}{
		"firebaseId": e.UID,
		"email":      e.Email,
		"createdAt":  e.Metadata.CreatedAt,
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("X-Webhook-Token", token)
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("call was not successful, status: %s", resp.Status)
	}
	return nil
}
