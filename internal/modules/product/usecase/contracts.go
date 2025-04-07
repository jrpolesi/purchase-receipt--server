package product_usecase

import "github.com/jrpolesi/purchase-receipt--server/internal/modules/product/domain"

type CreateProductInput struct {
	Code            string
	CategoryID      string
	Description     string
	Dimensions      string
	UnitsPerPackage int
	UnitPrice       float64
}

type CreateProductOutput struct {
	ID              string
	Code            string
	Category        domain.Category
	Description     string
	Dimensions      string
	UnitsPerPackage int
	UnitPrice       float64
}
