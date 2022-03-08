package product

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/kalsmic/inventory-service/cors"
)

const productsBasePath = "products"

func SetupRoutes(apiBasePath string) {
	handleProducts := http.HandlerFunc(productsHandler)
	handleProduct := http.HandlerFunc(productHandler)
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, productsBasePath), cors.Middleware(handleProducts))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, productsBasePath), cors.Middleware(handleProduct))

}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		productList := getProductList()
		productsJson, err := json.Marshal(productList)

		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		_, err = w.Write(productsJson)

		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return

	case http.MethodPost:
		var newProduct Product
		bodyBytes, err := ioutil.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(bodyBytes, &newProduct)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if newProduct.ProductID != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = addOrUpdateProduct(newProduct)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	}
}

func productHandler(w http.ResponseWriter, r *http.Request) {

	urlPathSegments := strings.Split(r.URL.Path, "products/")
	productID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	foundProduct := getProduct(productID)
	if foundProduct == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
		return
	}

	switch r.Method {
	case http.MethodGet:
		productsJSON, err := json.Marshal(foundProduct)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(productsJSON)
		return
	case http.MethodPut:
		//  update the product in the list
		var updatedProduct Product
		bodyBytes, err := ioutil.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Missing Body"))

			return
		}
		err = json.Unmarshal(bodyBytes, &updatedProduct)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Poorly formated Json"))
			return
		}

		if foundProduct.ProductID != productID {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Product with Specified Id does not exist"))
			return
		}
		addOrUpdateProduct(updatedProduct)
		w.WriteHeader(http.StatusOK)
		return
	case http.MethodDelete:
		removeProduct(productID)
		w.WriteHeader(http.StatusNoContent)
		return
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return

	}
}
