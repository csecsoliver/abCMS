lapis = require "lapis"
bcrypt = require "bcrypt"
import Users from require "models"
import slugify from require "lapis.util"
os = require "os"

import respond_to from require "lapis.application"


class ProtectedApplication extends lapis.Application
  @before_filter =>
    if (@session.user) and (Users\find username: @session.user) and (@session.expiry > os.time!)
      @current_user = Users\find username: @session.user
    else
      @write redirect_to: @url_for "login"

  [dashboard: "/dashboard"]: =>
    @page_title = "Dashboard"
    render: "dashboard"

  @include "applications.blog"
    