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

	tx, err := db.DB.Begin()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	queryProduct := `
		INSERT INTO products (id, name, price)
		VALUES (?, ?, ?)
	`

	_, err = tx.Exec(queryProduct, id, name, price)
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), 500)
		return
	}
	
	queryDetail := `
		INSERT INTO product_details (product_id, stock, is_active, created_at, image)
		VALUES (?, ?, ?, NOW(), ?)
	`

	_, err = tx.Exec(queryDetail, id, stock, true, imgBytes)
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), 500)
		return
	}

	err = tx.Commit()
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

	tx, err := db.DB.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// update tabel products
	queryProduct := `
	UPDATE products 
	SET name=?, price=? 
	WHERE id=?`

	_, err = tx.Exec(queryProduct, name, price, id)
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if imgBytes != nil {
		queryDetail := `
		UPDATE product_details 
		SET stock=?, is_active=?, image=? 
		WHERE product_id=?`

		_, err = tx.Exec(queryDetail, stock, isActive, imgBytes, id)
	} else {
		queryDetail := `
		UPDATE product_details 
		SET stock=?, is_active=? 
		WHERE product_id=?`

		_, err = tx.Exec(queryDetail, stock, isActive, id)
	}

	if err != nil {
	tx.Rollback()
	http.Error(w, err.Error(), http.StatusInternalServerError)
	return
}

	err = tx.Commit()
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