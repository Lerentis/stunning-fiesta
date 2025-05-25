package gitlab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/lerentis/stunning-fiesta/internal/config"
)

// --- Generic API Helper ---

func doGitlabAPIRequest(cfg config.Config, method, path string, payload interface{}) (*http.Response, error) {
	gitlabToken := os.Getenv("GITLAB_TOKEN")
	if gitlabToken == "" {
		return nil, fmt.Errorf("GITLAB_TOKEN environment variable not set")
	}

	apiURL := fmt.Sprintf("%s/api/v4/%s", cfg.GitlabURL, path)

	var body *bytes.Buffer
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}
		body = bytes.NewBuffer(b)
	} else {
		body = &bytes.Buffer{}
	}

	req, err := http.NewRequest(method, apiURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("PRIVATE-TOKEN", gitlabToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	return resp, nil
}

// --- Group Helpers ---

func groupExists(cfg config.Config, groupName string) (bool, error) {
	checkPath := fmt.Sprintf("groups?search=%s", groupName)
	resp, err := doGitlabAPIRequest(cfg, "GET", checkPath, nil)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var groups []struct {
		Name string `json:"name"`
		Path string `json:"path"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&groups); err != nil {
		return false, fmt.Errorf("failed to decode group search response: %w", err)
	}
	for _, group := range groups {
		if group.Path == groupName {
			return true, nil
		}
	}
	return false, nil
}

func createGroupAPI(cfg config.Config, path string, groupName string) error {
	payload := map[string]interface{}{
		"name":       groupName,
		"path":       groupName,
		"visibility": "private",
	}
	createResp, err := doGitlabAPIRequest(cfg, "POST", "groups", payload)
	if err != nil {
		return err
	}
	defer createResp.Body.Close()

	if createResp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create group: status %s", createResp.Status)
	}
	return nil
}

func CreateGroup(cfg config.Config, groupName string) error {
	exists, err := groupExists(cfg, groupName)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("group '%s' already exists", groupName)
	}
	return createGroupAPI(cfg, groupName, groupName)
}

// --- Repo Helpers ---

func getGroupID(cfg config.Config, groupName string) (int, error) {
	checkPath := fmt.Sprintf("groups?search=%s", groupName)
	resp, err := doGitlabAPIRequest(cfg, "GET", checkPath, nil)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var groups []struct {
		ID   int    `json:"id"`
		Path string `json:"path"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&groups); err != nil {
		return 0, fmt.Errorf("failed to decode group search response: %w", err)
	}
	for _, group := range groups {
		if group.Path == groupName {
			return group.ID, nil
		}
	}
	return 0, fmt.Errorf("group '%s' not found", groupName)
}

func getOrCreateSubgroup(cfg config.Config, parentGroupID int, subgroupName string) (int, error) {
	// Check if subgroup exists
	searchK8sPath := fmt.Sprintf("groups/%d/subgroups?search=%s", parentGroupID, subgroupName)
	resp, err := doGitlabAPIRequest(cfg, "GET", searchK8sPath, nil)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var subgroups []struct {
		ID   int    `json:"id"`
		Path string `json:"path"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&subgroups); err != nil {
		return 0, fmt.Errorf("failed to decode subgroup search response: %w", err)
	}
	for _, sg := range subgroups {
		if sg.Path == subgroupName {
			return sg.ID, nil
		}
	}

	// Create subgroup if not found
	payload := map[string]interface{}{
		"name":       subgroupName,
		"path":       subgroupName,
		"parent_id":  parentGroupID,
		"visibility": "private",
	}
	createResp, err := doGitlabAPIRequest(cfg, "POST", "groups", payload)
	if err != nil {
		return 0, err
	}
	defer createResp.Body.Close()
	if createResp.StatusCode != http.StatusCreated {
		return 0, fmt.Errorf("failed to create %s subgroup: status %s", subgroupName, createResp.Status)
	}
	var newGroup struct {
		ID int `json:"id"`
	}
	if err := json.NewDecoder(createResp.Body).Decode(&newGroup); err != nil {
		return 0, fmt.Errorf("failed to decode %s group creation response: %w", subgroupName, err)
	}
	return newGroup.ID, nil
}

func projectExists(cfg config.Config, fullPath string, serviceName string) (string, error) {
	searchProjectPath := fmt.Sprintf("projects?search=%s", serviceName)
	resp, err := doGitlabAPIRequest(cfg, "GET", searchProjectPath, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var projects []struct {
		PathWithNamespace string `json:"path_with_namespace"`
		HTTPURLToRepo     string `json:"http_url_to_repo"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		return "", fmt.Errorf("failed to decode project search response: %w", err)
	}
	for _, project := range projects {
		if project.PathWithNamespace == fullPath {
			return project.HTTPURLToRepo, nil
		}
	}
	return "", nil
}

func createProject(cfg config.Config, namespaceID int, serviceName string) (string, error) {
	payload := map[string]interface{}{
		"name":         serviceName,
		"path":         serviceName,
		"namespace_id": namespaceID,
		"visibility":   "private",
	}
	createResp, err := doGitlabAPIRequest(cfg, "POST", "projects", payload)
	if err != nil {
		return "", err
	}
	defer createResp.Body.Close()
	if createResp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("failed to create project: status %s", createResp.Status)
	}
	var createdProject struct {
		HTTPURLToRepo string `json:"http_url_to_repo"`
	}
	if err := json.NewDecoder(createResp.Body).Decode(&createdProject); err != nil {
		return "", fmt.Errorf("failed to decode project creation response: %w", err)
	}
	return createdProject.HTTPURLToRepo, nil
}

