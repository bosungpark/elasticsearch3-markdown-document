package main

import (
	"elastic/source/utills"
	"fmt"
)

func main() {
	elasticManager := utills.ElasticManager{}

	err := elasticManager.Connect()
	if err != nil {
		fmt.Println("client failed", err)

		return
	}

	res, err := elasticManager.GetDocuments("bosung", "frDAto0BFjoKUr_vtenp")
	fmt.Println("GetDocuments: ", string(res))
}
