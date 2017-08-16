#!/usr/bin/env tarantool

box.cfg{
    -- log_level
    -- 1 – SYSERROR
    -- 2 – ERROR
    -- 3 – CRITICAL
    -- 4 – WARNING
    -- 5 – INFO
    -- 6 – DEBUG
    log_level = 5,

    slab_alloc_arena = 1,
    -- wal_dir='xlog',
    -- snap_dir='snap',
}
local prefix = 'ff_'
local log = require("log")

log.info('Info %s', box.info.version)

--------------------
-- Users
--------------------

s = box.schema.create_space(prefix..'chats', {
    if_not_exists=true,
    })
s:create_index('primary', {
    if_not_exists=true,
    type = 'tree',
    unique = true,
    parts = {1, 'string'},
})