package main

import (
	"fmt"
	"net/http"
	"pert3_npm/db"
	"pert3_npm/handlers"
)

const PORT = 8000 // 4 digit npm belakang

func main(){
	db.ConnectDB()
	// Akses halaman home
	http.HandleFunc("/", handlers.HomeHandler)
	// Akses route untuk menampilkan gambar
	http.HandleFunc("/image", handlers.ImageHandler)
	// akses halaman edit
	http.HandleFunc("/edit", handlers.EditProductHandler)
	// update data
	http.HandleFunc("/update", handlers.UpdateProduct)
	// create data
	http.HandleFunc("/create", handlers.CreateProductHandler)
	// delete data
	http.HandleFunc("/delete", handlers.DeleteProduct)
	
	fmt.Printf("Server berjalan di http://localhost:%d", PORT)
	http.ListenAndServe(fmt.Sprintf(":%d", PORT),nil)
}