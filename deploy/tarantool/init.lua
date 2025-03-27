box.cfg {}

username = os.getenv("TARANTOOL_USER_NAME")
pass = os.getenv("TARANTOOL_USER_PASSWORD")

box.schema.space.create('kv', {if_not_exists = true})

box.space.kv:format({{ name = 'id', type = 'string'}, { name = 'band_name', type = 'any' }})

box.space.kv:create_index('primary', {parts = {'id'}, if_not_exists=true})

box.schema.user.create('user',{password = 'password',if_not_exists = true})
