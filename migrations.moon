import create_table, types, add_column from require "lapis.db.schema"

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
    [3]: =>
        add_column "posts", "thumbnail_path", types.text default: ""
    [4]: =>
        add_column "users", "coins", types.integer default: 0
        add_column "users", "social", types.text default: ""
        add_column "users", "upload_token", types.text default: ""
        add_column "posts", "color", types.text default: ""
}