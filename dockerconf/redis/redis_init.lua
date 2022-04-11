-- sets all fields for a hash from a dictionary
local hmset = function(key, dict)
    if next(dict) == nil then
        return nil
    end
    local bulk = {}
    for k, v in pairs(dict) do
        table.insert(bulk, k)
        table.insert(bulk, v)
    end
    return redis.call('HMSET', key, unpack(bulk))
end

local default_user = {
    { id = 1, name = 'mockUser1'},
    { id = 2, name = 'mockUser2'},
    { id = 3, name = 'mockUser3'},
    { id = 4, name = 'mockUser4'},
    { id = 5, name = 'mockUser5'},
    { id = 6, name = 'mockUser6'},
}

for i, v in ipairs(default_user) do
    hmset('user:user_id:' .. i, v)
end
