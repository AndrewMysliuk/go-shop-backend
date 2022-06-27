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
	Type         string    `json:"type" binding:"required"`
	Subtype      string    `json:"subtype" binding:"required"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateProductInput struct {
	Title        string `json:"title" binding:"required"`
	Image        string `json:"image"`
	Price        uint   `json:"price" binding:"required"`
	Sale         uint   `json:"sale"`
	SaleOldPrice uint   `json:"sale_old_price"`
	CategoryId   string `json:"category_id" binding:"required"`
	Type         string `json:"type" binding:"required"`
	Subtype      string `json:"subtype" binding:"required"`
	Description  string `json:"description"`
}

type UpdateProductInput struct {
	Title        *string `json:"title" binding:"required"`
	Image        *string `json:"image"`
	Price        *uint   `json:"price" binding:"required"`
	Sale         *uint   `json:"sale"`
	SaleOldPrice *uint   `json:"sale_old_price"`
	CategoryId   *string `json:"category_id" binding:"required"`
	Type         *string `json:"type" binding:"required"`
	Subtype      *string `json:"subtype" binding:"required"`
	Description  *string `json:"description"`
}
