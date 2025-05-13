package openlibraryapi

import (
	"encoding/json"
	"fmt"
	"log"
)

func prettyPrint(data any) {
	encodedData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatal("failed to pretty print map")
	}
	fmt.Println(string(encodedData))
}
