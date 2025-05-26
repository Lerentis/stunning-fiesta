package gitlab

// import (
// 	"context"
// 	"fmt"
// 	"os"
// 	"testing"
// 	"time"

// 	"github.com/lerentis/stunning-fiesta/internal/config"
// 	"github.com/testcontainers/testcontainers-go"
// 	"github.com/testcontainers/testcontainers-go/wait"
// )

// func setupGitlabContainer(t *testing.T) (testcontainers.Container, string, func()) {
// 	ctx := context.Background()
// 	req := testcontainers.ContainerRequest{
// 		Image:        "gitlab/gitlab-ce:latest",
// 		ExposedPorts: []string{"8080/tcp", "443/tcp", "22/tcp"},
// 		Env: map[string]string{
// 			"GITLAB_OMNIBUS_CONFIG": "external_url 'http://localhost:8080'; gitlab_rails['initial_root_password'] = 'phaeyah1EShashabeuyu3cii3oong1ei';",
// 		},
// 		WaitingFor: wait.ForHTTP("/users/sign_in").WithPort("8080/tcp").WithStartupTimeout(5 * time.Minute),
// 	}
// 	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
// 		ContainerRequest: req,
// 		Started:          true,
// 	})
// 	if err != nil {
// 		t.Fatalf("Failed to start GitLab container: %v", err)
// 	}

// 	host, err := container.Host(ctx)
// 	if err != nil {
// 		t.Fatalf("Failed to get container host: %v", err)
// 	}
// 	port, err := container.MappedPort(ctx, "8080")
// 	if err != nil {
// 		t.Fatalf("Failed to get mapped port: %v", err)
// 	}
// 	baseURL := fmt.Sprintf("http://%s:%s", host, port.Port())

// 	cleanup := func() {
// 		_ = container.Terminate(ctx)
// 	}

// 	return container, baseURL, cleanup
// }

// func TestGitlabAPIIntegration(t *testing.T) {
// 	// Skip in short mode or CI if needed
// 	if testing.Short() {
// 		t.Skip("skipping integration test in short mode")
// 	}

// 	_, baseURL, cleanup := setupGitlabContainer(t)
// 	defer cleanup()

// 	// Set up config
// 	cfg := config.Config{
// 		GitlabURL: baseURL,
// 	}

// 	// Set the root token for testing (in real use, fetch from container logs or API)
// 	os.Setenv("GITLAB_TOKEN", "phaeyah1EShashabeuyu3cii3oong1ei")

// 	// Wait for GitLab to be ready for API (it may take a while after HTTP is up)
// 	time.Sleep(60 * time.Second)

// 	// Test group creation
// 	groupName := "testgroup"
// 	err := CreateGroup(cfg, groupName)
// 	if err != nil {
// 		t.Fatalf("CreateGroup failed: %v", err)
// 	}

// 	// Test duplicate group creation
// 	err = CreateGroup(cfg, groupName)
// 	if err == nil {
// 		t.Errorf("Expected error for duplicate group creation, got nil")
// 	}

// 	// Test repo creation
// 	repoURL, err := CreateKubernetesRepo(cfg, groupName, "testservice")
// 	if err != nil {
// 		t.Fatalf("CreateKubernetesRepo failed: %v", err)
// 	}
// 	if repoURL == "" {
// 		t.Errorf("Expected non-empty repo URL")
// 	}
// }
