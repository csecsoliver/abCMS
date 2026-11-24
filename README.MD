# abCMS

abCMS is a lightweight content management system written in **Go 1.27** using **templ**, **SQLite** (`modernc.org/sqlite`), and **ImageMagick** (`magick`) for image processing.

> **Note:** Images are intentionally **not** stripped of EXIF/metadata on upload. If you do not want metadata published, strip it before uploading.

## Legacy Checkpoint

The previous MoonScript/Lua (Lapis + OpenResty + CopyParty) implementation is preserved on the `checkpoint/lapis-lua` branch. Legacy `v3-legacy-python` and `v4-rewrite` branches also remain in git history.

## Requirements

- Go `1.27+`
- ImageMagick (`magick` binary in `PATH`, with HEIC/WebP/PNG/JPEG delegates)

## Quick Start

```bash
# Generate templ templates
make generate

# Run tests
make test

# Build and run
make build
./abcms
```

By default, abCMS listens on `:8080`, uses `abcms.sqlite` in the working directory, and stores uploads in `static/uploads/`.

### Environment Variables

- `ABCMS_ADDR` / `PORT` — listen address (default `:8080`)
- `ABCMS_DB_PATH` — SQLite database path (default `abcms.sqlite`)
- `ABCMS_STATIC_DIR` — static asset directory (default `static`)
- `ABCMS_UPLOAD_DIR` — image upload directory (default `static/uploads`)
- `ABCMS_SESSION_SECRET` — HMAC signing key for session cookies (ephemeral random key if unset)
- `ABCMS_SECURE_COOKIES` — set `true` or `1` when serving over HTTPS

## Migrating an Existing Database

Existing `abcms.sqlite` databases created by the earlier Lapis version are automatically upgraded on startup via `internal/db/migrations.go` without altering existing user rows, posts, or `lapis_migrations`. Legacy Argon2id/Argon2i password hashes remain supported alongside bcrypt.
