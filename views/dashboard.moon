lapis = require "lapis"

import Widget from require "lapis.html"
import Posts from require "models"
PostForm = require "views/postform"
class Dashboard extends Widget
    content: =>
        nav class: "dashboard-nav", ->
            a href: "/dashboard", "Dashboard Home"
            a href: "/dashboard/posts", " My Posts"
            a href: "/dashboard/posts/add", "Add New Post"
        div ->
            switch @params.splat
                when "posts/add"
                    h2 "Add New Post"
                    div ->
                        widget PostForm
                else
                    h2 "Welcome to your Dashboard"
                    p "Use the navigation above to manage your posts."

