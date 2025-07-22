package homecontroller

import (
	
	"net/http"
	"text/template"
	"sateayammadura/models/productmodel"
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

