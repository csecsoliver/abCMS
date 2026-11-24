package handler

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/csecsoliver/abCMS/internal/markdown"
	"github.com/csecsoliver/abCMS/internal/view"
)

func (h *Handler) handleIndex(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	author := strings.TrimSpace(r.URL.Query().Get("author"))

	posts, err := h.store.ListPostsFiltered(r.Context(), query, author)
	if err != nil {
		log.Printf("list posts: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	for i := range posts {
		if posts[i].HasImage && posts[i].ThumbnailPath == "" && posts[i].Path != "" {
			thumbURL, err := h.media.GenerateThumbnail(r.Context(), posts[i].Path)
			if err != nil {
				log.Printf("lazy thumbnail generation for post %d failed: %v", posts[i].ID, err)
				continue
			}
			if err := h.store.UpdatePostThumbnail(r.Context(), posts[i].ID, thumbURL); err != nil {
				log.Printf("save lazy thumbnail for post %d failed: %v", posts[i].ID, err)
				continue
			}
			posts[i].ThumbnailPath = thumbURL
		}
	}

	cards := toPostCards(posts, 300, 200)
	filter := view.FeedFilter{Query: query, Author: author}

	if isPartialHTMX(r) && hxTarget(r) == "post-grid-region" {
		setHTMXVary(w)
		h.renderComponent(w, r, http.StatusOK, view.PostGridRegion(cards, filter))
		return
	}

	meta := h.pageMeta(r, "abCMS", "home")
	h.renderPage(w, r, http.StatusOK, meta, view.IndexContent(cards, filter))
}

func (h *Handler) handlePost(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseInt(r.PathValue("post_id"), 10, 64)
	if err != nil || postID <= 0 {
		h.handleNotFound(w, r)
		return
	}

	post, err := h.store.FindPostByID(r.Context(), postID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.handleNotFound(w, r)
			return
		}
		log.Printf("find post %d: %v", postID, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var imageSize, thumbSize int64
	var hasImageSize, hasThumbSize bool
	if post.HasImage {
		imageSize, hasImageSize = h.media.FileSizeBytes(post.Path)
		thumbSize, hasThumbSize = h.media.FileSizeBytes(post.ThumbnailPath)
	}

	curUser := h.authenticatedUser(r)
	isOwner := curUser != nil && curUser.ID == post.UserID

	detail := view.PostDetail{
		ID:                 post.ID,
		Title:              post.Title,
		AuthorUsername:     post.AuthorUsername,
		CreatedAt:          post.CreatedAt,
		HasImage:           post.HasImage,
		Path:               post.Path,
		ThumbnailPath:      post.ThumbnailPath,
		ContentHTML:        markdown.RenderHTML(post.Content),
		ImageSizeBytes:     imageSize,
		HasImageSize:       hasImageSize,
		ThumbnailSizeBytes: thumbSize,
		HasThumbnailSize:   hasThumbSize,
		IsOwner:            isOwner,
	}

	title := post.Title
	if title == "" {
		title = "Post"
	}
	meta := h.pageMeta(r, title, "home")
	h.renderPage(w, r, http.StatusOK, meta, view.PostContent(detail))
}

func (h *Handler) handleNotFound(w http.ResponseWriter, r *http.Request) {
	meta := h.pageMeta(r, "404 - Page Not Found", "")
	h.renderPage(w, r, http.StatusNotFound, meta, view.NotFoundContent(meta))
}
