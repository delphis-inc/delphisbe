package main

import (
	"fmt"

	"github.com/nedrocks/delphisbe/gql"
)

func main() {
	_, err := gql.NewServer()
	if err != nil {
		fmt.Printf("Failed\n")
		return
	}
	fmt.Printf("Success\n")
}
