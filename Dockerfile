FROM golang:latest

WORKDIR /app

COPY ./ /app



# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client


RUN go mod download
RUN go build -o films-api ./cmd/api


RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/api/main.go

ENTRYPOINT go run ./cmd/api