import create_table, types from require "lapis.db.schema"

{
    [1]: =>
        create_table "users", {
            {"id", types.integer primary_key: true}
            {"username", types.text unique: true}
            {"passhash", types.text}
            {"xp", types.integer default: 0}
        }
    [2]: =>
        create_table "posts", {
            {"id", types.integer primary_key: true}
            {"user_id", types.integer}
            {"created_at", types.integer}
            {"title", types.text default: ""}
            {"has_image", types.integer default: 0}
            {"path", types.text default: ""}
            {"content", types.text default: ""}
        }
}