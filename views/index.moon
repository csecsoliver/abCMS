lapis = require "lapis"

import Widget from require "lapis.html"
import Posts from require "models"
markdown = require "markdown"
string = require "string"
class HomePage extends Widget
    content: =>
        div class: "grid-container", ->
            for post in *Posts\select "order by created_at desc"
                style_attr = if post.color and post.color != "" then "border-color: #{post.color}; border-width: 2px; border-style: solid;" else ""
                div class: "grid-item", style: style_attr, ->
                    h3 post.title
                    author = post\get_user!.username
                    p "Author: #{author} on " .. os.date("%Y-%m-%d %H:%M:%S", post.created_at)
                    if post.has_image == 1 and post.thumbnail_path != ""
                        img src: post.thumbnail_path, ["data-fullsrc"]: post.path, alt: "Post Image", style: "width: 100%;"
                    else if post.has_image == 1 and post.thumbnail_path == ""
                        img src: post.path, alt: "Post Image", style: "width: 100%;"
                    if string.len(post.content) <= 200 and post.has_image == 0
                        raw markdown(post.content)
                    elseif string.len(post.content) > 200 and post.has_image == 0
                        content = string.sub(post.content, 1, 200)
                        raw markdown(content.."...")
                    elseif string.len(post.content) <= 50 and post.has_image == 1
                        raw markdown(post.content)
                    elseif string.len(post.content) > 50 and post.has_image == 1
                        content = string.sub(post.content, 1, 50)
                        raw markdown(content.."...")
                
                    a href: "/posts/#{post.id}", "Open post"
