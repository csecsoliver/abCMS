local lapis = require("lapis")
local slugify
slugify = require("lapis.util").slugify
local os = require("os")
local respond_to, render
do
  local _obj_0 = require("lapis.application")
  respond_to, render = _obj_0.respond_to, _obj_0.render
end
local bcrypt = require("bcrypt")
local Users
Users = require("models").Users
do
  local _class_0
  local _parent_0 = lapis.Application
  local _base_0 = {
    layout = require("views.layout"),
    ["/"] = function(self)
      self.page_title = "abCMS"
      return {
        render = "index"
      }
    end,
    [{
      login = "/login"
    }] = respond_to({
      GET = function(self)
        if self.session.user and (Users:find({
          username = self.session.user
        })) and (self.session.expiry > os.time()) then
          return {
            redirect_to = "/dashboard"
          }
        end
        self.page_title = "Login"
        return {
          render = "login"
        }
      end,
      POST = function(self)
        if self.params.option == "signup" then
          if Users:find({
            username = slugify(self.params.username)
          }) then
            self.error_message = "Username already exists"
            return {
              render = "login"
            }
          end
          local user = Users:create({
            username = slugify(self.params.username),
            passhash = bcrypt.digest(self.params.password, 12)
          })
          self.session.user = user.username
          self.session.expiry = os.time() + 360000
          print("User " .. tostring(self.session.user) .. " signed up.")
          local redirect_url
          if self.params.return_to and self.params.return_to ~= "" then
            redirect_url = self.params.return_to
          else
            redirect_url = "/dashboard"
          end
          return {
            redirect_to = redirect_url,
            layout = false
          }
        elseif self.params.option == "login" then
          local user = Users:find({
            username = slugify(self.params.username)
          })
          if user and bcrypt.verify(self.params.password, user.passhash) then
            self.session.user = user.username
            self.session.expiry = os.time() + 360000
            print("User " .. tostring(self.session.user) .. " logged in.")
            local redirect_url
            if self.params.return_to and self.params.return_to ~= "" then
              redirect_url = self.params.return_to
            else
              redirect_url = "/dashboard"
            end
            return {
              redirect_to = redirect_url,
              layout = false
            }
          else
            self.error_message = "Invalid username or password"
            return {
              render = "login"
            }
          end
        else
          self.error_message = "Please select login or signup"
          return {
            render = "login"
          }
        end
      end
    })
  }
  _base_0.__index = _base_0
  setmetatable(_base_0, _parent_0.__base)
  _class_0 = setmetatable({
    __init = function(self, ...)
      return _class_0.__parent.__init(self, ...)
    end,
    __base = _base_0,
    __name = nil,
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
  self:enable("etlua")
  self:include("applications.protected")
  if _parent_0.__inherited then
    _parent_0.__inherited(_parent_0, _class_0)
  end
  return _class_0
end
