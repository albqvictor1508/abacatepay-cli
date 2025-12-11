package utils

import (
	"os/exec"
	"strings"
)

func GetGitRepo() string {
	output, err := exec.Command("git", "remote", "get-url", "origin").Output()

	if err != nil {
		return ""
	}

	url := strings.TrimSuffix(strings.TrimSpace(string(output)), ".git")

	// SSH style: git@github.com:user/repo
	if idx := strings.LastIndex(url, ":"); idx != -1 && !strings.Contains(url[:idx], "/") {
		url = url[idx+1:]
	}

	parts := strings.Split(url, "/")

	return parts[len(parts)-1]
}
