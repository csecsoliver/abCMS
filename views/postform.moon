lapis = require "lapis"

import Widget from require "lapis.html"
import Posts from require "models"

class PostForm extends Widget
    content: =>
        post = @options.post or {}
        form_method = "POST" -- todo: make proper mod handling
        form_action = if post.id then "/formapi/posts/mod/" .. post.id else "/formapi/posts/add"
        if @params.error_message
            p class: "error-message", @params.error_message
        form action: form_action, method: form_method, class: "post-form", enctype: "multipart/form-data", ->
            div ->
                p ->
                    b "Note: Images are not sanitized of metadata. This is an intentional design choice. In case you don't want your data on here, deal with it yourself."
            div ->
                label for: "title", "Title:"
                br!
                input type: "text", name: "title", id: "title", value: post.title or ""
            div ->
                label for: "content", "Content:"
                br!
                textarea name: "content", id: "content", post.content or ""
            div ->
                label for: "image", "Upload image:"
                br!
                input type: "file", name: "image", id: "image"
            div ->
                label for: "color", "Post color:"
                br!
                input type: "color", name: "color", id: "color", value: post.color or "#ffffff"
            div ->
                label for: "category", "Category:"
                br!
                select name: "category", id: "category", ->
                    categories = {"General", "Devlog", "Cool Stuff", "Questions", "Opinions", "Vibes"}
                    for category in *categories
                        if post.category == category
                            option value: category, selected: "selected", category
                        else
                            option value: category, category
                
            button type: "submit", "Post"
