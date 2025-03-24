box.cfg {}

username = os.getenv("USERNAME")
pass = os.getenv("PASSWORD")


if box.schema.user.exists(username) == false then
    box.schema.user.create(username, { password = pass })
end

box.schema.user.grant(username, 'read,write,execute', 'universe')

if box.space.kv == nil then
    box.schema.space.create('kv', {
        format = {
            {name = 'id', type = 'string'},
            {name = 'value', type = 'any'}
        }
    })
    box.space.kv:create_index('primary', {parts = {'id'}})
end