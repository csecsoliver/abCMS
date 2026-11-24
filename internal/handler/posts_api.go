package handler

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/csecsoliver/abCMS/internal/db"
	"github.com/csecsoliver/abCMS/internal/media"
	"github.com/csecsoliver/abCMS/internal/view"
)

func (h *Handler) handlePostCreate(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, h.cfg.MaxUploadBytes)
	if err := r.ParseMultipartForm(h.cfg.MaxUploadBytes); err != nil {
		h.renderPostFormError(w, r, http.StatusRequestEntityTooLarge, "/dashboard/posts/add", "Upload is too large or invalid.")
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	content := r.FormValue("content")

	fileData, origFilename, hasUpload, err := readUploadedImage(r)
	if err != nil {
		log.Printf("read uploaded image: %v", err)
		h.renderPostFormError(w, r, http.StatusBadRequest, "/dashboard/posts/add", "Failed to read uploaded image.")
		return
	}

	if title == "" && strings.TrimSpace(content) == "" && !hasUpload {
		h.renderPostFormError(w, r, http.StatusUnprocessableEntity, "/dashboard/posts/add", "At least one of title, content, or image must be provided.")
		return
	}

	var hasImage bool
	var imagePath, thumbPath string
	if hasUpload {
		imagePath, thumbPath, err = h.media.SaveUpload(r.Context(), origFilename, fileData)
		if err != nil {
			msg := "Image processing failed."
			if errors.Is(err, media.ErrUnsupportedFormat) {
				msg = "Unsupported image format."
			}
			h.renderPostFormError(w, r, http.StatusUnprocessableEntity, "/dashboard/posts/add", msg)
			return
		}
		hasImage = true
	}

	author := h.authenticatedUser(r)
	now := time.Now().Unix()
	if _, err := h.store.CreatePost(r.Context(), &db.Post{
		UserID:        author.ID,
		CreatedAt:     now,
		UpdatedAt:     now,
		Title:         title,
		HasImage:      hasImage,
		Path:          imagePath,
		Content:       content,
		ThumbnailPath: thumbPath,
	}); err != nil {
		if hasImage {
			h.media.DeleteFile(imagePath)
			h.media.DeleteFile(thumbPath)
		}
		log.Printf("create post: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	h.broadcastFeedAndStats(r)

	if isPartialHTMX(r) {
		h.renderDashboardPostsAfterMutation(w, r, author.ID)
		return
	}
	http.Redirect(w, r, "/dashboard/posts", http.StatusSeeOther)
}

func (h *Handler) handlePostUpdate(w http.ResponseWriter, r *http.Request) {
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

	author := h.authenticatedUser(r)
	if post.UserID != author.ID {
		h.redirectOrRenderDashboardPostsError(w, r, "You are not authorized to edit this post.")
		return
	}

	editURL := fmt.Sprintf("/dashboard/posts/edit/%d", postID)
	r.Body = http.MaxBytesReader(w, r.Body, h.cfg.MaxUploadBytes)
	if err := r.ParseMultipartForm(h.cfg.MaxUploadBytes); err != nil {
		h.renderPostFormError(w, r, http.StatusRequestEntityTooLarge, editURL, "Upload is too large or invalid.")
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	content := r.FormValue("content")

	fileData, origFilename, hasUpload, err := readUploadedImage(r)
	if err != nil {
		log.Printf("read uploaded image on update: %v", err)
		h.renderPostFormError(w, r, http.StatusBadRequest, editURL, "Failed to read uploaded image.")
		return
	}

	if title == "" && strings.TrimSpace(content) == "" && !hasUpload && !post.HasImage {
		h.renderPostFormError(w, r, http.StatusUnprocessableEntity, editURL, "At least one of title, content, or image must be provided.")
		return
	}

	oldPath := post.Path
	oldThumbPath := post.ThumbnailPath
	if hasUpload {
		newPath, newThumb, err := h.media.SaveUpload(r.Context(), origFilename, fileData)
		if err != nil {
			msg := "Image processing failed."
			if errors.Is(err, media.ErrUnsupportedFormat) {
				msg = "Unsupported image format."
			}
			h.renderPostFormError(w, r, http.StatusUnprocessableEntity, editURL, msg)
			return
		}
		post.HasImage = true
		post.Path = newPath
		post.ThumbnailPath = newThumb
	}

	post.Title = title
	post.Content = content
	post.UpdatedAt = time.Now().Unix()

	if err := h.store.UpdatePost(r.Context(), post); err != nil {
		if hasUpload {
			h.media.DeleteFile(post.Path)
			h.media.DeleteFile(post.ThumbnailPath)
		}
		log.Printf("update post %d: %v", postID, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if hasUpload {
		h.media.DeleteFile(oldPath)
		h.media.DeleteFile(oldThumbPath)
	}

	h.broadcastFeedAndStats(r)

	if isPartialHTMX(r) {
		h.renderDashboardPostsAfterMutation(w, r, author.ID)
		return
	}
	http.Redirect(w, r, "/dashboard/posts", http.StatusSeeOther)
}

func (h *Handler) handlePostDelete(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseInt(r.PathValue("post_id"), 10, 64)
	if err != nil || postID <= 0 {
		h.renderDeleteError(w, r, "Post not found.")
		return
	}

	post, err := h.store.FindPostByID(r.Context(), postID)
	if err != nil {
		h.renderDeleteError(w, r, "Post not found.")
		return
	}

	author := h.authenticatedUser(r)
	if post.UserID != author.ID {
		h.renderDeleteError(w, r, "You are not authorized to delete this post.")
		return
	}

	if err := h.store.DeletePost(r.Context(), postID); err != nil {
		log.Printf("delete post %d: %v", postID, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if post.HasImage {
		h.media.DeleteFile(post.Path)
		h.media.DeleteFile(post.ThumbnailPath)
	}

	h.broadcastFeedAndStats(r)

	if isPartialHTMX(r) {
		from := r.FormValue("from")
		if from == "list" || strings.HasPrefix(hxTarget(r), "post-card-") {
			remaining, err := h.store.ListPostsByUser(r.Context(), author.ID)
			emptyFallback := err == nil && len(remaining) == 0
			setHTMXVary(w)
			if emptyFallback {
				w.Header().Set("HX-Retarget", "#user-posts-list")
				w.Header().Set("HX-Reswap", "outerHTML")
			}
			meta := h.pageMeta(r, "My Posts - Dashboard", "dashboard")
			h.renderComponent(w, r, http.StatusOK, view.HTMXDeleteCardResponse(meta, emptyFallback, ""))
			return
		}

		// Deleted from single post view -> transition to dashboard posts shell in #app-main.
		posts, err := h.store.ListPostsByUser(r.Context(), author.ID)
		if err != nil {
			log.Printf("list user %d posts after delete: %v", author.ID, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		setHTMXVary(w)
		w.Header().Set("HX-Push-Url", "/dashboard/posts")
		meta := h.pageMeta(r, "My Posts - Dashboard", "dashboard")
		h.renderComponent(w, r, http.StatusOK, view.HTMXPageFragment(meta, view.DashboardShell("posts", "", view.UserPosts(toPostCards(posts, 200, 50), ""))))
		return
	}

	http.Redirect(w, r, "/dashboard/posts?error_message="+url.QueryEscape("Post deleted successfully."), http.StatusSeeOther)
}

func (h *Handler) renderPostFormError(w http.ResponseWriter, r *http.Request, status int, redirectBase, msg string) {
	if isPartialHTMX(r) {
		setHTMXVary(w)
		h.renderComponent(w, r, status, view.FormErrorBanner(msg))
		return
	}
	http.Redirect(w, r, redirectBase+"?error_message="+url.QueryEscape(msg), http.StatusSeeOther)
}

func (h *Handler) renderDeleteError(w http.ResponseWriter, r *http.Request, msg string) {
	if isPartialHTMX(r) {
		setHTMXVary(w)
		w.Header().Set("HX-Reswap", "none")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	http.Redirect(w, r, "/dashboard/posts?error_message="+url.QueryEscape(msg), http.StatusSeeOther)
}

func (h *Handler) renderDashboardPostsAfterMutation(w http.ResponseWriter, r *http.Request, userID int64) {
	posts, err := h.store.ListPostsByUser(r.Context(), userID)
	if err != nil {
		log.Printf("list user %d posts after mutation: %v", userID, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("HX-Push-Url", "/dashboard/posts")
	meta := h.pageMeta(r, "My Posts - Dashboard", "dashboard")
	h.renderDashboardPage(w, r, http.StatusOK, meta, "posts", "", view.UserPosts(toPostCards(posts, 200, 50), ""))
}

func readUploadedImage(r *http.Request) ([]byte, string, bool, error) {
	file, header, err := r.FormFile("image")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return nil, "", false, nil
		}
		return nil, "", false, err
	}
	defer file.Close()

	if header.Filename == "" {
		return nil, "", false, nil
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, "", false, err
	}
	if len(data) == 0 {
		return nil, "", false, nil
	}
	return data, header.Filename, true, nil
}
