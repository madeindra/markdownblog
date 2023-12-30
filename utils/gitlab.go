package utils

import "fmt"

// fetchGitlabFiles reads all files from gitlab repository according to the parameters (repo, category, branch)
func fetchGitlabFiles(gitDetail GitDetail, branch, token string, isPrivate bool) ([]File, error) {
	// TODO: add Gitlab support
	// return error for now
	return nil, fmt.Errorf("error: gitlab is not supported yet")
}
