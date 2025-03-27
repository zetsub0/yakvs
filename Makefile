_COMPOSE= docker compose -f ./deploy/docker-compose.yaml --env-file ./deploy/compose.env

.PHONY:up
up:
	${_COMPOSE} up --build -d

.PHONY:dev
dev:
	${_COMPOSE} --profile deps up -d

.PHONY:down
down:
	${_COMPOSE} down -v --remove-orphans
