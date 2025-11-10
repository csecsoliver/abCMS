local lapis = require("lapis")
local Widget
Widget = require("lapis.html").Widget
local Posts
Posts = require("models").Posts
local PostForm = require("views/postform")
local Dashboard
do
  local _class_0
  local _parent_0 = Widget
  local _base_0 = {
    content = function(self)
      nav({
        class = "dashboard-nav"
      }, function()
        a({
          href = "/dashboard"
        }, "Dashboard Home")
        return a({
          href = "/dashboard/posts/add"
        }, "Add New Post")
      end)
      return div(function()
        local _exp_0 = self.params.splat
        if "posts/add" == _exp_0 then
          h2("Add New Post")
          return div(function()
            return widget(PostForm)
          end)
        else
          h2("Welcome to your Dashboard")
          return p("Use the navigation above to manage your posts.")
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
    __name = "Dashboard",
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
  Dashboard = _class_0
  return _class_0
end
