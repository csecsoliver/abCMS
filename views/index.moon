lapis = require "lapis"

import Widget from require "lapis.html"
import Posts from require "models"
markdown = require "markdown"
string = require "string"
class HomePage extends Widget
    content: =>
        div class: "grid-container", ->
            for post in *Posts\select!
                div class: "grid-item", ->
                    h3 post.title
                    content = string.sub(post.content, 1, 150)
                    raw markdown(content.."...")
                    a href: "/posts/#{post.id}", "Read more"
