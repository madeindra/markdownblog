package main

import (
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/madeindra/markdownblog/utils"
	"github.com/urfave/cli/v2"
)

const (
	flagRepo     = "repo"
	flagBranch   = "branch"
	flagToken    = "token"
	flagOut      = "out"
	flagTheme    = "theme"
	flagBlogName = "name"
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
			&cli.StringFlag{
				Name:    flagTheme,
				Aliases: []string{"th"},
				Usage:   "Theme for your blog (template for generation)",
				Value:   "examples",
			},
			&cli.StringFlag{
				Name:     flagBlogName,
				Aliases:  []string{"n"},
				Usage:    "Name of your blog (for title)",
				Required: true,
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
			themeName := strings.TrimSpace(c.String(flagTheme))
			blogName := strings.TrimSpace(c.String(flagBlogName))

			// check whether the git repository is private or public by the presence of git token
			isPrivate := !utils.IsEmptyString(gitToken)

			// read all files from github or gitlab repository according to the parameters (repo, category, branch)
			files, err := utils.FetchFiles(gitDetail, gitBranch, gitToken, isPrivate)
			if err != nil {
				return err
			}

			// ready to go, print welcome message
			fmt.Println("Welcome to Markdown Blog generator")

			return generateBlog(files, blogName, themeName, outDir)
		},
		After: func(c *cli.Context) error {
			// print success message
			fmt.Println("Blog generated successfully")

			return nil
		},
	}
}

// generateBlog is a function to generate blog from markdown files
func generateBlog(files []utils.File, blogName, themeName, outDir string) error {
	// remove the output directory if exists
	err := os.RemoveAll(outDir)
	if err != nil {
		return err
	}

	// create a new output directory
	err = os.MkdirAll(outDir, os.ModePerm)
	if err != nil {
		return err
	}

	// loop through all files
	posts := make([]utils.Post, 0)
	for _, file := range files {
		// read content of files by downloading content
		content, err := utils.DownloadContent(file.URL)
		if err != nil {
			return err
		}

		// convert each markdown into html
		result := utils.MarkdownToHTML(content)

		// put each parsed file into the templates according
		filepath := path.Join("themes", themeName, "templates", "template.html")
		tmpl, err := template.ParseFiles(filepath)
		if err != nil {
			return err
		}

		// create template data
		data := map[string]string{
			"title":    utils.CreateTitle(file.Name),
			"contents": utils.StringifyHTML(result),
			"name":     blogName,
		}

		// create file writer
		filename := fmt.Sprintf("%s.html", strings.TrimSuffix(file.Name, ".md"))
		newFile, err := os.Create(path.Join(outDir, filename))
		if err != nil {
			return err
		}
		defer newFile.Close()

		// execute template and write the result into the new file
		err = tmpl.Execute(newFile, data)
		if err != nil {
			return err
		}

		// append to posts
		posts = append(posts, utils.Post{
			Data:     result,
			Filename: filename,
		})
	}

	// copy assets (css, js, etc) into the output directory
	err = utils.CopyDirectory(path.Join("themes", themeName, "assets"), path.Join(outDir, "assets"))
	if err != nil {
		return err
	}

	// create index.html as homepage according to the parameters
	generateIndex(posts, blogName, themeName, outDir)

	return nil
}

// generateIndex is a function to generate index.html as homepage
func generateIndex(posts []utils.Post, blogName, themeName, outDir string) error {
	// put each parsed file into the templates according
	filepath := path.Join("themes", themeName, "templates", "template.html")
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		return err
	}

	summaries := utils.CreateSummaries(posts)

	// create template data
	data := map[string]interface{}{
		"title":    blogName,
		"name":     blogName,
		"contents": summaries,
	}

	// create file writer
	filename := "index.html"
	newFile, err := os.Create(path.Join(outDir, filename))
	if err != nil {
		return err
	}
	defer newFile.Close()

	// execute template and write the result into the new file
	err = tmpl.Execute(newFile, data)
	if err != nil {
		return err
	}

	return nil
}
