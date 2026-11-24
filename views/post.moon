import Widget from require "lapis.html"
markdown = require "markdown"
lfs = require "lfs"
class PostView extends Widget
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
                image_file = lfs.attributes "."..@post.path
                thumbnail_file = lfs.attributes "."..@post.thumbnail_path
                if image_file.size > 4100
                    image_size_bytes = image_file.size 
                    p "The original image takes up " .. image_size_bytes .. " bytes of space."
                if thumbnail_file.size > 4100
                    thumbnail_size_bytes = thumbnail_file.size
                    p "The thumbnail image takes up " .. thumbnail_size_bytes .. " bytes of space."