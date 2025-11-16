local config
config = require("lapis.config").config
return config("development", function()
  server("nginx")
  code_cache("off")
  num_workers("1")
  secret("ashdksduzawo1q28787triguw78zr32ughkjbb7fs8oi3h4")
  return sqlite(function()
    return database("abcms.sqlite")
  end)
end)
