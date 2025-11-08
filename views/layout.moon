import Widget from require "lapis.html"

class Layout extends Widget
  content: =>
    raw "<!DOCTYPE HTML>"
    html lang: "en", ->
      head ->
        meta charset: "UTF-8"
        meta name: "viewport", content: "width=device-width, initial-scale=1.0"
        title @page_title or "Page"
        
        script src: "./htmx.min.js"
        script src: "./htmx-ext-response-targets.js"
        script src: "./htmx-ext-sse.js"
        script src: "./htmx-ext-ws.js"
        link rel: "stylesheet", href: "./styles.css"
        link rel: "stylesheet", href: "./cozy.css"
        script src: "https://sfxr.me/riffwave.js"
        script src: "https://sfxr.me/sfxr.js"
        script src: "./script.js", defer: true
      
      body ->
        header ->
          h1 "abCMS.text"
          nav ->
            a href: "./", "Homepage"
            a href: "./admin.html", "Account"
        
        @content_for "inner"