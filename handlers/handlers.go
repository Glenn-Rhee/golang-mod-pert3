package handlers

import (
	"html/template"
	"net/http"
)

var tmpl *template.Template

type ProductView struct {
	Id        string
	Name      string
	Price     float64
	Stock     int
	IsActive  bool
	CreatedAt string
}

type HomePageData struct {
	Products []ProductView
	Error string
}

type EditPageData struct {
	Product ProductView
	Error string
}

/* TODO - Buat handler untuk membuat produk */
func CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	
}

/* TODO - Buat handler untuk mengupdate produk */
func UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
	
}

/* TODO - Buat handler untuk menghapus produk */
func DeleteProductHandler(w http.ResponseWriter, r *http.Request){
	
}