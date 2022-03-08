package main

import (
	"net/http"
	"os"
	"github.com/kalsmic/inventory-service/product"
)

const apiBasePath = "/api/v1"

func helloInventoryService(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
func main() {
	product.SetupRoutes(apiBasePath)
	http.HandleFunc("/health", helloInventoryService)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)

}
