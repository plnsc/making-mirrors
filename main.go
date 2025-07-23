package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

type Repository struct {
	Provider string
	Owner    string
	Name     string
	URL      string
}

func main() {
	fmt.Println("Making Mirrors for Git Repositories")
	fmt.Println("===")

	// Define CLI flags
	var registryFile = flag.String("input", "registry.csv", "Path to the registry CSV file")
	var mirrorsDir = flag.String("output", "$HOME/Code/mirrors", "Directory to store mirrors")
	flag.Parse()

	// Use the CLI flag values (with defaults if not provided)
	finalMirrorsDir := *mirrorsDir
	finalRegistryFile := *registryFile

	// Expand environment variables and tilde (~) to full paths
	finalMirrorsDir = expandPath(finalMirrorsDir)
	fmt.Printf("Output directory: %s\n", finalMirrorsDir)
	finalRegistryFile = expandPath(finalRegistryFile)
	fmt.Printf("Registry file: %s\n", finalRegistryFile)

	// Create mirrors directory if it doesn't exist
	if err := os.MkdirAll(finalMirrorsDir, 0755); err != nil {
		log.Fatalf("Failed to create mirrors directory: %v", err)
	}

	// Read repositories from registry file
	repos, err := readRegistry(finalRegistryFile)
	if err != nil {
		log.Fatalf("Failed to read registry: %v", err)
	}

	fmt.Printf("Found %d repositories to mirror\n", len(repos))

	// Set up worker pool with all available CPU cores
	numWorkers := runtime.NumCPU()
	fmt.Printf("Using %d workers (CPU cores)\n", numWorkers)

	// Create channels for work distribution
	repoChan := make(chan Repository, len(repos))
	resultChan := make(chan string, len(repos))

	// Start workers
	var wg sync.WaitGroup
	for range numWorkers {
		wg.Add(1)
		go worker(finalMirrorsDir, repoChan, resultChan, &wg)
	}

	// Send repositories to workers
	go func() {
		for _, repo := range repos {
			repoChan <- repo
		}
		close(repoChan)
	}()

	// Collect results
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Print results
	fmt.Println("\nMirroring progress:")
	successCount := 0
	for result := range resultChan {
		fmt.Println(result)
		if strings.Contains(result, "✓") {
			successCount++
		}
	}

	fmt.Printf("\nCompleted! Successfully mirrored %d/%d repositories\n", successCount, len(repos))
}

// expandPath expands environment variables and tilde (~) in file paths
func expandPath(path string) string {
	// First expand environment variables
	path = os.ExpandEnv(path)

	// Then handle tilde expansion
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Failed to get home directory: %v", err)
		}
		path = filepath.Join(homeDir, path[2:])
	}

	return path
}

func readRegistry(filename string) ([]Repository, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %v", filename, err)
	}
	defer file.Close()

	var repos []Repository
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue // Skip empty lines and comments
		}

		repo, err := parseRepositoryLine(line)
		if err != nil {
			log.Printf("Warning: Failed to parse line '%s': %v", line, err)
			continue
		}

		repos = append(repos, repo)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return repos, nil
}

func parseRepositoryLine(line string) (Repository, error) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return Repository{}, fmt.Errorf("invalid format: expected 'provider:owner/repo.git'")
	}

	provider := parts[0]
	repoPath := strings.TrimSuffix(parts[1], ".git")

	pathParts := strings.Split(repoPath, "/")
	if len(pathParts) != 2 {
		return Repository{}, fmt.Errorf("invalid repository path: expected 'owner/repo'")
	}

	owner := pathParts[0]
	name := pathParts[1]

	var url string
	switch provider {
	case "github":
		url = fmt.Sprintf("https://github.com/%s/%s.git", owner, name)
	case "gitlab":
		url = fmt.Sprintf("https://gitlab.com/%s/%s.git", owner, name)
	case "bitbucket":
		url = fmt.Sprintf("https://bitbucket.org/%s/%s.git", owner, name)
	default:
		return Repository{}, fmt.Errorf("unsupported provider: %s", provider)
	}

	return Repository{
		Provider: provider,
		Owner:    owner,
		Name:     name,
		URL:      url,
	}, nil
}

func worker(mirrorsDir string, repoChan <-chan Repository, resultChan chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for repo := range repoChan {
		result := mirrorRepository(mirrorsDir, repo)
		resultChan <- result
	}
}

func mirrorRepository(mirrorsDir string, repo Repository) string {
	repoDir := filepath.Join(mirrorsDir, repo.Provider, repo.Owner, repo.Name)

	// Check if repository already exists (check for refs directory in bare repository)
	if _, err := os.Stat(filepath.Join(repoDir, "refs")); err == nil {
		// Repository exists, pull latest changes
		return pullRepository(repoDir, repo)
	} else {
		// Repository doesn't exist, clone it
		return cloneRepository(mirrorsDir, repo)
	}
}

func cloneRepository(mirrorsDir string, repo Repository) string {
	repoDir := filepath.Join(mirrorsDir, repo.Provider, repo.Owner, repo.Name)

	// Create parent directory
	if err := os.MkdirAll(filepath.Dir(repoDir), 0755); err != nil {
		return fmt.Sprintf("✗ %s/%s: Failed to create directory: %v", repo.Owner, repo.Name, err)
	}

	// Clone the repository
	cmd := exec.Command("git", "clone", "--mirror", repo.URL, repoDir)
	if err := cmd.Run(); err != nil {
		return fmt.Sprintf("✗ %s/%s: Clone failed: %v", repo.Owner, repo.Name, err)
	}

	return fmt.Sprintf("✓ %s/%s: Cloned successfully", repo.Owner, repo.Name)
}

func pullRepository(repoDir string, repo Repository) string {
	// Get the current state of refs before update
	beforeCmd := exec.Command("git", "-C", repoDir, "show-ref")
	beforeOutput, beforeErr := beforeCmd.Output()

	// Perform remote update
	cmd := exec.Command("git", "-C", repoDir, "remote", "update")
	if err := cmd.Run(); err != nil {
		return fmt.Sprintf("✗ %s/%s: Remote update failed: %v", repo.Owner, repo.Name, err)
	}

	// Get the state of refs after update
	afterCmd := exec.Command("git", "-C", repoDir, "show-ref")
	afterOutput, afterErr := afterCmd.Output()

	// If we couldn't get refs info, assume update was successful
	if beforeErr != nil || afterErr != nil {
		return fmt.Sprintf("✓ %s/%s: Updated successfully", repo.Owner, repo.Name)
	}

	// Compare before and after refs to see if anything changed
	beforeRefs := strings.TrimSpace(string(beforeOutput))
	afterRefs := strings.TrimSpace(string(afterOutput))

	if beforeRefs == afterRefs {
		return fmt.Sprintf("✓ %s/%s: Already up to date", repo.Owner, repo.Name)
	}

	// Check if this might be a force push by looking for shortened ref lists
	// (indicating some refs were updated/overwritten)
	beforeLines := strings.Split(beforeRefs, "\n")
	afterLines := strings.Split(afterRefs, "\n")

	// If we have significantly different number of refs, it might be a complex update
	if len(beforeLines) > 0 && len(afterLines) > 0 &&
		abs(len(beforeLines)-len(afterLines)) > len(beforeLines)/10 {
		return fmt.Sprintf("✓ %s/%s: Updated (significant changes detected)", repo.Owner, repo.Name)
	}

	return fmt.Sprintf("✓ %s/%s: Updated successfully", repo.Owner, repo.Name)
}

// Helper function to calculate absolute difference
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
