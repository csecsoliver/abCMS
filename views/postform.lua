local lapis = require("lapis")
local Widget
Widget = require("lapis.html").Widget
local Posts
Posts = require("models").Posts
local PostForm
do
  local _class_0
  local _parent_0 = Widget
  local _base_0 = {
    content = function(self)
      local post = self.options.post or { }
      local form_method = "POST"
      local form_action
      if post.id then
        form_action = "/formapi/posts/mod/" .. post.id
      else
        form_action = "/formapi/posts/add"
      end
      return form({
        action = form_action,
        method = form_method,
        class = "post-form",
        enctype = "multipart/form-data"
      }, function()
        div(function()
          label({
            ["for"] = "title"
          }, "Title:")
          br()
          return input({
            type = "text",
            name = "title",
            id = "title",
            value = post.title or "",
            required = true
          })
        end)
        div(function()
          label({
            ["for"] = "content"
          }, "Content:")
          br()
          return textarea({
            name = "content",
            id = "content",
            required = true
          }, post.content or "")
        end)
        div(function()
          label({
            ["for"] = "image"
          }, "Upload image:")
          br()
          return input({
            type = "file",
            name = "image",
            id = "image"
          })
        end)
        return button({
          type = "submit"
        }, "Post")
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
    __name = "PostForm",
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
  PostForm = _class_0
  return _class_0
end
