local config
config = require("lapis.config").config
return config("development", function()
  server("nginx")
  code_cache("off")
  num_workers("1")
  port("8081")
  return sqlite(function()
    return database("abcms.sqlite")
  end)
end)
