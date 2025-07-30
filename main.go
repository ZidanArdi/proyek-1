package main

import (
	"log"
	"net/http"
	"sateayammadura/config"
	"sateayammadura/controllers/categorycontroller"
	"sateayammadura/controllers/homecontroller"
	"sateayammadura/controllers/productcontroller" // <--- Tambahkan ini
)

func main() {
	config.Connection() // Initialize the database connection
	//1. Home Page
	http.HandleFunc("/", homecontroller.Welcome)
	http.HandleFunc("/admin/login", homecontroller.AdminLogin)
	http.HandleFunc("/admin/check", homecontroller.AdminCheck)
	http.HandleFunc("/admin/logout", homecontroller.AdminLogout)

	//2. Category Routes
	http.HandleFunc("/categories", categorycontroller.Index)
	http.HandleFunc("/categories/about", categorycontroller.AboutUs)
	http.HandleFunc("/categories/add", categorycontroller.Add)
	http.HandleFunc("/categories/edit", categorycontroller.Edit)
	http.HandleFunc("/categories/delete", categorycontroller.Delete)

	//3. Product Routes
	http.HandleFunc("/products", productcontroller.Index)
	http.HandleFunc("/products/detail", productcontroller.Detail)
	http.HandleFunc("/products/add", productcontroller.Add)
	http.HandleFunc("/products/edit", productcontroller.Edit)
	http.HandleFunc("/products/delete", productcontroller.Delete)

	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("public/images"))))
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	log.Println("Server running on port 8080")
	http.ListenAndServe(":8080", nil)

}
