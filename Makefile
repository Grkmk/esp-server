build-server:
	docker build .

prod: build-server
	docker-compose -p espresso up

dev:
	docker-compose -p espresso run --service-ports postgres

go:
	go run main.go