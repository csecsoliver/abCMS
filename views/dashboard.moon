lapis = require "lapis"

import Widget from require "lapis.html"
import Posts from require "models"
PostForm = require "views/postform"
UserPosts = require "views/userposts"
class Dashboard extends Widget
    content: =>
        nav class: "dashboard-nav", ->
            a href: "/dashboard", "Dashboard Home"
            a href: "/dashboard/posts", "My Posts"
            -- a href: "/dashboard/posts", " My Posts"
            a href: "/dashboard/posts/add", "Add New Post"
        div ->
            if @params.error_message
                p class: "error-message", @params.error_message
            switch @params.splat
                when "posts/add"
                    h2 "Add New Post"
                    div ->
                        widget PostForm
                when "posts"
                    h2 "My Posts"
                    div ->
                        widget UserPosts
                else
                    h2 "Welcome to your Dashboard"
                    p "Use the navigation above to manage your posts."

