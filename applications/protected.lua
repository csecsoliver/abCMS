local lapis = require("lapis")
local bcrypt = require("bcrypt")
local Users
Users = require("models").Users
local slugify
slugify = require("lapis.util").slugify
local os = require("os")
local respond_to
respond_to = require("lapis.application").respond_to
local UserApplication
do
  local _class_0
  local _parent_0 = lapis.Application
  local _base_0 = {
    [{
      prottest = "/prottest"
    }] = function(self)
      return "cool protected content for " .. tostring(self.current_user.username) .. "!"
    end
  }
  _base_0.__index = _base_0
  setmetatable(_base_0, _parent_0.__base)
  _class_0 = setmetatable({
    __init = function(self, ...)
      return _class_0.__parent.__init(self, ...)
    end,
    __base = _base_0,
    __name = "UserApplication",
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
    print("Checking protected route for user " .. tostring(self.session.user))
    print("Current time: " .. tostring(os.time()))
    print("Session expiry time: " .. tostring(self.session.expiry - os.time()))
    print("User exists: " .. tostring((Users:find({
      username = self.session.user
    })).username))
    if (self.session.user) and (Users:find({
      username = self.session.user
    })) and (self.session.expiry > os.time()) then
      self.current_user = Users:find({
        username = self.session.user
      })
    else
      return self:write({
        redirect_to = self:url_for("login")
      })
    end
  end)
  self:include("applications.blog")
  if _parent_0.__inherited then
    _parent_0.__inherited(_parent_0, _class_0)
  end
  UserApplication = _class_0
  return _class_0
end
