package product_usecase

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jrpolesi/purchase-receipt--server/internal/modules/product/domain"
	"github.com/jrpolesi/purchase-receipt--server/internal/modules/product/domain/gateway"
)

var (
	ErrEmptyOrInvalidCategoryID = errors.New("empty or invalid category ID")
	ErrProductAlreadyExists		 = errors.New("product already exists")
)

type createProduct struct {
	productRepository  gateway.Product
	categoryRepository gateway.Category
}

func NewCreateProduct(
	productRepository gateway.Product,
	categoryRepository gateway.Category,
) *createProduct {
	return &createProduct{
		productRepository:  productRepository,
		categoryRepository: categoryRepository,
	}
}

func (c *createProduct) Execute(input CreateProductInput) (productOutput CreateProductOutput, err error) {
	categoryID, err := uuid.Parse(input.CategoryID)
	if err != nil {
		return CreateProductOutput{}, ErrEmptyOrInvalidCategoryID
	}

	category, err := c.categoryRepository.FindByID(categoryID)
	if err != nil {
		return CreateProductOutput{}, err
	}

	newProduct := domain.NewProduct(
		uuid.Nil,
		input.Code,
		input.Description,
		category,
		input.Dimensions,
		input.UnitsPerPackage,
		input.UnitPrice,
	)

	isValid, err := newProduct.IsValid()
	if !isValid {
		return CreateProductOutput{}, err
	}

	productExist, err := c.productRepository.ExistsByCode(input.Code)
	if err != nil {
		return CreateProductOutput{}, err
	}
	if productExist {
		return CreateProductOutput{}, ErrProductAlreadyExists
	}

	_, err = c.productRepository.Create(newProduct)
	if err != nil {
		return CreateProductOutput{}, err
	}

	return CreateProductOutput{
		ID:              newProduct.ID.String(),
		Code:            newProduct.Code,
		Description:     newProduct.Description,
		Category:        newProduct.Category,
		Dimensions:      newProduct.Dimensions,
		UnitsPerPackage: newProduct.UnitsPerPackage,
		UnitPrice:       newProduct.UnitPrice,
	}, err
}
