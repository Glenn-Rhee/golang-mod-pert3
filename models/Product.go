package models

import "time"

// Buat Struct Product
type Product struct {
	Id        string
	Name      string
	Price     float64
	Stock     int
	IsActive  bool
	CreatedAt time.Time
	Image     []byte
}