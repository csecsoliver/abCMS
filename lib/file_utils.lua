local lfs = require('lfs')
local uuid = require("lua-uuid")
local magick = require("magick")
local escape
escape = require("lapis.util").escape
local curl = require("cURL")
local lume = require("lume")
local UploadImage
UploadImage = function(self)
  local image = magick.load_image_from_blob(self.params.image.content)
  local format = image and image:get_format()
  local allowed = {
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
  if not (format and lume.find(allowed, string.upper(format))) then
    self:write({
      status = 418
    })
    return 
  end
  local filename = ""
  local fullsize_content = nil
  if format == "HEIC" or format == "HEIF" then
    image:set_format("JPEG")
    fullsize_content = image:get_blob()
    filename = tostring(uuid.new()) .. self.params.image.filename .. ".jpg"
  else
    fullsize_content = self.params.image.content
    filename = tostring(uuid.new()) .. self.params.image.filename
  end
  local file = assert(io.open('static/uploads/' .. filename, 'wb'))
  file:write(fullsize_content)
  print("Uploaded image to static/uploads/" .. filename)
  file:close()
  if format == "HEIC" or format == "HEIF" then
    image = magick.load_image_from_blob(self.params.image.content)
  end
  image:set_format("JPEG")
  image:thumb("1000x1000")
  local thumbfile = assert(io.open('static/uploads/thumb-' .. filename .. ".jpg", 'wb'))
  thumbfile:write(image:get_blob())
  print("Uploaded thumbnail to static/uploads/thumb-" .. filename .. ".jpg")
  thumbfile:close()
  image:destroy()
  return "/static/uploads/" .. filename, "/static/uploads/thumb-" .. filename .. ".jpg"
end
local GetFreeSpace
GetFreeSpace = function()
  local pieces = { }
  local r = curl.easy({
    url = "http://localhost:25511/img?pw=changeme&ls=t",
    writefunction = function(chunk)
      table.insert(pieces, chunk)
      return chunk
    end
  })
  r:perform()
  r:close()
  local response = table.concat(pieces)
  response = lume.split(response, "//")[2]
  response = lume.split(response, "2025")[1]
  return response
end
local GenThumb
GenThumb = function(path)
  local imgf = assert(io.open("." .. path, "rb"))
  local imgdata = imgf:read("*all")
  imgf:close()
  local imgimg = magick.load_image_from_blob(imgdata)
  imgimg:set_format("JPEG")
  imgimg:thumb("1000x1000")
  local newfilename = tostring(uuid.new())
  local thumbfile = assert(io.open('static/uploads/thumb-' .. newfilename .. ".jpg", 'wb'))
  thumbfile:write(imgimg:get_blob())
  print("Uploaded thumbnail to static/uploads/thumb-" .. newfilename .. ".jpg")
  thumbfile:close()
  return "/static/uploads/thumb-" .. newfilename .. ".jpg"
end
return {
  UploadImage = UploadImage,
  GetFreeSpace = GetFreeSpace,
  GenThumb = GenThumb
}
