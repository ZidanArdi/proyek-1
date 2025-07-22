package categorycontroller

import (
	"html/template"
	"net/http"
	"sateayammadura/entities"
	"sateayammadura/models/categorymodel"
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Index(w http.ResponseWriter, r *http.Request) {
    // Ambil data kategori dari model
    categories, err := categorymodel.GetAllCategories()
    if err != nil {
        http.Error(w, "Gagal mengambil data kategori", http.StatusInternalServerError)
        return
    }

    // Cek apakah ada query ?deleted=true di URL
    deleted := r.URL.Query().Get("deleted")

    // Siapkan data untuk dikirim ke template
    data := map[string]any{
        "categories": categories,
        "deleted":    deleted, // tambahkan untuk keperluan alert
    }

    // Parsing dan eksekusi template
    temp, err := template.ParseFiles("views/category/index.html")
    if err != nil {
        http.Error(w, "Template error", http.StatusInternalServerError)
        return
    }

    temp.Execute(w, data)
}


func AboutUs(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Ini adalah halaman About Us"))
}

func Add(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Tampilkan halaman form
		temp, err := template.ParseFiles("views/category/create.html")
		if err != nil {
			http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		temp.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		// WAJIB: baca form POST
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Gagal membaca form: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Ambil data dari form
		var category entities.Category
		category.Name = r.FormValue("name")
		category.CreatedAt = time.Now()
		category.UpdatedAt = time.Now()

		// Simpan ke MongoDB
		if ok := categorymodel.Create(category); !ok {
			// Jika gagal, render ulang dengan error
			temp, _ := template.ParseFiles("views/category/create.html")
			temp.Execute(w, struct {
				Error string
			}{
				Error: "Gagal menyimpan data kategori",
			})
			return
		}

		// Redirect kalau berhasil
		http.Redirect(w, r, "/categories", http.StatusSeeOther)
	}
}

func Edit(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        temp, err := template.ParseFiles("views/category/edit.html")
        if err != nil {
            http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
            return
        }

        idString := r.URL.Query().Get("id")
        objectID, err := primitive.ObjectIDFromHex(idString)
        if err != nil {
            http.Error(w, "ID tidak valid: "+err.Error(), http.StatusBadRequest)
            return
        }

        category := categorymodel.Detail(objectID)
        if category.ID.IsZero() {
            http.Error(w, "Data tidak ditemukan", http.StatusNotFound)
            return
        }

        data := map[string]any{
            "category": category,
        }

        temp.Execute(w, data)
    }

    if r.Method == http.MethodPost {
        // Ambil ID dari form
        idString := r.FormValue("id")
        objectID, err := primitive.ObjectIDFromHex(idString)
        if err != nil {
            http.Error(w, "ID tidak valid: "+err.Error(), http.StatusBadRequest)
            return
        }

        // Ambil data lama dari DB
        category := categorymodel.Detail(objectID)
        if category.ID.IsZero() {
            http.Error(w, "Data tidak ditemukan", http.StatusNotFound)
            return
        }

        // Perbarui field dari form
        category.Name = r.FormValue("name")
        category.UpdatedAt = time.Now()

        // Simpan perubahan ke database
        if ok := categorymodel.Update(objectID, category); !ok {
            http.Error(w, "Gagal menyimpan perubahan", http.StatusInternalServerError)
            return
        }

        // Redirect ke halaman daftar kategori
        http.Redirect(w, r, "/categories", http.StatusSeeOther)
    }
}



func Delete(w http.ResponseWriter, r *http.Request) {
    // Ambil ID dari query string
    idString := r.URL.Query().Get("id")
    objectID, err := primitive.ObjectIDFromHex(idString)
    if err != nil {
        http.Error(w, "ID tidak valid: "+err.Error(), http.StatusBadRequest)
        return
    }

    // Jalankan delete dari model
    if ok := categorymodel.Delete(objectID); !ok {
        http.Error(w, "Gagal menghapus data kategori", http.StatusInternalServerError)
        return
    }

    // Redirect ke halaman kategori dengan alert berhasil hapus
    http.Redirect(w, r, "/categories?deleted=success", http.StatusSeeOther)

}


