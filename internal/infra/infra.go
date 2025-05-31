package infra

import (
	"fmt"

	"github.com/lerentis/stunning-fiesta/internal/config"
	"github.com/lerentis/stunning-fiesta/internal/git"
	"github.com/lerentis/stunning-fiesta/internal/gitlab"
	"github.com/lerentis/stunning-fiesta/utils"
)

func CreateAndPushInfraRepo(cfg config.Config, groupName string, vars map[string]interface{}) error {
	cloneURL, err := gitlab.CreateInfrastructureRepo(cfg, groupName)
	if err != nil {
		return fmt.Errorf("failed to create infra repository: %w", err)
	}
	fmt.Printf("infra repository created successfully: %s\n", cloneURL)

	repo, err := git.CloneRepository(cloneURL)
	if err != nil {
		return fmt.Errorf("failed to clone infra repository: %w", err)
	}
	repo.ChangeBranch("main")

	templatePath, err := utils.GetInfrastructureTemplates(cfg)
	if err != nil {
		return fmt.Errorf("failed to get infra templates: %w", err)
	}

	err = utils.RenderTemplatesDir(templatePath, repo.Path, vars)
	if err != nil {
		return fmt.Errorf("failed to render templates: %w", err)
	}
	repo.AddChanges()
	repo.CommitChanges(fmt.Sprintf("Add infra files for %s", groupName))
	if err := repo.PushChanges(); err != nil {
		return fmt.Errorf("failed to push changes to infra repository: %w", err)
	}
	fmt.Printf("infra files rendered and pushed successfully for %s\n", groupName)

	return nil
}
