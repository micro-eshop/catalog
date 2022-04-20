package nats

import (
	"context"

	"github.com/micro-eshop/catalog/pkg/core/services"
	microeshop "github.com/micro-eshop/common-go"
)

const topic = "PRODUCTS.created"

type productsCreatedPublisher struct {
	messagePublisher microeshop.MessagePublisher[services.ProductCreated]
}

func NewPublisher(client *microeshop.NatsClient) (*productsCreatedPublisher, error) {
	publisher := microeshop.NewMessagePublisher[services.ProductCreated](client)
	return &productsCreatedPublisher{messagePublisher: publisher}, nil
}

func (c productsCreatedPublisher) Publish(ctx context.Context, event services.ProductCreated) error {
	headers := make(map[string]string)
	return c.messagePublisher.Publish(ctx, microeshop.NatsMessage[services.ProductCreated]{Data: event, MetaData: microeshop.NatsMessageMetaData{Topic: topic, Headers: headers}})
}
