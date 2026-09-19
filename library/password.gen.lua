---@meta password

---@class password
---@field DEFAULT_COST number recommended bcrypt cost for production use (tests should use a much lower cost, e.g. 4, to stay fast)
local password = {}

--- hashes a password with bcrypt at the given cost; raises if cost is out of range or the password exceeds 72 bytes
---@param plaintext string the password to hash; must not exceed 72 bytes (bcrypt's own limit)
---@param cost number bcrypt cost factor in [4, 31]; use password.DEFAULT_COST unless you have a specific reason not to
---@return string hash a self-describing bcrypt hash ($2a$<cost>$...) suitable for storage
function password.hash(plaintext, cost) end

--- verifies a password against a stored bcrypt hash in constant time; returns false on a wrong password, raises on a malformed hash
---@param plaintext string the password attempt
---@param hash string a bcrypt hash previously produced by password.hash; this is application data, not attacker input, so a malformed value raises rather than returning false
---@return boolean ok true if plaintext matches hash
function password.verify(plaintext, hash) end

--- returns the bcrypt cost embedded in a hash, for deciding whether to rehash on login; raises on a malformed hash
---@param hash string a bcrypt hash previously produced by password.hash
---@return number cost the cost factor the hash was created with
function password.cost(hash) end

return password
