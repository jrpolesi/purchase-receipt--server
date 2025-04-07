package domain

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrEmptyProductCode                    = errors.New("product code cannot be empty")
	ErrEmptyProductDescription             = errors.New("product description cannot be empty")
	ErrEmptyProductCategory                = errors.New("product category cannot be empty")
	ErrEmptyProductDimensions              = errors.New("product dimensions cannot be empty")
	ErrLowerThanZeroProductUnitsPerPackage = errors.New("product units per package cannot be lower than zero")
	ErrLowerThanZeroProductUnitPrice       = errors.New("product unit price cannot be lower than zero")
)

type Product struct {
	ID              uuid.UUID
	Code            string
	Category        Category
	Description     string
	Dimensions      string
	UnitsPerPackage int
	UnitPrice       float64
}

func NewProduct(
	id uuid.UUID,
	code string,
	description string,
	category *Category,
	dimensions string,
	unitsPerPackage int,
	unitPrice float64,
) *Product {
	if unitsPerPackage == 0 {
		unitsPerPackage = 1
	}

	if id == uuid.Nil {
		id = uuid.New()
	}

	return &Product{
		ID:              id,
		Code:            code,
		Category:        *category,
		Description:     description,
		Dimensions:      dimensions,
		UnitsPerPackage: unitsPerPackage,
		UnitPrice:       unitPrice,
	}
}

func (p *Product) IsValid() (isValid bool, err error) {
	var errorsList []error

	if p.Code == "" {
		errorsList = append(errorsList, ErrEmptyProductCode)
	}
	if p.Description == "" {
		errorsList = append(errorsList, ErrEmptyProductDescription)
	}
	if p.Category.ID == uuid.Nil {
		errorsList = append(errorsList, ErrEmptyProductCategory)
	}
	if p.Dimensions == "" {
		errorsList = append(errorsList, ErrEmptyProductDimensions)
	}
	if p.UnitsPerPackage <= 0 {
		errorsList = append(errorsList, ErrLowerThanZeroProductUnitsPerPackage)
	}
	if p.UnitPrice < 0 {
		errorsList = append(errorsList, ErrLowerThanZeroProductUnitPrice)
	}

	if len(errorsList) > 0 {
		return false, errors.Join(errorsList...)
	}

	return true, nil
}
