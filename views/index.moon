lapis = require "lapis"

import Widget from require "lapis.html"
import Posts from require "models"

class HomePage extends Widget
    content: =>
        div class: "grid", ->
            for post in *Posts\select!
                div class: "cardfield", ->
                    img src: "/img/"..post.path
