import Widget from require "lapis.html"

class Layout extends Widget
  content: =>
    raw "<!DOCTYPE HTML>"
    html lang: "en", ->
      head ->
        meta charset: "UTF-8"
        meta name: "viewport", content: "width=device-width, initial-scale=1.0"
        title @page_title or "Page"

        script src: "/static/js/htmx.min.js"
        script src: "/static/js/htmx-ext-response-targets.js"
        script src: "/static/js/htmx-ext-sse.js"
        script src: "/static/js/htmx-ext-ws.js"
        link rel: "stylesheet", href: "/static/css/styles.css"
        link rel: "stylesheet", href: "/static/css/cozy.css"
        script src: "https://sfxr.me/riffwave.js"
        script src: "https://sfxr.me/sfxr.js"
        script src: "/static/js/script.js", defer: true

      body ->
        header ->
          h1 "abCMS"
          nav ->
            a href: "/", "Homepage"
            a href: "/dashboard", "Dashboard"
        main ->
          @content_for "inner"