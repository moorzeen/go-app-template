# syntax=docker/dockerfile:1
FROM golang:latest
RUN go version
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

# make wait-for-postgres.sh executable
RUN chmod +x wait-for-postgres.sh

# build go app
RUN go build -v -o cmd/ cmd/main.go

EXPOSE 8080
CMD [ "cmd/main" ]
