local Widget
Widget = require("lapis.html").Widget
local Error404
do
  local _class_0
  local _parent_0 = Widget
  local _base_0 = {
    content = function(self)
      return div({
        class = "error-container"
      }, function()
        h1({
          class = "error-code"
        }, "404")
        h2({
          class = "error-message"
        }, "Page Not Found")
        p({
          class = "error-description"
        }, "The page you're looking for doesn't exist or has been moved.")
        div({
          class = "error-actions"
        }, function()
          a({
            href = "/",
            class = "btn"
          }, "Go Home")
          if self.current_user then
            return a({
              href = "/dashboard",
              class = "btn"
            }, "Dashboard")
          else
            return a({
              href = "/login",
              class = "btn"
            }, "Login")
          end
        end)
        return p({
          class = "error-footer"
        }, function()
          text("Need help? ")
          a({
            href = "/posts"
          }, "View all posts")
          return text(" or contact the site administrator.")
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
    __name = "Error404",
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
  Error404 = _class_0
  return _class_0
end
