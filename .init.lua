local fm = require("fullmoon")
local posts = require("posts")
local pages = require("pages")


fm.setTemplate({"/views/", tmpl = "fmt"})


---@param r table
---@return function
fm.setRoute("/", function (r)
    return fm.serveContent()
end)
fm.run()