local Widget
Widget = require("lapis.html").Widget
local Layout
do
  local _class_0
  local _parent_0 = Widget
  local _base_0 = {
    content = function(self)
      raw("<!DOCTYPE HTML>")
      return html({
        lang = "en"
      }, function()
        head(function()
          meta({
            charset = "UTF-8"
          })
          meta({
            name = "viewport",
            content = "width=device-width, initial-scale=1.0"
          })
          title(self.page_title or "Page")
          script({
            src = "/static/js/htmx.min.js"
          })
          script({
            src = "/static/js/htmx-ext-response-targets.js"
          })
          script({
            src = "/static/js/htmx-ext-sse.js"
          })
          script({
            src = "/static/js/htmx-ext-ws.js"
          })
          link({
            rel = "stylesheet",
            href = "/static/css/styles.css"
          })
          link({
            rel = "stylesheet",
            href = "/static/css/cozy.css"
          })
          script({
            src = "https://sfxr.me/riffwave.js"
          })
          script({
            src = "https://sfxr.me/sfxr.js"
          })
          return script({
            src = "/static/js/script.js",
            defer = true
          })
        end)
        return body(function()
          header(function()
            h1("abCMS.text")
            return nav(function()
              a({
                href = "./"
              }, "Homepage")
              return a({
                href = "./admin.html"
              }, "Account")
            end)
          end)
          return self:content_for("inner")
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
    __name = "Layout",
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
  Layout = _class_0
  return _class_0
end
