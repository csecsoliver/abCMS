lfs = require 'lfs'
uuid = require "lua-uuid"
UploadImage = =>
    file = assert(io.open 'filestore/img/' .. uuid.new! .. @params.files.image.filename , 'wb')
    file:write @params.files.image.content
    file:close!
    return file:getfd!
    