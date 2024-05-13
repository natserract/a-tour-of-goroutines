package product

import (
	"time"

	"github.com/gofrs/uuid"
)

type Product struct {
	Id          uuid.UUID
	Name        string
	Sku         string
	Category    string
	ImageUrl    string
	Notes       string
	Price       float64
	Stock       int
	Location    string
	IsAvailable bool

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
