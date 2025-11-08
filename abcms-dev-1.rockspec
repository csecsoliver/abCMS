package = "abcms"
version = "dev-1"
source = {
   url = "git://github.com/csecsoliver/abCMS.git"
}
description = {
   summary = "A simple CMS built with Lapis",
   detailed = [[
      abCMS is a simple content management system built with the Lapis web framework for Lua/MoonScript.
   ]],
   homepage = "https://github.com/csecsoliver/abCMS",
   license = "MIT"
}
dependencies = {
   "lua >= 5.1",
   "lapis >= 1.7.0",
   "moonscript >= 0.5.0",
   "lua-resty-jwt >= 0.2.0",
   "bcrypt >= 2.1",
   "luaossl >= 20181207"
}
build = {
   type = "none"
}
