package amqp

import (
	"context"
	"fmt"

	"github.com/go-godin/log"

	"github.com/go-godin/rabbitmq"
)

func Logging(logger log.Logger, routingKey string, handler rabbitmq.SubscriptionHandler) rabbitmq.SubscriptionHandler {
	return func(ctx context.Context, delivery *rabbitmq.Delivery) {
		logger.Info(
			"incoming AMQP message",
			"routing_key", routingKey,
			"redelivered", fmt.Sprint(delivery.Redelivered),
		)

		handler(ctx, delivery)
	}
}
