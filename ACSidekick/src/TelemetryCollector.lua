---@class DBLink
---@field dbHttp string
DBLink = {}
DBLink.__index = DBLink

---@param dbHttp string
function DBLink:new(dbHttp)
    local dbLink = {}
    setmetatable(dbLink, DBLink)
    dbLink.dbHttp = dbHttp
    return dbLink
end

---@class Packet
---@field dbHttp string