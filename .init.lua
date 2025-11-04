local fm = require("fullmoon")
fm.setRoute("/", function ()
    return "hello world"
end)
fm.run()