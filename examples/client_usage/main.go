package main

import (
	"fmt"
	"log"

	"github.com/jimmitjoo/ecom/src/client"
)

func main() {
	// Skapa ny klient
	c := client.NewClient("localhost:8080")

	// Hämta alla produkter
	products, err := c.Products.ListProducts()
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range products {
		fmt.Printf("Produkt: %s\n", p.BaseTitle)
	}
}
