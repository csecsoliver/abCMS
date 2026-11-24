package view

import (
	"fmt"
	"net/url"
	"time"
)

// PageMeta holds shared layout metadata rendered in the header and head.
type PageMeta struct {
	Title         string
	CurrentUser   string
	ActiveNav     string
	FreeSpaceText string
	DBSizeBytes   int64
	HasDBSize     bool
}

// FeedFilter holds active search and author filter parameters for the homepage feed.
type FeedFilter struct {
	Query  string
	Author string
}

// HasActive returns true if any filter criteria is set.
func (f FeedFilter) HasActive() bool {
	return f.Query != "" || f.Author != ""
}

// PostCard holds data for a post summary card on the homepage or dashboard.
type PostCard struct {
	ID             int64
	Title          string
	AuthorUsername string
	CreatedAt      int64
	HasImage       bool
	Path           string
	ThumbnailPath  string
	PreviewHTML    string
}

// DisplayImageURL returns the thumbnail URL when present, falling back to the full image URL.
func (p PostCard) DisplayImageURL() string {
	if p.ThumbnailPath != "" {
		return p.ThumbnailPath
	}
	return p.Path
}

// PostDetail holds full post details for the single-post page.
type PostDetail struct {
	ID                 int64
	Title              string
	AuthorUsername     string
	CreatedAt          int64
	HasImage           bool
	Path               string
	ThumbnailPath      string
	ContentHTML        string
	ImageSizeBytes     int64
	HasImageSize       bool
	ThumbnailSizeBytes int64
	HasThumbnailSize   bool
	IsOwner            bool
}

// FormPost holds pre-populated fields when creating or editing a post.
type FormPost struct {
	ID       int64
	Title    string
	Content  string
	HasImage bool
	Path     string
}

// FormatTimestamp formats a Unix timestamp in seconds as YYYY-MM-DD HH:MM:SS.
func FormatTimestamp(unixSec int64) string {
	return time.Unix(unixSec, 0).Local().Format("2006-01-02 15:04:05")
}

// PostURL builds the canonical URL path for viewing a post.
func PostURL(id int64) string {
	return fmt.Sprintf("/posts/%d", id)
}

// PostEditURL builds the dashboard edit URL for a post.
func PostEditURL(id int64) string {
	return fmt.Sprintf("/dashboard/posts/edit/%d", id)
}

// PostDeleteURL builds the POST action URL to delete a post.
func PostDeleteURL(id int64) string {
	return fmt.Sprintf("/formapi/posts/delete/%d", id)
}

// PostFormAction returns the POST action URL for creating or modifying a post.
func PostFormAction(id int64) string {
	if id > 0 {
		return fmt.Sprintf("/formapi/posts/mod/%d", id)
	}
	return "/formapi/posts/add"
}

// AuthorURL builds a homepage URL filtered to a specific author.
func AuthorURL(author string) string {
	return "/?author=" + url.QueryEscape(author)
}

// FeedURL builds a homepage URL preserving active search and author filters.
func FeedURL(query, author string) string {
	vals := url.Values{}
	if query != "" {
		vals.Set("q", query)
	}
	if author != "" {
		vals.Set("author", author)
	}
	if len(vals) == 0 {
		return "/"
	}
	return "/?" + vals.Encode()
}

// PostCardID returns the DOM element ID for a dashboard post card.
func PostCardID(id int64) string {
	return fmt.Sprintf("post-card-%d", id)
}

// NavActiveClass returns "nav-active" when active is true.
func NavActiveClass(active bool) string {
	if active {
		return "nav-active"
	}
	return ""
}
