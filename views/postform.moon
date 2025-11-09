lapis = require "lapis"

import Widget from require "lapis.html"
import Posts from require "models"

class PostForm extends Widget
    content: =>
        post = @options.post or {}
        form_method = if post.id then "PUT" else "POST"
        form_action = if post.id then "/formapi/posts/mod/" + post.id else "/formapi/posts/add"
        
        form action: form_action, method: form_method, class: "post-form", ->
            div ->
                label for: "title", "Title:"
                br!
                input type: "text", name: "title", id: "title", value: post.title or "", required: true
            div ->
                label for: "content", "Content:"
                br!
                textarea name: "content", id: "content", required: true, post.content or ""

            button type: "submit", "Post"
