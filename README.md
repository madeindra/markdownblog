# Markdown Blog
Generate static blog from markdown files in your git repository.

## Usage
1. Clone this repo
2. Install dependencies
```bash
go mod tidy
```
3. Run the CLI tool
```bash
go run main.go -r <repo-url> -t <token> -o <output-dir> -n <blog-name> -th <theme-name>
```

CLI Flags:
- `-h` or `--help` : Show help
- `-r` or `--repo` : The URL of the git repository (Github only, must be a public repo if no token is provided)
- `-t` or `--token` : The access token for the git repository
- `-o` or `--output` : The output directory for the generated blog
- `-n` or `--name` : The name of the blog, will be used as the title of the blog
- `-th` or `--theme` : The name of the theme to be used for the blog (will find this in the `themes` directory) 

## Adding Custom Theme
1. Create a new folder in `themes` directory. The name of the folder will be the name of the theme.
2. Inside the theme folder, create a `templates` folder and an `assets` folder.
3. Create a `template.html` file in the `templates` folder. (See `themes/examples/template.html` for reference)
4. Put all your static files (css, js, images) in the `assets` folder.

## TODO
- [X] Generate from a Github repo (public & private)
- [ ] Generate from a Gitlab repo (public & private)
- [ ] Branch selection
- [ ] Generate from a source directory
- [ ] Add more themes 