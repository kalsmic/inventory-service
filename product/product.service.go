package product

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		productsJson, err := json.Marshal(ProductList)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(productsJson)

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
		newProduct.ProductID = GetNextId()
		ProductList = append(ProductList, newProduct)
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

	foundProduct, listItemIndex := FindProductById(productID)
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
		var updateProduct Product
		bodyBytes, err := ioutil.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Missing Body"))

			return
		}
		err = json.Unmarshal(bodyBytes, &updateProduct)
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
		updateProduct.ProductID = productID
		foundProduct = &updateProduct
		ProductList[listItemIndex] = *foundProduct
		w.WriteHeader(http.StatusOK)
		return
	case http.MethodDelete:
		size := len(ProductList)

		if listItemIndex == size-1 {
			ProductList = ProductList[:listItemIndex]
		} else {
			ProductList = append(ProductList[:listItemIndex], ProductList[listItemIndex+1:]...)

		}
		w.WriteHeader(http.StatusNoContent)

		return
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return

	}
}
