# Architecture

## Entry Point
- `cmd/abcms/main.go`: Loads configuration, opens SQLite (running automatic migrations), wires dependencies, and starts `net/http` with graceful shutdown (closing active SSE streams before HTTP drain).

## Internal Packages
- `internal/config/`: Environment-based runtime configuration.
- `internal/auth/`: Username slugification (`Slugify`), password hashing/verification (`bcrypt` cost 12 + legacy `argon2id`/`argon2i` PHC verification), and HMAC-SHA256 signed session cookies.
- `internal/db/`: Pure-Go SQLite (`modernc.org/sqlite`) store, search/author filter queries, and idempotent schema migrations compatible with legacy Lapis `abcms.sqlite` databases.
- `internal/markdown/`: `goldmark` Markdown rendering sanitized with `bluemonday.UGCPolicy()` and rune-safe preview truncation.
- `internal/media/`: Image upload validation, HEIC/HEIF-to-JPEG conversion, metadata-preserving storage, `1000x1000` JPEG thumbnail generation via `magick`, and native `syscall.Statfs` disk space reporting.
- `internal/view/`: Type-safe `templ` templates (`layout.templ`, `partials.templ`, `index.templ`, `post.templ`, `login.templ`, `dashboard.templ`, `post_form.templ`, `user_posts.templ`, `not_found.templ`) supporting both full-page SSR and partial HTMX 4 SPA fragments with `<template hx type="partial">` updates.
- `internal/handler/`: HTTP routing (`net/http.ServeMux`), `CrossOriginProtection` CSRF middleware, auth middleware, dual-mode HTMX 4 SPA rendering (`htmx.go`), SSE live event broker (`events.go`), public routes, dashboard CRUD routes, and static file serving.

## HTMX 4 SPA Architecture
- **Persistent Shell**: `<body>` uses HTMX 4.0.2 explicit inheritance (`hx-boost:inherited="true"`, `hx-target:inherited="#app-main"`, `hx-swap:inherited="innerHTML show:top showTarget:window"`, `hx-indicator:inherited="#htmx-progress-bar"`) with `<main id="app-main" hx-history-elt>` and `Vary: HX-Request, HX-Target`, keeping the header, SSE stream, and `starry_night.js` background canvas alive across navigations.
- **Targeted Fragment Swaps & `<template hx type="partial">`**: Server routes detect `HX-Request`, `HX-Target` (`tag#id` format), and `HX-History-Restore-Request` to return full documents on initial load/history restore, `#app-main` fragments on page navigation, `#dashboard-content` fragments on dashboard tab switching, or `#post-grid-region`/`#user-posts-list` on live search—alongside native HTMX 4 `<template hx type="partial">` swaps for `#site-nav`, `#dashboard-nav`, and `#header-stats`.
- **Built-in `hx-status` & Official v4 Extensions**: Uses native `hx-status:4xx` / `hx-status:5xx` for inline validation errors without losing form state, `hx-preload.min.js` for zero-latency `mousedown`/`touchstart` link & tab preloading, and `hx-sse.min.js` (`hx-sse:connect="/events"`) streaming `<template hx type="partial">` header stats and `posts-changed` events.

## Routes
- `GET /` — Homepage post grid with optional `?q=` search & `?author=` filtering
- `GET /posts/{post_id}` — Single post detail
- `GET /events` — Server-Sent Events stream (`<template hx type="partial">` stats + `posts-changed` event)
- `GET /login`, `POST /login` — Login & signup (inline 4xx errors via `hx-status:4xx`)
- `POST /logout` — Clear session cookie
- `GET /dashboard`, `GET /dashboard/posts`, `GET /dashboard/posts/add`, `GET /dashboard/posts/edit/{post_id}` — Auth-gated tabbed dashboard
- `POST /formapi/posts/add`, `POST /formapi/posts/mod/{post_id}`, `POST /formapi/posts/delete/{post_id}` — Auth-gated post mutations (multipart uploads + in-place card deletion)
- `GET /static/` — Static assets & uploads
