lapis = require "lapis"

import Widget from require "lapis.html"
import Posts from require "models"
import GenThumb from require "lib/file_utils"
markdown = require "markdown"
string = require "string"
class HomePage extends Widget
    content: =>
        div class: "grid-container", ->
            for post in *Posts\select "order by created_at desc"
                div class: "grid-item", ->
                    h3 post.title
                    author = post\get_user!.username
                    p "Author: #{author} on " .. os.date("%Y-%m-%d %H:%M:%S", post.created_at)
                    if post.has_image == 1
                        if post.thumbnail_path == ""
                            post.thumbnail_path = GenThumb post.path
                            post\update!
                        img src: post.thumbnail_path, ["data-fullsrc"]: post.path, alt: "Post Image", style: "width: 100%;"
                        
                    if string.len(post.content) <= 300 and post.has_image == 0
                        raw markdown(post.content)
                    elseif string.len(post.content) > 300 and post.has_image == 0
                        content = string.sub(post.content, 1, 300)
                        raw markdown(content.."...")
                    elseif string.len(post.content) <= 200 and post.has_image == 1
                        raw markdown(post.content)
                    elseif string.len(post.content) > 200 and post.has_image == 1
                        content = string.sub(post.content, 1, 200)
                        raw markdown(content.."...")
                
                    a href: "/posts/#{post.id}", "Open post"
