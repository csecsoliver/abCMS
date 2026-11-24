package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "modernc.org/sqlite"
)

// Store wraps SQLite database access for abCMS.
type Store struct {
	db *sql.DB
}

// User represents a row in the users table.
type User struct {
	ID          int64
	Username    string
	PassHash    string
	XP          int64
	Coins       int64
	Social      string
	UploadToken string
}

// Post represents a row in the posts table joined with its author's username.
type Post struct {
	ID             int64
	UserID         int64
	AuthorUsername string
	CreatedAt      int64
	UpdatedAt      int64
	Title          string
	HasImage       bool
	Path           string
	Content        string
	ThumbnailPath  string
	Color          string
}

// Open opens the SQLite database at path, configures pragmas, and runs automatic migrations.
func Open(path string) (*Store, error) {
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite %q: %w", path, err)
	}
	conn.SetMaxOpenConns(1)

	pragmas := []string{
		"PRAGMA foreign_keys = ON;",
		"PRAGMA journal_mode = WAL;",
		"PRAGMA busy_timeout = 5000;",
	}
	for _, p := range pragmas {
		if _, err := conn.Exec(p); err != nil {
			_ = conn.Close()
			return nil, fmt.Errorf("apply pragma %q: %w", p, err)
		}
	}

	if err := Migrate(conn); err != nil {
		_ = conn.Close()
		return nil, err
	}

	return &Store{db: conn}, nil
}

// Close closes the underlying database connection.
func (s *Store) Close() error {
	return s.db.Close()
}

// FindUserByUsername retrieves a user by exact username.
func (s *Store) FindUserByUsername(ctx context.Context, username string) (*User, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, username, passhash, COALESCE(xp, 0), COALESCE(coins, 0), COALESCE(social, ''), COALESCE(upload_token, '')
		FROM users
		WHERE username = ?
	`, username)
	return scanUser(row)
}

// FindUserByID retrieves a user by ID.
func (s *Store) FindUserByID(ctx context.Context, id int64) (*User, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, username, passhash, COALESCE(xp, 0), COALESCE(coins, 0), COALESCE(social, ''), COALESCE(upload_token, '')
		FROM users
		WHERE id = ?
	`, id)
	return scanUser(row)
}

// CreateUser inserts a new user and returns the created record.
func (s *Store) CreateUser(ctx context.Context, username, passHash string) (*User, error) {
	res, err := s.db.ExecContext(ctx, `
		INSERT INTO users (username, passhash) VALUES (?, ?)
	`, username, passHash)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return s.FindUserByID(ctx, id)
}

const selectPostColumns = `
	SELECT
		p.id,
		p.user_id,
		COALESCE(u.username, 'unknown'),
		COALESCE(p.created_at, 0),
		COALESCE(p.updated_at, 0),
		COALESCE(p.title, ''),
		COALESCE(p.has_image, 0),
		COALESCE(p.path, ''),
		COALESCE(p.content, ''),
		COALESCE(p.thumbnail_path, ''),
		COALESCE(p.color, '')
	FROM posts p
	LEFT JOIN users u ON u.id = p.user_id
`

