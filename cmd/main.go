package main

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/lerentis/stunning-fiesta/internal/helm"
	"github.com/spf13/cobra"
)

func main() {
	var operation string
	var path string
	var interactive bool

	var rootCmd = &cobra.Command{
		Use:   "stunning-fiesta",
		Short: "A CLI helms to generate projects",
		Run: func(cmd *cobra.Command, args []string) {
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
				app := helm.ParseAppJson(path)
				helm.RenderTemplates(app)
			}
		},
	}

	rootCmd.Flags().StringVar(&operation, "operation", "", "What kind of operation do you want to perform")
	rootCmd.Flags().StringVar(&path, "path", "", "Path to app.json File")
	rootCmd.Flags().BoolVar(&interactive, "interactive", false, "Run in interactive mode")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
