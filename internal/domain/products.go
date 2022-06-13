package domain

import "time"

type ProductsList struct {
	Id           string    `json:"id"`
	Title        string    `json:"title" binding:"required"`
	Image        string    `json:"image"`
	Price        uint      `json:"price" binding:"required"`
	Sale         uint      `json:"sale"`
	SaleOldPrice uint      `json:"sale_old_price"`
	CategoryId   string    `json:"category_id" binding:"required"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateProductInput struct {
	Title        string `json:"title" binding:"required"`
	Image        string `json:"image"`
	Price        uint   `json:"price" binding:"required"`
	Sale         uint   `json:"sale"`
	SaleOldPrice uint   `json:"sale_old_price"`
	CategoryId   string `json:"category_id" binding:"required"`
}

type UpdateProductInput struct {
	Title        *string `json:"title" binding:"required"`
	Image        *string `json:"image"`
	Price        *uint   `json:"price" binding:"required"`
	Sale         *uint   `json:"sale"`
	SaleOldPrice *uint   `json:"sale_old_price"`
	CategoryId   *int    `json:"category_id" binding:"required"`
}
