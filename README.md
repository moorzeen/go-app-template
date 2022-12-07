# Go application template

[![CI/CD](https://github.com/moorzeen/go-app-template/actions/workflows/ci-cd.yml/badge.svg?branch=master)](https://github.com/moorzeen/go-app-template/actions/workflows/ci-cd.yml)

This is a simple Go application template using popular technologies and approaches.
It can be taken as the basis of your own project.

As an example, this application template implements REST and gRPC API for user registration and authorization.

#### Under the hood
- gRPC/REST server;
- Postgres, Redis, Kafka;
- Zerolog;
- Docker containerization with CI/CD;

### Launching the app locally
1. Clone the repository as usually and go to project directory:
```
git clone https://github.com/moorzeen/go-app-template
cd go-app-template
```

2. Run Postgres, Redis and Kafka as the Docker containers:
```
docker-compose -f docker-compose.local.yml up -d --build
````

3. Build and run Go application:
```
export GO111MODULE="on"
go run cmd/main.go
````

3. After successfully launching the application, you will see the lines:
````
03.11.22 18:27:01 INF internal/app/app.go:90 > grpc server started on port 8081
03.11.22 18:27:01 INF internal/app/app.go:92 > rest server started on port 8082
````

### HTTP API

There are two options for sending a request to the server: REST and gRPC

#### **User Registration**

Handler: `POST /register`

Registration is performed using a login/password pair. Each login is unique.

Request format:
````
POST /register HTTP/1.1
Content-Type: application/json
...
{
    "login": "<login>",
    "password": "<password>"
}
````

As a response, the handler returns the registered user in the format:
````
{
    "statusCode": 201,
    "message": "user successfully registered",
    "user": {
        "id": "1e3165e5-a958-41bd-89eb-c2233178cbe5",
        "login": "user",
        "active": false,
        "createdAt": "22-11-03 12:24:31.116",
        "updatedAt": "22-11-03 12:24:31.116",
        "lastLoginAt": "01-01-01 00:00:00.000"
    }
}
````

Possible response codes:
- `201` — the user has been successfully registered;
- `400` — invalid request format;
- `409` — login is already occupied;
- `500` — internal server error.

## To Do
1. [ ] complete documentation; 
2. [ ] implement authorization methods;
3. [ ] full cover with tests including benchmark test;
3. [ ] implement a load performance testing tool.

## Contributing
We would love for you to contribute to `moorzeen/go-app-template`, pull requests are welcome!