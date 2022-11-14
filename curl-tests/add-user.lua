request = function()
    headers = {}
    headers["Content-Type"] = "application/json"
    body   = '{"username": "user-' .. tostring(math.random(-1000000000, 1000000000)) ..  '"}'

    return wrk.format("POST", "/api/v1/users", headers, body)
end

