//go:generate mockgen  -package=gatewaymocks -source=product.go -destination=mocks/product.go
package gateway

import "github.com/jrpolesi/purchase-receipt--server/internal/modules/product/domain"

type Product interface {
	Create(product *domain.Product) (id string, err error)
	ExistsByCode(code string) (exists bool, err error)
}
