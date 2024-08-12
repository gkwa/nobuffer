-- Load the socket module
local socket = require("socket")

-- Define the URL for the GET request
local url = "https://httpbin.org/uuid"

-- Create a socket and connect to the server
local host = "httpbin.org"
local sock = assert(socket.connect(host, 443))

-- Upgrade the connection to SSL/TLS
local ssl = require("ssl")
sock = assert(ssl.wrap(sock, {mode="client", protocol="tlsv1_2"}))
assert(sock:dohandshake())

-- Prepare the GET request
local request = string.format(
    "GET %s HTTP/1.1\r\n" ..
    "Host: %s\r\n" ..
    "Accept: application/json\r\n" ..
    "Connection: close\r\n\r\n",
    "/uuid",
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

-- Simple JSON parsing (for demonstration purposes)
local function parse_json(json_str)
    local t = {}
    for k, v in json_str:gmatch('"([^"]+)"%s*:%s*"([^"]+)"') do
        t[k] = v
    end
    return t
end

-- Parse the JSON response
local data = parse_json(json_response)

-- Fetch the UUID from the dictionary
local uuid = data["uuid"]

print("Fetched UUID: " .. tostring(uuid))
