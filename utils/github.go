package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GithubFile holds github file detail
type GithubFile struct {
	Name string `json:"name"`
	URL  string `json:"download_url"`
	Type string `json:"type"`
}

// fetchGithubFiles reads all files from github repository according to the parameters (repo, category, branch)
func fetchGithubFiles(gitDetail GitDetail, branch, token string, isPrivate bool) ([]File, error) {
	// static github api url
	url := "https://api.github.com/repos/%s/%s/contents"

	// prepare http get request
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(url, gitDetail.Username, gitDetail.RepoName), nil)
	if err != nil {
		return nil, fmt.Errorf("error: failed to prepare github http request")
	}

	// set authorization header if the repository is private
	if isPrivate {
		// make sure the token is not empty
		if IsEmptyString(token) {
			return nil, fmt.Errorf("error: missing authorization token for private github repository")
		}

		// add authorization header
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	// send http request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error: failed to send github http request")
	}

	// check whether the response is ok
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: failed to fetch github files")
	}

	// parse the response body to json
	files := []GithubFile{}
	err = json.NewDecoder(resp.Body).Decode(&files)
	if err != nil {
		return nil, fmt.Errorf("error: failed to parse github response body")
	}

	return translateGithubFile(files), nil
}

// translateGithubFile translates github file type to file type for easier use
func translateGithubFile(files []GithubFile) []File {
	// return nil if files is nil
	if files == nil {
		return nil
	}

	// translate github file to file
	result := []File{}
	for _, f := range files {
		result = append(result, File(f))
	}

	return result
}
