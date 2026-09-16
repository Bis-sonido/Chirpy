# Chirpy
BootDev - Chirpy Project

A Golang micro blogging API for posting and managing short messages.

## Installation

Make sure you have the latest [Go toolchain] (https://golang.org/dl/) installed as well a local PostgresSQL database. You can always check if Go or PostgresSQL were installed by checking their latest version:

Golang(Go)
```go version```

PostgresSQL
```psql --version```

You will also need to install Goose to run migrations.

`go install github.com/pressly/goose/v3/cmd/goose@latest`
check version:
`goose -version`

## How to Run

Boot the server
```go run .```

Server starts on http://localhost:8080

## Use / Interact with API

Below are a few examples to interact with the server using `curl` or any other API client of your choosing.

**Example: Register a User**
```sh
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "securepassword"}'
```

**Example: Log In**
```sh
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "securepassword"}'
```

**Example: Creating an Authenticated Chirp**
```sh
curl -X POST http://localhost:8080/api/chirps \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"body": "Hello world!"}'
```

## Cheat Sheet Table
| Method | Endpoint | Description | Auth Required |
| :----- | :-------: | :---------: | :------------: |
|`GET`| `/api/healthz` | Readiness Check | No |
|`POST` | `api/users` | Create User | No |
|`POST` | `api/login` | Authenticate & get JWT | No |
|`POST` | `api/chirps` | Create a chirp | Yes |
|`GET` | `api/chirps` | List all chirps | No |
|`DELETE` | `api/chirps/{id}` | Delete a chirp | Yes |