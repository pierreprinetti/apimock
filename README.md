# apimock

This is a very basic fake API server. I use it to build the front-end of web applications, without the need for the backend to be ready.

It is an in-memory, non-persistent key-value store you can fill with `PUT` requests, where the request path is the key and the request body is the value.
Retrieve the saved value with a subsequent `GET` request at the same endpoint.

It is meant for prototyping. **Please do NOT use _apimock_ in production**.

## Example:

    $ HOST=localhost:8800 apimock &
    $ curl -X PUT -d '{"message": "This is not a pipe"}' localhost:8800/my/endpoint
    > {"message": "This is not a pipe"}
    $ curl -X GET localhost:8800/my/endpoint
    > {"message": "This is not a pipe"}
    $ curl -X DELETE localhost:8800/my/endpoint
    $ curl -X GET localhost:8800/my/endpoint
    $

## Docker container

    docker run --name apimock -p 8800:80 -d pierreprinetti/apimock:latest

## Features

It currently supports:
- [x] `PUT`
- [x] `GET`
- [x] `DELETE`

What it might support in the future:
- [ ] `Content-Type` headers
- [ ] `POST` to an endpoint with fake ID generator (e.g. `POST` to `example.com/items` would result in the storage of the element in `example.com/items/1`

