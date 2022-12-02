# syntax=docker/dockerfile:1
FROM golang:latest

# move to working directory /app
WORKDIR /app

# copy, download and verify dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download && go mod verify

# copy the code into the container
COPY . .

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

# make wait-for-postgres.sh executable
RUN chmod +x wait-for-postgres.sh

# compile application
RUN go build -v -o cmd/ cmd/main.go

EXPOSE 8080

CMD [ "cmd/main" ]
