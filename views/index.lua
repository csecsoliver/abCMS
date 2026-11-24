local Widget
Widget = require("lapis.html").Widget
local Posts
Posts = require("models").Posts
local generate_thumbnail
generate_thumbnail = require("lib.file_utils").generate_thumbnail
local markdown = require("markdown")
local string = require("string")
local IndexView
do
  local _class_0
  local _parent_0 = Widget
  local _base_0 = {
    content = function(self)
      return div({
        class = "grid-container"
      }, function()
        local _list_0 = Posts:select("order by created_at desc")
        for _index_0 = 1, #_list_0 do
          local post = _list_0[_index_0]
          div({
            class = "grid-item"
          }, function()
            h3(post.title)
            local author = post:get_user().username
            p("Author: " .. tostring(author) .. " on " .. os.date("%Y-%m-%d %H:%M:%S", post.created_at))
            if post.has_image == 1 then
              if post.thumbnail_path == "" or post.thumbnail_path == nil then
                local thumbnail_path = generate_thumbnail(post.path)
                post:update({
                  thumbnail_path = thumbnail_path
                })
                post.thumbnail_path = thumbnail_path
              end
              img({
                src = post.thumbnail_path,
                ["data-full-src"] = post.path,
                alt = "Post Image",
                style = "width: 100%;"
              })
            end
            if string.len(post.content) <= 300 and post.has_image == 0 then
              raw(markdown(post.content))
            elseif string.len(post.content) > 300 and post.has_image == 0 then
              local preview_content = string.sub(post.content, 1, 300)
              raw(markdown(preview_content .. "..."))
            elseif string.len(post.content) <= 200 and post.has_image == 1 then
              raw(markdown(post.content))
            elseif string.len(post.content) > 200 and post.has_image == 1 then
              local preview_content = string.sub(post.content, 1, 200)
              raw(markdown(preview_content .. "..."))
            end
            return a({
              href = "/posts/" .. tostring(post.id)
            }, "Open post")
          end)
        end
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
    __name = "IndexView",
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
  IndexView = _class_0
  return _class_0
end
