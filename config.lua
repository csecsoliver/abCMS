local config
config = require("lapis.config").config
return config("development", function()
  server("nginx")
  code_cache("off")
  num_workers("1")
  session_name("abcms_session")
  secret(os.getenv("SESSION_SECRET") or "change-this-secret-in-production")
  return sqlite(function()
    return database("abcms.sqlite")
  end)
end)
