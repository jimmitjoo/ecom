package main

import (
	"fmt"
	"log"

	"github.com/jimmitjoo/ecom/src/client"
	"github.com/jimmitjoo/ecom/src/client/products"
)

func main() {
	// Skapa ny klient
	c := client.NewClient("localhost:8080")

	// Skapa parametrar för GetProducts
	params := products.NewGetProductsParams()

	// Hämta alla produkter
	result, err := c.Products.GetProducts(params)
	if err != nil {
		log.Fatal(err)
	}

	// Skriv ut produkterna
	for _, p := range result.Payload {
		title := "Ingen titel"
		if p.BaseTitle != nil {
			title = *p.BaseTitle
		}
		fmt.Printf("Produkt: %s\n", title)
	}
}
