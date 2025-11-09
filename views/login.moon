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
        action: "/login"
      }, ->
        -- Hidden field to preserve return URL
        input {
          type: "hidden"
          name: "return_to"
          value: @return_to or ""
        }
        
        div ->
          label for: "username", "Username:"
          input {
            type: "text"
            id: "username"
            name: "username"
            required: true
          }
        
        div ->
          label for: "password", "Password:"
          input {
            type: "password"
            id: "password"
            name: "password"
            required: true
          }
        div ->
          label for: "option-login", "Login"
          label for: "option-signup", "Sign Up"
          input {
            id: "option-login"
            type: "radio"
            name: "option"
            value: "login"
            -- style: "width: 0; height: 0;"
          }
          input {
            id: "option-signup"
            type: "radio"
            name: "option"
            value: "signup"
            -- style: "width: 0; height: 0;"
          }
        div ->
          button type: "submit", "Go"
