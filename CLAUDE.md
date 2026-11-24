# CLAUDE.md

This file provides guidance when working with this repository.

## Project Overview

abCMS is a content management system written in **Go 1.27** with **templ** templates, **SQLite** (`modernc.org/sqlite`), and **ImageMagick** (`magick`) for image validation and thumbnailing.

The previous MoonScript/Lua (Lapis + OpenResty + CopyParty) implementation is preserved on the `checkpoint/lapis-lua` branch.

## Development Workflow

- `go tool templ generate` (or `make generate`) — Generate `*_templ.go` files from `internal/view/*.templ`
- `go test ./...` (or `make test`) — Run unit and integration tests
- `go build -o abcms ./cmd/abcms` (or `make build`) — Build binary
- `go run ./cmd/abcms` (or `make run`) — Start server (default `:8080`)

## Architecture

- `cmd/abcms/main.go`: Application entrypoint and graceful HTTP shutdown
- `internal/config/config.go`: Environment configuration (`ABCMS_ADDR`, `ABCMS_DB_PATH`, `ABCMS_STATIC_DIR`, `ABCMS_UPLOAD_DIR`, `ABCMS_SESSION_SECRET`, `ABCMS_SECURE_COOKIES`)
- `internal/auth/`: Slugify, bcrypt + legacy Argon2 password verification, HMAC-SHA256 cookie sessions
- `internal/db/`: SQLite store (`modernc.org/sqlite`) + automatic migrations compatible with legacy Lapis `abcms.sqlite` databases
- `internal/markdown/`: `goldmark` + `bluemonday` HTML-sanitized Markdown rendering
- `internal/media/`: ImageMagick (`magick`) upload/thumbnail pipeline (metadata intentionally preserved) and `syscall.Statfs` free-space helper
- `internal/view/`: `templ` UI templates and generated `*_templ.go` files
- `internal/handler/`: HTTP handlers, auth middleware, and `http.NewCrossOriginProtection()` CSRF middleware
- `static/`: CSS (`union.css`), JS (`image_hover.js`, `starry_night.js`, vendored HTMX), images, and `uploads/`

## Database Schema (`abcms.sqlite`)

### `users`
- `id` (integer, primary key)
- `username` (text, unique)
- `passhash` (text)
- `xp` (integer, default 0)
- `coins` (integer, default 0)
- `social` (text, default `""`)
- `upload_token` (text, default `""`)

### `posts`
- `id` (integer, primary key)
- `user_id` (integer, foreign key to `users`)
- `created_at` (integer timestamp)
- `updated_at` (integer timestamp, default 0)
- `title` (text, default `""`)
- `has_image` (integer, default 0)
- `path` (text, default `""`)
- `content` (text, default `""`)
- `thumbnail_path` (text, default `""`)
- `color` (text, default `""`)

## Naming Conventions

- **Go packages & files**: `snake_case` filenames (e.g., `post_form.templ`, `user_posts.templ`, `posts_api.go`).
- **HTML Data Attributes**: `kebab-case` `data-*` attributes (e.g., `data-full-src`).
- **JavaScript**: `snake_case` filenames (`image_hover.js`, `starry_night.js`), `camelCase` symbols (`initImageHover`, `starryNight`).

## Important Notes

- **Image Metadata**: Images are intentionally NOT sanitized of EXIF/metadata.
- **Session Expiry**: Sessions expire after 100 hours (`360000` seconds).
- **CSRF Protection**: State-changing routes use `POST` protected by Go's `http.NewCrossOriginProtection()`.
- **Legacy DB Compatibility**: `internal/db/migrations.go` upgrades older Lapis `abcms.sqlite` tables automatically without touching `lapis_migrations`.
