import create_table, types from require "lapis.db.schema"

{
    [1]: =>
        create_table "users", {
            {"username", types.text unique: true}
            {"userid", types.integer primary_key: true}
            {"passhash", types.text}
            {"xp", types.integer default: 0}
        }
}