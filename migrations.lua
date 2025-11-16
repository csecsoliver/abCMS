local create_table, types, add_column
do
  local _obj_0 = require("lapis.db.schema")
  create_table, types, add_column = _obj_0.create_table, _obj_0.types, _obj_0.add_column
end
return {
  [1] = function(self)
    return create_table("users", {
      {
        "id",
        types.integer({
          primary_key = true
        })
      },
      {
        "username",
        types.text({
          unique = true
        })
      },
      {
        "passhash",
        types.text
      },
      {
        "xp",
        types.integer({
          default = 0
        })
      }
    })
  end,
  [2] = function(self)
    return create_table("posts", {
      {
        "id",
        types.integer({
          primary_key = true
        })
      },
      {
        "user_id",
        types.integer
      },
      {
        "created_at",
        types.integer
      },
      {
        "title",
        types.text({
          default = ""
        })
      },
      {
        "has_image",
        types.integer({
          default = 0
        })
      },
      {
        "path",
        types.text({
          default = ""
        })
      },
      {
        "content",
        types.text({
          default = ""
        })
      }
    })
  end,
  [3] = function(self)
    return add_column("posts", "thumbnail_path", types.text({
      default = ""
    }))
  end,
  [4] = function(self)
    add_column("users", "coins", types.integer({
      default = 0
    }))
    add_column("users", "social", types.text({
      default = ""
    }))
    add_column("users", "upload_token", types.text({
      default = ""
    }))
    return add_column("posts", "color", types.text({
      default = ""
    }))
  end
}
