// Package main provides a tool for creating mirrors of Git repositories.
//
// Making Mirrors is a command-line application that reads a registry of Git repositories
// and creates local mirrors of them. It supports concurrent operations and handles
// repositories from various Git hosting providers.
//
// Usage:
//
//	making-mirrors [flags]
//
// Flags:
//
//	-input string
//	  	Path to the registry CSV file (default "$HOME/Code/mirrors/registry.txt")
//	-output string
//	  	Directory to store mirrors (default "$HOME/Code/mirrors")
//
// The registry file should contain repository information in a supported format,
// and the tool will create bare Git mirrors in the specified output directory.
//
// Example:
//
//	making-mirrors -input ./repos.txt -output ./mirrors
//
// Author: Paulo Nascimento <paulornasc@gmail.com>
// License: MIT
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

// Package metadata and constants
const (
	// AppName is the name of the application
	AppName = "making-mirrors"

	// AppVersion is the current version of the application
	AppVersion = "0.0.3"

	// AppDescription is a brief description of what the application does
	AppDescription = "A Go application for creating mirrors of Git repositories"

	// AppAuthor is the author of the application
	AppAuthor = "Paulo Nascimento <paulornasc@gmail.com>"

	// AppLicense is the license under which the application is distributed
	AppLicense = "MIT"

	// AppRepository is the URL to the source code repository
	AppRepository = "https://github.com/plnsc/making-mirrors"

	// DefaultRegistryFile is the default path for the registry file
	DefaultRegistryFile = "$HOME/Code/mirrors/registry.txt"

	// DefaultMirrorsDir is the default directory for storing mirrors
	DefaultMirrorsDir = "$HOME/Code/mirrors"
)

// BuildInfo contains build-time information
type BuildInfo struct {
	Version   string
	GitCommit string
	BuildTime string
	GoVersion string
}

type Repository struct {
	Provider string
	Owner    string
	Name     string
	URL      string
}

func main() {
	fmt.Printf("%s v%s\n", AppName, AppVersion)
	fmt.Println(AppDescription)
	fmt.Println("===")

	// Define CLI flags
	var registryFile = flag.String("input", DefaultRegistryFile, "Path to the registry CSV file")
	var mirrorsDir = flag.String("output", DefaultMirrorsDir, "Directory to store mirrors")
	var version = flag.Bool("version", false, "Show version information")
	flag.Parse()

	// Handle version flag
	if *version {
		fmt.Printf("%s version %s\n", AppName, AppVersion)
		fmt.Printf("Author: %s\n", AppAuthor)
		fmt.Printf("License: %s\n", AppLicense)
		fmt.Printf("Repository: %s\n", AppRepository)
		return
	}

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
	for i := 0; i < numWorkers; i++ {
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

	// Wait for all workers to finish and close the result channel
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results from workers
	fmt.Println("\nMirroring repositories...")
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
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Warning: failed to close file: %v", err)
		}
	}()

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
	case "gitea":
		url = fmt.Sprintf("https://gitea.com/%s/%s.git", owner, name)
	case "codecommit":
		awsParts := strings.Split(owner, "-")
		if len(awsParts) == 2 {
			region := awsParts[1]
			url = fmt.Sprintf("https://git-codecommit.%s.amazonaws.com/v1/repos/%s", region, name)
		} else {
			url = fmt.Sprintf("https://git-codecommit.%s.amazonaws.com/v1/repos/%s", owner, name)
		}
	case "azure":
		url = fmt.Sprintf("https://dev.azure.com/%s/%s/_git/%s", owner, name, name)
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
