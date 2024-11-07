package main

import (
	"fmt"
	"log"

	"github.com/jimmitjoo/ecom/src/client"
	"github.com/jimmitjoo/ecom/src/client/products"
)

func main() {
	// Create a new client
	c := client.NewClient("localhost:8080")

	// Create parameters for GetProducts
	params := products.NewGetProductsParams()

	// Get all products
	result, err := c.Products.GetProducts(params)
	if err != nil {
		log.Fatal(err)
	}

	// Print the products
	for _, p := range result.Payload {
		title := "No title"
		if p.BaseTitle != nil {
			title = *p.BaseTitle
		}
		fmt.Printf("Product: %s\n", title)
	}
}
