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
      -- Not authenticated, redirect to login page
      return redirect_to: @url_for "protected_login_page"
    
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
        @error_message = "Username and password required"
        return render: "login"
      
      user = Auth\authenticate(username, password)
      
      unless user
        @error_message = "Invalid username or password"
        return render: "login"
      
      -- Generate JWT token
      token = Auth\generate_token(username)
      
      -- Set cookie
      @cookies.jwt_token = token
      
      -- Redirect to protected area
      redirect_to: @url_for "protected_dashboard"
  }

  -- Signup endpoint
  [protected_signup: "/protected/signup"]: respond_to {
    POST: =>
      username = @params.username
      password = @params.password
      confirm_password = @params.confirm_password
      
      unless username and password and confirm_password
        @error_message = "All fields required"
        return render: "login"
      
      unless password == confirm_password
        @error_message = "Passwords do not match"
        return render: "login"
      
      user, err = Auth\create_user(username, password)
      
      unless user
        @error_message = err or "Failed to create user"
        return render: "login"
      
      -- Generate JWT token
      token = Auth\generate_token(username)
      
      -- Set cookie
      @cookies.jwt_token = token
      
      -- Redirect to protected area
      redirect_to: @url_for "protected_dashboard"
  }

  -- Sample protected route
  [protected_dashboard: "/protected/dashboard"]: =>
    "Welcome, #{@current_user.username}! This is a protected page. Your XP: #{@current_user.xp}"

  -- Logout endpoint
  [protected_logout: "/protected/logout"]: =>
    @cookies.jwt_token = ""
    redirect_to: @url_for "protected_login_page"