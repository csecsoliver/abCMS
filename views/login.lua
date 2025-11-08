local Widget
Widget = require("lapis.html").Widget
local Login
do
  local _class_0
  local _parent_0 = Widget
  local _base_0 = {
    content = function(self)
      return self:div({
        class = "login-container"
      }, function()
        self:h2("Login")
        if self.error_message then
          self:div({
            class = "error"
          }, self.error_message)
        end
        self:form({
          method = "POST",
          action = "/protected/login"
        }, function()
          self:input({
            type = "hidden",
            name = "return_to",
            value = self.return_to or ""
          })
          self:div(function()
            self:label({
              ["for"] = "login-username"
            }, "Username:")
            return self:input({
              type = "text",
              id = "login-username",
              name = "username",
              required = true
            })
          end)
          self:div(function()
            self:label({
              ["for"] = "login-password"
            }, "Password:")
            return self:input({
              type = "password",
              id = "login-password",
              name = "password",
              required = true
            })
          end)
          return self:div(function()
            return self:button({
              type = "submit"
            }, "Login")
          end)
        end)
        self:hr()
        self:h2("Sign Up")
        return self:form({
          method = "POST",
          action = "/protected/signup"
        }, function()
          self:input({
            type = "hidden",
            name = "return_to",
            value = self.return_to or ""
          })
          self:div(function()
            self:label({
              ["for"] = "signup-username"
            }, "Username:")
            return self:input({
              type = "text",
              id = "signup-username",
              name = "username",
              required = true
            })
          end)
          self:div(function()
            self:label({
              ["for"] = "signup-password"
            }, "Password:")
            return self:input({
              type = "password",
              id = "signup-password",
              name = "password",
              required = true
            })
          end)
          self:div(function()
            self:label({
              ["for"] = "signup-confirm"
            }, "Confirm Password:")
            return self:input({
              type = "password",
              id = "signup-confirm",
              name = "confirm_password",
              required = true
            })
          end)
          return self:div(function()
            return self:button({
              type = "submit"
            }, "Sign Up")
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
