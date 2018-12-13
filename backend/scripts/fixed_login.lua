counter = 1

request = function()
    uid = counter + 1
    username = "admin" .. uid
    password = "111111"
    wrk.method = "POST"
    wrk.body = "username=" .. username .. "&password=" .. password
    wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"
    return wrk.format()
end