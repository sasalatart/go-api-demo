# Go API Demo

This repository contains a demo API built with [go][golang]. Its main purposes are:

1. Show how the standard library may be used to build an API, and
2. Show how Docker image sizes may be minimized by using multi-stage builds.

The API listens on port 3000 and exposes a `/ping` endpoint for `GET`, `POST`, `PUT`, `PATCH` and
`DELETE` HTTP requests. It responds with the values of the request's HTTP method, query parameters
and body, in case there was one.

Example request:

```sh
curl --location --request POST 'localhost:3000/ping?qs1-key=qs1-value&qs2-key=qs2-value' \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --data-urlencode 'body1-key=body1-value' \
  --data-urlencode 'body2-key=body2-value'
```

Example response:

```json
{
  "method": "POST",
  "queryParams": {
    "qs1-key": ["qs1-value"],
    "qs2-key": ["qs2-value"]
  },
  "body": {
    "body1-key": ["body1-value"],
    "body2-key": ["body2-value"]
  }
}
```

## Running locally

```sh
go build -o app && ./app
```

## Docker

Three Dockerfiles have been added to the repository:

- `Dockerfile.full`: contains a basic setup that uses the `golang:1.14.3` base image.
- `Dockerfile.alpine`: same setup as `Dockerfile.full`, but based on a golang alpine image.
- `Dockerfile`: uses a multi-stage build to optimize image size.

The resulting image sizes are:

- `Dockerfile.full`: 817MB
- `Dockerfile.alpine`: 378MB
- `Dockerfile`: 13.3MB

In order to try out the multi-stage image, you must run the following:

```sh
# Build the image
docker build -t sasalatart/go-api-demo:latest .

# Run a container
docker run -p 3000:3000 sasalatart/go-api-demo:latest
```

The API will then be listening on your machine's port 3000.

## Tests

```sh
go test
```

[golang]: https://golang.org/
