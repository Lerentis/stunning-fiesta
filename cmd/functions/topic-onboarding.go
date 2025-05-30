package functions

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/lerentis/stunning-fiesta/internal/config"
	"github.com/lerentis/stunning-fiesta/utils"
	"github.com/spf13/cobra"
)

var topicOnboarding = &cobra.Command{
	Use:   "topic-onboarding",
	Short: "Onboard a new topic",
	Long:  `Onboard a new topic by providing necessary details and configurations.`,
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

		var selectedCluster string
		prompt := &survey.Select{
			Message: "Choose a cluster:",
			Options: clusterNames,
		}
		err = survey.AskOne(prompt, &selectedCluster)
		if err != nil {
			fmt.Printf("Error during selection: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("You selected cluster: %s\n", selectedCluster)

	},
}

func init() {
	rootCmd.AddCommand(topicOnboarding)
}
