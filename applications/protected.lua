local lapis = require("lapis")
local respond_to
respond_to = require("lapis.application").respond_to
local Auth
Auth = require("lib.auth").Auth
local Users
Users = require("models").Users
local ProtectedApplication
do
  local _class_0
  local _parent_0 = lapis.Application
  local _base_0 = {
    [{ protected_login_page = "/protected/login" }] = function(self)
      return {
        render = "login"
      }
    end,
    [{ protected_login = "/protected/login" }] = respond_to({
      POST = function(self)
        local username = self.params.username
        local password = self.params.password
        if not (username and password) then
          self.error_message = "Username and password required"
          return {
            render = "login"
          }
        end
        local user = Auth:authenticate(username, password)
        if not user then
          self.error_message = "Invalid username or password"
          return {
            render = "login"
          }
        end
        local token = Auth:generate_token(username)
        self.cookies.jwt_token = token
        return {
          redirect_to = self:url_for("protected_dashboard")
        }
      end
    }),
    [{ protected_signup = "/protected/signup" }] = respond_to({
      POST = function(self)
        local username = self.params.username
        local password = self.params.password
        local confirm_password = self.params.confirm_password
        if not (username and password and confirm_password) then
          self.error_message = "All fields required"
          return {
            render = "login"
          }
        end
        if not (password == confirm_password) then
          self.error_message = "Passwords do not match"
          return {
            render = "login"
          }
        end
        local user, err = Auth:create_user(username, password)
        if not user then
          self.error_message = err or "Failed to create user"
          return {
            render = "login"
          }
        end
        local token = Auth:generate_token(username)
        self.cookies.jwt_token = token
        return {
          redirect_to = self:url_for("protected_dashboard")
        }
      end
    }),
    [{ protected_dashboard = "/protected/dashboard" }] = function(self)
      return "Welcome, " .. tostring(self.current_user.username) .. "! This is a protected page. Your XP: " .. tostring(self.current_user.xp)
    end,
    [{ protected_logout = "/protected/logout" }] = function(self)
      self.cookies.jwt_token = ""
      return {
        redirect_to = self:url_for("protected_login_page")
      }
    end
  }
  _base_0.__index = _base_0
  setmetatable(_base_0, _parent_0.__base)
  _class_0 = setmetatable({
    __init = function(self, ...)
      return _class_0.__parent.__init(self, ...)
    end,
    __base = _base_0,
    __name = "ProtectedApplication",
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
  local self = _class_0
  self:before_filter(function(self)
    if self.route_name == "protected_login" or self.route_name == "protected_signup" or self.route_name == "protected_login_page" then
      return
    end
    local token = self.cookies.jwt_token
    local username = Auth:verify_token(token)
    if not username then
      return {
        redirect_to = self:url_for("protected_login_page")
      }
    end
    self.current_user = Users:find({
      username = username
    })
  end)
  if _parent_0.__inherited then
    _parent_0.__inherited(_parent_0, _class_0)
  end
  ProtectedApplication = _class_0
  return _class_0
end
