# JWT Authentication System

This project implements a simple JWT-based authentication system with middleware protection.

## Features

- JWT token-based authentication
- Bcrypt password hashing
- Login and signup functionality
- Protected routes with middleware
- HTMX-powered forms (no page reloads)

## Routes

### Public Routes
- `GET /protected/login` - Login page with login and signup forms

### Authentication Endpoints
- `POST /protected/login` - Authenticate user and return JWT token
- `POST /protected/signup` - Create new user and return JWT token

### Protected Routes
- `GET /protected/dashboard` - Sample protected page (requires authentication)
- `GET /protected/logout` - Clear JWT token and redirect to login

## How It Works

1. **Authentication Flow**:
   - User submits username/password via HTMX form
   - Server validates credentials against database
   - On success, JWT token is generated and set as cookie
   - User is redirected to protected area

2. **Middleware Protection**:
   - Every route in the protected module checks for JWT token
   - If token is missing or invalid, user sees login page
   - If token is valid, user info is loaded into `@current_user`

3. **JWT Tokens**:
   - Tokens expire after 7 days
   - Stored in `jwt_token` cookie
   - Signed with JWT_SECRET environment variable

## Usage Example

### Creating a Protected Route

```moonscript
-- In applications/protected.moon
[protected_myroute: "/protected/myroute"]: =>
  -- @current_user is available here
  "Hello #{@current_user.username}!"
```

### Accessing Protected Route

1. Visit `/protected/login` to sign up or log in
2. After authentication, you can access `/protected/dashboard` or any protected route
3. If you try to access a protected route without authentication, you'll be redirected to login

## Environment Variables

- `JWT_SECRET` - Secret for signing JWT tokens (defaults to insecure value if not set)
- `SESSION_SECRET` - Secret for session management (defaults to insecure value if not set)

**Important**: Always set these in production!

## Security Considerations

- Passwords are hashed with bcrypt (cost factor: 12)
- JWT tokens are httpOnly cookies (when properly configured)
- Default secrets should never be used in production
- Tokens expire after 7 days
- Username must be unique in database

## Database Schema

The `users` table has the following structure:
- `userid` (INTEGER, PRIMARY KEY)
- `username` (TEXT, UNIQUE)
- `passhash` (TEXT)
- `xp` (INTEGER, DEFAULT 0)

Run migrations to create the table:
```bash
lapis migrate
```
