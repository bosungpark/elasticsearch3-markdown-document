package main

import (
	"elastic/source/stores"
	"fmt"
)

func main() {
	elasticStore := stores.ElasticStore{}

	err := elasticStore.Connect()
	if err != nil {
		fmt.Println("connect client failed", err)

		return
	}

	err = elasticStore.DeleteDocument("bosung", "frDAto0BFjoKUr_vtenp")
}
