package htmltemplate

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	blackfriday "github.com/russross/blackfriday/v2"
	"jaytaylor.com/html2text"
)

// PlaintextBlurb returns a plaintext blurb from markdown content.
func PlaintextBlurb(s string) string {
	unsafeHTML := blackfriday.Run([]byte(s))
	plaintext, err := html2text.FromString(string(unsafeHTML))
	if err != nil {
		plaintext = s
	}
	period := strings.Index(plaintext, ". ")
	if period > 0 {
		plaintext = plaintext[:period+1]
	}

	return plaintext
}

// sanitizedContent returns a sanitized content block or an error is one occurs.
func (te *Engine) sanitizedContent(t *template.Template, content string) (*template.Template, error) {
	b := []byte(content)
	// Ensure unit line endings are used when pulling out of JSON.
	markdownWithUnixLineEndings := strings.Replace(string(b), "\r\n", "\n", -1)
	htmlCode := blackfriday.Run([]byte(markdownWithUnixLineEndings))

	// Sanitize by removing HTML if true.
	if !te.allowUnsafeHTML {
		htmlCode = bluemonday.UGCPolicy().SanitizeBytes(htmlCode)
	}

	// Change delimiters temporarily so code samples can use Go blocks.
	safeContent := fmt.Sprintf(`[{[{define "content"}]}]%s[{[{end}]}]`, string(htmlCode))
	t = t.Delims("[{[{", "}]}]")
	var err error
	t, err = t.Parse(safeContent)
	if err != nil {
		return nil, err
	}
	// Reset delimiters
	t = t.Delims("{{", "}}")
	return t, nil
}

// sanitized returns a sanitized content block or an error is one occurs.
func (te *Engine) sanitized(content string) string {
	b := []byte(content)
	// Ensure unit line endings are used when pulling out of JSON.
	markdownWithUnixLineEndings := strings.Replace(string(b), "\r\n", "\n", -1)
	htmlCode := blackfriday.Run([]byte(markdownWithUnixLineEndings))

	// Sanitize by removing HTML if true.
	if !te.allowUnsafeHTML {
		htmlCode = bluemonday.UGCPolicy().SanitizeBytes(htmlCode)
	}

	return string(htmlCode)
}
