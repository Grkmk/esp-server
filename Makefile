build:
	docker-compose build

server: build
	docker-compose up