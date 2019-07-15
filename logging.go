package middleware

import (
	"context"
	"time"

	"github.com/go-godin/log"
	"github.com/go-kit/kit/endpoint"
)

func Logging(logger log.Logger, endpointName string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				logger.Info("", "endpoint", endpointName, "took", time.Since(begin))
			}(time.Now())

			return next(ctx, request)
		}
	}
}
