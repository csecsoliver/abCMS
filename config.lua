local config
config = require("lapis.config").config
config("development", function()
  server("nginx")
  code_cache("off")
  num_workers("1")
  return sqlite(function()
    return database("abcms.sqlite")
  end)
end)
config("test", function()
  server("nginx")
  code_cache("off")
  num_workers("1")
  return sqlite(function()
    return database("abcms.test.sqlite")
  end)
end)
return nil
