lapis = require "lapis"

import Widget from require "lapis.html"
import Posts from require "models"

class PostForm extends Widget
    content: =>
        post = @options.post or {}
        form_method = "POST" -- todo: make proper mod handling
        form_action = if post.id then "/formapi/posts/mod/" .. post.id else "/formapi/posts/add"

        form action: form_action, method: form_method, class: "post-form", enctype: "multipart/form-data", ->
            div ->
                label for: "title", "Title:"
                br!
                input type: "text", name: "title", id: "title", value: post.title or "", required: true
            div ->
                label for: "content", "Content:"
                br!
                textarea name: "content", id: "content", required: true, post.content or ""
            div ->
                label for: "image", "Upload image:"
                br!
                input type: "file", name: "image", id: "image", required: true
                
            button type: "submit", "Post"
