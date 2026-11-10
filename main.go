package main

import (
	"fmt"
	"log"
	"vran/blogaggre/internal/config"
)

func main() {
	c, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	config.SetUser(&c, "alex")

	c2, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", c2)
}
