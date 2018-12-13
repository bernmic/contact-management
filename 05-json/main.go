package main

import (
	"encoding/json"
	"fmt"
	"kontaktmanagement/05-json/model"
	"log"
)

func main() {
	contact := model.Contact{1, "Klaus", "Schmitz"}
	json, err := json.Marshal(contact)
	if err != nil {
		log.Fatalf("Error marshalling data: %v\n", err)
	}
	fmt.Println(string(json))
}
