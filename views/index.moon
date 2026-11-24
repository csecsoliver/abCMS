import Widget from require "lapis.html"
import Posts from require "models"
import generate_thumbnail from require "lib.file_utils"
markdown = require "markdown"
string = require "string"
class IndexView extends Widget
    content: =>
        p ->
            a href: "/dashboard/posts/add", "Add post"
        div class: "grid-container", ->
            for post in *Posts\select "order by created_at desc"
                div class: "grid-item", ->
                    h3 post.title
                    author = post\get_user!.username
                    p "Author: #{author} on " .. os.date("%Y-%m-%d %H:%M:%S", post.created_at)
                    if post.has_image == 1
                        if post.thumbnail_path == "" or post.thumbnail_path == nil
                            thumbnail_path = generate_thumbnail post.path
                            post\update thumbnail_path: thumbnail_path
                            post.thumbnail_path = thumbnail_path
                        img src: post.thumbnail_path, ["data-full-src"]: post.path, alt: "Post Image", style: "width: 100%;"
                        
                    if string.len(post.content) <= 300 and post.has_image == 0
                        raw markdown(post.content)
                    elseif string.len(post.content) > 300 and post.has_image == 0
                        preview_content = string.sub(post.content, 1, 300)
                        raw markdown(preview_content.."...")
                    elseif string.len(post.content) <= 200 and post.has_image == 1
                        raw markdown(post.content)
                    elseif string.len(post.content) > 200 and post.has_image == 1
                        preview_content = string.sub(post.content, 1, 200)
                        raw markdown(preview_content.."...")
                
                    a href: "/posts/#{post.id}", "Open post"
