package messaging

import (
	"context"
	"encoding/json"

	"github.com/micro-eshop/catalog/internal/core/services"
	"github.com/nats-io/nats.go"
)

const topic = "PRODUCTS.created"

type natsClient struct {
	Connection *nats.Conn
}

func NewPublisher(url string) (*natsClient, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &natsClient{Connection: nc}, nil
}

func (c natsClient) Close() {
	c.Connection.Drain()
	c.Connection.Close()
}

func (c natsClient) Publish(ctx context.Context, event *services.ProductCreated) error {
	json, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return c.Connection.Publish(topic, json)
}
