package helper

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
)

func RedirectError(w http.ResponseWriter, r *http.Request, msg string) {
	errMsg := url.QueryEscape(msg)
	http.Redirect(w, r, "/?error="+errMsg, http.StatusSeeOther)
}

func ParsePrice(r *http.Request) (float64, error) {
	return strconv.ParseFloat(r.FormValue("price"), 64)
}

func ParseStock(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.FormValue("stock"), 10, 64)
}

func ReadAndValidateImage(file multipart.File) ([]byte, error) {
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed read file")
	}

	if len(fileBytes) > 1<<20 {
		return nil, fmt.Errorf("file size more than 1MB")
	}

	fileType := http.DetectContentType(fileBytes)

	if fileType != "image/jpeg" && fileType != "image/png" {
		return nil, fmt.Errorf("JPEG or PNG type only")
	}

	return fileBytes, nil
}