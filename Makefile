_COMPOSE= docker compose -f ./deploy/docker-compose.yaml --env-file ./deploy/compose.env

.PHONY:up
up:
	${_COMPOSE} up --build

.PHONY:dev
dev:
	${_COMPOSE} --profile deps up

.PHONY:down
down:
	${_COMPOSE} down -v --remove-orphans
