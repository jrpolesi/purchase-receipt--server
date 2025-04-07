package product_usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jrpolesi/purchase-receipt--server/internal/modules/product/domain"
	gatewayMocks "github.com/jrpolesi/purchase-receipt--server/internal/modules/product/domain/gateway/mocks"
	productUsecase "github.com/jrpolesi/purchase-receipt--server/internal/modules/product/usecase"
	testUtils "github.com/jrpolesi/purchase-receipt--server/internal/testutils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_CreateProduct_Execute_Success(t *testing.T) {
	var categoryID = uuid.New()

	type mocks struct {
		productRepository  *gatewayMocks.MockProduct
		categoryRepository *gatewayMocks.MockCategory
	}

	tests := map[string]struct {
		input     productUsecase.CreateProductInput
		expected  productUsecase.CreateProductOutput
		mockSetup func(ctx context.Context, mocks *mocks)
	}{
		"Should create a product successfully": {
			input: productUsecase.CreateProductInput{
				Code:            "ABC123",
				Description:     "Test Product",
				CategoryID:      categoryID.String(),
				Dimensions:      "10x10x10",
				UnitPrice:       10.0,
				UnitsPerPackage: 10,
			},
			expected: productUsecase.CreateProductOutput{
				Code:        "ABC123",
				Description: "Test Product",
				Category: domain.Category{
					ID:   categoryID,
					Name: "Test Category",
				},
				Dimensions:      "10x10x10",
				UnitPrice:       10.0,
				UnitsPerPackage: 10,
			},
			mockSetup: func(ctx context.Context, mocks *mocks) {
				category := &domain.Category{
					ID:   categoryID,
					Name: "Test Category",
				}

				mocks.categoryRepository.
					EXPECT().
					FindByID(categoryID).
					Return(category, nil).
					Times(1)

				mocks.productRepository.
					EXPECT().
					ExistsByCode("ABC123").
					Return(false, nil).
					Times(1)

				mocks.productRepository.
					EXPECT().
					Create(testUtils.MatchedByFn(func(value *domain.Product) bool {
						return assert.NotEmpty(t, value.ID) &&
							assert.Equal(t, "ABC123", value.Code) &&
							assert.Equal(t, "Test Product", value.Description) &&
							assert.Equal(t, *category, value.Category) &&
							assert.Equal(t, "10x10x10", value.Dimensions) &&
							assert.Equal(t, 10.0, value.UnitPrice) &&
							assert.Equal(t, 10, value.UnitsPerPackage)
					})).
					Return("", nil).
					Times(1)
			},
		},
		"Should create a product with 1 unit per package when it is not provided": {
			input: productUsecase.CreateProductInput{
				Code:            "ABC123",
				Description:     "Test Product",
				CategoryID:      categoryID.String(),
				Dimensions:      "10x10x10",
				UnitPrice:       10.99,
				UnitsPerPackage: 0,
			},
			expected: productUsecase.CreateProductOutput{
				Code:        "ABC123",
				Description: "Test Product",
				Category: domain.Category{
					ID:   categoryID,
					Name: "Test Category",
				},
				Dimensions:      "10x10x10",
				UnitPrice:       10.99,
				UnitsPerPackage: 1,
			},
			mockSetup: func(ctx context.Context, mocks *mocks) {
				category := &domain.Category{
					ID:   categoryID,
					Name: "Test Category",
				}

				mocks.categoryRepository.
					EXPECT().
					FindByID(categoryID).
					Return(category, nil).
					Times(1)

				mocks.productRepository.
					EXPECT().
					ExistsByCode("ABC123").
					Return(false, nil).
					Times(1)

				mocks.productRepository.
					EXPECT().
					Create(testUtils.MatchedByFn(func(value *domain.Product) bool {
						return assert.NotEmpty(t, value.ID) &&
							assert.Equal(t, "ABC123", value.Code) &&
							assert.Equal(t, "Test Product", value.Description) &&
							assert.Equal(t, *category, value.Category) &&
							assert.Equal(t, "10x10x10", value.Dimensions) &&
							assert.Equal(t, 10.99, value.UnitPrice) &&
							assert.Equal(t, 1, value.UnitsPerPackage)
					})).
					Return("", nil).
					Times(1)
			},
		},
		"Should create a product with 0 unit price": {
			input: productUsecase.CreateProductInput{
				Code:            "ABC123",
				Description:     "Test Product",
				CategoryID:      categoryID.String(),
				Dimensions:      "10x10x10",
				UnitPrice:       0,
				UnitsPerPackage: 15,
			},
			expected: productUsecase.CreateProductOutput{
				Code:        "ABC123",
				Description: "Test Product",
				Category: domain.Category{
					ID:   categoryID,
					Name: "Test Category",
				},
				Dimensions:      "10x10x10",
				UnitPrice:       0,
				UnitsPerPackage: 15,
			},
			mockSetup: func(ctx context.Context, mocks *mocks) {
				category := &domain.Category{
					ID:   categoryID,
					Name: "Test Category",
				}

				mocks.categoryRepository.
					EXPECT().
					FindByID(categoryID).
					Return(category, nil).
					Times(1)

				mocks.productRepository.
					EXPECT().
					ExistsByCode("ABC123").
					Return(false, nil).
					Times(1)

				mocks.productRepository.
					EXPECT().
					Create(testUtils.MatchedByFn(func(value *domain.Product) bool {
						return assert.NotEmpty(t, value.ID) &&
							assert.Equal(t, "ABC123", value.Code) &&
							assert.Equal(t, "Test Product", value.Description) &&
							assert.Equal(t, *category, value.Category) &&
							assert.Equal(t, "10x10x10", value.Dimensions) &&
							assert.Equal(t, 0.0, value.UnitPrice) &&
							assert.Equal(t, 15, value.UnitsPerPackage)
					})).
					Return("", nil).
					Times(1)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Arrange
			ctx := context.Background()
			productRepoMock := gatewayMocks.NewMockProduct(ctrl)
			categoryRepoMock := gatewayMocks.NewMockCategory(ctrl)

			test.mockSetup(ctx, &mocks{
				productRepository:  productRepoMock,
				categoryRepository: categoryRepoMock,
			})

			// Act
			usecase := productUsecase.NewCreateProduct(productRepoMock, categoryRepoMock)
			output, err := usecase.Execute(test.input)

			// Assert
			assert.Nil(t, err)
			assert.NotEmpty(t, output.ID)

			_, err = uuid.Parse(output.ID)
			assert.NoError(t, err, "Expected a valid UUID")

			expectedOutput := test.expected
			expectedOutput.ID = output.ID
			assert.Equal(t, expectedOutput, output)
		})
	}
}

