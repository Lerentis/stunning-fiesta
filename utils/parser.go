package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func ParseAppJson(jsonPath string) App {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}
	var app App
	err = json.Unmarshal(data, &app)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON data: %v", err)
	}
	fmt.Printf("App: %+v\n", app)
	return app
}
