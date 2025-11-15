local lapis = require("lapis")
local escape
escape = require("lapis.html").escape
local UploadImage
UploadImage = require("lib.file_utils").UploadImage
local Posts, Users
do
  local _obj_0 = require("models")
  Posts, Users = _obj_0.Posts, _obj_0.Users
end
local BlogApplication
do
  local _class_0
  local _parent_0 = lapis.Application
  local _base_0 = {
    ["/formapi/posts/add"] = function(self)
      local author = Users:find({
        username = self.session.user
      })
      local title = escape(self.params.title or "")
      local content = escape(self.params.content or "")
      local has_image = 0
      local path = ""
      local thumbnail_path = ""
      print(self.params.image.filename)
      if self.params.image and self.params.image.filename ~= "" then
        has_image = 1
        path, thumbnail_path = UploadImage(self)
      end
      Posts:create({
        user_id = author.id,
        title = title,
        content = content,
        created_at = os.time(),
        has_image = has_image,
        path = path,
        thumbnail_path = thumbnail_path
      })
      return self:write({
        redirect_to = "/dashboard/posts"
      })
    end,
    ["/posts/:postid"] = function(self)
      self.post = Posts:find({
        id = self.params.postid
      })
      if self.post then
        return self:write({
          render = "post"
        })
      else
        return self:write("Post not found", {
          status = 404
        })
      end
    end
  }
  _base_0.__index = _base_0
  setmetatable(_base_0, _parent_0.__base)
  _class_0 = setmetatable({
    __init = function(self, ...)
      return _class_0.__parent.__init(self, ...)
    end,
    __base = _base_0,
    __name = "BlogApplication",
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
  BlogApplication = _class_0
  return _class_0
end
