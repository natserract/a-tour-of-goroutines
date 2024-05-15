package response

import (
	"goroutines/internal/product"
	"time"
)

type ProductShow struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Sku         string    `json:"sku"`
	Category    string    `json:"category"`
	Notes       string    `json:"notes"`
	ImageUrl    string    `json:"imageUrl"`
	Stock       int       `json:"stock"`
	Price       float64   `json:"price"`
	Location    string    `json:"location"`
	IsAvailable bool      `json:"isAvailable"`
	CreatedAt   time.Time `json:"createdAt"`
}

type ProductCreateResponse struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateProductResponse struct {
	Message string                `json:"message"`
	Data    ProductCreateResponse `json:"data"`
}

const ProductsCreateSuccMessage = "Successfully create products"

func ProductToCreateResponse(data *product.Product) *CreateProductResponse {
	return &CreateProductResponse{
		Message: ProductsCreateSuccMessage,
		Data: ProductCreateResponse{
			Id:        data.Id.String(),
			CreatedAt: data.CreatedAt,
		},
	}
}
