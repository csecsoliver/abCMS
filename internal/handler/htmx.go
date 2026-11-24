package handler

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/a-h/templ"

	"github.com/csecsoliver/abCMS/internal/view"
)

func isHTMX(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

func isHistoryRestore(r *http.Request) bool {
	return r.Header.Get("HX-History-Restore-Request") == "true"
}

func isPartialHTMX(r *http.Request) bool {
	return isHTMX(r) && !isHistoryRestore(r)
}

// hxTarget extracts the element ID from HTMX 4's HX-Target header ("tag#id", e.g.
// "div#dashboard-content"), while also accepting bare IDs for compatibility.
func hxTarget(r *http.Request) string {
	raw := strings.TrimSpace(r.Header.Get("HX-Target"))
	if idx := strings.LastIndex(raw, "#"); idx >= 0 {
		return raw[idx+1:]
	}
	return raw
}

func setHTMXVary(w http.ResponseWriter) {
	w.Header().Add("Vary", "HX-Request")
	w.Header().Add("Vary", "HX-Target")
}

// renderPage renders content inside Layout for normal browser requests (and HTMX history
// restores), or as an SPA fragment with <template hx type="partial"> updates for HTMX requests.
func (h *Handler) renderPage(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	meta view.PageMeta,
	content templ.Component,
) {
	setHTMXVary(w)
	if isPartialHTMX(r) {
		h.renderComponent(w, r, status, view.HTMXPageFragment(meta, content))
		return
	}
	h.renderComponent(w, r, status, view.FullPage(meta, content))
}

// renderDashboardPage renders either only the inner #dashboard-content (when switching tabs
// inside the dashboard), the #app-main dashboard shell (when navigating from outside the
// dashboard via HTMX), or the complete HTML document for standard requests.
func (h *Handler) renderDashboardPage(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	meta view.PageMeta,
	activeTab string,
	errorMessage string,
	inner templ.Component,
) {
	setHTMXVary(w)
	if isPartialHTMX(r) {
		if hxTarget(r) == "dashboard-content" {
			h.renderComponent(w, r, status, view.HTMXDashboardFragment(meta, activeTab, errorMessage, inner))
			return
		}
		h.renderComponent(w, r, status, view.HTMXPageFragment(meta, view.DashboardShell(activeTab, errorMessage, inner)))
		return
	}
	h.renderComponent(w, r, status, view.FullPage(meta, view.DashboardShell(activeTab, errorMessage, inner)))
}

// broadcastFeedAndStats pushes live header stats as an HTMX 4 <template hx type="partial">
// message and emits a named posts-changed event over SSE.
func (h *Handler) broadcastFeedAndStats(r *http.Request) {
	if h.events == nil {
		return
	}
	meta := h.pageMeta(r, "", "")
	var buf bytes.Buffer
	if err := view.HeaderStatsPartial(meta).Render(context.Background(), &buf); err != nil {
		log.Printf("render SSE stats partial: %v", err)
	} else {
		// Unnamed SSE messages are swapped directly by HTMX 4 hx-sse via <template hx type="partial">.
		h.events.Publish("", buf.String())
	}
	// Named SSE events are dispatched as DOM CustomEvents on #sse-bridge (bubbling to body).
	h.events.Publish("posts-changed", "refresh")
}
