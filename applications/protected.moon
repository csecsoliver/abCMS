lapis = require "lapis"
bcrypt = require "bcrypt"
import Users from require "models"
import slugify from require "lapis.util"
os = require "os"

import respond_to from require "lapis.application"


class UserApplication extends lapis.Application
  @before_filter =>
    print "Checking protected route for user #{@session.user}"
    print "Current time: #{os.time!}"
    print "Session expiry time: #{@session.expiry - os.time!}"
    print "User exists: #{(Users\find username: @session.user).username}"
    if (@session.user) and (Users\find username: @session.user) and (@session.expiry > os.time!)
      @current_user = Users\find username: @session.user
    else
      @write redirect_to: @url_for "login"
  
  [prottest: "/d"]: =>
    "cool protected content for #{@current_user.username}!"
  
  @include "applications.blog"
    