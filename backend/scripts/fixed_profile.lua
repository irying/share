local counter = 1

request = function()
   uid = counter + 1
   -- path = "/backend/index?token=test&uid=" .. uid
   path = "/api/users/profiles?token=test&uid=" .. uid
   return wrk.format(nil, path)
end