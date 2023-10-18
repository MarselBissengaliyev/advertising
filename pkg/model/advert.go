package model

import (
	"time"
)

type Advert struct {
	Id          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Photos      string    `json:"photos" db:"photos" binding:"required"`
	Price       float64   `json:"price" db:"price" binding:"required"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at" binding:"required"`
	MainPhoto   string    `json:"main_photo" db:"main_photo"`
}

type GetAllAdvertsResponse struct {
	Id        int     `json:"id"`
	Title     string  `json:"title"`
	MainPhoto string  `json:"main_photo" db:"main_photo"`
	Price     float64 `json:"price"`
}

type GetAdvertByIdResponse struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	MainPhoto   string    `json:"main_photo" db:"main_photo"`
	Price       float64   `json:"price"`
	Photos      *[]string `json:"photos"`
	Description *string   `json:"description"`
}

type CreateAdvertsBody struct {
	Title       string   `json:"title" db:"title" binding:"required"`
	Photos      []string `json:"photos" db:"photos" binding:"required"`
	Price       float64  `json:"price" db:"price" binding:"required"`
	Description string   `json:"description" db:"description"`
}
