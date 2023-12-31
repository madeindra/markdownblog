package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
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

// DownloadContent downloads file from url and return the content
func DownloadContent(url string) ([]byte, error) {
	// download file from url
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error: failed to download file")
	}
	defer resp.Body.Close()

	// read the content of the file
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error: failed to read file")
	}

	return content, nil
}

// CopyDirectory copies a directory and all its contents to a new directory
func CopyDirectory(src, dest string) error {
	// open the source directory
	srcDir, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("error: failed to open source directory: %v", err)
	}
	defer srcDir.Close()

	// create the destination directory
	err = os.MkdirAll(dest, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error: failed to create destination directory: %v", err)
	}

	// read all files and subdirectories in the source directory
	fileInfos, err := srcDir.Readdir(-1)
	if err != nil {
		return fmt.Errorf("error: failed to read source directory: %v", err)
	}

	// copy each file and subdirectory to the destination directory
	for _, fileInfo := range fileInfos {
		srcPath := path.Join(src, fileInfo.Name())
		destPath := path.Join(dest, fileInfo.Name())

		if fileInfo.IsDir() {
			// recursively copy subdirectories
			err = CopyDirectory(srcPath, destPath)
			if err != nil {
				return fmt.Errorf("error: failed to copy subdirectory: %v", err)
			}
		} else {
			// copy regular files
			err = CopyFile(srcPath, destPath)
			if err != nil {
				return fmt.Errorf("error: failed to copy file: %v", err)
			}
		}
	}

	return nil
}

// CopyFile copies a file from source to destination
func CopyFile(src, dest string) error {
	// open the source file
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("error: failed to open source file: %v", err)
	}
	defer srcFile.Close()

	// create the destination file
	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("error: failed to create destination file: %v", err)
	}
	defer destFile.Close()

	// copy the contents of the source file to the destination file
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return fmt.Errorf("error: failed to copy file contents: %v", err)
	}

	return nil
}

// CreateTitle convert from file name to title
func CreateTitle(filename string) string {
	// early return if filename is empty
	if filename == "" {
		return ""
	}

	// strip .md from the end of the filename
	title := strings.TrimSuffix(filename, ".md")

	// replace all hyphens or underscore with space
	title = strings.ReplaceAll(title, "-", " ")
	title = strings.ReplaceAll(title, "_", " ")

	// capitalize each word
	return strings.Title(title)
}
