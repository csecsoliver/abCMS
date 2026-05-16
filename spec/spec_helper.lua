local use_test_env = require("lapis.spec").use_test_env

use_test_env()

local migrations = require("lapis.db.migrations")
local db = require("lapis.db")

migrations.run_migrations(require("migrations"))

local function truncate_tables()
  db.query("delete from posts")
  db.query("delete from users")
end

return {
  truncate_tables = truncate_tables
}
