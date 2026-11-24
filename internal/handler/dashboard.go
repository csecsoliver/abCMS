package handler

import (
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/csecsoliver/abCMS/internal/view"
)

func (h *Handler) handleDashboardHome(w http.ResponseWriter, r *http.Request) {
	meta := h.pageMeta(r, "Dashboard", "dashboard")
	h.renderDashboardPage(w, r, http.StatusOK, meta, "home", r.URL.Query().Get("error_message"), view.DashboardHomeSection())
}

func (h *Handler) handleDashboardPosts(w http.ResponseWriter, r *http.Request) {
	user := h.authenticatedUser(r)
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	posts, err := h.store.ListPostsByUserFiltered(r.Context(), user.ID, query)
	if err != nil {
		log.Printf("list user %d posts: %v", user.ID, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	cards := toPostCards(posts, 200, 50)
	if isPartialHTMX(r) && hxTarget(r) == "user-posts-list" {
		setHTMXVary(w)
		h.renderComponent(w, r, http.StatusOK, view.UserPostsList(cards, query))
		return
	}

	meta := h.pageMeta(r, "My Posts - Dashboard", "dashboard")
	h.renderDashboardPage(w, r, http.StatusOK, meta, "posts", r.URL.Query().Get("error_message"), view.UserPosts(cards, query))
}

func (h *Handler) handleDashboardPostAdd(w http.ResponseWriter, r *http.Request) {
	meta := h.pageMeta(r, "Add Post - Dashboard", "dashboard")
	h.renderDashboardPage(w, r, http.StatusOK, meta, "form", "", view.PostForm(view.FormPost{}, r.URL.Query().Get("error_message")))
}

func (h *Handler) handleDashboardPostEdit(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseInt(r.PathValue("post_id"), 10, 64)
	if err != nil || postID <= 0 {
		h.redirectOrRenderDashboardPostsError(w, r, "Post not found.")
		return
	}

	post, err := h.store.FindPostByID(r.Context(), postID)
	if err != nil {
		h.redirectOrRenderDashboardPostsError(w, r, "Post not found.")
		return
	}

	user := h.authenticatedUser(r)
	if post.UserID != user.ID {
		h.redirectOrRenderDashboardPostsError(w, r, "You are not authorized to edit this post.")
		return
	}

	formPost := view.FormPost{
		ID:       post.ID,
		Title:    post.Title,
		Content:  post.Content,
		HasImage: post.HasImage,
		Path:     post.Path,
	}
	meta := h.pageMeta(r, "Edit Post - Dashboard", "dashboard")
	h.renderDashboardPage(w, r, http.StatusOK, meta, "form", "", view.PostForm(formPost, r.URL.Query().Get("error_message")))
}

func (h *Handler) redirectOrRenderDashboardPostsError(w http.ResponseWriter, r *http.Request, msg string) {
	if !isPartialHTMX(r) {
		http.Redirect(w, r, "/dashboard/posts?error_message="+url.QueryEscape(msg), http.StatusSeeOther)
		return
	}
	user := h.authenticatedUser(r)
	posts, err := h.store.ListPostsByUser(r.Context(), user.ID)
	if err != nil {
		log.Printf("list user %d posts on edit error: %v", user.ID, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("HX-Push-Url", "/dashboard/posts")
	meta := h.pageMeta(r, "My Posts - Dashboard", "dashboard")
	h.renderDashboardPage(w, r, http.StatusOK, meta, "posts", msg, view.UserPosts(toPostCards(posts, 200, 50), ""))
}
