local key = KEYS[1]
local inputCode = ARGV[1]
local code = redis.call("get", key)

-- 重试次数
local cntKey = key..":cnt"
local cnt = tonumber(redis.call("get", cntKey))

if cnt == nil or cnt <= 0 then
    -- 重试次数超过 3 次
    return -1
elseif code == inputCode then
    -- 验证码正确
    redis.call("set", cntKey, -1)
    return 0
else 
    -- 验证码输入错误，重试次数减 1
    redis.call("decr", cntKey)
    return -2
end

