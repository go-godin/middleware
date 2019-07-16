package middleware

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"

	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

var requestDuration = prometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
	Name:    "endpoint_request_duration_ms",
	Help:    "Request duration in milliseconds",
	Buckets: []float64{50, 100, 250, 500, 1000},
}, []string{"method"})

var requestsCurrent = prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
	Name: "endpoint_requests_current",
	Help: "The current number of gRPC requests by endpoint",
}, []string{"method"})

var requestStatus = prometheus.NewCounterFrom(stdprometheus.CounterOpts{
	Name: "endpoint_requests_total",
	Help: "The total number of gRPC requests and whether the business failed or not",
}, []string{"method", "success"})

var amqpInbound = prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
	Name: "amqp_inbound",
	Help: "Increased on incoming message, decreased on ack/nack",
}, []string{"routing_key"})

var amqpTransportError = prometheus.NewCounterFrom(stdprometheus.CounterOpts{
	Name: "amqp_transport_error",
	Help: "Increased when a message could not be decoded or necessary content is missing",
}, []string{"routing_key"})

// Prometheus adds basic RED metrics on all endpoints. The transport layer (gRPC, AMQP, HTTP, ...) should also have metrics attached and
// will then take care of monitoring gRPC endpoints including their status.
func Prometheus(methodName string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			requestsCurrent.With("method", methodName).Add(1)

			defer func(begin time.Time) {
				requestDuration.With("method", methodName).Observe(time.Since(begin).Seconds())
				requestsCurrent.With("method", methodName).Add(-1)
				requestStatus.With("method", methodName, "success", fmt.Sprint(err == nil))
			}(time.Now())

			return next(ctx, request)
		}
	}
}

func InstrumentRabbitMQ(topic string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			amqpInbound.With("routing_key", topic).Add(1)

			defer func(begin time.Time) {
				amqpInbound.With("routing_key", topic).Add(-1)
				if err != nil {
					amqpTransportError.With("routing_key", topic).Add(1)
				}
			}(time.Now())

			return next(ctx, request)
		}
	}
}
