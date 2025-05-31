package app

import (
	"fmt"

	"github.com/lerentis/stunning-fiesta/internal/config"
	"github.com/lerentis/stunning-fiesta/internal/git"
	"github.com/lerentis/stunning-fiesta/internal/gitlab"
	"github.com/lerentis/stunning-fiesta/utils"
)

func CreateAndPushAppRepo(cfg config.Config, groupName string, serviceName string, vars map[string]interface{}) error {
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
	repo.AddChanges()
	repo.CommitChanges(fmt.Sprintf("Add App files for %s", serviceName))
	if err := repo.PushChanges(); err != nil {
		return fmt.Errorf("failed to push changes to App repository: %w", err)
	}
	fmt.Printf("App files rendered and pushed successfully for %s\n", serviceName)

	return nil
}
