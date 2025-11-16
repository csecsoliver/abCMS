import Widget from require "lapis.html"
file_utils = require "lib/file_utils"
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
        link rel: "stylesheet", href: "/static/css/union.css"
        -- link rel: "stylesheet", href: "/static/css/styles.css"
        -- link rel: "stylesheet", href: "/static/css/cozy.css"
        script src: "https://sfxr.me/riffwave.js"
        script src: "https://sfxr.me/sfxr.js"
        script src: "/static/js/misc.js", defer: true
        script src: "/static/js/script.js", defer: true
        -- link rel: "stylesheet", href: "https://cdn.jsdelivr.net/npm/easymde/dist/easymde.min.css"
        -- script src: "https://cdn.jsdelivr.net/npm/easymde/dist/easymde.min.js"


      body ->
        header ->
          h1 "abCMS"
          p file_utils\GetFreeSpace!
          nav ->
            a href: "/", "Homepage"
            a href: "/dashboard", "Dashboard"
        main ->
          @content_for "inner"
        -- script ->
        --   raw "const easymde = new EasyMDE({element: document.getElementById('content'),});"