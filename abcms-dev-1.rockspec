package = "abcms"
version = "dev-1"

source = {
   url = "git://github.com/csecsoliver/abCMS.git",
   branch = "v4-rewrite"
}

description = {
   summary = "A CMS built with Lapis framework",
   detailed = [[
      abCMS is a content management system built with the Lapis web framework
      for OpenResty/nginx. It features blog functionality, user authentication,
      and markdown support.
   ]],
   homepage = "https://github.com/csecsoliver/abCMS"
}

dependencies = {
   "lua >= 5.1",
   -- Web framework
   "lapis >= 1.16.0",
   "moonscript >= 0.5.0",
   
   -- Database
   "lsqlite3 >= 0.9.6",
   "pgmoon >= 1.16.0",
   
   -- Security & Auth
   "bcrypt >= 2.3",
   "lua-resty-jwt >= 0.2.3",
   "lua-resty-openssl >= 1.7.0",
   "luaossl >= 20250929",
   
   -- Markdown parsers
   "markdown >= 0.33",
   "lua-discount >= 1.2.10",
   "lua-markdown-extra >= 0.8",
   "cmark >= 0.31.1",
   
   -- Utilities
   "lua-cjson >= 2.1.0",
   "etlua >= 1.3.0",
   "date >= 2.2.1",
   "luafilesystem >= 1.8.0",
   "luasocket >= 3.1.0",
   "lpeg >= 1.1.0",
   "ansicolors >= 1.0.2",
   "argparse >= 0.7.1",
   "alt-getopt >= 0.8.0",
   "loadkit >= 1.1.0",
   "lub >= 1.1.0",
   "yaml >= 1.1.2",
}

build = {
   type = "builtin",
   modules = {
      ["app"] = "app.lua",
      ["config"] = "config.lua",
      ["migrations"] = "migrations.lua",
      ["models"] = "models.lua",
      
      -- Applications
      ["applications.blog"] = "applications/blog.lua",
      ["applications.protected"] = "applications/protected.lua",
      
      -- Models
      ["models.posts"] = "models/posts.lua",
      ["models.users"] = "models/users.lua",
      
      -- Views
      ["views.dashboard"] = "views/dashboard.lua",
      ["views.index"] = "views/index.lua",
      ["views.layout"] = "views/layout.lua",
      ["views.login"] = "views/login.lua",
   },
   copy_directories = {
      "static",
      "views"
   }
}
