package app

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"

	"github.com/lerentis/stunning-fiesta/internal/config"
	"github.com/lerentis/stunning-fiesta/internal/git"
	"github.com/lerentis/stunning-fiesta/internal/gitlab"
	"github.com/lerentis/stunning-fiesta/utils"
)

func CreateAndPushAppRepo(cfg config.Config, groupName string, serviceName string, buildTool string, vars map[string]interface{}) error {
	cloneURL, err := gitlab.CreateServiceRepo(cfg, groupName, serviceName)
	if err != nil {
		return fmt.Errorf("failed to create App repository: %w", err)
	}
	fmt.Printf("App repository created successfully: %s\n", cloneURL)

	repo, err := git.CloneRepository(cloneURL)
	if err != nil {
		return fmt.Errorf("failed to clone App repository: %w", err)
	}
	repo.ChangeBranch("main")

	templatePath, err := utils.GetApplicationTemplates(cfg)
	if err != nil {
		return fmt.Errorf("failed to get App templates: %w", err)
	}

	err = utils.RenderTemplatesDir(templatePath, repo.Path, vars)
	if err != nil {
		return fmt.Errorf("failed to render templates: %w", err)
	}
	if err := AddDependenciesToAppRepo(cfg, repo.Path, serviceName, buildTool); err != nil {
		return fmt.Errorf("failed to add dependencies to App repository: %w", err)
	}
	repo.AddChanges()
	repo.CommitChanges(fmt.Sprintf("Add App files for %s", serviceName))
	if err := repo.PushChanges(); err != nil {
		return fmt.Errorf("failed to push changes to App repository: %w", err)
	}
	fmt.Printf("App files rendered and pushed successfully for %s\n", serviceName)

	return nil
}

func AddDependenciesToAppRepo(cfg config.Config, repoPath string, serviceName string, buildTool string) error {

	dependencies, err := utils.GetDependencies(cfg)
	if err != nil {
		return fmt.Errorf("failed to get dependencies: %w", err)
	}

	if buildTool == "maven" {
		if err := createPom(repoPath, serviceName, dependencies); err != nil {
			return fmt.Errorf("failed to create pom.xml: %w", err)
		}
	} else if buildTool == "gradle" {
		if err := createGradleBuild(repoPath, dependencies); err != nil {
			return fmt.Errorf("failed to create pom.xml: %w", err)
		}
	}

	return nil
}

func createGradleBuild(repoPath string, dependencies map[string]interface{}) error {
	buildGradlePath := filepath.Join(repoPath, "build.gradle")
	file, err := os.Create(buildGradlePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("plugins {\n    id 'java'\n}\n\ngroup = 'com.example'\nversion = '1.0.0'\n\n"))
	if err != nil {
		return err
	}

	_, err = file.WriteString("repositories {\n    mavenCentral()\n}\n\ndependencies {\n")
	if err != nil {
		return err
	}

	if depList, ok := dependencies["dependencies"].([]interface{}); ok {
		for _, dep := range depList {
			if depMap, ok := dep.(map[string]interface{}); ok {
				group := toString(depMap["groupId"])
				artifact := toString(depMap["artifactId"])
				version := toString(depMap["version"])
				_, err = file.WriteString(fmt.Sprintf("    implementation '%s:%s:%s'\n", group, artifact, version))
				if err != nil {
					return err
				}
			}
		}
	}

	_, err = file.WriteString("}\n")
	if err != nil {
		return err
	}

	return nil
}

func createPom(repoPath string, serviceName string, dependencies map[string]interface{}) error {
	var deps []utils.PomDependency
	if depList, ok := dependencies["dependencies"].([]interface{}); ok {
		for _, dep := range depList {
			if depMap, ok := dep.(map[string]interface{}); ok {
				deps = append(deps, utils.PomDependency{
					GroupID:    toString(depMap["groupId"]),
					ArtifactID: toString(depMap["artifactId"]),
					Version:    toString(depMap["version"]),
				})
			}
		}
	}

	var plugins []utils.PomPlugin
	if pluginList, ok := dependencies["plugins"].([]interface{}); ok {
		for _, plugin := range pluginList {
			if pluginMap, ok := plugin.(map[string]interface{}); ok {
				var executions []utils.PomPluginExecution
				if execList, ok := pluginMap["executions"].([]interface{}); ok {
					for _, exec := range execList {
						if execMap, ok := exec.(map[string]interface{}); ok {
							var goals []string
							if goalsList, ok := execMap["goals"].([]interface{}); ok {
								for _, g := range goalsList {
									goals = append(goals, toString(g))
								}
							}
							executions = append(executions, utils.PomPluginExecution{
								Phase: toString(execMap["phase"]),
								Goals: goals,
							})
						}
					}
				}
				plugins = append(plugins, utils.PomPlugin{
					GroupID:    toString(pluginMap["groupId"]),
					ArtifactID: toString(pluginMap["artifactId"]),
					Version:    toString(pluginMap["version"]),
					Executions: executions,
					Configuration: func() *utils.PomPluginConfiguration {
						if cfg, ok := pluginMap["configuration"].(*utils.PomPluginConfiguration); ok {
							return cfg
						}
						return nil
					}(),
				})
			}
		}
	}

	pom := utils.Pom{
		Xmlns:        "http://maven.apache.org/POM/4.0.0",
		XmlnsXsi:     "http://www.w3.org/2001/XMLSchema-instance",
		XsiSchema:    "http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd",
		ModelVersion: "4.0.0",
		GroupID:      "com.example",
		ArtifactID:   serviceName,
		Version:      "1.0.0",
		Dependencies: deps,
		Plugins:      plugins,
	}

	pomPath := filepath.Join(repoPath, "pom.xml")
	file, err := os.Create(pomPath)
	if err != nil {
		return err
	}
	defer file.Close()

	file.WriteString(xml.Header)
	enc := xml.NewEncoder(file)
	enc.Indent("", "  ")
	if err := enc.Encode(pom); err != nil {
		return err
	}

	return nil
}

// Helper to safely convert interface{} to string
func toString(val interface{}) string {
	if val == nil {
		return ""
	}
	if s, ok := val.(string); ok {
		return s
	}
	return ""
}
