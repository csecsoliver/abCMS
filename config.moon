import config from require "lapis.config"

config "development", ->
  server "nginx"
  code_cache "off"
  num_workers "1"
  sqlite ->
    database "abcms.sqlite"
  session_name "abcms_session"
  secret os.getenv("SESSION_SECRET") or "change-this-secret-in-production"
