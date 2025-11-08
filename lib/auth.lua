local jwt = require("resty.jwt")
local bcrypt = require("bcrypt")
local Users
Users = require("models").Users
local Auth
do
  local _class_0
  local _base_0 = { }
  _base_0.__index = _base_0
  _class_0 = setmetatable({
    __init = function() end,
    __base = _base_0,
    __name = "Auth"
  }, {
    __index = _base_0,
    __call = function(cls, ...)
      local _self_0 = setmetatable({}, _base_0)
      cls.__init(_self_0, ...)
      return _self_0
    end
  })
  _base_0.__class = _class_0
  local self = _class_0
  self.get_secret = function(self)
    local secret = os.getenv("JWT_SECRET")
    if not secret then
      secret = "default-insecure-secret-change-in-production"
      print("WARNING: Using default JWT secret. Set JWT_SECRET environment variable!")
    end
    return secret
  end
  self.generate_token = function(self, username)
    local secret = self:get_secret()
    local jwt_token = jwt:sign(secret, {
      header = {
        typ = "JWT",
        alg = "HS256"
      },
      payload = {
        username = username,
        exp = os.time() + (7 * 24 * 60 * 60)
      }
    })
    return jwt_token
  end
  self.verify_token = function(self, token)
    if not token then
      return nil
    end
    local secret = self:get_secret()
    local jwt_obj = jwt:verify(secret, token)
    if not (jwt_obj and jwt_obj.verified) then
      return nil
    end
    return jwt_obj.payload.username
  end
  self.hash_password = function(self, password)
    return bcrypt.digest(password, 12)
  end
  self.verify_password = function(self, password, hash)
    return bcrypt.verify(password, hash)
  end
  self.authenticate = function(self, username, password)
    local user = Users:find({
      username = username
    })
    if not user then
      return nil
    end
    if not self:verify_password(password, user.passhash) then
      return nil
    end
    return user
  end
  self.create_user = function(self, username, password)
    local existing = Users:find({
      username = username
    })
    if existing then
      return nil, "User already exists"
    end
    local hash = self:hash_password(password)
    local user = Users:create({
      username = username,
      passhash = hash,
      xp = 0
    })
    return user
  end
  Auth = _class_0
end
return {
  Auth = Auth
}
