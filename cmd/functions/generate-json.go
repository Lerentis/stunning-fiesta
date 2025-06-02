package functions

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/lerentis/stunning-fiesta/internal/config"
	"github.com/lerentis/stunning-fiesta/utils"
	"github.com/spf13/cobra"
)

var output string

var generateJSON = &cobra.Command{
	Use:   "generate-json",
	Short: "Generate an application JSON file",
	Long:  `Generate an application JSON file based on the provided templates and configurations.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(defaultConfig)
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}

		clusterList, err := utils.ListAvailableClusters(*cfg)
		if err != nil {
			fmt.Printf("Error listing clusters: %v\n", err)
			os.Exit(1)
		}

		var clusterNames []string
		for _, c := range clusterList.Clusters {
			clusterNames = append(clusterNames, c.Name)
		}

		var app utils.App

		questions := []*survey.Question{
			{
				Name:     "Name",
				Prompt:   &survey.Input{Message: "What is the name of the application?"},
				Validate: survey.Required,
			},
			{
				Name:     "Registry",
				Prompt:   &survey.Input{Message: "What is the docker registry Path?"},
				Validate: survey.Required,
			},
			{
				Name:     "Project",
				Prompt:   &survey.Input{Message: "What is the project name?"},
				Validate: survey.Required,
			},
			{
				Name:     "IngressEnabeld",
				Prompt:   &survey.Input{Message: "Is ingress enabled? (yes/no)"},
				Validate: survey.Required,
			},
			{
				Name:     "DatabaseEnabled",
				Prompt:   &survey.Input{Message: "Is database enabled? (yes/no)"},
				Validate: survey.Required,
			},
			{
				Name: "Cluster",
				Prompt: &survey.Select{
					Message: "Select the cluster:",
					Options: clusterNames,
				},
				Validate: survey.Required,
			},
			{
				Name:     "ServiceNowOffering",
				Prompt:   &survey.Input{Message: "What is the ServiceNow offering?"},
				Validate: survey.Required,
			},
			{
				Name:     "ProjectSSH",
				Prompt:   &survey.Input{Message: "What is the project clone url (ssh)?"},
				Validate: survey.Required,
			},
		}

		err = survey.Ask(questions, &app)
		if err != nil {
			fmt.Printf("Error during survey: %v\n", err)
			os.Exit(1)
		}

		file, err := os.Create(output)
		if err != nil {
			fmt.Printf("Error creating file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		err = encoder.Encode(app)
		if err != nil {
			fmt.Printf("Error writing JSON to file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Application JSON successfully written to %s\n", output)
	},
}

func init() {
	rootCmd.AddCommand(generateJSON)
	generateJSON.PersistentFlags().StringVarP(&output, "output-file", "o", "", "Path to the app.json file to generate")

	generateJSON.MarkFlagRequired("output-file")
}
