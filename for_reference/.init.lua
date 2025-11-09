local fm = require("for_reference.lua.fullmoon")
local posts = require("for_reference.lua.posts")
local pages = require("for_reference.lua.pages")


fm.setTemplate({"/views/", tmpl = "fmt"})


---@param r table
---@return function
fm.setRoute("/:file", function (r)
    return fm.serveAsset("resources/"..r.params.file)
end)
fm.setRoute("/", function (r)
    return fm.serveContent("index")
end)
fm.setRoute("/post", function (r)
    
end)
fm.setRoute("/cozy", function (r)
    
end)
fm.setRoute("/cozy/post", function (r)
    
end)
fm.run()