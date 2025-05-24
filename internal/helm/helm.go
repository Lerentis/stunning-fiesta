package helm

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/lerentis/stunning-fiesta/utils"
)

func ParseAppJson(jsonPath string) utils.App {
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
	return app
}

func GetTemplate(name string) (string, error) {
	//data, err := templates.FS.ReadFile(name)
	//if err != nil {
	//	return "", err
	//}
	return "", nil
}

func RenderTemplates(app utils.App) {
	for _, templateName := range []string{"requirements.yaml.tmpl", "values.yaml.tmpl"} {
		template, err := GetTemplate(templateName)
		if err != nil {
			log.Fatalf("Error getting template: %v", err)
		}
		fmt.Println(template)
	}
}
