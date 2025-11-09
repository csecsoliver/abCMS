import Model from require "lapis.db.model"

class Posts extends Model
    @relations: {
        {"user", belongs_to: "Users"}
    }