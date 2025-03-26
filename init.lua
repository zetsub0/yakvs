box.cfg {}

username = os.getenv("USERNAME")
pass = os.getenv("PASSWORD")

if box and box.user then
    if not box.user.exists(username) then
        box.user.create(username, { password = pass })
    end

    local function safe_grant(username, privileges, object)
        local ok, err = pcall(function()
            box.user.grant(username, privileges, object)
        end)
        if not ok and err.code ~= ER_PRIV_GRANTED then
            error(err)
        end
    end

    safe_grant(username, 'read,write,execute', 'universe')
else
    print("Error: box.user is not available.")
end

if box and box.space then
    if box.space.kv == nil then
        box.space.create('kv', {
            format = {
                {name = 'id', type = 'string'},
                {name = 'value', type = 'any'}
            }
        })
        box.space.kv:create_index('primary', {parts = {'id'}})
    end
else
    print("Error: box.space is not available.")
end