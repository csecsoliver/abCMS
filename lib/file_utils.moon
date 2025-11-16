lfs = require 'lfs'
uuid = require "lua-uuid"
magick = require "magick"
import escape from require "lapis.util"
curl = require "cURL"
lume = require "lume"
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

GetFreeSpace = ->
    pieces = {}
    r = curl.easy
      url: "http://localhost:25511/img?pw=changeme&ls=t"--..os.getenv("COPYPARTY_PASS") currently hardcoded, so don't expose the 25511 port publicly
    --   httpheader:
    --     "X-Test-Header1: Header-Data1"
    --     "X-Test-Header2: Header-Data2"
      writefunction: (chunk) ->
        table.insert pieces, chunk
        chunk
    r\perform!
    r\close!
    response = table.concat pieces
    response = lume.split(response, "//")[2]
    response = lume.split(response, "2025")[1]
    response

{ :UploadImage, :GetFreeSpace }