uid = math.random(1, 10000000)
username = "admin" .. uid
password = "111111"
wrk.method = "POST"
wrk.body = "username=" .. username .. "&password=" .. password
wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"