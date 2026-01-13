# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

abCMS is a content management system built with MoonScript/Lua and the Lapis web framework. It's designed for managing blog posts with image support, user authentication, and file storage capabilities.

## Development Environment Setup

### Prerequisites
- Ubuntu/Debian-based system (setup script targets Ubuntu)
- Python 3 for CopyParty file server
- Git

### Initial Setup

1. **Install system dependencies and Lua packages:**
   ```bash
   ./setup.sh
   ```
   This script installs OpenResty, LuaRocks, and all required Lua packages.

2. **Compile MoonScript files:**
   ```bash
   moonc .
   ```

3. **Run database migrations:**
   ```bash
   lapis migrate
   ```

### Development Workflow

1. **Start the main application (requires two terminals):**

   - Terminal 1: Start Lapis server
     ```bash
     lapis serve
     ```
     This starts the web server on port 8080 (default).

   - Terminal 2: Start CopyParty file server
     ```bash
     ./setup_files.sh
     ```
     This starts the file server on port 3923.

2. **Common development commands:**
   - `moonc .` - Compile all MoonScript files to Lua
   - `lapis migrate` - Run database migrations
   - `lapis serve` - Start development server

3. **Ports to remember:**
   - 8080: Main application (can be exposed via reverse proxy)
   - 3923: CopyParty file server (DO NOT expose publicly)
   - 25511: CopyParty API endpoint (internal only)

## Architecture

### Technology Stack
- **Backend**: MoonScript (compiles to Lua)
- **Web Framework**: Lapis (runs on OpenResty/Nginx)
- **Database**: SQLite with Lapis ORM
- **Image Processing**: ImageMagick via LuaMagick
- **File Server**: CopyParty (Python-based)

### Directory Structure

- `app.moon` / `app.lua`: Main application entry point with core routes
- `config.lua`: Lapis configuration (development mode, SQLite)
- `migrations.lua`: Database schema definitions
- `models/`: Lapis models (Users, Posts)
- `applications/`: Application modules
  - `protected.moon`: Authentication middleware for protected routes
  - `blog.moon`: Post management (create, edit, delete)
- `views/`: moonscript HTML builder templates
- `lib/`: Utility libraries
  - `file_utils.moon`: Image upload and storage management
- `static/`: Static assets (CSS, JS, uploaded files)
- `filestore/`: CopyParty configuration and scripts

### Key Components

1. **Authentication System**:
   - Session-based authentication with expiration
   - Password hashing with bcrypt/argon2
   - Login/signup routes in `app.moon`
   - Protected routes middleware in `applications/protected.moon`

2. **Post Management**:
   - Create, view, and delete posts
   - Image uploads with thumbnail generation
   - HEIC/HEIF to JPEG conversion
   - Handled in `applications/blog.moon`

3. **Image Handling**:
   - Supports multiple formats (JPEG, PNG, GIF, WEBP, etc.)
   - Automatic thumbnail generation (1000x1000 max)
   - HEIC/HEIF conversion to JPEG
   - Files stored in `static/uploads/`

4. **File Storage**:
   - CopyParty integration for file serving
   - Space monitoring via API endpoint
   - Configuration in `filestore/abcms.conf`

## Database Schema

### Users Table
- `id` (primary key)
- `username` (unique text)
- `passhash` (text)
- `xp` (integer, default 0)

### Posts Table
- `id` (primary key)
- `user_id` (foreign key to users)
- `created_at` (integer timestamp)
- `title` (text)
- `content` (text)
- `filename` (text)
- `thumbname` (text)

## Important Notes

- **Image Metadata**: Images are NOT sanitized of metadata - this is intentional
- **File Paths**: Always use absolute paths when accessing files
- **Session Expiry**: Sessions expire after 360000 seconds (100 hours)
- **Port Safety**: Never expose ports 3923 or 25511 publicly
- **Compilation**: Always run `moonc .` after modifying .moon files
- **Database**: SQLite database file is `abcms.sqlite`

## Code Style

- MoonScript uses Python-like syntax with significant whitespace
- Indentation-based blocks (tabs or spaces, but be consistent)
- Functions defined with `->` for single expressions, `=>` for blocks
- Class definitions use `class` keyword
- Import statements use `import` keyword

## Nix Development Environment

For Nix-based systems (including NixOS), use the provided `shell.nix` for development:

### Nix Shell Setup

```bash
nix-shell
```

This will automatically install all required Lua packages and system dependencies.

### Key Nix Configuration Details

1. **Pre-built packages from nixpkgs**:
   - `lua51Packages.lua`, `lua51Packages.luarocks`, `lua51Packages.moonscript`
   - `lua51Packages.markdown`, `lua51Packages.luafilesystem`
   - `lua51Packages.lua-curl`, `lua51Packages.luasec`, `lua51Packages.luaossl`

2. **Packages installed via luarocks** (missing or broken in nixpkgs):
   - `bcrypt`, `argon2` - Password hashing libraries
   - `magick` - ImageMagick bindings (broken in nixpkgs)
   - `lume` - Utility library (not in nixpkgs)
   - `lua-uuid` - UUID generation
   - `lsqlite3` - SQLite bindings

3. **Critical environment variables**:
   - `LD_LIBRARY_PATH="${pkgs.imagemagick.out}/lib:$LD_LIBRARY_PATH"` - Required for magick package
   - `MAGICK_WAND_LIBRARY="${pkgs.imagemagick.out}/lib/libMagickWand-7.Q16HDRI.so"` - MagickWand library path

4. **Important flags for luarocks installations**:
   - `CRYPTO_DIR`, `CRYPTO_INCDIR`, `OPENSSL_DIR`, `OPENSSL_INCDIR` for lapis/openssl
   - `ARGON2_DIR`, `ARGON2_INCDIR`, `ARGON2_LIBDIR` for argon2
   - `MAGICK_DIR`, `MAGICK_INCDIR`, `MAGICK_LIBDIR` for magick
   - `UUID_DIR`, `UUID_INCDIR`, `UUID_LIBDIR` for lua-uuid
   - `SQLITE_DIR`, `SQLITE_INCDIR`, `SQLITE_LIBDIR` for lsqlite3

### Notes for Future Modifications

- Always check if a lua package exists in `lua51Packages` before adding to shellHook
- Broken packages (like magick) must be installed via luarocks with proper paths
- System libraries (libargon2, libuuid, etc.) need both buildInputs and luarocks flags
- ImageMagick requires LD_LIBRARY_PATH to be set for the magick Lua package

## Security Considerations

- Passwords are hashed using bcrypt (default) or argon2
- User input is escaped using Lapis's `escape` function
- File uploads are validated for allowed formats
- Session-based authentication with expiration
- Authorization checks ensure users can only modify their own posts