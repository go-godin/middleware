package middleware

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/openzipkin/zipkin-go"
)

func InstrumentZipkin() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			span := zipkin.SpanFromContext(ctx)
			span.Annotate(time.Now(), "endpoint.start")

			defer func() {
				if err := response.(endpoint.Failer).Failed(); err != nil {
					span.Tag("error", err.Error())
				}
				span.Annotate(time.Now(), "endpoint.end")
			}()

			return next(ctx, request)
		}
	}
}
