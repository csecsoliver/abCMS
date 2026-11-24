import Widget from require "lapis.html"
PostFormView = require "views.post_form"
UserPostsView = require "views.user_posts"
class DashboardView extends Widget
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
                        widget PostFormView
                when "posts"
                    h2 "My Posts"
                    div ->
                        widget UserPostsView
                else
                    h2 "Welcome to your Dashboard"
                    p "Use the navigation above to manage your posts."

