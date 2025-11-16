lapis = require "lapis"

import Widget from require "lapis.html"
import Posts, Users from require "models"
PostForm = require "views/postform"
markdown = require "markdown"
class UserPosts extends Widget
    content: =>
        posts = Posts\select "where user_id = ? order by created_at desc", (Users\find username: @session.user).id
        if #posts == 0
            p "You have no posts yet."
        else
            div class: "grid-container", ->
                for post in *posts
                    div class: "grid-item", ->
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
                        br!
                        a href: "/posts/#{post.id}", "Open post"
                        br!
                        a href: "/formapi/posts/delete/#{post.id}", " Delete Post"