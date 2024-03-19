build:
	docker-compose build film-api

run:
	docker-compose up film-api -d

migrate:
	migrate -path ./db/migration -database 'postgres://api_tester:testing@0.0.0.0:5436/film_api?sslmode=disable' up

add_keys:
	psql -h localhost -p 5436 -d film_api -U api_tester  -W
	testing
	INSERT INTO public.users (id, role, api_key) VALUES (1, 'admin', 'root');

swag:
	swag init -g cmd/main.go