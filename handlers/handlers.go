package handlers

import (
	"html/template"
	"net/http"
	"pert3_npm/db"
	"pert3_npm/helper"
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

func CreateProductHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := r.ParseMultipartForm(1 << 20)
	if err != nil {
		helper.RedirectError(w, r, "File is too large. Max 1MB")
		return
	}

	id := r.FormValue("id")
	name := r.FormValue("name")

	price, err := helper.ParsePrice(r)
	if err != nil {
		helper.RedirectError(w, r, "Price must be a number")
		return
	}

	stock, err := helper.ParseStock(r)
	if err != nil {
		helper.RedirectError(w, r, "Stock must be a number")
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		helper.RedirectError(w, r, "Failed read file")
		return
	}

	imgBytes, err := helper.ReadAndValidateImage(file)
	if err != nil {
		helper.RedirectError(w, r, err.Error())
		return
	}

	query := `
	INSERT INTO products (id, name, price, stock, is_active, created_at, image)
	VALUES (?, ?, ?, ?, ?, NOW(), ?)
	`

	_, err = db.DB.Exec(query, id, name, price, stock, true, imgBytes)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func UpdateProductHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := r.ParseMultipartForm(1 << 20)
	if err != nil {
		helper.RedirectError(w, r, "File is too large. Max 1MB")
		return
	}

	id := r.FormValue("id")
	name := r.FormValue("name")

	price, err := helper.ParseStock(r)
	if err != nil {
		helper.RedirectError(w, r, "Price must be a number")
		return
	}

	stock, err := helper.ParseStock(r)
	if err != nil {
		helper.RedirectError(w, r, "Stock must be a number")
		return
	}

	isActive := r.FormValue("is_active") == "on"

	var imgBytes []byte

	file, _, err := r.FormFile("image")

	if err == nil {
		imgBytes, err = helper.ReadAndValidateImage(file)
		if err != nil {
			helper.RedirectError(w, r, err.Error())
			return
		}
	}

	var query string

	if imgBytes != nil {
		query = `
		UPDATE products 
		SET name=?, price=?, stock=?, is_active=?, image=? 
		WHERE id=?`
		_, err = db.DB.Exec(query, name, price, stock, isActive, imgBytes, id)

	} else {
		query = `
		UPDATE products 
		SET name=?, price=?, stock=?, is_active=? 
		WHERE id=?`
		_, err = db.DB.Exec(query, name, price, stock, isActive, id)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func DeleteProductHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	r.ParseForm()
	id := r.FormValue("id")
	_, err := db.DB.Exec("DELETE FROM products WHERE id=?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}