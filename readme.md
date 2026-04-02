# 📘 Modul 3 - Activity 1

## Integrasi CRUD & Database (SQL) dengan Go

---

## 🧩 Langkah Pengerjaan

### 1. Koneksi ke MySQL

Gunakan perintah berikut untuk melakukan koneksi:

```bash
# Koneksi untuk Kampus F4 :
mysql -u root -p

# Koneksi untuk Kampus F8 :
mysql -h dbms.lepkom.f4.com -u APCx -p
```

> 📌 _Catatan:_ `APCx` x disesuaikan dengan no PC masing - masing

---

### 2. Membuat Database & Menggunakannya

```sql
CREATE DATABASE tokolepkom_npm;

USE tokolepkom_npm;
```

> 📌 _Catatan:_ `npm` disesuaikan dengan NPM praktikan

---

### 3. Membuat Tabel Database

#### Tabel `products`

```sql
CREATE TABLE products (
    id VARCHAR(100) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(15,2) NOT NULL
);
```

#### Tabel `product_details`

```sql
CREATE TABLE product_details (
    product_id VARCHAR(100) PRIMARY KEY,
    stock INT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    image LONGBLOB,
    FOREIGN KEY (product_id) REFERENCES products(id)
);
```

---

### 4. Import File SQL

Copy seluruh isi file `product.sql` lalu paste ke MySQL CLI.

---

### 5. Membuat Server Go

- Gunakan **port 4 digit terakhir NPM**
- Pastikan koneksi database sudah benar

Contoh koneksi:

```go
/* KONEKSI DI BAWAH UNTUK KAMPUS F4 */
username := "root"
password := ""
server := "127.0.0.1:3306"
dbName := "tokolepkom_npm"
dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
    username, password, server, dbName)


/* KONEKSI DI BAWAH UNTUK KAMPUS F8 */
username := "APCx"
password := "lepkom@123"
server := "dbms.lepkom.f4.com:3306"
dbName := "tokolepkom_npm"

dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
    username, password, server, dbName)
```

---

### 6. Membuat Model Product

File: `/models/Product.go`

```go
type Product struct {
    Id        string
    Name      string
    Price     float64
    Stock     int
    IsActive  bool
    CreatedAt time.Time
    Image     []byte
}
```

---

### 7. Setup Routing Server

Pastikan server memiliki route berikut:

| Route     | Fungsi                    |
| --------- | ------------------------- |
| `/`       | Menampilkan semua produk  |
| `/image`  | Menampilkan gambar produk |
| `/edit`   | Halaman edit produk       |
| `/update` | Update data produk        |
| `/create` | Tambah produk             |
| `/delete` | Hapus produk              |

---

### 8. Membuat Handler View

#### a. HomeView

- Menampilkan semua data produk
- Mengirim data ke template
- Menampilkan error jika ada

#### b. ImageView

- Mengambil gambar berdasarkan `id`
- Return sebagai response (image/jpeg atau png)
- Jika tidak ada → 404

#### c. EditView

- Mengambil data berdasarkan `id`
- Ditampilkan ke form edit

> 📌 _Catatan:_ isi dari file html dapat dicopy dari folder `templates/`

---

### 9. Membuat Helper Functions

#### a. RedirectError

```go
func RedirectError(w http.ResponseWriter, r *http.Request, err error) {
    http.Redirect(w, r, "/?error="+err.Error(), http.StatusSeeOther)
}
```

---

#### b. ParsePrice

```go
func ParsePrice(r *http.Request) (float64, error) {
    return strconv.ParseFloat(r.FormValue("price"), 64)
}
```

---

#### c. ParseStock

```go
func ParseStock(r *http.Request) (int, error) {
    return strconv.Atoi(r.FormValue("stock"))
}
```

---

#### d. ReadAndValidateImage

```go
func ReadAndValidateImage(r *http.Request) ([]byte, error) {
    file, _, err := r.FormFile("image")
    if err != nil {
        return nil, err
    }
    defer file.Close()

    data, err := io.ReadAll(file)
    if err != nil {
        return nil, err
    }

    if len(data) > 1<<20 {
        return nil, fmt.Errorf("file terlalu besar (maks 1MB)")
    }

    return data, nil
}
```

---

### 10. Membuat Handler CRUD

#### a. CreateProductHandler

- Ambil data dari form
- Validasi input
- Insert ke database

---

#### b. UpdateProductHandler

- Update data produk berdasarkan `id`
- Validasi input

---

#### c. DeleteProductHandler

- Hapus produk berdasarkan `id`

---

### 11. Testing Aplikasi

- Jalankan server dengan menggunakan perintah `go run -mod=vendor main.go`
- Akses melalui browser
- Pastikan semua fitur CRUD berjalan

---

## 📸 Dokumentasi

- Screenshot:
  - Koneksi DB
  - Tabel
  - Tampilan Web
  - CRUD berjalan

---
