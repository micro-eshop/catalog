package rabbitmq

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/micro-eshop/catalog/pkg/core/services"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func (client *RabbitMqStreamClient) PublishProductCreated(ctx context.Context, event services.ProductCreated) error {
	json, err := json.Marshal(event)
	if err != nil {
		return err
	}
	msg := amqp.NewMessage(json)
	return client.Publish(ctx, msg)
}
