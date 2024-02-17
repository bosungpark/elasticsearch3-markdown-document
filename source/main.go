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

	res, err := elasticStore.GetDocument("bosung", "frDAto0BFjoKUr_vtenp")
	fmt.Println(string(res))
}
