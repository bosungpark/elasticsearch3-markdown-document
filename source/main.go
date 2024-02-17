package main

import (
	"elastic/source/utills"
	"elastic/source/schemas"
	"encoding/json"
	"fmt"
)

func main() {
	elasticManager := utills.ElasticManager{}

	err := elasticManager.Connect()
	if err != nil {
		fmt.Println("client failed", err)

		return
	}

	document := schemas.TestDocument{Message: "테스트가 잘 될까요?"}
	data, _ := json.Marshal(document)
	
	elasticManager.IndexDocuments("bosung", data)
}
