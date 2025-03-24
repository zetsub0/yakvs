up:
	docker compose up -d

down:
	docker compose down
	docker volume rm yatc_tarantool_data
	docker rmi yatc-tarantool