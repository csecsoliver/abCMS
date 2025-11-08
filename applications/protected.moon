lapis = require "lapis"
import respond_to, json_params from require "lapis.application"
import Auth from require "lib.auth"
import Users from require "models"

class ProtectedApplication extends lapis.Application
  @before_filter =>
    -- Skip authentication for login and signup routes
    if @route_name == "protected_login" or @route_name == "protected_signup" or @route_name == "protected_login_page"
      return
    
    -- Check for JWT token in cookie
    token = @cookies.jwt_token
    username = Auth\verify_token(token)
    
    if not username
      -- Not authenticated, redirect to login page
      if @req.headers["hx-request"]
        -- HTMX request
        @res.headers["HX-Redirect"] = "/protected/login"
        return status: 401
      else
        -- Regular request, show login page
        return render: "login"
    
    -- Store authenticated user in request context
    @current_user = Users\find username: username

  -- Login page
  [protected_login_page: "/protected/login"]: =>
    render: "login"

  -- Login endpoint
  [protected_login: "/protected/login"]: respond_to {
    POST: =>
      username = @params.username
      password = @params.password
      
      unless username and password
        return status: 400, layout: false, "Username and password required"
      
      user = Auth\authenticate(username, password)
      
      unless user
        return status: 401, layout: false, "Invalid username or password"
      
      -- Generate JWT token
      token = Auth\generate_token(username)
      
      -- Set cookie
      @cookies.jwt_token = token
      
      -- Redirect to protected area
      @res.headers["HX-Redirect"] = "/protected/dashboard"
      return layout: false, "Login successful"
  }

  -- Signup endpoint
  [protected_signup: "/protected/signup"]: respond_to {
    POST: =>
      username = @params.username
      password = @params.password
      confirm_password = @params.confirm_password
      
      unless username and password and confirm_password
        return status: 400, layout: false, "All fields required"
      
      unless password == confirm_password
        return status: 400, layout: false, "Passwords do not match"
      
      user, err = Auth\create_user(username, password)
      
      unless user
        return status: 400, layout: false, err or "Failed to create user"
      
      -- Generate JWT token
      token = Auth\generate_token(username)
      
      -- Set cookie
      @cookies.jwt_token = token
      
      -- Redirect to protected area
      @res.headers["HX-Redirect"] = "/protected/dashboard"
      return layout: false, "Signup successful"
  }

  -- Sample protected route
  [protected_dashboard: "/protected/dashboard"]: =>
    "Welcome, #{@current_user.username}! This is a protected page. Your XP: #{@current_user.xp}"

  -- Logout endpoint
  [protected_logout: "/protected/logout"]: =>
    @cookies.jwt_token = ""
    redirect_to: @url_for "protected_login_page"