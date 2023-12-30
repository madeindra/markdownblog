package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/madeindra/markdownblog/utils"
	"github.com/urfave/cli/v2"
)

const (
	flagRepo   = "repo"
	flagBranch = "branch"
	flagToken  = "token"
	flagOut    = "out"
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
			&cli.StringFlag{
				Name:    flagOut,
				Aliases: []string{"o"},
				Usage:   "Output directory for result",
				Value:   "generated",
			},
		},
		// validate required params
		Before: func(c *cli.Context) error {
			// ensure all required parameters are not empty
			err := utils.EnsureParams(c.String(flagRepo), c.String(flagOut))
			if err != nil {
				return err
			}

			return nil
		},
		// run the app
		Action: func(c *cli.Context) error {
			// get git username and repo name from git repository url
			gitDetail, err := utils.GetGitDetail(c.String(flagRepo))
			if err != nil {
				return err
			}

			// trim space from parameters
			gitBranch := strings.TrimSpace(c.String(flagBranch))
			gitToken := strings.TrimSpace(c.String(flagToken))
			outDir := strings.TrimSpace(c.String(flagOut))

			// check whether the git repository is private or public by the presence of git token
			isPrivate := !utils.IsEmptyString(gitToken)

			// read all files from github or gitlab repository according to the parameters (repo, category, branch)
			files, err := utils.FetchFiles(gitDetail, gitBranch, gitToken, isPrivate)
			if err != nil {
				return err
			}

			// ready to go, print welcome message
			fmt.Println("Welcome to Markdown Blog generator")

			return generateBlog(files, outDir)
		},
	}
}

func generateBlog(files []utils.File, outDir string) error {
	// loop through all files
	for _, file := range files {
		// read content of files by downloading content
		content, err := utils.DownloadContent(file.URL)
		if err != nil {
			return err
		}

		// convert each markdown into html
		result := utils.MarkdownToHTML(content)

		// !DELETE THIS!
		fmt.Println(result)

		// TODO: put each parsed file into the templates according to the parameters (theme)
	}

	// TODO: create a directory for the result (outdir)

	// TODO: put all as html files into directory according to the parameters (outdir, category)

	// TODO: create index.html as homepage according to the parameters (theme, title, etc)
	return nil
}