func Test_CreateProduct_Execute_Failure(t *testing.T) {
	var categoryID = uuid.New()

	type mocks struct {
		productRepository  *gatewayMocks.MockProduct
		categoryRepository *gatewayMocks.MockCategory
	}

	tests := map[string]struct {
		input       productUsecase.CreateProductInput
		expectedErr error
		mockSetup   func(ctx context.Context, mocks *mocks)
	}{
		"Should return error when code is empty": {
			input: productUsecase.CreateProductInput{
				Description:     "Test Product",
				CategoryID:      categoryID.String(),
				Dimensions:      "10x10x10",
				UnitPrice:       10.0,
				UnitsPerPackage: 10,
			},
			expectedErr: domain.ErrEmptyProductCode,
			mockSetup: func(ctx context.Context, mocks *mocks) {
				mocks.categoryRepository.
					EXPECT().
					FindByID(categoryID).
					Return(&domain.Category{
						ID:   categoryID,
						Name: "Test Category",
					}, nil).
					Times(1)
			},
		},
		"Should return error when description is empty": {
			input: productUsecase.CreateProductInput{
				Code:            "ABC123",
				CategoryID:      categoryID.String(),
				Dimensions:      "10x10x10",
				UnitPrice:       10.0,
				UnitsPerPackage: 10,
			},
			expectedErr: domain.ErrEmptyProductDescription,
			mockSetup: func(ctx context.Context, mocks *mocks) {
				mocks.categoryRepository.
					EXPECT().
					FindByID(categoryID).
					Return(&domain.Category{
						ID:   categoryID,
						Name: "Test Category",
					}, nil).
					Times(1)
			},
		},
		"Should return error when category is empty": {
			input: productUsecase.CreateProductInput{
				Code:            "ABC123",
				Description:     "Test Product",
				Dimensions:      "10x10x10",
				UnitPrice:       10.0,
				UnitsPerPackage: 10,
			},
			expectedErr: productUsecase.ErrEmptyOrInvalidCategoryID,
			mockSetup:   func(ctx context.Context, mocks *mocks) {},
		},
		"Should return error when dimensions is empty": {
			input: productUsecase.CreateProductInput{
				Code:            "ABC123",
				Description:     "Test Product",
				CategoryID:      categoryID.String(),
				UnitPrice:       10.0,
				UnitsPerPackage: 10,
			},
			expectedErr: domain.ErrEmptyProductDimensions,
			mockSetup: func(ctx context.Context, mocks *mocks) {
				mocks.categoryRepository.
					EXPECT().
					FindByID(categoryID).
					Return(&domain.Category{
						ID:   categoryID,
						Name: "Test Category",
					}, nil).
					Times(1)
			},
		},
		"Should return error when units per package is lower than 0": {
			input: productUsecase.CreateProductInput{
				Code:            "ABC123",
				Description:     "Test Product",
				CategoryID:      categoryID.String(),
				Dimensions:      "10x10x10",
				UnitPrice:       10.0,
				UnitsPerPackage: -1,
			},
			expectedErr: domain.ErrLowerThanZeroProductUnitsPerPackage,
			mockSetup: func(ctx context.Context, mocks *mocks) {
				mocks.categoryRepository.
					EXPECT().
					FindByID(categoryID).
					Return(&domain.Category{
						ID:   categoryID,
						Name: "Test Category",
					}, nil).
					Times(1)
			},
		},
		"Should return error when unit price is lower than 0": {
			input: productUsecase.CreateProductInput{
				Code:            "ABC123",
				Description:     "Test Product",
				CategoryID:      categoryID.String(),
				Dimensions:      "10x10x10",
				UnitPrice:       -1,
				UnitsPerPackage: 10,
			},
			expectedErr: domain.ErrLowerThanZeroProductUnitPrice,
			mockSetup: func(ctx context.Context, mocks *mocks) {
				mocks.categoryRepository.
					EXPECT().
					FindByID(categoryID).
					Return(&domain.Category{
						ID:   categoryID,
						Name: "Test Category",
					}, nil).
					Times(1)
			},
		},
		"Should return error when product already exists": {
			input: productUsecase.CreateProductInput{
				Code:            "ABC123",
				Description:     "Test Product",
				CategoryID:      categoryID.String(),
				Dimensions:      "10x10x10",
				UnitPrice:       10.0,
				UnitsPerPackage: 10,
			},
			expectedErr: productUsecase.ErrProductAlreadyExists,
			mockSetup: func(ctx context.Context, mocks *mocks) {
				mocks.categoryRepository.
					EXPECT().
					FindByID(categoryID).
					Return(&domain.Category{
						ID:   categoryID,
						Name: "Test Category",
					}, nil).
					Times(1)

				mocks.productRepository.
					EXPECT().
					ExistsByCode("ABC123").
					Return(true, nil).
					Times(1)
			},
		},
		"Should return error when categoryRepository.FindByID fails": {
			input: productUsecase.CreateProductInput{
				Code:            "ABC123",
				Description:     "Test Product",
				CategoryID:      categoryID.String(),
				Dimensions:      "10x10x10",
				UnitPrice:       10.0,
				UnitsPerPackage: 10,
			},
			expectedErr: errors.New("any error"),
			mockSetup: func(ctx context.Context, mocks *mocks) {
				mocks.categoryRepository.
					EXPECT().
					FindByID(categoryID).
					Return(nil, errors.New("any error")).
					Times(1)
			},
		},
		"Should return error when productRepository.ExistsByCode fails": {
			input: productUsecase.CreateProductInput{
				Code:            "ABC123",
				Description:     "Test Product",
				CategoryID:      categoryID.String(),
				Dimensions:      "10x10x10",
				UnitPrice:       10.0,
				UnitsPerPackage: 10,
			},
			expectedErr: errors.New("product repository error"),
			mockSetup: func(ctx context.Context, mocks *mocks) {
				mocks.categoryRepository.
					EXPECT().
					FindByID(categoryID).
					Return(&domain.Category{
						ID:   categoryID,
						Name: "Test Category",
					}, nil).
					Times(1)

				mocks.productRepository.
					EXPECT().
					ExistsByCode("ABC123").
					Return(false, errors.New("product repository error")).
					Times(1)
			},
		},
		"Should return error when productRepository.Create fails": {
			input: productUsecase.CreateProductInput{
				Code:            "ABC123",
				Description:     "Test Product",
				CategoryID:      categoryID.String(),
				Dimensions:      "10x10x10",
				UnitPrice:       10.0,
				UnitsPerPackage: 10,
			},
			expectedErr: errors.New("any product repository error"),
			mockSetup: func(ctx context.Context, mocks *mocks) {
				category := &domain.Category{
					ID:   categoryID,
					Name: "Test Category",
				}

				mocks.categoryRepository.
					EXPECT().
					FindByID(categoryID).
					Return(category, nil).
					Times(1)

				mocks.productRepository.
					EXPECT().
					ExistsByCode("ABC123").
					Return(false, nil).
					Times(1)

				mocks.productRepository.
					EXPECT().
					Create(testUtils.MatchedByFn(func(value *domain.Product) bool {
						return assert.NotEmpty(t, value.ID) &&
							assert.Equal(t, "ABC123", value.Code) &&
							assert.Equal(t, "Test Product", value.Description) &&
							assert.Equal(t, *category, value.Category) &&
							assert.Equal(t, "10x10x10", value.Dimensions) &&
							assert.Equal(t, 10.0, value.UnitPrice) &&
							assert.Equal(t, 10, value.UnitsPerPackage)
					})).
					Return("", errors.New("any product repository error")).
					Times(1)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Arrange
			ctx := context.Background()
			productRepoMock := gatewayMocks.NewMockProduct(ctrl)
			categoryRepoMock := gatewayMocks.NewMockCategory(ctrl)

			test.mockSetup(ctx, &mocks{
				productRepository:  productRepoMock,
				categoryRepository: categoryRepoMock,
			})

			// Act
			createProductUsecase := productUsecase.NewCreateProduct(
				productRepoMock,
				categoryRepoMock,
			)
			_, err := createProductUsecase.Execute(test.input)

			// Assert
			assert.NotNil(t, err)
			assert.EqualError(t, err, test.expectedErr.Error())
		})
	}
}
