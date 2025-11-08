lapis = require "lapis"
import respond_to from require "lapis.application"
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
      -- Not authenticated, redirect to login page with return URL
      return_to = @req.parsed_url.path
      if @req.parsed_url.query
        return_to = return_to .. "?" .. @req.parsed_url.query
      return redirect_to: @url_for("protected_login_page", {}, return_to: return_to)
    
    -- Store authenticated user in request context
    @current_user = Users\find username: username

  -- Login page
  [protected_login_page: "/protected/login"]: =>
    @return_to = @params.return_to or @url_for "protected_dashboard"
    render: "login"

  -- Login endpoint
  [protected_login: "/protected/login"]: respond_to {
    POST: =>
      username = @params.username
      password = @params.password
      return_to = @params.return_to or @url_for "protected_dashboard"
      
      unless username and password
        @error_message = "Username and password required"
        @return_to = return_to
        return render: "login"
      
      user = Auth\authenticate(username, password)
      
      unless user
        @error_message = "Invalid username or password"
        @return_to = return_to
        return render: "login"
      
      -- Generate JWT token
      token = Auth\generate_token(username)
      
      -- Set cookie
      @cookies.jwt_token = token
      
      -- Redirect to original destination or dashboard
      redirect_to: return_to
  }

  -- Signup endpoint
  [protected_signup: "/protected/signup"]: respond_to {
    POST: =>
      username = @params.username
      password = @params.password
      confirm_password = @params.confirm_password
      return_to = @params.return_to or @url_for "protected_dashboard"
      
      unless username and password and confirm_password
        @error_message = "All fields required"
        @return_to = return_to
        return render: "login"
      
      unless password == confirm_password
        @error_message = "Passwords do not match"
        @return_to = return_to
        return render: "login"
      
      user, err = Auth\create_user(username, password)
      
      unless user
        @error_message = err or "Failed to create user"
        @return_to = return_to
        return render: "login"
      
      -- Generate JWT token
      token = Auth\generate_token(username)
      
      -- Set cookie
      @cookies.jwt_token = token
      
      -- Redirect to original destination or dashboard
      redirect_to: return_to
  }

  -- Sample protected route
  [protected_dashboard: "/protected/dashboard"]: =>
    "Welcome, #{@current_user.username}! This is a protected page. Your XP: #{@current_user.xp}"

  -- Logout endpoint
  [protected_logout: "/protected/logout"]: =>
    @cookies.jwt_token = ""
    redirect_to: @url_for "protected_login_page"