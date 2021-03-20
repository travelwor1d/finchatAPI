package messaging

import (
	"context"
	"fmt"

	"github.com/finchatapp/finchat-api/internal/appconfig"
	pubnub "github.com/pubnub/go"
)

type Messager interface {
	User(id int) Messager
	Register(ctx context.Context, firstName, lastName, email string) error
}

type Client struct {
	Server *pubnub.PubNub
	pn     *pubnub.PubNub
	pubKey string
	subKey string
}

func UUID(id int) string {
	return fmt.Sprintf("chat-user-%d", id)
}

func New(c appconfig.Pubnub) *Client {
	config := pubnub.NewConfig()
	config.PublishKey = c.PubKey
	config.SubscribeKey = c.SubKey
	config.UUID = c.ServerUUID
	pn := pubnub.NewPubNub(config)
	return &Client{pn, pn, c.PubKey, c.SubKey}
}

func (c *Client) User(id int) Messager {
	config := pubnub.NewConfig()
	config.PublishKey = c.pubKey
	config.SubscribeKey = c.subKey
	config.UUID = UUID(id)
	pn := pubnub.NewPubNub(config)
	return &Client{c.pn, pn, c.pubKey, c.subKey}
}

func (c *Client) Register(ctx context.Context, firstName, lastName, email string) error {
	_, _, err := c.pn.SetUUIDMetadataWithContext(ctx).Name(firstName + " " + lastName).Email(email).ProfileURL("/").Execute()
	if err != nil {
		return err
	}
	return nil
}
