local socket = require("socket")
local ssl = require("ssl")
local dkjson = require("dkjson")

-- Define the URL for the GET request
local host = "httpbin.org"
local path = "/uuid"

-- Create a socket and connect to the server
local sock = assert(socket.connect(host, 443))

-- Upgrade the connection to SSL/TLS
sock = assert(ssl.wrap(sock, {mode="client", protocol="tlsv1_2"}))
assert(sock:dohandshake())

-- Prepare the GET request
local request = string.format(
    "GET %s HTTP/1.1\r\n" ..
    "Host: %s\r\n" ..
    "Accept: application/json\r\n" ..
    "Connection: close\r\n\r\n",
    path,
    host
)

-- Send the request
sock:send(request)

-- Receive the response
local response = {}
while true do
    local s, status = sock:receive()
    if status == "closed" then break end
    response[#response + 1] = s
end

-- Close the socket
sock:close()

-- Find the JSON part of the response (assuming it's after the headers)
local json_start = 1
while response[json_start] ~= "" do
    json_start = json_start + 1
end
local json_response = table.concat(response, "\n", json_start + 1)

-- Parse the JSON response using dkjson
local data = dkjson.decode(json_response)

-- Fetch the UUID from the parsed data
local uuid = data.uuid

-- Assert that the UUID matches the expected format
local uuid_pattern = "^%x%x%x%x%x%x%x%x%-%x%x%x%x%-%x%x%x%x%-%x%x%x%x%-%x%x%x%x%x%x%x%x%x%x%x%x$"
assert(uuid:match(uuid_pattern), "UUID does not match expected format")

print("UUID format validation passed")
