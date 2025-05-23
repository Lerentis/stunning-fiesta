package helm

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/lerentis/stunning-fiesta/pkg/utils"
)

func ParseAppJson(jsonPath string) {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}
	var app utils.App
	err = json.Unmarshal(data, &app)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON data: %v", err)
	}
	fmt.Printf("App: %+v\n", app)
}
