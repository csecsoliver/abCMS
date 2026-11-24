local uuid = require("lua-uuid")
local magick = require("magick")
local curl = require("cURL")
local lume = require("lume")
local upload_image
upload_image = function(self)
  local image = magick.load_image_from_blob(self.params.image.content)
  local format = image and image:get_format()
  local allowed_formats = {
    "JPEG",
    "JPG",
    "PNG",
    "GIF",
    "WEBP",
    "BMP",
    "TIFF",
    "TIF",
    "ICO",
    "CUR",
    "SVG",
    "SVGZ",
    "APNG",
    "AVIF",
    "JXL",
    "HEIC",
    "HEIF"
  }
  if not (format and lume.find(allowed_formats, string.upper(format))) then
    self:write({
      status = 418
    })
    return 
  end
  local filename = ""
  local full_size_content = nil
  if format == "HEIC" or format == "HEIF" then
    image:set_format("JPEG")
    full_size_content = image:get_blob()
    filename = tostring(uuid.new()) .. self.params.image.filename .. ".jpg"
  else
    full_size_content = self.params.image.content
    filename = tostring(uuid.new()) .. self.params.image.filename
  end
  local output_file = assert(io.open('static/uploads/' .. filename, 'wb'))
  output_file:write(full_size_content)
  print("Uploaded image to static/uploads/" .. filename)
  output_file:close()
  if format == "HEIC" or format == "HEIF" then
    image = magick.load_image_from_blob(self.params.image.content)
  end
  image:set_format("JPEG")
  image:thumb("1000x1000")
  local thumbnail_file = assert(io.open('static/uploads/thumb-' .. filename .. ".jpg", 'wb'))
  thumbnail_file:write(image:get_blob())
  print("Uploaded thumbnail to static/uploads/thumb-" .. filename .. ".jpg")
  thumbnail_file:close()
  image:destroy()
  return "/static/uploads/" .. filename, "/static/uploads/thumb-" .. filename .. ".jpg"
end
local get_free_space
get_free_space = function()
  local response_chunks = { }
  local request = curl.easy({
    url = "http://localhost:25511/img?pw=changeme&ls=t",
    writefunction = function(chunk)
      table.insert(response_chunks, chunk)
      return chunk
    end
  })
  request:perform()
  request:close()
  local response = table.concat(response_chunks)
  response = lume.split(response, "//")[2]
  response = lume.split(response, "2025")[1]
  return response
end
local generate_thumbnail
generate_thumbnail = function(path)
  local image_file = assert(io.open("." .. path, "rb"))
  local image_data = image_file:read("*all")
  image_file:close()
  local image = magick.load_image_from_blob(image_data)
  image:set_format("JPEG")
  image:thumb("1000x1000")
  local thumbnail_id = tostring(uuid.new())
  local thumbnail_file = assert(io.open('static/uploads/thumb-' .. thumbnail_id .. ".jpg", 'wb'))
  thumbnail_file:write(image:get_blob())
  print("Uploaded thumbnail to static/uploads/thumb-" .. thumbnail_id .. ".jpg")
  thumbnail_file:close()
  return "/static/uploads/thumb-" .. thumbnail_id .. ".jpg"
end
return {
  upload_image = upload_image,
  get_free_space = get_free_space,
  generate_thumbnail = generate_thumbnail
}
