package models

type Product struct {
	Id        string
	Name      string
	Price     float64
	Stock     int
 IsActive  bool
 CreatedAt time.Time
 Image     []byte
}