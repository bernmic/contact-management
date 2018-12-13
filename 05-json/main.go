package main

import (
	"encoding/json"
	"fmt"
	"kontaktmanagement/05-json/model"
	"log"
)

func main() {
	contact := model.Contact{1, "Klaus", "Schmitz"}
	jsondata, err := json.Marshal(contact)
	if err != nil {
		log.Fatalf("Error marshalling data: %v\n", err)
	}
	fmt.Println(string(jsondata))

	ucontact := model.Contact{}
	err = json.Unmarshal(jsondata, &ucontact)
	if err != nil {
		log.Fatalf("Error unmarshalling data: %v\n", err)
	}
	fmt.Println(ucontact)
}
