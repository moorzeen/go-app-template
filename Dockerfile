# syntax=docker/dockerfile:1
FROM golang:latest

RUN curl -fsSLO https://get.docker/builds/Linux/x86_64/docker-17.04.0-ce.tgz \
  && tar xzvf docker-17.04.0-ce.tgz \
  && mv docker/docker /usr/local/bin \
  && rm -r docker docker-17.04.0-ce.tgz

RUN go version
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .

# build go app
RUN go build -v -o cmd/ cmd/main.go

EXPOSE 8080
CMD [ "cmd/main" ]
