local Widget
Widget = require("lapis.html").Widget
local markdown = require("markdown")
local lfs = require("lfs")
local PostView
do
  local _class_0
  local _parent_0 = Widget
  local _base_0 = {
    content = function(self)
      return div({
        class = "post-container"
      }, function()
        return div({
          class = "grid-item post-big"
        }, function()
          h2(self.post.title)
          local author = self.post:get_user().username
          p("Author: " .. tostring(author) .. " on " .. os.date("%Y-%m-%d %H:%M:%S", self.post.created_at))
          if self.post.has_image == 1 then
            img({
              src = self.post.path,
              alt = "Post Image",
              style = "width: 100%;"
            })
          end
          div(function()
            return raw(markdown(self.post.content))
          end)
          local image_file = lfs.attributes("." .. self.post.path)
          local thumbnail_file = lfs.attributes("." .. self.post.thumbnail_path)
          if image_file.size > 4100 then
            local image_size_bytes = image_file.size
            p("The original image takes up " .. image_size_bytes .. " bytes of space.")
          end
          if thumbnail_file.size > 4100 then
            local thumbnail_size_bytes = thumbnail_file.size
            return p("The thumbnail image takes up " .. thumbnail_size_bytes .. " bytes of space.")
          end
        end)
      end)
    end
  }
  _base_0.__index = _base_0
  setmetatable(_base_0, _parent_0.__base)
  _class_0 = setmetatable({
    __init = function(self, ...)
      return _class_0.__parent.__init(self, ...)
    end,
    __base = _base_0,
    __name = "PostView",
    __parent = _parent_0
  }, {
    __index = function(cls, name)
      local val = rawget(_base_0, name)
      if val == nil then
        local parent = rawget(cls, "__parent")
        if parent then
          return parent[name]
        end
      else
        return val
      end
    end,
    __call = function(cls, ...)
      local _self_0 = setmetatable({}, _base_0)
      cls.__init(_self_0, ...)
      return _self_0
    end
  })
  _base_0.__class = _class_0
  if _parent_0.__inherited then
    _parent_0.__inherited(_parent_0, _class_0)
  end
  PostView = _class_0
  return _class_0
end
