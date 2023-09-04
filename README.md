# Gin JWT Boilerplate

This is a basic and easy to use boilerplate for a Golang project using Gin and JWT.

## Features

- [x] JWT Authentication
- [x] JWT Authorization
- [x] Bearer Token
- [x] Refresh Token
- [x] User CRUD
- [x] User Login
- [x] User Logout
- [x] MongoDB Database
- [x] Docker TestContainers

## Requirements

- [Golang](https://golang.org/)
- [Docker](https://www.docker.com/)

Each package contains a readme with more information about the package.

## Usage

Create a `.env` file in the internal/service directory and add the following environment variables:

```env
MONGODB_URI={YOUR_MONGODB_URI}
DB_NAME={YOUR_DB_NAME}
JWT_SECRET={YOUR_JWT_SECRET}
JWT_EXPIRES_IN_HOUR={YOUR_JWT_EXPIRES_IN_HOUR}
JWT_REFRESH_EXPIRES_IN_HOUR={YOUR_JWT_REFRESH_EXPIRES_IN_HOUR}
```