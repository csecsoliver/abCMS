import Widget from require "lapis.html"

class Login extends Widget
  content: =>
    div class: "login-container", ->
      h2 "Login"
      
      if @error_message
        div class: "error", @error_message
      
      form {
        method: "POST"
        action: "/login"
      }, ->
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
        div class: "radio-group", ->
          input {
            id: "option-login"
            type: "radio"
            name: "option"
            value: "login"
            checked: true
            style: "position: absolute; opacity: 0; width: 0; height: 0;"
          }
          label {
            for: "option-login"
            class: "radio-label"
          }, "Login"
          
          input {
            id: "option-signup"
            type: "radio"
            name: "option"
            value: "signup"
            style: "position: absolute; opacity: 0; width: 0; height: 0;"
          }
          label {
            for: "option-signup"
            class: "radio-label"
          }, "Sign Up"
        div ->
          button type: "submit", "Go"
