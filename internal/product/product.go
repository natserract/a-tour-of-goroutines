package product

import (
	"time"
)

type Product struct {
	Id          string  `json:"id"`
	Name        string  `json:"name" validate:"required,min=1,max=30"`
	Sku         string  `json:"sku" validate:"required,min=1,max=30"`
	Category    string  `json:"category" validate:"required"`
	ImageUrl    string  `json:"imageUrl" validate:"required,dive,url"`
	Notes       string  `json:"notes" validate:"required,min=1,max=200"`
	Price       float64 `json:"price" validate:"required,min=1"`
	Stock       int     `json:"stock" validate:"required,min=0,max=100000"`
	Location    string  `json:"location" validate:"required,min=1,max=200"`
	IsAvailable bool    `json:"isAvailable"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
