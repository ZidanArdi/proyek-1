package homecontroller

import (
	"encoding/json"
	"net/http"
	"sateayammadura/models/adminmodel"
	"sateayammadura/models/productmodel"
	"text/template"
)

// Welcome Handler untuk halaman utama
func Welcome(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("views/home/index.html")
	if err != nil {
		panic(err)
	}
	products, _ := productmodel.GetAll()
	data := map[string]interface{}{
		"Products": products,
	}
	temp.Execute(w, data)
}

// Handler login admin
func AdminLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	admin, err := adminmodel.GetAdminByUsername(req.Username)
	if err != nil || admin == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if req.Password != admin.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// Set cookie session sederhana
	http.SetCookie(w, &http.Cookie{
		Name:     "admin_session",
		Value:    "loggedin",
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
	})
	w.WriteHeader(http.StatusOK)
}

// Handler cek status login admin
func AdminCheck(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("admin_session")
	if err != nil || cookie.Value != "loggedin" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Handler logout admin
func AdminLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "admin_session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	w.WriteHeader(http.StatusOK)
}
