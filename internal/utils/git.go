package utils

import (
	"os/exec"
	"strings"
)

func GetGitRepo() string {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	url := strings.TrimSpace(string(output))
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		repoName := parts[len(parts)-1]
		repoName = strings.TrimSuffix(repoName, ".git")

		return repoName
	}

	return ""
}
