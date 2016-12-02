package main

import (
	"fmt"
	"log"
	"os"

	"github.com/funayman/jadict-core/parser"
)

func main() {
	//get the file
	data, err := os.Open("./data/test-data.xml")
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	err = parser.Parse(data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Current count: ", parser.Count())
}
