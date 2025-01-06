local key = KEYS[1]
-- 验证次数，格式为：phone_code:login:152xxxxxxxx:cnt
local cntKey = key..":cnt"
local val = ARGV[1]

local ttl = tonumber(redis.call("ttl", key))
if ttl == -1 then
    -- key 存在，但过期时间未设置
    return -2
elseif ttl == -2 or ttl < 540 then
    -- ttl 小于 540 秒(过了 1 分钟，则重新设置验证码)
    redis.call("set", key, val)
    redis.call("expire", key, 600)
    redis.call("set", cntKey, 3)
    redis.call("expire", cntKey, 600)
    return 0
else
    -- 操作太频繁
    return -1
end

