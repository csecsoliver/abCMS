# Quick Reference: Adding Protected Routes

## Example 1: Simple Protected Route

Add this to `/applications/protected.moon`:

```moonscript
[protected_profile: "/protected/profile"]: =>
  "User Profile: #{@current_user.username}, XP: #{@current_user.xp}"
```

This route will be automatically protected by the middleware.

## Example 2: Protected Route with HTML

```moonscript
[protected_settings: "/protected/settings"]: =>
  @html ->
    h1 "Settings for #{@current_user.username}"
    p "XP: #{@current_user.xp}"
    form {
      ["hx-post"]: "/protected/update-settings"
      ["hx-target"]: "#result"
    }, ->
      button type: "submit", "Save Settings"
    div id: "result"
```

## Example 3: Protected API Endpoint

```moonscript
import json_params from require "lapis.application"

[protected_api_user: "/protected/api/user"]: json_params =>
  json: {
    username: @current_user.username
    xp: @current_user.xp
    userid: @current_user.userid
  }
```

## Example 4: Protected POST Endpoint

```moonscript
import respond_to from require "lapis.application"

[protected_update_xp: "/protected/update-xp"]: respond_to {
  POST: =>
    amount = tonumber(@params.amount) or 0
    @current_user\update xp: @current_user.xp + amount
    layout: false, "XP updated! New XP: #{@current_user.xp}"
}
```

## Middleware Behavior

The `@before_filter` in `protected.moon` runs before every route and:

1. **Skips authentication** for these routes:
   - `/protected/login` (GET)
   - `/protected/login` (POST)
   - `/protected/signup` (POST)

2. **Checks JWT token** for all other routes:
   - If no token or invalid token:
     - Shows login page (or redirects if HTMX request)
   - If valid token:
     - Loads user into `@current_user`
     - Allows route to execute

## Adding Routes That Don't Need Protection

If you want to add a route in the protected module that doesn't need authentication, update the `@before_filter`:

```moonscript
@before_filter =>
  -- Skip authentication for these routes
  if @route_name == "protected_login" or 
     @route_name == "protected_signup" or 
     @route_name == "protected_login_page" or
     @route_name == "protected_public_route"  -- Add your route here
    return
  
  -- ... rest of middleware code
```

## Example: Complete Protected Module

```moonscript
lapis = require "lapis"
import respond_to from require "lapis.application"
import Auth from require "lib.auth"
import Users from require "models"

class ProtectedApplication extends lapis.Application
  @before_filter =>
    -- Skip auth for public routes
    if @route_name == "protected_login" or @route_name == "protected_signup" or @route_name == "protected_login_page"
      return
    
    -- Check JWT token
    token = @cookies.jwt_token
    username = Auth\verify_token(token)
    
    unless username
      if @req.headers["hx-request"]
        @res.headers["HX-Redirect"] = "/protected/login"
        return status: 401
      else
        return render: "login"
    
    @current_user = Users\find username: username

  -- ... your protected routes here ...
  
  [protected_my_route: "/protected/my-route"]: =>
    "Hello #{@current_user.username}!"
```
