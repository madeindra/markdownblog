package utils

import (
	"fmt"
	"strings"
)

// GitDetail holds git username and repo name
type GitDetail struct {
	Username string
	Type     Git
	RepoName string
}

// Git holds git type
type Git string

// Git types
const (
	Github Git = "github.com"
	Gitlab Git = "gitlab.com"
)

// File holds file detail
type File struct {
	Name string
	URL  string
	Type string
}

// GetGitDetail return git username and repo name from git repository url
// this function accepts github and gitlab repository url
func GetGitDetail(repo string) (GitDetail, error) {
	// trim space
	repo = strings.TrimSpace(repo)

	// remove slash from the end of the url
	repo = strings.TrimSuffix(repo, "/")

	// remove .git from the end of the url
	repo = strings.TrimSuffix(repo, ".git")

	// remove http:// and https:// from the beginning of the url
	repo = strings.TrimPrefix(repo, "http://")
	repo = strings.TrimPrefix(repo, "https://")

	// split the url by slash, should have at 3 parts "github.com/username/repo" or "gitlab.com/username/repo"
	s := strings.Split(repo, "/")
	if len(s) != 3 {
		return GitDetail{}, fmt.Errorf("error: invalid git repository url")
	}

	// check whether the first part contains github or gitlab
	var gitType Git
	switch s[0] {
	case string(Github):
		gitType = Github
	case string(Gitlab):
		gitType = Gitlab
	default:
		return GitDetail{}, fmt.Errorf("error: repository is not supported")
	}

	return GitDetail{
		Type:     gitType,
		Username: s[1],
		RepoName: s[2],
	}, nil
}

// fetchFiles reads all files from github or gitlab repository according to the parameters (repo, category, branch)
func FetchFiles(gitDetail GitDetail, branch, token string, isPrivate bool) ([]File, error) {
	switch gitDetail.Type {
	case Github:
		return fetchGithubFiles(gitDetail, branch, token, isPrivate)
	case Gitlab:
		return fetchGitlabFiles(gitDetail, branch, token, isPrivate)
	default:
		return nil, fmt.Errorf("error: repository is not supported")
	}
}

// EnsureParams checks whether all required parameters are provided
func EnsureParams(v ...string) error {
	for _, s := range v {
		if strings.TrimSpace(s) == "" {
			return fmt.Errorf("error: missing required parameter(s)")
		}
	}
	return nil
}

// IsEmptyString checks whether a string is empty or not
func IsEmptyString(s string) bool {
	return strings.TrimSpace(s) == ""
}

// IsMarkdownFile checks whether a file is markdown or not
func IsMarkdownFile(filename string) bool {
	return strings.HasSuffix(filename, ".md")
}
