lapis = require "lapis"

class extends lapis.Application
  layout:require "views.layout"
  @enable "etlua"
  @include "applications.protected"
  "/": =>
    render: "index"
  
  