// ListPosts returns all posts ordered newest-first.
func (s *Store) ListPosts(ctx context.Context) ([]Post, error) {
	rows, err := s.db.QueryContext(ctx, selectPostColumns+` ORDER BY p.created_at DESC, p.id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanPosts(rows)
}

// ListPostsFiltered returns posts filtered by optional search query and author username.
func (s *Store) ListPostsFiltered(ctx context.Context, query, author string) ([]Post, error) {
	q := strings.TrimSpace(query)
	a := strings.TrimSpace(author)
	if q == "" && a == "" {
		return s.ListPosts(ctx)
	}

	var clauses []string
	var args []any
	if q != "" {
		like := "%" + q + "%"
		clauses = append(clauses, `(p.title LIKE ? OR p.content LIKE ? OR COALESCE(u.username, '') LIKE ?)`)
		args = append(args, like, like, like)
	}
	if a != "" {
		clauses = append(clauses, `u.username = ?`)
		args = append(args, a)
	}

	sqlQuery := selectPostColumns + ` WHERE ` + strings.Join(clauses, ` AND `) + ` ORDER BY p.created_at DESC, p.id DESC`
	rows, err := s.db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanPosts(rows)
}

// ListPostsByUser returns all posts authored by userID ordered newest-first.
func (s *Store) ListPostsByUser(ctx context.Context, userID int64) ([]Post, error) {
	rows, err := s.db.QueryContext(ctx, selectPostColumns+` WHERE p.user_id = ? ORDER BY p.created_at DESC, p.id DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanPosts(rows)
}

// ListPostsByUserFiltered returns posts authored by userID filtered by an optional search query.
func (s *Store) ListPostsByUserFiltered(ctx context.Context, userID int64, query string) ([]Post, error) {
	q := strings.TrimSpace(query)
	if q == "" {
		return s.ListPostsByUser(ctx, userID)
	}
	like := "%" + q + "%"
	rows, err := s.db.QueryContext(
		ctx,
		selectPostColumns+` WHERE p.user_id = ? AND (p.title LIKE ? OR p.content LIKE ?) ORDER BY p.created_at DESC, p.id DESC`,
		userID, like, like,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanPosts(rows)
}

// FindPostByID fetches a single post by ID.
func (s *Store) FindPostByID(ctx context.Context, id int64) (*Post, error) {
	row := s.db.QueryRowContext(ctx, selectPostColumns+` WHERE p.id = ?`, id)
	var p Post
	var hasImage int
	if err := row.Scan(
		&p.ID,
		&p.UserID,
		&p.AuthorUsername,
		&p.CreatedAt,
		&p.UpdatedAt,
		&p.Title,
		&hasImage,
		&p.Path,
		&p.Content,
		&p.ThumbnailPath,
		&p.Color,
	); err != nil {
		return nil, err
	}
	p.HasImage = hasImage == 1
	return &p, nil
}

// CreatePost inserts a new post and returns its ID.
func (s *Store) CreatePost(ctx context.Context, post *Post) (int64, error) {
	hasImage := 0
	if post.HasImage {
		hasImage = 1
	}
	res, err := s.db.ExecContext(ctx, `
		INSERT INTO posts (user_id, created_at, updated_at, title, has_image, path, content, thumbnail_path, color)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, post.UserID, post.CreatedAt, post.UpdatedAt, post.Title, hasImage, post.Path, post.Content, post.ThumbnailPath, post.Color)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// UpdatePost updates an existing post's editable fields.
func (s *Store) UpdatePost(ctx context.Context, post *Post) error {
	hasImage := 0
	if post.HasImage {
		hasImage = 1
	}
	_, err := s.db.ExecContext(ctx, `
		UPDATE posts
		SET updated_at = ?, title = ?, has_image = ?, path = ?, content = ?, thumbnail_path = ?, color = ?
		WHERE id = ?
	`, post.UpdatedAt, post.Title, hasImage, post.Path, post.Content, post.ThumbnailPath, post.Color, post.ID)
	return err
}

// UpdatePostThumbnail persists a lazily generated thumbnail path for an existing post.
func (s *Store) UpdatePostThumbnail(ctx context.Context, id int64, thumbnailPath string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE posts SET thumbnail_path = ? WHERE id = ?`, thumbnailPath, id)
	return err
}

// DeletePost removes a post row by ID.
func (s *Store) DeletePost(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM posts WHERE id = ?`, id)
	return err
}

func scanUser(row *sql.Row) (*User, error) {
	var u User
	if err := row.Scan(&u.ID, &u.Username, &u.PassHash, &u.XP, &u.Coins, &u.Social, &u.UploadToken); err != nil {
		return nil, err
	}
	return &u, nil
}

func scanPosts(rows *sql.Rows) ([]Post, error) {
	var posts []Post
	for rows.Next() {
		var p Post
		var hasImage int
		if err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.AuthorUsername,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.Title,
			&hasImage,
			&p.Path,
			&p.Content,
			&p.ThumbnailPath,
			&p.Color,
		); err != nil {
			return nil, err
		}
		p.HasImage = hasImage == 1
		posts = append(posts, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}
