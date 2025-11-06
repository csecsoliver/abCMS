local lapis = require("lapis")
local app = lapis.Application()

app:enable("etlua")
app.layout = require "views.layout"

app:get("/", function()
  return 
end)



return app
