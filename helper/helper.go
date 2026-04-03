package helper

import (
	"mime/multipart"
	"net/http"
)


/* TODO - Buat helper untuk redirect error */
func RedirectError(w http.ResponseWriter, r *http.Request, msg string) {
}

/* TODO - Buat helper untuk Parsing Price */
func ParsePrice(r *http.Request) (float64, error) {
}

/* TODO - Buat helper untuk Parsing Stock */
func ParseStock(r *http.Request) (int64, error) {
}

/* TODO - Buat helper untuk Membaca dan validasi gambar */
func ReadAndValidateImage(file multipart.File) ([]byte, error) {
	
}