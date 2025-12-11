import Widget from require "lapis.html"

class Error404 extends Widget
  content: =>
    div class: "error-container", ->
      h1 class: "error-code", "404"
      h2 class: "error-message", "Page Not Found"
      p class: "error-description", "The page you're looking for doesn't exist or has been moved."

      div class: "error-actions", ->
        a href: "/", class: "btn", "Go Home"
        if @current_user
          a href: "/dashboard", class: "btn", "Dashboard"
        else
          a href: "/login", class: "btn", "Login"

      p class: "error-footer", ->
        text "Need help? "
        a href: "/posts", "View all posts"
        text " or contact the site administrator."