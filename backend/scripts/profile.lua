request = function()
   uid = math.random(1, 10000000)
   -- path = "/backend/index?token=test&uid=" .. uid
   path = "/api/users/profiles?token=test&uid=" .. uid
   return wrk.format(nil, path)
end