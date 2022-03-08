package product

import (
	"encoding/json"
	"log"
)

var ProductList []Product

func init() {
	productsJSON := `
	[
		{
			"productId": 1,
			"manufacturer": "Johns-Jenkins",
			"sku": "p5z343vdS",
			"upc": "939581000000",
			"pricePerUnit": "497.45",
			"quantityOnHand": 9703,
			"productName": "sticky note"
		  },
		  {
			"productId": 2,
			"manufacturer": "Hessel, Schimmel and Feeney",
			"sku": "i7v300kmx",
			"upc": "740979000000",
			"pricePerUnit": "282.29",
			"quantityOnHand": 9217,
			"productName": "leg warmers"
		  }
	]`
	err := json.Unmarshal([]byte(productsJSON), &ProductList)
	if err != nil {
		log.Fatal(err)
	}
}

func GetNextId() int {
	highestID := -1

	for _, product := range ProductList {
		if highestID < product.ProductID {
			highestID = product.ProductID
		}
	}
	return highestID + 1
}

func FindProductById(productID int) (*Product, int) {
	for i, product := range ProductList {
		if product.ProductID == productID {
			return &product, i
		}
	}
	return nil, 0
}
