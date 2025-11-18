lapis = require "lapis"

import Widget from require "lapis.html"
import Posts from require "models"
markdown = require "markdown"
string = require "string"
lfs = require "lfs"
class PostPage extends Widget
    content: =>
        div class: "post-container", ->
            div class: "grid-item post-big", ->
                h2 @post.title
                author = @post\get_user!.username
                p "Author: #{author} on " .. os.date("%Y-%m-%d %H:%M:%S", @post.created_at)
                if @post.has_image == 1
                    img src: @post.path, alt: "Post Image", style: "width: 100%;"
                div ->
                    raw markdown(@post.content)
                imagefile = lfs.attributes "."..@post.path
                thumbfile = lfs.attributes "."..@post.thumbnail_path
                if imagefile.size > 4100
                    size_kb = imagefile.size 
                    p "The original image takes up " .. size_kb .. " bytes of space."
                if thumbfile > 4100
                    thumb_size_kb = thumbfile.size
                    p "The thumbnail image takes up " .. thumb_size_kb .. " bytes of space."