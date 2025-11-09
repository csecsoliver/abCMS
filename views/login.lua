local Widget
Widget = require("lapis.html").Widget
local Login
do
  local _class_0
  local _parent_0 = Widget
  local _base_0 = {
    content = function(self)
      return div({
        class = "login-container"
      }, function()
        h2("Login")
        if self.error_message then
          div({
            class = "error"
          }, self.error_message)
        end
        return form({
          method = "POST",
          action = "/login"
        }, function()
          input({
            type = "hidden",
            name = "return_to",
            value = self.return_to or ""
          })
          div(function()
            label({
              ["for"] = "username"
            }, "Username:")
            return input({
              type = "text",
              id = "username",
              name = "username",
              required = true
            })
          end)
          div(function()
            label({
              ["for"] = "password"
            }, "Password:")
            return input({
              type = "password",
              id = "password",
              name = "password",
              required = true
            })
          end)
          div(function()
            label({
              ["for"] = "option-login"
            }, "Login")
            label({
              ["for"] = "option-signup"
            }, "Sign Up")
            input({
              id = "option-login",
              type = "radio",
              name = "option",
              value = "login"
            })
            return input({
              id = "option-signup",
              type = "radio",
              name = "option",
              value = "signup"
            })
          end)
          return div(function()
            return button({
              type = "submit"
            }, "Go")
          end)
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
    __name = "Login",
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
  Login = _class_0
  return _class_0
end