func CreateKubernetesRepo(cfg config.Config, groupName string, serviceName string) (string, error) {
	// 1. Get group ID
	groupID, err := getGroupID(cfg, groupName)
	if err != nil {
		return "", err
	}

	// 2. Get or create k8s subgroup
	k8sGroupID, err := getOrCreateSubgroup(cfg, groupID, "k8s")
	if err != nil {
		return "", err
	}

	// 3. Check if project exists
	projectPath := fmt.Sprintf("%s/k8s/%s", groupName, serviceName)
	if url, err := projectExists(cfg, projectPath, serviceName); err != nil {
		return "", err
	} else if url != "" {
		return url, fmt.Errorf("repository '%s' already exists", projectPath)
	}

	// 4. Create project
	return createProject(cfg, k8sGroupID, serviceName)
}

func CreateServiceRepo(cfg config.Config, groupName string, serviceName string) (string, error) {
	topicGroupID, err := getGroupID(cfg, groupName)
	if err != nil {
		return "", err
	}

	projectPath := fmt.Sprintf("%s/%s", groupName, serviceName)
	if url, err := projectExists(cfg, projectPath, serviceName); err != nil {
		return "", err
	} else if url != "" {
		return url, fmt.Errorf("repository '%s' already exists", projectPath)
	}
	return createProject(cfg, topicGroupID, serviceName)
}

func CreateInfrastructureRepo(cfg config.Config, groupName string) (string, error) {
	infraProjectName := "infrastructure-clz"
	topicGroupID, err := getGroupID(cfg, groupName)
	if err != nil {
		return "", err
	}

	projectPath := fmt.Sprintf("%s/%s", groupName, infraProjectName)
	if url, err := projectExists(cfg, projectPath, infraProjectName); err != nil {
		return "", err
	} else if url != "" {
		return url, fmt.Errorf("repository '%s' already exists", projectPath)
	}
	return createProject(cfg, topicGroupID, infraProjectName)
}

func CreateNamespaceRepo(cfg config.Config, serviceName string) (string, error) {
	topicGroupID, err := getGroupID(cfg, cfg.NamespacesRepo)
	if err != nil {
		return "", err
	}

	projectPath := fmt.Sprintf("%s", cfg.NamespacesRepo)
	if url, err := projectExists(cfg, projectPath, serviceName); err != nil {
		return "", err
	} else if url != "" {
		return url, fmt.Errorf("repository '%s' already exists", projectPath)
	}
	return createProject(cfg, topicGroupID, serviceName)
}
