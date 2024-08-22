# Go Backend REST API

This repository is a showcase of my personal project for building a Golang REST API. It demonstrates practices in structuring a Go project and handling various tasks such as dependency management, API handling, and database migrations.

Additionally, this project implements integration with the [Midtrans](https://midtrans.com/) payment gateway and Firebase Cloud Messaging (FCM) for push notifications.

Idiomatic Way Format Your Code

```shell script
go fmt ./...
```

Run REST Server

```shell script
go run cmd/semesta-ban/*.go rest
```

Running Local Server

```shell script
docker-compose up -d
```

Running migrations

```shell script
go run cmd/semesta-ban/*.go migrate-up
```

## Directory Structure

- `bootstrap`

  This folder manages dependency.

- `cmd/semesta-ban`

  The main package.

- `cmd/semesta-ban/commands`

  For sub-commands of the main package, using cobra command.

- `internal/api`

  For API related handling

- `pkg`

  For any packages that doesn't have dependencies to other packages in this repository (sub-packages are exception).

- `files/db_migrations`

  For database related changes. Table creation/alteration and data insertion/modification goes here.
  Use `golang-migrate`.

References Api Documentation:

- https://github.com/david1312/go_application_structure/wiki/documentation-api-semesta-ban
