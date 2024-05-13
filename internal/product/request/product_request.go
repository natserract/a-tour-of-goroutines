package request

type ProductCreateRequest struct {
	Name        string  `form:"name" binding:"required,min=1,max=30"`
	Sku         string  `form:"sku" binding:"required,min=1,max=30"`
	Category    string  `form:"category"`
	ImageUrl    string  `form:"imageUrl" binding:"required,url"`
	Notes       string  `form:"notes" binding:"required,min=1,max=200"`
	Price       float64 `form:"price" binding:"required,min=1"`
	Stock       *int    `form:"stock" binding:"min=0,max=100000"`
	Location    string  `form:"location" binding:"required"`
	IsAvailable *bool   `form:"isAvailable" binding:"required"`
}
