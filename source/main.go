package main

import (
	"elastic/source/utills"
	"fmt"
)

func main() {
	elasticManager := utills.ElasticManager{}

	err := elasticManager.Connect()
	if err != nil {
		fmt.Println("connect client failed", err)

		return
	}

	query := `{ "query": { "match_all": {} } }`
	res, err := elasticManager.SearchDocument("bosung", query)
	fmt.Println("SearchDocument: ", string(res))
}
