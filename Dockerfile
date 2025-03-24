FROM tarantool/tarantool:2.10

WORKDIR /opt/tarantool

COPY init.lua /opt/tarantool/

CMD ["tarantool", "init.lua"]
