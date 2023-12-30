package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

const (
	flagRepo   = "repo"
	flagBranch = "branch"
	flagToken  = "token"
)

func main() {
	if err := App().Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func App() *cli.App {
	// create cli app
	return &cli.App{
		// set app name and usage
		Name:    "Markdown Blog",
		Version: "1.0.0",
		Usage:   "Generate static blog from markdown files in a git repository",
		Authors: []*cli.Author{
			{
				Name:  "Made Indra",
				Email: "made.indra@pm.me",
			},
		},
		// accept params (repo, branch, token) from flags
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     flagRepo,
				Aliases:  []string{"r"},
				Usage:    "URL of the git repository",
				Required: true,
			},
			&cli.StringFlag{
				Name:    flagBranch,
				Aliases: []string{"b"},
				Usage:   "Branch of your git repository",
				Value:   "master",
			},
			&cli.StringFlag{
				Name:    flagToken,
				Aliases: []string{"t"},
				Usage:   "Token for private repository",
			},
		},
		// validate required params
		Before: func(c *cli.Context) error {
			// ensure all required parameters are provided
			err := ensureParams(c.String(flagRepo))
			if err != nil {
				return err
			}

			return nil
		},
		// run the app
		Action: func(c *cli.Context) error {
			// get git username and repo name from git repository url
			gitDetail, err := getGitDetail(c.String(flagRepo))
			if err != nil {
				return err
			}

			// trim space from parameters
			gitBranch := strings.TrimSpace(c.String(flagBranch))
			gitToken := strings.TrimSpace(c.String(flagToken))

			// check whether the git repository is private or public by the presence of git token
			isPrivate := !isEmptyString(gitToken)

			// read all files from github or gitlab repository according to the parameters (repo, category, branch)
			files, err := fetchFiles(gitDetail, gitBranch, gitToken, isPrivate)
			if err != nil {
				return err
			}

			// ready to go, print welcome message
			fmt.Println("Welcome to Markdown Blog generator")

			// !DELETE LATER
			fmt.Println(files)

			// TODO: parse each markdown into html using gomarkdown

			// TODO: put each parsed file into the templates according to the parameters (theme)

			// TODO: create a directory for the result (outdir)

			// TODO: put all as html files into directory according to the parameters (outdir, category)

			// TODO: create index.html as homepage according to the parameters (theme, title, etc)
			return nil
		},
	}
}

// ensureParams checks whether all required parameters are provided
func ensureParams(v ...string) error {
	for _, s := range v {
		if strings.TrimSpace(s) == "" {
			return fmt.Errorf("error: missing required parameter(s)")
		}
	}
	return nil
}

// isEmptyString checks whether a string is empty or not
func isEmptyString(s string) bool {
	return strings.TrimSpace(s) == ""
}

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

// getGitDetails return git username and repo name from git repository url
// accepts github and gitlab repository url
func getGitDetail(repo string) (GitDetail, error) {
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
func fetchFiles(gitDetail GitDetail, branch, token string, isPrivate bool) ([]File, error) {
	switch gitDetail.Type {
	case Github:
		return fetchGithubFiles(gitDetail, branch, token, isPrivate)
	case Gitlab:
		return fetchGitlabFiles(gitDetail, branch, token, isPrivate)
	default:
		return nil, fmt.Errorf("error: repository is not supported")
	}
}

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
		if isEmptyString(token) {
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

// fetchGitlabFiles reads all files from gitlab repository according to the parameters (repo, category, branch)
func fetchGitlabFiles(gitDetail GitDetail, branch, token string, isPrivate bool) ([]File, error) {
	// TODO: add Gitlab support
	// return error for now
	return nil, fmt.Errorf("error: gitlab is not supported yet")
}

// File holds file detail
type File struct {
	Name string
	URL  string
	Type string
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
