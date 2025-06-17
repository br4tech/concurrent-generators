package ports

import (
	"context"

	"github.com/br4tech/concurrent-generators/internal/core/domain"
)

type OrderProcessor interface {
	Process(ctx context.Context, order domain.Order) error
}
