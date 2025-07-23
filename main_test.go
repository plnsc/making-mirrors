package main

import (
	"testing"
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
