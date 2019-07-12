package middleware

import (
	"context"
	"github.com/google/uuid"
	grpcMetadata "google.golang.org/grpc/metadata"

	metadata "github.com/go-godin/grpc-metadata"
	"github.com/go-kit/kit/endpoint"
)

func RequestIDMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			requestID := metadata.GetRequestID(ctx)

			// add the requestId to ctx if missing
			if requestID == "" {
				ctx = grpcMetadata.AppendToOutgoingContext(ctx, string(metadata.RequestID), uuid.New().String())
			}

			return next(ctx, request)
		}
	}
}
