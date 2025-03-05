-- redis 的 map 结构对应的 key
local key = KEYS[1]
-- cntKey：map 中的操作字段
local cntKey = ARGV[1]
-- delta：表示增量（-1 、+1）
local delta = tonumber(ARGV[2])

-- 检查 key 是否存在
local exists = redis.call("EXISTS", key)
if exists == 1 then
    -- HINCRBY 命令会自动处理字段不存在的情况，
    -- 如果字段不存在，它会先将字段设置为 0，然后再递增。
    redis.call("HINCRBY", key, cntKey, delta)
    return 1
else
    -- Key 不存在或缓存过期
    return 0
end
