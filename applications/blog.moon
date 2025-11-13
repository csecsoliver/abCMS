lapis = require "lapis"
import escape from require "lapis.html"
file_utils = require "lib.file_utils"
import Posts, Users from require "models"
class BlogApplication extends lapis.Application
    "/formapi/posts/add": =>
        author = Users\find username: @session.user
        title = escape(@params.title or "")
        content = escape(@params.content or "")
        if @params.files
            if @params.files.image
                file_utils.UploadImage!
                

            
        
        Posts\create
            user_id: author.id
            title: title
            content: content
            created_at: os.time!
        @write redirect_to: "/dashboard/posts"
    "/posts/:postid": =>
        @post = Posts\find id: @params.postid
        if @post
            @write render: "post"
        else
            @write "Post not found", status: 404
