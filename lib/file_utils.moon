uuid = require "lua-uuid"
magick = require "magick"
curl = require "cURL"
lume = require "lume"
upload_image = =>

    image = magick.load_image_from_blob @params.image.content
    format = image and image\get_format!
    allowed_formats = {"JPEG", "JPG", "PNG", "GIF", "WEBP", "BMP", "TIFF", "TIF", "ICO", "CUR", "SVG", "SVGZ", "APNG", "AVIF", "JXL", "HEIC", "HEIF"}
    unless format and lume.find(allowed_formats, string.upper format)
        @write status: 418
        return
    filename = ""
    full_size_content = nil
    -- Convert HEIC/HEIF to JPEG for full-size image
    if format == "HEIC" or format == "HEIF"
        image\set_format "JPEG"
        full_size_content = image\get_blob!
        filename = tostring(uuid.new!) .. @params.image.filename .. ".jpg"
    else
        full_size_content = @params.image.content
        filename = tostring(uuid.new!) .. @params.image.filename
    
    output_file = assert(io.open('static/uploads/' .. filename, 'wb'))
    output_file\write(full_size_content)
    print("Uploaded image to static/uploads/" .. filename)
    output_file\close()
    
    -- Reset image for thumbnail if we converted HEIC/HEIF
    if format == "HEIC" or format == "HEIF"
        image = magick.load_image_from_blob @params.image.content
    
    image\set_format "JPEG"
    image\thumb "1000x1000"

    thumbnail_file = assert(io.open('static/uploads/thumb-' .. filename.. ".jpg", 'wb'))
    thumbnail_file\write(image\get_blob!)
    print("Uploaded thumbnail to static/uploads/thumb-" .. filename.. ".jpg")
    thumbnail_file\close!
    
    image\destroy!
    return "/static/uploads/" .. filename, "/static/uploads/thumb-" .. filename.. ".jpg"

get_free_space = ->
    response_chunks = {}
    request = curl.easy
      url: "http://localhost:25511/img?pw=changeme&ls=t"--..os.getenv("COPYPARTY_PASS") currently hardcoded, so don't expose the 25511 port publicly
    --   httpheader:
    --     "X-Test-Header1: Header-Data1"
    --     "X-Test-Header2: Header-Data2"
      writefunction: (chunk) ->
        table.insert response_chunks, chunk
        chunk
    request\perform!
    request\close!
    response = table.concat response_chunks
    response = lume.split(response, "//")[2]
    response = lume.split(response, "2025")[1]
    response

generate_thumbnail = (path) ->
    image_file = assert io.open("." .. path, "rb")
    image_data = image_file\read "*all"
    image_file\close!
    image = magick.load_image_from_blob image_data

    image\set_format "JPEG"
    image\thumb "1000x1000"
    thumbnail_id = tostring uuid.new!
    thumbnail_file = assert(io.open('static/uploads/thumb-' .. thumbnail_id.. ".jpg", 'wb'))
    thumbnail_file\write(image\get_blob!)
    print("Uploaded thumbnail to static/uploads/thumb-" .. thumbnail_id.. ".jpg")
    thumbnail_file\close!
    return "/static/uploads/thumb-" .. thumbnail_id .. ".jpg"

{ :upload_image, :get_free_space, :generate_thumbnail }