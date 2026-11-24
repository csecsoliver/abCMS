package db

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func TestOpenAndCRUD(t *testing.T) {
	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "test.sqlite")

	store, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open error: %v", err)
	}
	defer store.Close()

	u, err := store.CreateUser(ctx, "alice", "hash123")
	if err != nil {
		t.Fatalf("CreateUser error: %v", err)
	}
	if u.Username != "alice" || u.ID <= 0 {
		t.Fatalf("unexpected user: %+v", u)
	}

	foundUser, err := store.FindUserByUsername(ctx, "alice")
	if err != nil || foundUser.ID != u.ID {
		t.Fatalf("FindUserByUsername = (%+v, %v)", foundUser, err)
	}

	postID, err := store.CreatePost(ctx, &Post{
		UserID:    u.ID,
		CreatedAt: 1000,
		UpdatedAt: 1000,
		Title:     "Hello",
		Content:   "World",
	})
	if err != nil {
		t.Fatalf("CreatePost error: %v", err)
	}

	p, err := store.FindPostByID(ctx, postID)
	if err != nil {
		t.Fatalf("FindPostByID error: %v", err)
	}
	if p.Title != "Hello" || p.AuthorUsername != "alice" || p.HasImage {
		t.Fatalf("unexpected post: %+v", p)
	}

	p.Title = "Updated"
	p.Content = "New body"
	p.UpdatedAt = 2000
	if err := store.UpdatePost(ctx, p); err != nil {
		t.Fatalf("UpdatePost error: %v", err)
	}

	if err := store.UpdatePostThumbnail(ctx, postID, "/static/uploads/thumb-x.jpg"); err != nil {
		t.Fatalf("UpdatePostThumbnail error: %v", err)
	}

	posts, err := store.ListPosts(ctx)
	if err != nil || len(posts) != 1 || posts[0].Title != "Updated" {
		t.Fatalf("ListPosts = (%+v, %v)", posts, err)
	}

	userPosts, err := store.ListPostsByUser(ctx, u.ID)
	if err != nil || len(userPosts) != 1 {
		t.Fatalf("ListPostsByUser = (%+v, %v)", userPosts, err)
	}

	if err := store.DeletePost(ctx, postID); err != nil {
		t.Fatalf("DeletePost error: %v", err)
	}
	if _, err := store.FindPostByID(ctx, postID); err != sql.ErrNoRows {
		t.Fatalf("expected ErrNoRows after delete, got %v", err)
	}
}

func TestLegacyLapisDBCompatibility(t *testing.T) {
	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "legacy.sqlite")

	// Simulate an older Lapis SQLite database prior to migration 3/4.
	raw, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open raw sqlite: %v", err)
	}
	_, err = raw.Exec(`
		CREATE TABLE lapis_migrations (name TEXT PRIMARY KEY);
		INSERT INTO lapis_migrations (name) VALUES ('1'), ('2');
		CREATE TABLE users (
			id INTEGER PRIMARY KEY,
			username TEXT UNIQUE,
			passhash TEXT,
			xp INTEGER DEFAULT 0
		);
		CREATE TABLE posts (
			id INTEGER PRIMARY KEY,
			user_id INTEGER,
			created_at INTEGER,
			title TEXT DEFAULT '',
			has_image INTEGER DEFAULT 0,
			path TEXT DEFAULT '',
			content TEXT DEFAULT ''
		);
		INSERT INTO users (id, username, passhash, xp) VALUES (1, 'legacyuser', '$argon2id$v=19$m=16384,t=2,p=1$c2FsdA$aGFzaA', 42);
		INSERT INTO posts (id, user_id, created_at, title, has_image, path, content) VALUES (10, 1, 1700000000, 'Legacy Post', 0, NULL, 'Body');
	`)
	if err != nil {
		t.Fatalf("seed legacy schema: %v", err)
	}
	_ = raw.Close()

	store, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open on legacy DB failed: %v", err)
	}
	defer store.Close()

	u, err := store.FindUserByUsername(ctx, "legacyuser")
	if err != nil {
		t.Fatalf("FindUserByUsername on legacy user failed: %v", err)
	}
	if u.XP != 42 || u.Coins != 0 {
		t.Fatalf("unexpected legacy user fields: %+v", u)
	}

	p, err := store.FindPostByID(ctx, 10)
	if err != nil {
		t.Fatalf("FindPostByID on legacy post failed: %v", err)
	}
	if p.Title != "Legacy Post" || p.AuthorUsername != "legacyuser" || p.Path != "" || p.ThumbnailPath != "" {
		t.Fatalf("unexpected legacy post: %+v", p)
	}

	// Ensure lapis_migrations was preserved.
	var lapisCount int
	if err := store.db.QueryRow(`SELECT COUNT(*) FROM lapis_migrations`).Scan(&lapisCount); err != nil || lapisCount != 2 {
		t.Fatalf("expected 2 lapis_migrations rows preserved, got %d (err=%v)", lapisCount, err)
	}
}
