# abCMS Dependencies

## Lua Dependencies

This project uses LuaRocks to manage dependencies. The dependencies are tracked in `abcms-dev-1.rockspec`.

### Core Dependencies

- **lua** (>= 5.1) - Lua programming language
- **lapis** (>= 1.7.0) - Web framework for Lua/MoonScript
- **moonscript** (>= 0.5.0) - MoonScript compiler
- **lua-resty-jwt** (>= 0.2.0) - JWT implementation for Lua/OpenResty
- **bcrypt** (>= 2.1) - Password hashing library
- **luaossl** (>= 20181207) - OpenSSL bindings for Lua

### Installation

To install all dependencies:

```bash
luarocks install --only-deps abcms-dev-1.rockspec
```

Or install individually:

```bash
luarocks install lapis
luarocks install moonscript
luarocks install lua-resty-jwt
luarocks install bcrypt
luarocks install luaossl
```

### Environment Variables

The following environment variables should be set for production:

- `JWT_SECRET` - Secret key for JWT token signing (required for security)
- `SESSION_SECRET` - Secret key for session management (required for security)

### Development

For local development, default secrets are used (with warnings). Make sure to set proper secrets in production.

## Compiling MoonScript

If you modify `.moon` files, compile them to `.lua`:

```bash
moonc .
```

Or compile specific files:

```bash
moonc models/users.moon
moonc applications/protected.moon
moonc views/login.moon
moonc lib/auth.moon
```
