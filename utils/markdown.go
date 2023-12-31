package utils

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"

	gohtml "html"
)

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
