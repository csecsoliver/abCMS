local lfs = require('lfs')
local uuid = require("lua-uuid")
local magick = require("magick")
local escape
escape = require("lapis.util").escape
local UploadImage
UploadImage = function(self)
  local filename = tostring(uuid.new()) .. self.params.image.filename
  local file = assert(io.open('static/uploads/' .. filename, 'wb'))
  file:write(self.params.image.content)
  print("Uploaded image to static/uploads/" .. filename)
  file:close()
  local image = magick.load_image_from_blob(self.params.image.content)
  image:set_format("JPEG")
  image:thumb("1000x1000")
  local thumbfile = assert(io.open('static/uploads/thumb-' .. filename .. ".jpg", 'wb'))
  thumbfile:write(image:get_blob())
  print("Uploaded thumbnail to static/uploads/thumb-" .. filename .. ".jpg")
  thumbfile:close()
  image:destroy()
  return "/static/uploads/" .. filename, "/static/uploads/thumb-" .. filename .. ".jpg"
end
local GenerateThumbnail
GenerateThumbnail = function(image_path, filename)
  magick = require("magick")
  magick.thumb(image_path, "100x100", "static/uploads/thumb-" .. filename .. ".jpg")
  print("Generated thumbnail at static/uploads/thumb-" .. filename .. ".jpg")
  return "/static/uploads/thumb-" .. filename .. ".jpg"
end
return {
  UploadImage = UploadImage,
  GenerateThumbnail = GenerateThumbnail
}
