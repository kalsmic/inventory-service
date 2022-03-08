package main

import (
	"net/http"

	"github.com/kalsmic/inventory-service/product"
)

func main() {
	http.HandleFunc("/products", product.ProductsHandler)
	http.HandleFunc("/products/", product.ProductHandler)
	http.ListenAndServe(":5000", nil)

}
