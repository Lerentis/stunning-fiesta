package main

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/lerentis/stunning-fiesta/internal/config"
	"github.com/lerentis/stunning-fiesta/internal/k8s"
	"github.com/lerentis/stunning-fiesta/utils"
	"github.com/spf13/cobra"
)

func main() {
	var operation string
	var path string
	var configPath string
	var interactive bool
	home, _ := os.UserHomeDir()
	defaultConfig := home + "/.stunning-fiesta.yaml"

	var version = "dev"

	var rootCmd = &cobra.Command{
		Use:   "stunning-fiesta",
		Short: "A CLI helms to generate projects",
		Run: func(cmd *cobra.Command, args []string) {
			// Load configuration
			config, err := config.LoadConfig(configPath)
			if err != nil {
				fmt.Printf("Error loading config: %v\n", err)
				os.Exit(1)
			}
			// Check for updates
			if !utils.UpdateCheck(config.Endpoints.Update, version) {
				fmt.Println("You are not using the latest version of Stunning Fiesta. Please update to the latest version.")
				os.Exit(1)
			}
			// If interactive or missing args, prompt
			if interactive || operation == "" || path == "" {
				survey.AskOne(&survey.Input{
					Message: "What kind of operation do you want to perform?",
				}, &operation)
			}

			if operation == "" || path == "" {
				fmt.Println("operation and path are required.")
				cmd.Usage()
				os.Exit(1)
			}

			if operation == "helm" {
				fmt.Println("Performing helm operation...")
				app := utils.ParseAppJson(path)
				k8s.RenderTemplates(app)
			}
		},
	}

	rootCmd.Flags().StringVar(&operation, "operation", "", "What kind of operation do you want to perform")
	rootCmd.Flags().StringVar(&path, "path", "", "Path to app.json File")
	rootCmd.Flags().StringVar(&configPath, "config", defaultConfig, "Path to config file")
	rootCmd.Flags().BoolVar(&interactive, "interactive", false, "Run in interactive mode")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
