package contract

import "context"

type ServiceHealthCheck interface {
	Ping(ctx context.Context) error
}
