lapis = require "lapis"

import Widget from require "lapis.html"
import Posts from require "models"
markdown = require "markdown"
string = require "string"
class PostPage extends Widget
    content: =>
        div class: "grid-item", ->
            h2 @post.title
            author = @post\get_user!.username
            p "Author: #{author} on " .. os.date("%Y-%m-%d %H:%M:%S", @post.created_at)
            div ->
                raw markdown(@post.content)