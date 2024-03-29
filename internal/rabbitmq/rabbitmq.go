package rabbitmq

import (
	"context"
	"fmt"
	"time"

	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const productsImportStream = "product-imported"

type RabbitMqStreamClient struct {
	env      *stream.Environment
	producer *stream.Producer
}

func NewRabbitMqStreamClient(uri string) (*RabbitMqStreamClient, error) {
	env, err := stream.NewEnvironment(
		stream.NewEnvironmentOptions().SetUri(uri))

	if err != nil {
		return nil, fmt.Errorf("failed to create stream environment: %w", err)
	}
	err = env.DeclareStream(productsImportStream,
		stream.NewStreamOptions().
			SetMaxLengthBytes(stream.ByteCapacity{}.GB(2)).SetMaxAge(time.Hour*24))
	if err != nil {
		return nil, fmt.Errorf("failed to declare stream: %w", err)
	}

	producer, err := env.NewProducer(productsImportStream, stream.NewProducerOptions())
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %s", err)
	}
	return &RabbitMqStreamClient{env: env, producer: producer}, nil
}

func (c *RabbitMqStreamClient) Close() {
	if err := c.producer.Close(); err != nil {
		log.WithError(err).Fatalf("failed to close producer")
	}
	if err := c.env.Close(); err != nil {
		log.WithError(err).Fatalf("failed to close environment")
	}
}

func (c *RabbitMqStreamClient) Publish(ctx context.Context, msg *amqp.AMQP10) error {
	span := c.startRabbitMqSpan(ctx, c.producer.GetName(), msg)
	defer span.End()
	err := c.producer.Send(msg)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}

// xD

func (publisher RabbitMqStreamClient) startRabbitMqSpan(ctx context.Context, queue string, msg *amqp.AMQP10) trace.Span {
	carrier := newRabbitMqCarrier(msg)
	propagator := otel.GetTextMapPropagator()
	ctx = propagator.Extract(ctx, carrier)
	tracerPrivider := otel.GetTracerProvider()
	ctx, span := tracerPrivider.Tracer("rabbitmq").Start(ctx, "publish")
	span.SetAttributes(
		attribute.KeyValue{Key: "messaging.destination", Value: attribute.StringValue(queue)},
		attribute.KeyValue{Key: "messaging.system", Value: attribute.StringValue("rabbitmq")},
		attribute.KeyValue{Key: "messaging.destination_kind", Value: attribute.StringValue("stream")},
		attribute.KeyValue{Key: "messaging.protocol", Value: attribute.StringValue("AMQP")},
		attribute.KeyValue{Key: "messaging.protocol_version", Value: attribute.StringValue("1.0.0")},
	)
	propagator.Inject(ctx, carrier)
	return span
}

type rabbitMqCarrier struct {
	msg *amqp.AMQP10
}

func newRabbitMqCarrier(msg *amqp.AMQP10) rabbitMqCarrier {
	setHeaderIfEmpty(msg)
	return rabbitMqCarrier{msg: msg}
}

func setHeaderIfEmpty(msg *amqp.AMQP10) {
	if msg.Annotations == nil {
		msg.Annotations = make(map[interface{}]interface{})
	}
}

// Get retrieves a single value for a given key.
func (c rabbitMqCarrier) Get(key string) string {
	if val, ok := c.msg.Annotations[key]; ok {
		return fmt.Sprint(val)
	}
	return ""
}

// Set sets a header.
func (c rabbitMqCarrier) Set(key, val string) {

	delete(c.msg.Annotations, key)
	c.msg.Annotations[key] = val
}

func (c rabbitMqCarrier) Keys() []string {
	keys := make([]string, 0)
	for k := range c.msg.Annotations {
		keys = append(keys, fmt.Sprint(k))
	}
	return keys
}
