package productcontroller

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"io"
	"strconv"
	"time"
	"fmt"
    

	"go.mongodb.org/mongo-driver/bson/primitive"

	"sateayammadura/entities"
	"sateayammadura/models/categorymodel"
	"sateayammadura/models/productmodel"
)

func Index(w http.ResponseWriter, r *http.Request) {
	products, err := productmodel.GetAll()
	if err != nil {
		log.Println("Error get products:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"products": products,
	}

	temp, err := template.ParseFiles("views/product/index.html")
	if err != nil {
		log.Println("Template error:", err)
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	if err := temp.Execute(w, data); err != nil {
		log.Println("Template execution error:", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}

// productcontroller.go
func Detail(w http.ResponseWriter, r *http.Request) {
    productID := r.URL.Query().Get("id")
    if productID == "" {
        http.Error(w, "Product ID is required", http.StatusBadRequest)
        return
    }

    product, err := productmodel.GetByID(productID)
    if err != nil {
        log.Println("Error getting product:", err)
        http.Error(w, "Produk tidak ditemukan", http.StatusNotFound)
        return
    }

    data := map[string]any{
        "product": product,
    }

    tmpl, err := template.ParseFiles("views/product/detail.html")
    if err != nil {
        log.Println("Template error:", err)
        http.Error(w, "Template error", http.StatusInternalServerError)
        return
    }

    if err := tmpl.Execute(w, data); err != nil {
        log.Println("Template execution error:", err)
        http.Error(w, "Error rendering page", http.StatusInternalServerError)
    }
}


func Add(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		temp, err := template.ParseFiles("views/product/create.html")
		if err != nil {
			log.Println("Template error:", err)
			http.Error(w, "Template error", http.StatusInternalServerError)
			return
		}

		categories, err := categorymodel.GetAllCategories()
		if err != nil {
			log.Println("Gagal mengambil kategori:", err)
			http.Error(w, "Gagal mengambil kategori", http.StatusInternalServerError)
			return
		}

		data := map[string]any{
			"categories": categories,
		}
		temp.Execute(w, data)
		return
	}

	// ========== HANDLE POST ==========
	r.Body = http.MaxBytesReader(w, r.Body, 50<<20)
	err := r.ParseMultipartForm(50 << 20)
	if err != nil {
		http.Error(w, "File terlalu besar", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("gambar")
	if err != nil {
		http.Error(w, "File gambar wajib diunggah", http.StatusBadRequest)
		return
	}
	defer file.Close()

	name := r.FormValue("name")
	categoryIDStr := r.FormValue("category_id")
	stockStr := r.FormValue("stock")
	description := r.FormValue("description")
	priceStr := r.FormValue("price")

	if name == "" || categoryIDStr == "" || stockStr == "" || description == "" || priceStr == "" {
		http.Error(w, "Semua field wajib diisi", http.StatusBadRequest)
		return
	}

	stock, err := strconv.Atoi(stockStr)
	if err != nil || stock < 0 {
		http.Error(w, "Stock harus berupa angka positif", http.StatusBadRequest)
		return
	}

	harga, err := strconv.Atoi(priceStr)
	if err != nil || harga < 0 {
		http.Error(w, "Harga harus berupa angka positif", http.StatusBadRequest)
		return
	}

	categoryID, err := primitive.ObjectIDFromHex(categoryIDStr)
	if err != nil {
		http.Error(w, "Kategori tidak valid", http.StatusBadRequest)
		return
	}

	// ✅ Ambil nama kategori dari koleksi
	category, err := categorymodel.GetCategoryByID(categoryID)
	if err != nil {
		http.Error(w, "Kategori tidak ditemukan", http.StatusBadRequest)
		return
	}

	// ✅ Simpan gambar
	filename := time.Now().Format("20060102150405") + "_" + handler.Filename
	filePath := "public/images/" + filename
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Gagal menyimpan file gambar", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	io.Copy(dst, file)

	// ✅ Buat entitas produk
	product := entities.Product{
		ID:           primitive.NewObjectID(),
		Name:         name,
		CategoryID:   categoryID,
		CategoryName: category.Name, // sinkronisasi nama kategori
		Stock:        stock,
		Description:  description,
		GambarURL:    "/images/" + filename,
		Harga:        harga,
		CreatedAt:    primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt:    primitive.NewDateTimeFromTime(time.Now()),
	}

	if ok := productmodel.Create(product); ok {
		http.Redirect(w, r, "/products", http.StatusSeeOther)
	} else {
		http.Error(w, "Gagal menyimpan produk", http.StatusInternalServerError)
	}
}

func Edit(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("id")
	if productID == "" {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodGet {
		product, err := productmodel.GetByID(productID)
		if err != nil {
			http.Error(w, "Produk tidak ditemukan", http.StatusNotFound)
			return
		}

		categories, err := categorymodel.GetAllCategories()
		if err != nil {
			http.Error(w, "Gagal mengambil kategori", http.StatusInternalServerError)
			return
		}

		data := map[string]any{
			"product":    product,
			"categories": categories,
		}

		tmpl, err := template.ParseFiles("views/product/edit.html")
		if err != nil {
			http.Error(w, "Template error", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, data)
		return
	}

	// Handle POST (update)
	err := r.ParseMultipartForm(50 << 20)
	if err != nil {
		http.Error(w, "Permintaan tidak valid", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")
	stockStr := r.FormValue("stock")
	priceStr := r.FormValue("price")
	categoryIDStr := r.FormValue("category_id")

	stock, _ := strconv.Atoi(stockStr)
	price, _ := strconv.Atoi(priceStr)
	categoryID, _ := primitive.ObjectIDFromHex(categoryIDStr)

	// Ambil nama kategori
	category, _ := categorymodel.GetCategoryByID(categoryID)

	// Cek apakah ada file gambar baru
	var imageURL string
	file, handler, err := r.FormFile("gambar")
	if err == nil {
		defer file.Close()
		filename := time.Now().Format("20060102150405") + "_" + handler.Filename
		filePath := "public/images/" + filename
		dst, _ := os.Create(filePath)
		defer dst.Close()
		io.Copy(dst, file)
		imageURL = "/images/" + filename
	}

	// Buat entitas produk update
	update := entities.Product{
		Name:         name,
		CategoryID:   categoryID,
		CategoryName: category.Name,
		Stock:        stock,
		Description:  description,
		Harga:        price,
		UpdatedAt:    primitive.NewDateTimeFromTime(time.Now()),
	}
	if imageURL != "" {
		update.GambarURL = imageURL
	}

	if err := productmodel.Update(productID, update); err != nil {
		http.Error(w, "Gagal mengupdate produk", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/products", http.StatusSeeOther)
}


func Delete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID produk tidak ditemukan", http.StatusBadRequest)
		return
	}

	err := productmodel.Delete(id)
	if err != nil {
		log.Println("Gagal menghapus produk:", err)
		http.Error(w, "Gagal menghapus produk", http.StatusInternalServerError)
		return
	}

	// ✅ Redirect ke halaman /products dan tambahkan query string
	http.Redirect(w, r, "/products?deleted=success", http.StatusSeeOther)
}

// Fungsi helper untuk format harga
func FormatRupiah(price int) string {
    result := fmt.Sprintf("Rp %s", formatThousands(price))
    return result
}

// Fungsi bantu untuk kasih titik pemisah ribuan
func formatThousands(n int) string {
    in := fmt.Sprintf("%d", n)
    out := ""
    for i, digit := range in {
        if i != 0 && (len(in)-i)%3 == 0 {
            out += "."
        }
        out += string(digit)
    }
    return out
}

