# Go application template

[![CI/CD](https://github.com/moorzeen/go-app-template/actions/workflows/ci-cd.yml/badge.svg?branch=master)](https://github.com/moorzeen/go-app-template/actions/workflows/ci-cd.yml)

### Description
**Template** is a simple Go application template using popular technologies and approaches.
It can be taken as the basis of your own project.

As an example, this application template implements REST and gRPC API for user registration and authorization.

#### Under the hood
- gRPC/REST server;
- error handling and logging using Zerolog;
- working with environment variables and env files;
- JWT authorization;
- working with Postgres, Redis, Kafka;
- autotests;
- launch the application in Docker.

### Launching the app
1. Clone the repository as usually:
```
git clone...
```

2. Assemble and run the Docker container with the command:
```
docker-compose...
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
After successful registration, the user is automatically authenticated.

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

## Suggest improvements
...