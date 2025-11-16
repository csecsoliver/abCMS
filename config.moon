import config from require "lapis.config"

config "development", ->
  server "nginx"
  code_cache "off"
  num_workers "1"
  secret "ashdksduzawo1q28787triguw78zr32ughkjbb7fs8oi3h4"
  sqlite ->
    database "abcms.sqlite"
