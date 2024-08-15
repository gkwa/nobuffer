local lu = require('luaunit')
local httpbin = require('httpbin')

TestUUID = {}

function TestUUID:testUUIDFormat()
   local uuid = httpbin.get_uuid()
   local uuid_pattern = "^%x%x%x%x%x%x%x%x%-%x%x%x%x%-%x%x%x%x%-%x%x%x%x%-%x%x%x%x%x%x%x%x%x%x%x%x$"
   lu.assertStrMatches(uuid, uuid_pattern)
end

function TestUUID:testUUIDNotIncorrectFormat()
   local uuid = httpbin.get_uuid()
   local incorrect_uuid_pattern = "^9999%x%x%x%x%x%x%x%x%-%x%x%x%x%-%x%x%x%x%-%x%x%x%x%-%x%x%x%x%x%x%x%x%x%x%x%x$"
   lu.assertNotEquals(uuid:match(incorrect_uuid_pattern), uuid)
end

local runner = lu.LuaUnit.new()
runner:setOutputType("tap")
local success, failures = runner:runSuite()
os.exit(0)
