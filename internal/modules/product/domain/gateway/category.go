//go:generate mockgen -package=gatewaymocks -source=category.go -destination=mocks/category.go
package gateway

import (
	"github.com/google/uuid"
	"github.com/jrpolesi/purchase-receipt--server/internal/modules/product/domain"
)

type Category interface {
	FindByID(id uuid.UUID) (*domain.Category, error)
}
