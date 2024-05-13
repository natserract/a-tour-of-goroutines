package category

import (
	"time"

	"github.com/gofrs/uuid"
)

// Category is an entity that represents a category of product
type Category struct {
	ID   uuid.UUID
	Name string

	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
