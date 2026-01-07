lfs = require 'lfs'
uuid = require "lua-uuid"
magick = require "magick"
import escape from require "lapis.util"
curl = require "cURL"
lume = require "lume"
UploadImage = =>

    image = magick.load_image_from_blob @params.image.content
    format = image and image\get_format!
    allowed = {"JPEG", "JPG", "PNG", "GIF", "WEBP", "BMP", "TIFF", "TIF", "ICO", "CUR", "SVG", "SVGZ", "APNG", "AVIF", "JXL", "HEIC", "HEIF"}
    unless format and lume.find(allowed, string.upper format)
        @write status: 418
        return
    filename = ""
    fullsize_content = nil
    -- Convert HEIC/HEIF to JPEG for full-size image
    if format == "HEIC" or format == "HEIF"
        image\set_format "JPEG"
        fullsize_content = image\get_blob!
        filename = tostring(uuid.new!) .. @params.image.filename .. ".jpg"
    else
        fullsize_content = @params.image.content
        filename = tostring(uuid.new!) .. @params.image.filename
    
    file = assert(io.open('static/uploads/' .. filename, 'wb'))
    file\write(fullsize_content)
    print("Uploaded image to static/uploads/" .. filename)
    file\close()
    
    -- Reset image for thumbnail if we converted HEIC/HEIF
    if format == "HEIC" or format == "HEIF"
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
    year = os.date("%Y")
    response = lume.split(response, year)[1]
    response

{ :UploadImage, :GetFreeSpace }