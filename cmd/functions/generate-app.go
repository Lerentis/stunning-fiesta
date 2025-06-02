package functions

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/lerentis/stunning-fiesta/internal/app"
	"github.com/lerentis/stunning-fiesta/internal/config"
	"github.com/lerentis/stunning-fiesta/internal/k8s"
	"github.com/lerentis/stunning-fiesta/utils"
	"github.com/spf13/cobra"
)

var generateApp = &cobra.Command{
	Use:   "generate-app",
	Short: "Onboard a new application or service",
	Long:  `Generate the Application and k8s repository with a default skeleton.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(defaultConfig)
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}

		var appName string
		if err := survey.AskOne(&survey.Input{
			Message: "Enter the name of your new service:",
		}, &appName, survey.WithValidator(survey.Required)); err != nil {
			fmt.Printf("Error entering ServiceNow offering: %v\n", err)
			os.Exit(1)
		}

		var buildTool string

		availableBuildTools := []string{"Maven", "Gradle"}

		if err := survey.AskOne(&survey.Select{
			Message: "Select the cluster:",
			Options: availableBuildTools,
		}, &buildTool, survey.WithValidator(survey.Required)); err != nil {
			fmt.Printf("Error selecting cluster: %v\n", err)
			os.Exit(1)
		}

		// 1. Select Cluster
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
		if err := survey.AskOne(&survey.Select{
			Message: "Select the cluster:",
			Options: clusterNames,
		}, &selectedCluster, survey.WithValidator(survey.Required)); err != nil {
			fmt.Printf("Error selecting cluster: %v\n", err)
			os.Exit(1)
		}

		// 2. Select Application Type
		appTypesList, err := utils.ListAvailableApplicationTypes(*cfg)
		if err != nil {
			fmt.Printf("Error listing application types: %v\n", err)
			os.Exit(1)
		}
		var appTypeNames []string
		for _, t := range appTypesList.ApplicationTypes {
			appTypeNames = append(appTypeNames, t.Name)
		}
		var selectedAppType string
		if err := survey.AskOne(&survey.Select{
			Message: "Select the application type:",
			Options: appTypeNames,
		}, &selectedAppType, survey.WithValidator(survey.Required)); err != nil {
			fmt.Printf("Error selecting application type: %v\n", err)
			os.Exit(1)
		}

		// 3. Ask if database should be enabled
		var databaseEnabled bool
		if err := survey.AskOne(&survey.Confirm{
			Message: "Enable database?",
			Default: false,
		}, &databaseEnabled); err != nil {
			fmt.Printf("Error selecting database enabled: %v\n", err)
			os.Exit(1)
		}

		// 4. Select DBMS if database is enabled
		var selectedDBMS string
		if databaseEnabled {
			dbmsList, err := utils.ListAvailableDBMS(*cfg)
			if err != nil {
				fmt.Printf("Error listing DBMS: %v\n", err)
				os.Exit(1)
			}
			var dbmsNames []string
			for _, d := range dbmsList.DBMS {
				dbmsNames = append(dbmsNames, d.Name)
			}
			if err := survey.AskOne(&survey.Select{
				Message: "Select the DBMS:",
				Options: dbmsNames,
			}, &selectedDBMS, survey.WithValidator(survey.Required)); err != nil {
				fmt.Printf("Error selecting DBMS: %v\n", err)
				os.Exit(1)
			}
		}

		// 5. Ask if ingress should be enabled
		var ingressEnabled bool
		if err := survey.AskOne(&survey.Confirm{
			Message: "Enable ingress?",
			Default: false,
		}, &ingressEnabled); err != nil {
			fmt.Printf("Error selecting ingress enabled: %v\n", err)
			os.Exit(1)
		}

		// 6. Select Team
		teamList, err := utils.ListAvailableTeams(*cfg)
		if err != nil {
			fmt.Printf("Error listing teams: %v\n", err)
			os.Exit(1)
		}
		var teamNames []string
		for _, t := range teamList.Teams {
			teamNames = append(teamNames, t.Name)
		}
		if len(teamNames) == 0 {
			fmt.Println("No teams available to select.")
			os.Exit(1)
		}
		var selectedTeam string
		if err := survey.AskOne(&survey.Select{
			Message: "Select the team:",
			Options: teamNames,
		}, &selectedTeam, survey.WithValidator(survey.Required)); err != nil {
			fmt.Printf("Error selecting team: %v\n", err)
			os.Exit(1)
		}

		// 7. Select Project
		projectList, err := utils.ListAvailableProjects(*cfg)
		if err != nil {
			fmt.Printf("Error listing projects: %v\n", err)
			os.Exit(1)
		}
		var projectNames []string
		for _, p := range projectList.Projects {
			projectNames = append(projectNames, p.Name)
		}

		var selectedTopic string
		if err := survey.AskOne(&survey.Select{
			Message: "Select the topic for your new app:",
			Options: projectNames,
		}, &selectedTopic, survey.WithValidator(survey.Required)); err != nil {
			fmt.Printf("Error selecting topic: %v\n", err)
			os.Exit(1)
		}
		// 8. Ask for ServiceNow offering
		var serviceNowOffering string
		if err := survey.AskOne(&survey.Input{
			Message: "Enter the ServiceNow offering:",
		}, &serviceNowOffering, survey.WithValidator(survey.Required)); err != nil {
			fmt.Printf("Error entering ServiceNow offering: %v\n", err)
			os.Exit(1)
		}

		vars := map[string]interface{}{
			"AppName":            appName,
			"BuildTool":          buildTool,
			"SelectedCluster":    selectedCluster,
			"SelectedAppType":    selectedAppType,
			"DatabaseEnabled":    databaseEnabled,
			"SelectedDBMS":       selectedDBMS,
			"IngressEnabled":     ingressEnabled,
			"SelectedTeam":       selectedTeam,
			"SelectedTopic":      selectedTopic,
			"ServiceNowOffering": serviceNowOffering,
		}

		if err := app.CreateAndPushAppRepo(
			*cfg,
			selectedTopic,
			appName,
			buildTool,
			vars,
		); err != nil {
			fmt.Printf("Error creating and pushing app repo: %v\n", err)
			os.Exit(1)
		}
		if err := k8s.CreateAndPushKubernetesRepo(*cfg, selectedTopic, appName, vars); err != nil {
			fmt.Printf("Error creating and pushing k8s repo: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Application and Kubernetes repositories created and pushed successfully.")
	},
}

func init() {
	rootCmd.AddCommand(generateApp)
}
