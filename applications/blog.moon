lapis = require "lapis"
import escape from require "lapis.html"
import UploadImage from require "lib.file_utils"
import Posts, Users from require "models"
class BlogApplication extends lapis.Application
    "/formapi/posts/add": =>
        author = Users\find username: @session.user
        title = escape(@params.title or "")
        content = escape(@params.content or "")
        has_image = 0
        path = ""
        thumbnail_path = ""
        print @params.image.filename
        if @params.image and @params.image.filename != "" 
            has_image = 1
            path, thumbnail_path = UploadImage self

        if title == "" and content == "" and @params.image.filename == "" and @params.image
            @write redirect_to: "/dashboard/posts/add?error_message=At least one of title, content, or image must be provided."
            return
        
        Posts\create
            user_id: author.id
            title: title
            content: content
            created_at: os.time!
            has_image: has_image
            path: path
            thumbnail_path: thumbnail_path

        @write redirect_to: "/dashboard/posts"
    "/posts/:postid": =>
        @post = Posts\find id: @params.postid
        if @post
            @write render: "post"
        else
            @write "Post not found", status: 404
