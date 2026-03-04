package handlers

import (
	"database/sql"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"pert3_npm/db"
	"pert3_npm/models"
	"strconv"
)

var tmpl *template.Template

type HomePageData struct {
	Products []models.Product
	Error string
}

type EditPageData struct {
	Product models.Product
	Error string
}

func HomeHandler(w http.ResponseWriter, r *http.Request){
	tmpl, _ = template.ParseFiles("templates/index.html")
	errorMsg := r.URL.Query().Get("error")
	rows, err := db.DB.Query("SELECT id, name, price, stock FROM products")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()
	var products []models.Product
	for rows.Next(){
		var product models.Product
		err := rows.Scan(&product.Id, &product.Name, &product.Price, &product.Stock)
		if err !=nil {
			http.Error(w, err.Error(), 500)
			return
		}

		products = append(products, product)
	}

	data := HomePageData{
		Products: products,
		Error: errorMsg,
	}

	tmpl.Execute(w, data)
}

func CreateProductHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	err := r.ParseMultipartForm(1 << 20)
	if err != nil {
		http.Redirect(w, r, "/?error=File is to large. Max 1MB", http.StatusBadRequest)
		return
	}

	r.ParseForm()

	id := r.FormValue("id")
	name := r.FormValue("name")
	// Mengecek apakah field price merupakan angka
	price, err := strconv.ParseInt(r.FormValue("price"), 10, 64)
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
		http.Redirect(w, r, "/?error=Failed read file", http.StatusBadRequest)
		return
	}
	
	defer file.Close()

	// Membaca file gambar, apabila ada error mengembalikan response
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Redirect(w, r, "/?error=Failed read fil", http.StatusInternalServerError)
		return
	}

	// Mengecek apakah ukuran gambar itu lebih kecil dari 1mb
	if len(fileBytes) > 1<<20 {
		http.Redirect(w, r, "?error=File size more than 1MB", http.StatusBadRequest)
		return
	}

	// Mengecek tipe file apakah jpeg atau png, jika bukan akan mengembalikan error
	fileType := http.DetectContentType(fileBytes)
	if fileType != "image/jpeg" && fileType != "image/png" {
		http.Redirect(w, r, "?error=JPEG or PNG type only", http.StatusBadRequest)
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

func EditProductHandler(w http.ResponseWriter, r *http.Request){
	tmpl, _ = template.ParseFiles("templates/edit.html")
	id := r.URL.Query().Get("id")
	if id == ""{
		http.Redirect(w, r, "?error=Id is required!", http.StatusSeeOther)
	}
	query := "SELECT id, name, price, stock FROM products WHERE id = ?"
	var product models.Product
	err := db.DB.QueryRow(query, id).Scan( 
		&product.Id,
  		&product.Name,
  		&product.Price,
  		&product.Stock,
	)
	var errorMsg string
	
	if err != nil {
		if err == sql.ErrNoRows {
			errorMsg = "Data is not found!"
		} else {
			errorMsg = "Internal server error"
		}
	}

	data := EditPageData {
		Product: product,
		Error: errorMsg,
	}

	tmpl.Execute(w, data)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
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

	query := "UPDATE products SET name=?, price=?, stock=? WHERE id=?"
	_, err = db.DB.Exec(query, name, price, stock, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request){
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