local socket = require("socket")
local ssl = require("ssl")
local dkjson = require("dkjson")

local function get_uuid()
   local host = "httpbin.org"
   local path = "/uuid"

   local sock = assert(socket.connect(host, 443))
   sock = assert(ssl.wrap(sock, { mode = "client", protocol = "tlsv1_2" }))
   assert(sock:dohandshake())

   local request = string.format(
       "GET %s HTTP/1.1\r\n" ..
       "Host: %s\r\n" ..
       "Accept: application/json\r\n" ..
       "Connection: close\r\n\r\n",
       path,
       host
   )

   sock:send(request)

   local response = {}
   while true do
       local s, status = sock:receive()
       if status == "closed" then break end
       response[#response + 1] = s
   end

   sock:close()

   local json_start = 1
   while response[json_start] ~= "" do
       json_start = json_start + 1
   end
   local json_response = table.concat(response, "\n", json_start + 1)

   local data = dkjson.decode(json_response)

   return data.uuid
end

return {
   get_uuid = get_uuid
}
