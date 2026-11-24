package markdown

import (
	"bytes"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
)

var (
	mdEngine  = goldmark.New()
	sanitizer = bluemonday.UGCPolicy()
)

// RenderHTML converts Markdown source into sanitized HTML safe for direct rendering.
func RenderHTML(src string) string {
	var buf bytes.Buffer
	if err := mdEngine.Convert([]byte(src), &buf); err != nil {
		return sanitizer.Sanitize(src)
	}
	return sanitizer.Sanitize(buf.String())
}

// RenderPreviewHTML truncates Markdown safely on rune boundaries and renders sanitized HTML.
func RenderPreviewHTML(src string, maxRunes int) string {
	if maxRunes > 0 {
		runes := []rune(src)
		if len(runes) > maxRunes {
			src = string(runes[:maxRunes]) + "..."
		}
	}
	return RenderHTML(src)
}
