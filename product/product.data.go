package product

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"sync"
)

var productMap = struct {
	sync.RWMutex
	items map[int]Product
}{items: make(map[int]Product)}

func init() {
	fmt.Println("loading products....")
	prodMap, err := loadProductMap()
	productMap.items = prodMap
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d products loaded....\n", len(productMap.items))
}

func loadProductMap() (map[int]Product, error) {
	filename := "products.json"

	// check if file exists
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file [%s] does not exist", filename)
	}

	file, _ := ioutil.ReadFile(filename)
	ProductList := make([]Product, 0)

	err = json.Unmarshal([]byte(file), &ProductList)
	if err != nil {
		log.Fatal(err)
	}

	prodMap := make(map[int]Product)

	for i := 0; i < len(ProductList); i++ {
		prodMap[ProductList[i].ProductID] = ProductList[i]
	}
	return prodMap, nil
}

// Returns a product with specified id from the Product Map
func getProduct(productID int) *Product {

	productMap.RLock()         // prevents another thread from getting a write lock on the struct while we read it
	defer productMap.RUnlock() // release lock from mutex
	if product, ok := productMap.items[productID]; ok {
		return &product
	}
	return nil
}

func removeProduct(productID int) {
	productMap.Lock()
	defer productMap.Unlock()
	delete(productMap.items, productID)
}

func getProductList() []Product {
	productMap.RLock()
	products := make([]Product, 0, len(productMap.items))
	for _, value := range productMap.items {
		products = append(products, value)
	}
	productMap.RUnlock()
	return products
}

/**
This method extracts product ids from the product map keys
@return list of product ids sorted in ascending order
**/
func getProductIds() []int {
	productMap.RLock()
	productIds := []int{}
	for key := range productMap.items {
		productIds = append(productIds, key)
	}
	productMap.RUnlock()
	sort.Ints(productIds)
	return productIds
}

/**
This method returns the next ID for a new product
**/
func getNextProductID() int {
	productIDs := getProductIds()
	size := len(productIDs)
	if size == 0 {
		return 1
	} else {
		return size
	}
}

// Create a new product if it does not already exist has an Id of Zero.
// Updates the product if it already exists
func addOrUpdateProduct(product Product) (int, error) {
	// if the product is set, update otherwise add
	addOrUpdateID := -1
	if product.ProductID > 0 {
		oldProduct := getProduct(product.ProductID)
		// if it does not exist, return error
		if oldProduct == nil {
			return 0, fmt.Errorf("Product id [%d] does not exist", product.ProductID)
		}
		// Otherwise replace it
		addOrUpdateID = product.ProductID
	} else {
		addOrUpdateID = getNextProductID()
		product.ProductID = addOrUpdateID
	}

	productMap.Lock()
	productMap.items[addOrUpdateID] = product
	productMap.Unlock()
	return addOrUpdateID, nil
}
