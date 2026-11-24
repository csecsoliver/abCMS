package handler

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/a-h/templ"

	"github.com/csecsoliver/abCMS/internal/auth"
	"github.com/csecsoliver/abCMS/internal/config"
	"github.com/csecsoliver/abCMS/internal/db"
	"github.com/csecsoliver/abCMS/internal/markdown"
	"github.com/csecsoliver/abCMS/internal/media"
	"github.com/csecsoliver/abCMS/internal/view"
)

type contextKey int

const userContextKey contextKey = 1

// Handler wires HTTP endpoints to the database, session manager, media service, and SSE broker.
type Handler struct {
	cfg      config.Config
	store    *db.Store
	sessions *auth.SessionManager
	media    *media.Service
	events   *EventBroker
}

// New constructs a Handler with all required dependencies.
func New(cfg config.Config, store *db.Store, sessions *auth.SessionManager, mediaSvc *media.Service) *Handler {
	return &Handler{
		cfg:      cfg,
		store:    store,
		sessions: sessions,
		media:    mediaSvc,
		events:   NewEventBroker(),
	}
}

// Close shuts down active background resources such as SSE streams.
func (h *Handler) Close() {
	if h.events != nil {
		h.events.Close()
	}
}

// Routes returns the configured HTTP handler wrapped with CSRF protection and panic recovery.
func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", h.handleIndex)
	mux.HandleFunc("GET /posts/{post_id}", h.handlePost)
	mux.HandleFunc("GET /events", h.handleEvents)
	mux.HandleFunc("GET /login", h.handleLoginPage)
	mux.HandleFunc("POST /login", h.handleLoginSubmit)
	mux.HandleFunc("GET /logout", h.handleLogout)
	mux.HandleFunc("POST /logout", h.handleLogout)

	mux.HandleFunc("GET /dashboard", h.requireAuth(h.handleDashboardHome))
	mux.HandleFunc("GET /dashboard/", h.requireAuth(h.handleDashboardHome))
	mux.HandleFunc("GET /dashboard/posts", h.requireAuth(h.handleDashboardPosts))
	mux.HandleFunc("GET /dashboard/posts/add", h.requireAuth(h.handleDashboardPostAdd))
	mux.HandleFunc("GET /dashboard/posts/edit/{post_id}", h.requireAuth(h.handleDashboardPostEdit))

	mux.HandleFunc("POST /formapi/posts/add", h.requireAuth(h.handlePostCreate))
	mux.HandleFunc("POST /formapi/posts/mod/{post_id}", h.requireAuth(h.handlePostUpdate))
	mux.HandleFunc("POST /formapi/posts/delete/{post_id}", h.requireAuth(h.handlePostDelete))

	staticFS := http.FileServer(http.Dir(h.cfg.StaticDir))
	mux.Handle("GET /static/", http.StripPrefix("/static/", staticFS))
	mux.HandleFunc("GET /favicon.ico", h.handleFavicon)
	mux.HandleFunc("/", h.handleNotFound)

	csrfProtected := http.NewCrossOriginProtection().Handler(mux)
	return h.recoverMiddleware(csrfProtected)
}

func (h *Handler) handleFavicon(w http.ResponseWriter, r *http.Request) {
	faviconPath := filepath.Join(h.cfg.StaticDir, "favicon.ico")
	if _, err := os.Stat(faviconPath); err != nil {
		h.handleNotFound(w, r)
		return
	}
	http.ServeFile(w, r, faviconPath)
}

func (h *Handler) recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic handling %s %s: %v\n%s", r.Method, r.URL.Path, rec, debug.Stack())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := h.authenticatedUser(r)
		if user == nil {
			reqURI := r.URL.RequestURI()
			target := "/login"
			returnTo := ""
			if r.Method == http.MethodGet && reqURI != "" && reqURI != "/dashboard" {
				returnTo = reqURI
				target += "?return_to=" + url.QueryEscape(reqURI)
			}
			if isPartialHTMX(r) {
				w.Header().Set("HX-Retarget", "#app-main")
				w.Header().Set("HX-Reswap", "innerHTML show:top showTarget:window")
				w.Header().Set("HX-Push-Url", target)
				h.renderPage(w, r, http.StatusOK, h.pageMeta(r, "Login", "login"), view.LoginContent(returnTo, ""))
				return
			}
			http.Redirect(w, r, target, http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(), userContextKey, user)
		next(w, r.WithContext(ctx))
	}
}

func (h *Handler) authenticatedUser(r *http.Request) *db.User {
	if u, ok := r.Context().Value(userContextKey).(*db.User); ok && u != nil {
		return u
	}
	username, ok := h.sessions.GetUsername(r)
	if !ok {
		return nil
	}
	user, err := h.store.FindUserByUsername(r.Context(), username)
	if err != nil {
		return nil
	}
	return user
}

func (h *Handler) pageMeta(r *http.Request, title, activeNav string) view.PageMeta {
	var currentUsername string
	if u := h.authenticatedUser(r); u != nil {
		currentUsername = u.Username
	}
	return h.pageMetaForUser(currentUsername, title, activeNav)
}

func (h *Handler) pageMetaForUser(username, title, activeNav string) view.PageMeta {
	dbSize, hasDBSize := media.DBSizeBytes(h.cfg.DBPath)
	return view.PageMeta{
		Title:         title,
		CurrentUser:   username,
		ActiveNav:     activeNav,
		FreeSpaceText: media.FreeSpaceText(h.cfg.UploadDir),
		DBSizeBytes:   dbSize,
		HasDBSize:     hasDBSize,
	}
}

func (h *Handler) renderComponent(w http.ResponseWriter, r *http.Request, status int, c templ.Component) {
	var buf bytes.Buffer
	if err := c.Render(r.Context(), &buf); err != nil {
		log.Printf("render error on %s: %v", r.URL.Path, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(buf.Bytes())
}

func safeReturnTo(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" || !strings.HasPrefix(trimmed, "/") || strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, `/\`) {
		return "/dashboard"
	}
	return trimmed
}

func toPostCards(posts []db.Post, maxTextRunes, maxImageRunes int) []view.PostCard {
	cards := make([]view.PostCard, 0, len(posts))
	for _, p := range posts {
		limit := maxTextRunes
		if p.HasImage {
			limit = maxImageRunes
		}
		cards = append(cards, view.PostCard{
			ID:             p.ID,
			Title:          p.Title,
			AuthorUsername: p.AuthorUsername,
			CreatedAt:      p.CreatedAt,
			HasImage:       p.HasImage,
			Path:           p.Path,
			ThumbnailPath:  p.ThumbnailPath,
			PreviewHTML:    markdown.RenderPreviewHTML(p.Content, limit),
		})
	}
	return cards
}
