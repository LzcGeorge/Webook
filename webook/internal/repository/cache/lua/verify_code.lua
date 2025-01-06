local key = KEYS[1]
local expectedCode = ARGV[1]
local code = redis.call("get", key)

-- 重试次数
local cntKey = key..":cnt"
local cnt = tonumber(redis.call("get", cntKey))

if code <= 0 then
    -- 重试次数超过 3 次
    return -1
elseif code == expectedCode then
    -- 验证码正确
    redis.call("set", cntKey, -1)
    return 0
else 
    -- 验证码输入错误，重试次数减 1
    redis.call("incr", cntKey)
    return -2
end

