package main

import (
	"net/http"

	"github.com/kalsmic/inventory-service/product"
)

const apiBasePath = "/api/v1"

func main() {
	product.SetupRoutes(apiBasePath)
	http.ListenAndServe(":5000", nil)

}
