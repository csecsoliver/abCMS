import Model from require "lapis.db.model"

class Users extends Model
    @relations: {
        {"posts", has_many: "Posts"}
    }

