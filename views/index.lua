local lapis = require("lapis")
local Widget
Widget = require("lapis.html").Widget
local Posts
Posts = require("models").Posts
local markdown = require("markdown")
local string = require("string")
local HomePage
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
          local style_attr
          if post.color and post.color ~= "" then
            style_attr = "border-color: " .. tostring(post.color) .. "; border-width: 2px; border-style: solid;"
          else
            style_attr = ""
          end
          div({
            class = "grid-item",
            style = style_attr
          }, function()
            h3(post.title)
            local author = post:get_user().username
            p("Author: " .. tostring(author) .. " on " .. os.date("%Y-%m-%d %H:%M:%S", post.created_at))
            if post.has_image == 1 and post.thumbnail_path ~= "" then
              img({
                src = post.thumbnail_path,
                ["data-fullsrc"] = post.path,
                alt = "Post Image",
                style = "width: 100%;"
              })
            else
              if post.has_image == 1 and post.thumbnail_path == "" then
                img({
                  src = post.path,
                  alt = "Post Image",
                  style = "width: 100%;"
                })
              end
            end
            if string.len(post.content) <= 200 and post.has_image == 0 then
              raw(markdown(post.content))
            elseif string.len(post.content) > 200 and post.has_image == 0 then
              local content = string.sub(post.content, 1, 200)
              raw(markdown(content .. "..."))
            elseif string.len(post.content) <= 50 and post.has_image == 1 then
              raw(markdown(post.content))
            elseif string.len(post.content) > 50 and post.has_image == 1 then
              local content = string.sub(post.content, 1, 50)
              raw(markdown(content .. "..."))
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
    __name = "HomePage",
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
  HomePage = _class_0
  return _class_0
end
