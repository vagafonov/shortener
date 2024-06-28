package contract

import "context"

// ServiceHealthCheck abstract interface for HealthCheck.
type ServiceHealthCheck interface {
	Ping(ctx context.Context) error
}
