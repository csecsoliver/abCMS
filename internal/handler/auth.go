package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/csecsoliver/abCMS/internal/auth"
	"github.com/csecsoliver/abCMS/internal/db"
	"github.com/csecsoliver/abCMS/internal/view"
)

func (h *Handler) handleLoginPage(w http.ResponseWriter, r *http.Request) {
	if h.authenticatedUser(r) != nil {
		if isPartialHTMX(r) {
			w.Header().Set("HX-Retarget", "#app-main")
			w.Header().Set("HX-Reswap", "innerHTML show:top showTarget:window")
			w.Header().Set("HX-Push-Url", "/dashboard")
			meta := h.pageMeta(r, "Dashboard", "dashboard")
			h.renderPage(w, r, http.StatusOK, meta, view.DashboardShell("home", "", view.DashboardHomeSection()))
			return
		}
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	returnTo := r.URL.Query().Get("return_to")
	if isPartialHTMX(r) && hxTarget(r) != "app-main" {
		w.Header().Set("HX-Retarget", "#app-main")
		w.Header().Set("HX-Reswap", "innerHTML show:top showTarget:window")
	}
	meta := h.pageMeta(r, "Login", "login")
	h.renderPage(w, r, http.StatusOK, meta, view.LoginContent(returnTo, ""))
}

func (h *Handler) handleLoginSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderLoginError(w, r, http.StatusBadRequest, "", "Invalid form submission")
		return
	}

	option := r.FormValue("option")
	rawUsername := r.FormValue("username")
	password := r.FormValue("password")
	returnTo := r.FormValue("return_to")
	username := auth.Slugify(rawUsername)

	switch option {
	case "signup":
		if username == "" || password == "" {
			h.renderLoginError(w, r, http.StatusBadRequest, returnTo, "Username and password are required")
			return
		}
		if _, err := h.store.FindUserByUsername(r.Context(), username); err == nil {
			h.renderLoginError(w, r, http.StatusConflict, returnTo, "Username already exists")
			return
		} else if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("check user %q: %v", username, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		passHash, err := auth.HashPassword(password)
		if err != nil {
			log.Printf("hash password: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		user, err := h.store.CreateUser(r.Context(), username, passHash)
		if err != nil {
			log.Printf("create user %q: %v", username, err)
			h.renderLoginError(w, r, http.StatusConflict, returnTo, "Username already exists")
			return
		}

		h.sessions.Set(w, user.Username)
		log.Printf("User %s signed up.", user.Username)
		h.completeAuthSuccess(w, r, user, safeReturnTo(returnTo))

	case "login":
		if username == "" || password == "" {
			h.renderLoginError(w, r, http.StatusUnauthorized, returnTo, "Invalid username or password")
			return
		}
		user, err := h.store.FindUserByUsername(r.Context(), username)
		if err != nil || !auth.VerifyPassword(password, user.PassHash) {
			h.renderLoginError(w, r, http.StatusUnauthorized, returnTo, "Invalid username or password")
			return
		}

		h.sessions.Set(w, user.Username)
		log.Printf("User %s logged in.", user.Username)
		h.completeAuthSuccess(w, r, user, safeReturnTo(returnTo))

	default:
		h.renderLoginError(w, r, http.StatusBadRequest, returnTo, "Please select login or signup")
	}
}

func (h *Handler) renderLoginError(w http.ResponseWriter, r *http.Request, status int, returnTo, msg string) {
	setHTMXVary(w)
	if isPartialHTMX(r) {
		h.renderComponent(w, r, status, view.FormErrorBanner(msg))
		return
	}
	meta := h.pageMeta(r, "Login", "login")
	h.renderComponent(w, r, status, view.LoginPage(meta, returnTo, msg))
}

func (h *Handler) completeAuthSuccess(w http.ResponseWriter, r *http.Request, user *db.User, dest string) {
	if !isPartialHTMX(r) {
		http.Redirect(w, r, dest, http.StatusSeeOther)
		return
	}

	rWithUser := r.WithContext(context.WithValue(r.Context(), userContextKey, user))
	setHTMXVary(w)

	switch dest {
	case "/dashboard", "/dashboard/":
		w.Header().Set("HX-Push-Url", "/dashboard")
		meta := h.pageMetaForUser(user.Username, "Dashboard", "dashboard")
		h.renderComponent(w, rWithUser, http.StatusOK, view.HTMXPageFragment(meta, view.DashboardShell("home", "", view.DashboardHomeSection())))
	case "/dashboard/posts":
		posts, err := h.store.ListPostsByUser(r.Context(), user.ID)
		if err != nil {
			log.Printf("list user %d posts after auth: %v", user.ID, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("HX-Push-Url", "/dashboard/posts")
		meta := h.pageMetaForUser(user.Username, "Dashboard", "dashboard")
		h.renderComponent(w, rWithUser, http.StatusOK, view.HTMXPageFragment(meta, view.DashboardShell("posts", "", view.UserPosts(toPostCards(posts, 200, 50), ""))))
	case "/dashboard/posts/add":
		w.Header().Set("HX-Push-Url", "/dashboard/posts/add")
		meta := h.pageMetaForUser(user.Username, "Dashboard", "dashboard")
		h.renderComponent(w, rWithUser, http.StatusOK, view.HTMXPageFragment(meta, view.DashboardShell("form", "", view.PostForm(view.FormPost{}, ""))))
	default:
		locPayload, _ := json.Marshal(map[string]string{
			"path":   dest,
			"target": "#app-main",
			"swap":   "innerHTML show:top showTarget:window",
		})
		w.Header().Set("HX-Location", string(locPayload))
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) handleLogout(w http.ResponseWriter, r *http.Request) {
	h.sessions.Clear(w)
	if !isPartialHTMX(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	posts, err := h.store.ListPosts(r.Context())
	if err != nil {
		log.Printf("list posts on logout: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	setHTMXVary(w)
	w.Header().Set("HX-Push-Url", "/")
	meta := h.pageMetaForUser("", "abCMS", "home")
	h.renderComponent(w, r, http.StatusOK, view.HTMXPageFragment(meta, view.IndexContent(toPostCards(posts, 300, 200), view.FeedFilter{})))
}
