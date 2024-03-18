FROM golang:latest

WORKDIR /app

COPY ./ /app

RUN go mod download

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

ENTRYPOINT go run ./cmd/api