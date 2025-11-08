import Widget from require "lapis.html"

class Login extends Widget
  content: =>
    div class: "login-container", ->
      h2 "Login"
      
      -- Show error message if present
      if @error_message
        div class: "error", @error_message
      
      -- Login form
      form {
        method: "POST"
        action: "/protected/login"
      }, ->
        div ->
          label for: "login-username", "Username:"
          input {
            type: "text"
            id: "login-username"
            name: "username"
            required: true
          }
        
        div ->
          label for: "login-password", "Password:"
          input {
            type: "password"
            id: "login-password"
            name: "password"
            required: true
          }
        
        div ->
          button type: "submit", "Login"
      
      hr!
      
      h2 "Sign Up"
      
      -- Signup form
      form {
        method: "POST"
        action: "/protected/signup"
      }, ->
        div ->
          label for: "signup-username", "Username:"
          input {
            type: "text"
            id: "signup-username"
            name: "username"
            required: true
          }
        
        div ->
          label for: "signup-password", "Password:"
          input {
            type: "password"
            id: "signup-password"
            name: "password"
            required: true
          }
        
        div ->
          label for: "signup-confirm", "Confirm Password:"
          input {
            type: "password"
            id: "signup-confirm"
            name: "confirm_password"
            required: true
          }
        
        div ->
          button type: "submit", "Sign Up"
