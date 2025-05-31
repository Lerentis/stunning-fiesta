package k8s

import (
	"fmt"

	"github.com/lerentis/stunning-fiesta/internal/config"
	"github.com/lerentis/stunning-fiesta/internal/git"
	"github.com/lerentis/stunning-fiesta/internal/gitlab"
	"github.com/lerentis/stunning-fiesta/utils"
)

func CreateAndPushKubernetesRepo(cfg config.Config, groupName string, serviceName string, vars map[string]interface{}) error {
	cloneURL, err := gitlab.CreateKubernetesRepo(cfg, groupName, serviceName)
	if err != nil {
		return fmt.Errorf("failed to create k8s repository: %w", err)
	}
	fmt.Printf("k8s repository created successfully: %s\n", cloneURL)

	repo, err := git.CloneRepository(cloneURL)
	if err != nil {
		return fmt.Errorf("failed to clone k8s repository: %w", err)
	}
	repo.ChangeBranch("main")

	templatePath, err := utils.GetHelmTemplates(cfg)
	if err != nil {
		return fmt.Errorf("failed to get k8s templates: %w", err)
	}

	err = utils.RenderTemplatesDir(templatePath, repo.Path, vars)
	if err != nil {
		return fmt.Errorf("failed to render templates: %w", err)
	}
	repo.AddChanges()
	repo.CommitChanges(fmt.Sprintf("Add k8s files for %s", serviceName))
	if err := repo.PushChanges(); err != nil {
		return fmt.Errorf("failed to push changes to k8s repository: %w", err)
	}
	fmt.Printf("k8s files rendered and pushed successfully for %s\n", serviceName)

	return nil
}
