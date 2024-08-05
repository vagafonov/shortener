package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values := md.Get("userID")
		if len(values) == 0 {
			return nil, status.Error(codes.Unauthenticated, "userID is empty in metadata")
		}
	}

	// вызываем RPC-метод
	return handler(ctx, req)
}
