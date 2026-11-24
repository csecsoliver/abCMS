import Widget from require "lapis.html"
file_utils = require "lib.file_utils"
lfs = require "lfs"
class LayoutView extends Widget
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
        script src: "/static/js/image_hover.js", defer: true
        script src: "/static/js/starry_night.js", defer: true
        -- link rel: "stylesheet", href: "https://cdn.jsdelivr.net/npm/easymde/dist/easymde.min.css"
        -- script src: "https://cdn.jsdelivr.net/npm/easymde/dist/easymde.min.js"


      body ->
        header ->
          h1 "abCMS"
          p file_utils.get_free_space!
          db_file = lfs.attributes "abcms.sqlite"
          if db_file
              db_size_bytes = db_file.size
              p "The database takes up " .. db_size_bytes .. " bytes of space."
          nav ->
            a href: "/", "Homepage"
            a href: "/dashboard", "Dashboard"
        main ->
          @content_for "inner"
        -- script ->
        --   raw "const easymde = new EasyMDE({element: document.getElementById('content'),});"