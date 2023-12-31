package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"

	gohtml "html"
)

// Post is a struct to store parsed markdown content
type Post struct {
	Data     []byte
	Filename string
}

// MarkdownToHTML convert markdown content into html
func MarkdownToHTML(content []byte) []byte {
	// initialize parser
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	mdParser := parser.NewWithExtensions(extensions)

	// parse markdown content
	res := mdParser.Parse(content)

	// initialize HTML renderer
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	htmlRenderer := html.NewRenderer(opts)

	return markdown.Render(res, htmlRenderer)
}

// StringifyHTML convert html content into string and unescape it
func StringifyHTML(content []byte) string {
	return gohtml.UnescapeString((string(content)))
}

// Cleanup remove all headings (including contents) and strip all html tags from html string
func Cleanup(content string) string {
	content = removeHeadings(content)
	content = stripTags(content)

	return content
}

// removeHeadings remove all headings from html string
func removeHeadings(input string) string {
	pattern := regexp.MustCompile("<h[1-6][^>]*>.*?</h[1-6]>")
	return pattern.ReplaceAllString(input, "")
}

// stripTags strip all html tags from html string
func stripTags(input string) string {
	pattern := regexp.MustCompile("<[^>]*>")
	return pattern.ReplaceAllString(input, "")
}

// CreateSummaries is a function to create summaries elements for posts for homepage
func CreateSummaries(posts []Post) string {
	var builder strings.Builder
	for i, post := range posts {
		// get title from the post
		title := getTitle(string(post.Data))

		// get summary from the post
		summary := getSummary(string(post.Data))

		// create summaries element
		builder.WriteString(fmt.Sprintf(`
			<h2>%s</h2>
			<p>%s</p>
			<a href="%s">Read more...</a>
		`, title, summary, post.Filename))

		// add <hr> element if not the last post
		if i != len(posts)-1 {
			builder.WriteString("<hr>")
		}
	}

	return builder.String()
}

// getTitle return the title from first heading recursively (h1 to h4), if not found return Untitled
func getTitle(content string) string {
	// unescape html content
	content = gohtml.UnescapeString(content)

	// find heading recursively until found (h1 to h4)
	for i := 1; i <= 4; i++ {
		// create heading tag
		tag := fmt.Sprintf("h%d", i)

		// find opening heading tag (regardless of attributes)
		openStart := strings.Index(content, fmt.Sprintf("<%s", tag))
		openEnd := strings.Index(content[openStart:], ">") + openStart

		// find closing heading tag
		end := strings.Index(content[openEnd:], fmt.Sprintf("</%s>", tag)) + openEnd

		// if found, return the title
		if openStart != -1 && openEnd != -1 && end != -1 {
			return content[openEnd+1 : end]
		}
	}

	// if not found, return Untitled
	return "Untitled"
}

// getSummary return the first paragraph from html string, if not found return No summary
func getSummary(content string) string {
	// unescape html content
	content = gohtml.UnescapeString(content)

	// find first paragraph
	openStart := strings.Index(content, "<p>")
	openEnd := strings.Index(content[openStart:], "</p>") + openStart

	// if found, return the title
	if openStart != -1 && openEnd != -1 {
		return content[openStart+3 : openEnd]
	}

	// if not found, return Untitled
	return "No summary"
}
