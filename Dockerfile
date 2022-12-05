# syntax=docker/dockerfile:1
FROM golang:latest
RUN go version
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .

# build go app
RUN go build -v -o cmd/ cmd/main.go

EXPOSE 8080
CMD [ "cmd/main" ]
