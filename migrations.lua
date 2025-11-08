local create_table, types
do
  local _obj_0 = require("lapis.db.schema")
  create_table, types = _obj_0.create_table, _obj_0.types
end
return {
  [1] = function(self)
    return create_table("users", {
      {
        "username",
        types.text({
          unique = true
        })
      },
      {
        "userid",
        types.integer({
          primary_key = true
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
  end
}
