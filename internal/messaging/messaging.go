package messaging

import (
	"context"
	"fmt"
	"log"

	"github.com/finchatapp/finchat-api/internal/appconfig"
	"github.com/finchatapp/finchat-api/internal/model"
	pubnub "github.com/pubnub/go"
)

type Messager interface {
	User(id int) Messager
	Register(ctx context.Context, firstName, lastName, email string) error
	SendMessage(ctx context.Context, threadID, senderID int, message string) (int64, error)
}

type MessageSaver interface {
	SaveMessage(ctx context.Context, msg *model.Message) error
}

type Client struct {
	Server *pubnub.PubNub
	pn     *pubnub.PubNub
	ms     MessageSaver
	pubKey string
	subKey string
}

func UUID(id int) string {
	return fmt.Sprintf("chat-user-%d", id)
}

func Channel(id int) string {
	return fmt.Sprintf("channel-%d", id)
}

func New(c appconfig.Pubnub, ms MessageSaver) *Client {
	config := pubnub.NewConfig()
	config.PublishKey = c.PubKey
	config.SubscribeKey = c.SubKey
	config.UUID = c.ServerUUID
	pn := pubnub.NewPubNub(config)
	return &Client{pn, pn, ms, c.PubKey, c.SubKey}
}

func (c *Client) User(id int) Messager {
	config := pubnub.NewConfig()
	config.PublishKey = c.pubKey
	config.SubscribeKey = c.subKey
	config.UUID = UUID(id)
	pn := pubnub.NewPubNub(config)
	return &Client{c.pn, pn, c.ms, c.pubKey, c.subKey}
}

func (c *Client) Register(ctx context.Context, firstName, lastName, email string) error {
	_, _, err := c.pn.SetUUIDMetadataWithContext(ctx).Name(firstName + " " + lastName).Email(email).ProfileURL("/").Execute()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) SendMessage(ctx context.Context, threadID, senderID int, message string) (int64, error) {
	resp, _, err := c.pn.PublishWithContext(ctx).Channel(Channel(threadID)).Message(message).Execute()
	if err != nil {
		return 0, err
	}
	go func() {
		if err := c.ms.SaveMessage(ctx, &model.Message{
			ThreadID:  threadID,
			SenderID:  senderID,
			Type:      "TEXT",
			Message:   message,
			Timestamp: resp.Timestamp,
		}); err != nil {
			log.Printf("[ERROR] failed to save message: %v", err)
		}
	}()
	return resp.Timestamp, nil
}
