package functions

import (
	"fmt"
	"os"

	"github.com/lerentis/stunning-fiesta/internal/config"
	"github.com/spf13/cobra"
)

var (
	operation     string
	path          string
	configPath    string
	interactive   bool
	home, _       = os.UserHomeDir()
	defaultConfig = home + "/.stunning-fiesta.yaml"

	rootCmd = &cobra.Command{
		Use:   "stunning-fiesta",
		Short: "A CLI helps to generate projects",
	}
)

func Execute() {
	_, err := config.LoadConfig(defaultConfig)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}
	// Check for updates
	/*if !utils.UpdateCheck(config.Endpoints.Update, rootCmd.Version) {
		fmt.Println("You are not using the latest version of Stunning Fiesta. Please update to the latest version.")
		os.Exit(1)
	}*/
	fmt.Println("Welcome to Stunning Fiesta CLI!")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func init() {

}
