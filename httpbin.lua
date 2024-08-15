local socket = require("socket")
local ssl = require("ssl")
local dkjson = require("dkjson")

local function get_uuid()
    -- Define the URL for the GET request
    local host = "httpbin.org"
    local path = "/uuid"

    -- Create a socket and connect to the server
    local sock = assert(socket.connect(host, 443))

    -- Upgrade the connection to SSL/TLS
    sock = assert(ssl.wrap(sock, { mode = "client", protocol = "tlsv1_2" }))
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

    -- Return the UUID from the parsed data
    return data.uuid
end

local function assert_match(pattern, value, message)
    if not string.match(value, pattern) then
        error(message or string.format("'%s' does not match pattern '%s'", value, pattern))
    end
end

local function assert_no_match(pattern, value, message)
    if string.match(value, pattern) then
        error(message or string.format("'%s' unexpectedly matches pattern '%s'", value, pattern))
    end
end

local function test_uuid_format_correct()
    local uuid = get_uuid()
    print("Received UUID:", uuid)
    local uuid_pattern = "^%x%x%x%x%x%x%x%x%-%x%x%x%x%-%x%x%x%x%-%x%x%x%x%-%x%x%x%x%x%x%x%x%x%x%x%x$"
    assert_match(uuid_pattern, uuid, "UUID does not match expected format")
    print("Correct UUID format test passed")
end

local function test_uuid_format_incorrect()
    local uuid = get_uuid()
    print("Received UUID:", uuid)
    local incorrect_uuid_pattern = "^9999%x%x%x%x%x%x%x%x%-%x%x%x%x%-%x%x%x%x%-%x%x%x%x%-%x%x%x%x%x%x%x%x%x%x%x%x$"
    assert_no_match(incorrect_uuid_pattern, uuid, "UUID unexpectedly matches incorrect pattern")
    print("Incorrect UUID format test passed")
end

local function run_tests(...)
    local tests = { ... }
    local pass_count, fail_count = 0, 0

    print("Starting tests...")

    for i, test_func in ipairs(tests) do
        local status, error_msg = pcall(test_func)
        if status then
            pass_count = pass_count + 1
            print(string.format("Test %d passed", i))
        else
            fail_count = fail_count + 1
            print(string.format("Test %d failed: %s", i, error_msg))
        end
    end

    print(string.format("Tests completed. Passed: %d, Failed: %d", pass_count, fail_count))

    -- Uncomment these lines if you want the script to exit with a non-zero status on test failures
    -- if fail_count > 0 then
    --     os.exit(1)
    -- else
    --     os.exit(0)
    -- end
end

run_tests(test_uuid_format_correct, test_uuid_format_incorrect)
