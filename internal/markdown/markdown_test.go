package markdown

import (
	"strings"
	"testing"
)

func TestRenderHTMLSanitizesXSS(t *testing.T) {
	src := "# Hello\n\n<script>alert('xss')</script>\n\n[ok](https://example.com)"
	out := RenderHTML(src)
	if strings.Contains(out, "<script") {
		t.Fatalf("expected <script> to be stripped, got: %s", out)
	}
	if !strings.Contains(out, "<h1>Hello</h1>") {
		t.Fatalf("expected <h1>Hello</h1>, got: %s", out)
	}
	if !strings.Contains(out, `href="https://example.com"`) {
		t.Fatalf("expected safe link preserved, got: %s", out)
	}
}

func TestRenderPreviewHTMLRuneSafe(t *testing.T) {
	src := "áéíóúáéíóú"
	out := RenderPreviewHTML(src, 4)
	if !strings.Contains(out, "áéíó...") {
		t.Fatalf("expected rune-safe truncation 'áéíó...', got: %s", out)
	}
}
