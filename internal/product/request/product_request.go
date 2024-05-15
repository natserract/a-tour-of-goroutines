package request

import (
	"goroutines/internal/product/errs"
	"strings"
)

type ProductCreateRequest struct {
	Name        string  `form:"name" binding:"required,min=1,max=30"`
	Sku         string  `form:"sku" binding:"required,min=1,max=30"`
	Category    string  `form:"category"`
	ImageUrl    string  `form:"imageUrl" binding:"required,url"`
	Notes       string  `form:"notes" binding:"required,min=1,max=200"`
	Price       float64 `form:"price" binding:"required,min=1"`
	Stock       *int    `form:"stock" binding:"required,min=0,max=100000"`
	Location    string  `form:"location" binding:"required"`
	IsAvailable bool    `form:"isAvailable" binding:"required"`
}

var ImageFormats = []string{".jpg", ".jpeg", ".png", ".webp"}

func (pr *ProductCreateRequest) ValidateProductCreate() error {
	var err error = nil

	// Check image url format
	if pr.ImageUrl != "" {
		for _, imageFormat := range ImageFormats {
			if strings.HasSuffix(pr.ImageUrl, imageFormat) {
				break
			}

			if imageFormat == ImageFormats[len(ImageFormats)-1] {
				err = errs.ProductErrsImageUrlInvalid
				break
			}
		}

	}

	return err
}
