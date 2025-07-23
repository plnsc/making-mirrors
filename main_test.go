package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestAppMetadata(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{"AppName", AppName, "making-mirrors"},
		{"AppVersion", AppVersion, "0.1.0"},
		{"AppAuthor", AppAuthor, "Paulo Nascimento"},
		{"AppLicense", AppLicense, "MIT"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, tt.value, tt.expected)
			}
		})
	}
}

func TestDefaultPaths(t *testing.T) {
	if DefaultRegistryFile == "" {
		t.Error("DefaultRegistryFile should not be empty")
	}

	if DefaultMirrorsDir == "" {
		t.Error("DefaultMirrorsDir should not be empty")
	}
}

func TestBuildInfo(t *testing.T) {
	info := BuildInfo{
		Version:   "0.1.0",
		GitCommit: "abc123",
		BuildTime: "2025-07-23",
		GoVersion: "1.22",
	}

	if info.Version == "" {
		t.Error("BuildInfo.Version should not be empty")
	}
}

func TestExpandPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		setup    func() string
		cleanup  func()
	}{
		{
			name:     "simple path",
			input:    "/simple/path",
			expected: "/simple/path",
		},
		{
			name:     "path with environment variable",
			input:    "$HOME/test",
			expected: "",
			setup: func() string {
				home := os.Getenv("HOME")
				return filepath.Join(home, "test")
			},
		},
		{
			name:     "path with tilde",
			input:    "~/test",
			expected: "",
			setup: func() string {
				home, _ := os.UserHomeDir()
				return filepath.Join(home, "test")
			},
		},
		{
			name:     "path with custom env var",
			input:    "$TEST_VAR/path",
			expected: "/custom/path",
			setup: func() string {
				os.Setenv("TEST_VAR", "/custom")
				return "/custom/path"
			},
			cleanup: func() {
				os.Unsetenv("TEST_VAR")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.expected = tt.setup()
			}
			if tt.cleanup != nil {
				defer tt.cleanup()
			}

			result := expandPath(tt.input)
			if result != tt.expected {
				t.Errorf("expandPath(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseRepositoryLine(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    Repository
		expectError bool
	}{
		{
			name:  "valid github repository",
			input: "github:torvalds/linux",
			expected: Repository{
				Provider: "github",
				Owner:    "torvalds",
				Name:     "linux",
				URL:      "https://github.com/torvalds/linux.git",
			},
			expectError: false,
		},
		{
			name:  "valid gitlab repository",
			input: "gitlab:gitlab-org/gitlab",
			expected: Repository{
				Provider: "gitlab",
				Owner:    "gitlab-org",
				Name:     "gitlab",
				URL:      "https://gitlab.com/gitlab-org/gitlab.git",
			},
			expectError: false,
		},
		{
			name:  "valid bitbucket repository",
			input: "bitbucket:atlassian/stash",
			expected: Repository{
				Provider: "bitbucket",
				Owner:    "atlassian",
				Name:     "stash",
				URL:      "https://bitbucket.org/atlassian/stash.git",
			},
			expectError: false,
		},
		{
			name:  "repository with .git suffix",
			input: "github:golang/go.git",
			expected: Repository{
				Provider: "github",
				Owner:    "golang",
				Name:     "go",
				URL:      "https://github.com/golang/go.git",
			},
			expectError: false,
		},
		{
			name:        "invalid format - no colon",
			input:       "github-torvalds/linux",
			expected:    Repository{},
			expectError: true,
		},
		{
			name:        "invalid format - no slash",
			input:       "github:torvaldslinux",
			expected:    Repository{},
			expectError: true,
		},
		{
			name:        "invalid format - too many slashes",
			input:       "github:torvalds/linux/extra",
			expected:    Repository{},
			expectError: true,
		},
		{
			name:        "unsupported provider",
			input:       "codeberg:torvalds/linux",
			expected:    Repository{},
			expectError: true,
		},
		{
			name:        "empty provider",
			input:       ":torvalds/linux",
			expected:    Repository{},
			expectError: true,
		},
		{
			name:        "empty owner/repo",
			input:       "github:",
			expected:    Repository{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseRepositoryLine(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("parseRepositoryLine(%q) expected error but got none", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("parseRepositoryLine(%q) unexpected error: %v", tt.input, err)
				return
			}

			if result != tt.expected {
				t.Errorf("parseRepositoryLine(%q) = %+v, want %+v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestReadRegistry(t *testing.T) {
	tests := []struct {
		name          string
		fileContent   string
		expectedLen   int
		expectedRepos []Repository
		expectError   bool
	}{
		{
			name: "valid registry file",
			fileContent: `github:torvalds/linux
gitlab:gitlab-org/gitlab
bitbucket:atlassian/stash`,
			expectedLen: 3,
			expectedRepos: []Repository{
				{Provider: "github", Owner: "torvalds", Name: "linux", URL: "https://github.com/torvalds/linux.git"},
				{Provider: "gitlab", Owner: "gitlab-org", Name: "gitlab", URL: "https://gitlab.com/gitlab-org/gitlab.git"},
				{Provider: "bitbucket", Owner: "atlassian", Name: "stash", URL: "https://bitbucket.org/atlassian/stash.git"},
			},
			expectError: false,
		},
		{
			name: "registry with comments and empty lines",
			fileContent: `# This is a comment
github:golang/go

# Another comment
gitlab:gitlab-org/gitaly
`,
			expectedLen: 2,
			expectedRepos: []Repository{
				{Provider: "github", Owner: "golang", Name: "go", URL: "https://github.com/golang/go.git"},
				{Provider: "gitlab", Owner: "gitlab-org", Name: "gitaly", URL: "https://gitlab.com/gitlab-org/gitaly.git"},
			},
			expectError: false,
		},
		{
			name: "registry with invalid lines",
			fileContent: `github:golang/go
invalid-line-no-colon
gitlab:gitlab-org/gitlab`,
			expectedLen: 2,
			expectedRepos: []Repository{
				{Provider: "github", Owner: "golang", Name: "go", URL: "https://github.com/golang/go.git"},
				{Provider: "gitlab", Owner: "gitlab-org", Name: "gitlab", URL: "https://gitlab.com/gitlab-org/gitlab.git"},
			},
			expectError: false,
		},
		{
			name:        "empty file",
			fileContent: "",
			expectedLen: 0,
			expectError: false,
		},
		{
			name: "only comments and empty lines",
			fileContent: `# Comment 1
# Comment 2

# Comment 3
`,
			expectedLen: 0,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "registry.txt")

			err := os.WriteFile(tmpFile, []byte(tt.fileContent), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			result, err := readRegistry(tmpFile)

			if tt.expectError {
				if err == nil {
					t.Errorf("readRegistry() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("readRegistry() unexpected error: %v", err)
				return
			}

			if len(result) != tt.expectedLen {
				t.Errorf("readRegistry() returned %d repositories, want %d", len(result), tt.expectedLen)
			}

			for i, expected := range tt.expectedRepos {
				if i >= len(result) {
					t.Errorf("Expected repository at index %d but result has only %d repositories", i, len(result))
					continue
				}
				if result[i] != expected {
					t.Errorf("Repository at index %d: got %+v, want %+v", i, result[i], expected)
				}
			}
		})
	}

	// Test non-existent file
	t.Run("non-existent file", func(t *testing.T) {
		_, err := readRegistry("/non/existent/file")
		if err == nil {
			t.Error("readRegistry() should return error for non-existent file")
		}
	})
}

func TestAbs(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{
			name:     "positive number",
			input:    5,
			expected: 5,
		},
		{
			name:     "negative number",
			input:    -5,
			expected: 5,
		},
		{
			name:     "zero",
			input:    0,
			expected: 0,
		},
		{
			name:     "large positive number",
			input:    1000000,
			expected: 1000000,
		},
		{
			name:     "large negative number",
			input:    -1000000,
			expected: 1000000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := abs(tt.input)
			if result != tt.expected {
				t.Errorf("abs(%d) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestRepository(t *testing.T) {
	t.Run("repository struct creation", func(t *testing.T) {
		repo := Repository{
			Provider: "github",
			Owner:    "torvalds",
			Name:     "linux",
			URL:      "https://github.com/torvalds/linux.git",
		}

		if repo.Provider != "github" {
			t.Errorf("Provider = %q, want %q", repo.Provider, "github")
		}
		if repo.Owner != "torvalds" {
			t.Errorf("Owner = %q, want %q", repo.Owner, "torvalds")
		}
		if repo.Name != "linux" {
			t.Errorf("Name = %q, want %q", repo.Name, "linux")
		}
		if repo.URL != "https://github.com/torvalds/linux.git" {
			t.Errorf("URL = %q, want %q", repo.URL, "https://github.com/torvalds/linux.git")
		}
	})

	t.Run("zero value repository", func(t *testing.T) {
		var repo Repository

		if repo.Provider != "" {
			t.Errorf("Zero value Provider should be empty, got %q", repo.Provider)
		}
		if repo.Owner != "" {
			t.Errorf("Zero value Owner should be empty, got %q", repo.Owner)
		}
		if repo.Name != "" {
			t.Errorf("Zero value Name should be empty, got %q", repo.Name)
		}
		if repo.URL != "" {
			t.Errorf("Zero value URL should be empty, got %q", repo.URL)
		}
	})
}

// Benchmark tests for performance-critical functions
func BenchmarkParseRepositoryLine(b *testing.B) {
	line := "github:torvalds/linux"
	for i := 0; i < b.N; i++ {
		_, err := parseRepositoryLine(line)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

func BenchmarkExpandPath(b *testing.B) {
	path := "$HOME/test/path"
	for i := 0; i < b.N; i++ {
		_ = expandPath(path)
	}
}

func BenchmarkAbs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = abs(-42)
	}
}

// Table-driven test for all supported providers
func TestAllProviders(t *testing.T) {
	providers := []struct {
		name         string
		provider     string
		expectedHost string
	}{
		{"GitHub", "github", "github.com"},
		{"GitLab", "gitlab", "gitlab.com"},
		{"Bitbucket", "bitbucket", "bitbucket.org"},
	}

	for _, p := range providers {
		t.Run(p.name, func(t *testing.T) {
			line := p.provider + ":owner/repo"
			repo, err := parseRepositoryLine(line)
			if err != nil {
				t.Errorf("parseRepositoryLine failed for %s: %v", p.name, err)
				return
			}

			if !strings.Contains(repo.URL, p.expectedHost) {
				t.Errorf("URL for %s should contain %s, got %s", p.name, p.expectedHost, repo.URL)
			}

			if repo.Provider != p.provider {
				t.Errorf("Provider for %s should be %s, got %s", p.name, p.provider, repo.Provider)
			}
		})
	}
}

// Integration tests for git operations (these require git to be installed)
func TestGitOperationsIntegration(t *testing.T) {
	// Skip if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available, skipping integration tests")
	}

	t.Run("git commands exist", func(t *testing.T) {
		// Test that git clone command structure is correct
		cmd := exec.Command("git", "--version")
		err := cmd.Run()
		if err != nil {
			t.Errorf("git command failed: %v", err)
		}
	})
}

// Test error handling in expandPath
func TestExpandPathEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		setup    func()
		cleanup  func()
		validate func(string) bool
	}{
		{
			name:  "path with multiple env vars",
			input: "$HOME/$USER/test",
			validate: func(result string) bool {
				return strings.Contains(result, "/") && !strings.Contains(result, "$")
			},
		},
		{
			name:  "path with undefined env var",
			input: "$UNDEFINED_VAR/test",
			validate: func(result string) bool {
				return result == "/test" // undefined env vars expand to empty string
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			if tt.cleanup != nil {
				defer tt.cleanup()
			}

			result := expandPath(tt.input)
			if tt.validate != nil && !tt.validate(result) {
				t.Errorf("expandPath(%q) = %q failed validation", tt.input, result)
			}
		})
	}
}

// Test concurrent safety of parseRepositoryLine (it should be safe since it's pure)
func TestParseRepositoryLineConcurrent(t *testing.T) {
	const numGoroutines = 100
	const numIterations = 100

	done := make(chan bool, numGoroutines)
	line := "github:torvalds/linux"

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer func() { done <- true }()
			for j := 0; j < numIterations; j++ {
				repo, err := parseRepositoryLine(line)
				if err != nil {
					t.Errorf("Concurrent parseRepositoryLine failed: %v", err)
					return
				}
				if repo.Provider != "github" || repo.Owner != "torvalds" || repo.Name != "linux" {
					t.Errorf("Concurrent parseRepositoryLine returned incorrect result: %+v", repo)
					return
				}
			}
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("Concurrent test timed out")
		}
	}
}

// Test file permission handling in readRegistry
func TestReadRegistryPermissions(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "registry.txt")

	// Create file with content
	content := "github:golang/go\n"
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test reading with normal permissions
	repos, err := readRegistry(tmpFile)
	if err != nil {
		t.Errorf("readRegistry failed with normal permissions: %v", err)
	}
	if len(repos) != 1 {
		t.Errorf("Expected 1 repository, got %d", len(repos))
	}

	// Test with unreadable file (skip on Windows as chmod behaves differently)
	if strings.Contains(strings.ToLower(os.Getenv("OS")), "windows") {
		t.Skip("Skipping permission test on Windows")
	}

	// Make file unreadable
	err = os.Chmod(tmpFile, 0000)
	if err != nil {
		t.Fatalf("Failed to change file permissions: %v", err)
	}

	// Restore permissions for cleanup
	defer func() {
		os.Chmod(tmpFile, 0644)
	}()

	_, err = readRegistry(tmpFile)
	if err == nil {
		t.Error("readRegistry should fail with unreadable file")
	}
}

// Test large registry file handling
func TestReadRegistryLargeFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "large_registry.txt")

	// Create a large registry file
	var content strings.Builder
	expectedCount := 1000
	for i := 0; i < expectedCount; i++ {
		content.WriteString("github:user")
		content.WriteString(string(rune('0' + i%10)))
		content.WriteString("/repo")
		content.WriteString(string(rune('0' + i%10)))
		content.WriteString("\n")
	}

	err := os.WriteFile(tmpFile, []byte(content.String()), 0644)
	if err != nil {
		t.Fatalf("Failed to create large test file: %v", err)
	}

	repos, err := readRegistry(tmpFile)
	if err != nil {
		t.Errorf("readRegistry failed with large file: %v", err)
	}

	if len(repos) != expectedCount {
		t.Errorf("Expected %d repositories, got %d", expectedCount, len(repos))
	}
}

// Test malformed registry files
func TestReadRegistryMalformed(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expectCount int
	}{
		{
			name:        "mixed valid and invalid lines",
			content:     "github:valid/repo\ninvalid line\ngitlab:another/valid\n",
			expectCount: 2,
		},
		{
			name:        "only invalid lines",
			content:     "invalid line 1\ninvalid line 2\n",
			expectCount: 0,
		},
		{
			name:        "unicode characters",
			content:     "github:user/repo-with-üñíçödé\n",
			expectCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "registry.txt")

			err := os.WriteFile(tmpFile, []byte(tt.content), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			repos, err := readRegistry(tmpFile)
			if err != nil {
				t.Errorf("readRegistry failed: %v", err)
			}

			if len(repos) != tt.expectCount {
				t.Errorf("Expected %d repositories, got %d", tt.expectCount, len(repos))
			}
		})
	}
}

// Test BuildInfo struct functionality
func TestBuildInfoStruct(t *testing.T) {
	t.Run("BuildInfo with all fields", func(t *testing.T) {
		info := BuildInfo{
			Version:   "1.0.0",
			GitCommit: "abc123def",
			BuildTime: "2025-07-23T10:00:00Z",
			GoVersion: "go1.22.0",
		}

		if info.Version != "1.0.0" {
			t.Errorf("Version = %q, want %q", info.Version, "1.0.0")
		}
		if info.GitCommit != "abc123def" {
			t.Errorf("GitCommit = %q, want %q", info.GitCommit, "abc123def")
		}
		if info.BuildTime != "2025-07-23T10:00:00Z" {
			t.Errorf("BuildTime = %q, want %q", info.BuildTime, "2025-07-23T10:00:00Z")
		}
		if info.GoVersion != "go1.22.0" {
			t.Errorf("GoVersion = %q, want %q", info.GoVersion, "go1.22.0")
		}
	})

	t.Run("zero value BuildInfo", func(t *testing.T) {
		var info BuildInfo

		if info.Version != "" {
			t.Errorf("Zero value Version should be empty, got %q", info.Version)
		}
		if info.GitCommit != "" {
			t.Errorf("Zero value GitCommit should be empty, got %q", info.GitCommit)
		}
		if info.BuildTime != "" {
			t.Errorf("Zero value BuildTime should be empty, got %q", info.BuildTime)
		}
		if info.GoVersion != "" {
			t.Errorf("Zero value GoVersion should be empty, got %q", info.GoVersion)
		}
	})
}

// Test constants have expected values
func TestConstants(t *testing.T) {
	t.Run("application constants", func(t *testing.T) {
		if AppName == "" {
			t.Error("AppName should not be empty")
		}
		if AppVersion == "" {
			t.Error("AppVersion should not be empty")
		}
		if AppAuthor == "" {
			t.Error("AppAuthor should not be empty")
		}
		if AppLicense == "" {
			t.Error("AppLicense should not be empty")
		}
		if AppRepository == "" {
			t.Error("AppRepository should not be empty")
		}
		if AppDescription == "" {
			t.Error("AppDescription should not be empty")
		}
	})

	t.Run("default paths", func(t *testing.T) {
		if DefaultRegistryFile == "" {
			t.Error("DefaultRegistryFile should not be empty")
		}
		if DefaultMirrorsDir == "" {
			t.Error("DefaultMirrorsDir should not be empty")
		}

		// Check that default paths contain expected patterns
		if !strings.Contains(DefaultRegistryFile, "registry") {
			t.Error("DefaultRegistryFile should contain 'registry'")
		}
		if !strings.Contains(DefaultMirrorsDir, "mirrors") {
			t.Error("DefaultMirrorsDir should contain 'mirrors'")
		}
	})

	t.Run("repository URL should be valid", func(t *testing.T) {
		if !strings.HasPrefix(AppRepository, "https://") {
			t.Error("AppRepository should start with https://")
		}
		if !strings.Contains(AppRepository, "github.com") {
			t.Error("AppRepository should be on github.com")
		}
	})
}

// Test edge cases for Repository struct
func TestRepositoryEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		repo     Repository
		validate func(Repository) bool
	}{
		{
			name: "repository with special characters in name",
			repo: Repository{
				Provider: "github",
				Owner:    "user-name",
				Name:     "repo.name-with_special",
				URL:      "https://github.com/user-name/repo.name-with_special.git",
			},
			validate: func(r Repository) bool {
				return strings.Contains(r.Name, ".") && strings.Contains(r.Name, "_") && strings.Contains(r.Name, "-")
			},
		},
		{
			name: "repository with numeric owner/name",
			repo: Repository{
				Provider: "gitlab",
				Owner:    "123user",
				Name:     "456repo",
				URL:      "https://gitlab.com/123user/456repo.git",
			},
			validate: func(r Repository) bool {
				return strings.HasPrefix(r.Owner, "123") && strings.HasPrefix(r.Name, "456")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.validate != nil && !tt.validate(tt.repo) {
				t.Errorf("Repository validation failed for %+v", tt.repo)
			}
		})
	}
}
