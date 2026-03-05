package handlers

import (
	"html/template"
	"io"
	"net/http"
	"net/url"
	"pert3_npm/db"
	"strconv"
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

func CreateProductHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := r.ParseMultipartForm(1 << 20)
	if err != nil {
		http.Redirect(w, r, "/?error=File is to large. Max 1MB", http.StatusSeeOther)
		return
	}

	r.ParseForm()

	id := r.FormValue("id")
	name := r.FormValue("name")
	// Mengecek apakah field price merupakan angka
	price, err := strconv.ParseFloat(r.FormValue("price"), 10)
	if err != nil {
		errMsg := url.QueryEscape("Price must be a number")
		http.Redirect(w, r, "/?error=" + errMsg, http.StatusSeeOther)
		return
	}
	// Mengecek apakah field stock merupakan angka
	stock, err := strconv.ParseInt(r.FormValue("stock"), 10, 64)
	if err != nil {
		errMsg := url.QueryEscape("Stock must be a number")
		http.Redirect(w, r, "/?error=" + errMsg, http.StatusSeeOther)
		return
	}

	// Mengecek apakah image ada
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Redirect(w, r, "/?error=Failed read file", http.StatusSeeOther)
		return
	}
	
	defer file.Close()

	// Membaca file gambar, apabila ada error mengembalikan response
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Redirect(w, r, "/?error=Failed read file", http.StatusInternalServerError)
		return
	}

	// Mengecek apakah ukuran gambar itu lebih kecil dari 1mb
	if len(fileBytes) > 1<<20 {
		http.Redirect(w, r, "?error=File size more than 1MB", http.StatusSeeOther)
		return
	}

	// Mengecek tipe file apakah jpeg atau png, jika bukan akan mengembalikan error
	fileType := http.DetectContentType(fileBytes)
	if fileType != "image/jpeg" && fileType != "image/png" {
		http.Redirect(w, r, "?error=JPEG or PNG type only", http.StatusSeeOther)
		return
	}

	query := `
		INSERT INTO products (id, name, price, stock, is_active, created_at, image)
		VALUES (?, ?, ?, ?, ?, NOW(), ?)
	`
	_, err = db.DB.Exec(query, id, name, price, stock, true, fileBytes)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func UpdateProductHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := r.ParseMultipartForm(1 << 20)
	if err != nil {
		http.Redirect(w, r, "/?error=File is to large. Max 1MB", http.StatusSeeOther)
		return
	}

	r.ParseForm()

	id := r.FormValue("id")
	name := r.FormValue("name")
	price, err := strconv.ParseInt(r.FormValue("price"), 10, 64)
	if err != nil {
		errMsg := url.QueryEscape("Price must be a number")
		http.Redirect(w, r, "/?error=" + errMsg, http.StatusSeeOther)
		return
	}
	stock, err := strconv.ParseInt(r.FormValue("stock"), 10, 64)
	if err != nil {
		errMsg := url.QueryEscape("Stock must be a number")
		http.Redirect(w, r, "/?error=" + errMsg, http.StatusSeeOther)
		return
	}

	isActive := r.FormValue("is_active") == "on"
	var imgFile []byte
	file, _, err := r.FormFile("image")

	if err == nil {
		defer file.Close()

		// Membaca file gambar, apabila ada error mengembalikan response
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			http.Redirect(w, r, "/?error=Failed read file", http.StatusInternalServerError)
			return
		}

		// Mengecek apakah ukuran gambar itu lebih kecil dari 1mb
		if len(fileBytes) > 1<<20 {
			http.Redirect(w, r, "?error=File size more than 1MB", http.StatusSeeOther)
			return
		}

		// Mengecek tipe file apakah jpeg atau png, jika bukan akan mengembalikan error
		fileType := http.DetectContentType(fileBytes)
		if fileType != "image/jpeg" && fileType != "image/png" {
			http.Redirect(w, r, "?error=JPEG or PNG type only", http.StatusSeeOther)
			return
		}

		imgFile = fileBytes
	}
	
	var query string

	if imgFile != nil {
		query = `
			UPDATE products 
			SET name=?, price=?, stock=?, is_active=?, image=? 
			WHERE id=?`
	} else {
		query = `
			UPDATE products 
			SET name=?, price=?, stock=?, is_active=? 
			WHERE id=?`
	}
	
	if imgFile != nil {
		_, err = db.DB.Exec(query, name, price, stock, isActive, imgFile, id)
	} else {
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