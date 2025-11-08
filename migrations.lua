local create_table, types
do
  local _obj_0 = require("lapis.db.schema")
  create_table, types = _obj_0.create_table, _obj_0.types
end
return {
  [1] = function(self)
    return create_table("users", {
      {
        "userid",
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
        "postid",
        types.integer({
          primary_key = true
        })
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
  end
}
