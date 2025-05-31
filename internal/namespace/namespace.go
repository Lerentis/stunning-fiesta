package namespace

import (
	"fmt"

	"github.com/lerentis/stunning-fiesta/internal/config"
	"github.com/lerentis/stunning-fiesta/internal/git"
	"github.com/lerentis/stunning-fiesta/internal/gitlab"
	"github.com/lerentis/stunning-fiesta/utils"
)

func CreateAndPushNamespace(cfg config.Config, serviceName string, vars map[string]interface{}) error {
	cloneURL, err := gitlab.CreateNamespaceRepo(cfg, serviceName)
	if err != nil {
		return fmt.Errorf("failed to create namespace repository: %w", err)
	}
	fmt.Printf("Namespace repository created successfully: %s\n", cloneURL)

	repo, err := git.CloneRepository(cloneURL)
	if err != nil {
		return fmt.Errorf("failed to clone namespace repository: %w", err)
	}
	repo.ChangeBranch("main")

	templatePath, err := utils.GetNamespaceTemplates(cfg)
	if err != nil {
		return fmt.Errorf("failed to get namespace templates: %w", err)
	}

	err = utils.RenderTemplatesDir(templatePath, repo.Path, vars)
	if err != nil {
		return fmt.Errorf("failed to render templates: %w", err)
	}
	repo.AddChanges()
	repo.CommitChanges(fmt.Sprintf("Add namespace files for %s", serviceName))
	if err := repo.PushChanges(); err != nil {
		return fmt.Errorf("failed to push changes to namespace repository: %w", err)
	}
	fmt.Printf("Namespace files rendered and pushed successfully for %s\n", serviceName)

	return nil
}
