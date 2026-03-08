package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
	"pert3_npm/db"
	"pert3_npm/models"
)

func HomeView(w http.ResponseWriter, r *http.Request) {
	tmpl, _ = template.ParseFiles("templates/index.html")
	errorMsg := r.URL.Query().Get("error")
	rows, err := db.DB.Query(`
		SELECT 
			p.id,
			p.name,
			p.price,
			pd.stock,
			pd.is_active,
			pd.created_at
		FROM products p
		JOIN product_details pd ON p.id = pd.product_id
	`)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var products []ProductView

	for rows.Next() {
		var product models.Product

		err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Price,
			&product.Stock,
			&product.IsActive,
			&product.CreatedAt,
		)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		formattedDate := product.CreatedAt.Format("02 Jan 2006 15:04")

		products = append(products, ProductView{
			Id:        product.Id,
			Name:      product.Name,
			Price:     product.Price,
			Stock:     product.Stock,
			IsActive:  product.IsActive,
			CreatedAt: formattedDate,
		})
	}
	data := HomePageData{
		Products: products,
		Error:    errorMsg,
	}

	tmpl.Execute(w, data)
}

func ImageView(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	var image []byte

	err := db.DB.QueryRow("SELECT image FROM products WHERE id = ?", id).Scan(&image)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	contentType := http.DetectContentType(image)
	w.Header().Set("Content-Type", contentType)
	w.Write(image)
}

func EditView(w http.ResponseWriter, r *http.Request) {
	tmpl, _ = template.ParseFiles("templates/edit.html")
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Redirect(w, r, "?error=Id is required!", http.StatusSeeOther)
	}
	query := "SELECT id, name, price, stock, is_active, created_at FROM products WHERE id = ?"
	var product models.Product
	err := db.DB.QueryRow(query, id).Scan(
		&product.Id,
		&product.Name,
		&product.Price,
		&product.Stock,
		&product.IsActive,
		&product.CreatedAt,
	)
	var errorMsg string

	if err != nil {
		if err == sql.ErrNoRows {
			errorMsg = "Data is not found!"
		} else {
			errorMsg = "Internal server error"
		}
	}
	formattedDate := product.CreatedAt.Format("02 Jan 2006 15:04")

	data := EditPageData{
		Product: ProductView{
			Id:        product.Id,
			Name:      product.Name,
			Price:     product.Price,
			Stock:     product.Stock,
			IsActive:  product.IsActive,
			CreatedAt: formattedDate,
		},
		Error: errorMsg,
	}

	tmpl.Execute(w, data)
}