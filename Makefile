build:
	docker-compose build film-api

run:
	docker-compose up film-api -d

migrate:
	migrate -path ./db/migration -database 'postgres://api_tester:testing@0.0.0.0:5436/film_api?sslmode=disable' up

add_keys:
	go run ./db

swag:
	swag init -g ./cmd/api/main.go