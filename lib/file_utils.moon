lfs = require 'lfs'
uuid = require "lua-uuid"
magick = require "magick"
import escape from require "lapis.util"

UploadImage = =>
    filename = tostring(uuid.new!) .. @params.image.filename
    file = assert(io.open('static/uploads/' ..  filename, 'wb'))
    file\write(@params.image.content)
    print("Uploaded image to static/uploads/" .. filename)
    file\close()

    image = magick.load_image_from_blob @params.image.content
    image\set_format "JPEG"
    image\thumb "1000x1000"

    thumbfile = assert(io.open('static/uploads/thumb-' .. filename.. ".jpg", 'wb'))
    thumbfile\write(image\get_blob!)
    print("Uploaded thumbnail to static/uploads/thumb-" .. filename.. ".jpg")
    thumbfile\close!

    image\destroy!

    return "/static/uploads/" .. filename, "/static/uploads/thumb-" .. filename.. ".jpg"

GenerateThumbnail = (image_path, filename) ->
    magick = require "magick"
    magick.thumb image_path, "100x100", "static/uploads/thumb-" .. filename.. ".jpg"
    print("Generated thumbnail at static/uploads/thumb-" .. filename.. ".jpg")
    return "/static/uploads/thumb-" .. filename.. ".jpg"

{ :UploadImage, :GenerateThumbnail }