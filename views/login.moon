import Widget from require "lapis.html"

class Login extends Widget
  content: =>
    div class: "login-container", ->
      h2 "Login"
      
      -- Login form
      form {
        ["hx-post"]: "/protected/login"
        ["hx-target"]: "#message"
        ["hx-swap"]: "innerHTML"
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
      
      div id: "message"
      
      hr!
      
      h2 "Sign Up"
      
      -- Signup form
      form {
        ["hx-post"]: "/protected/signup"
        ["hx-target"]: "#signup-message"
        ["hx-swap"]: "innerHTML"
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
      
      div id: "signup-message"
