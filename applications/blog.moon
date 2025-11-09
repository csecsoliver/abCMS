lapis = require "lapis"
import escape from require "lapis.html"

import Posts, Users from require "models"
class BlogApplication extends lapis.Application
   "/formapi/posts/add": =>
        author = Users\find username: @session.user
        title = escape(@params.title or "")
        content = escape(@params.content or "")
        
        Posts\create
            user_id: author.id
            title: title
            content: content
            created_at: os.time!
        @write redirect_to: "/dashboard/posts"
