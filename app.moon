lapis = require "lapis"
import slugify from require "lapis.util"
os = require "os"
import respond_to, render from require "lapis.application"
bcrypt = require "bcrypt"
import Users from require "models"


class extends lapis.Application
  layout:require "views.layout"
  @enable "etlua"
  @include "applications.protected"
  "/": =>
    render: "index"
  [login: "/login"]: respond_to {
    GET: => 
      if @session.user
        return redirect_to: "/dashboard"
      @page_title = "Login"
      render: "login"
    POST: =>
        if @params.option == "signup"
          if Users\find username: slugify @params.username
            @error_message = "Username already exists"
            return render: "login"
          user = Users\create {
            username: slugify @params.username
            passhash: bcrypt.digest(@params.password, 12)
          }
          @session.user = user.username
          @session.expiry = os.time! + 360000
          print "User #{@session.user} signed up."
          redirect_url = if @params.return_to and @params.return_to != "" then @params.return_to else "/dashboard"
          return { redirect_to: redirect_url, layout: false }
        elseif @params.option == "login"
          user = Users\find username: slugify @params.username
          if user and bcrypt.verify(@params.password, user.passhash)
            @session.user = user.username
            @session.expiry = os.time! + 360000
            print "User #{@session.user} logged in."
            redirect_url = if @params.return_to and @params.return_to != "" then @params.return_to else "/dashboard"
            return { redirect_to: redirect_url, layout: false }
            -- return "User #{@session.user} logged in."
          else
            @error_message = "Invalid username or password"
            return render: "login"
        else
          @error_message = "Please select login or signup"
          return render: "login"
  }
  