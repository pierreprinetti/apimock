# apimock
[![Go Report Card](https://goreportcard.com/badge/github.com/pierreprinetti/apimock)](https://goreportcard.com/report/github.com/pierreprinetti/apimock)
[![Build Status](https://travis-ci.org/pierreprinetti/apimock.svg?branch=master)](https://travis-ci.org/pierreprinetti/apimock)

This is a very basic fake API server. I use it to build the front-end of web applications, without the need for the backend to be ready.

It is an in-memory, non-persistent key-value store you can fill with `PUT` requests, where the request path is the key and the request body is the value.
Retrieve the saved value with a subsequent `GET` request at the same endpoint.

_apimock_ will serve back the same `Content-Type` is has received. If no `Content-Type` header was sent with the `PUT` request, the `DEFAULT_CONTENT_TYPE` environment variable will be sent.

_apimock_ is meant for prototyping. **Please do NOT use it in production**.

## Example:

    $ HOST=localhost:8800 apimock & # same as PORT=8800 apimock &
    $ curl -X PUT -d '{"message": "This is not a pipe"}' localhost:8800/my/endpoint
    > {"message": "This is not a pipe"}
    $ curl -X GET localhost:8800/my/endpoint
    > {"message": "This is not a pipe"}
    $ curl -X DELETE localhost:8800/my/endpoint
    $ curl -X GET localhost:8800/my/endpoint
    $

## Content-Type
Apimock will remember the `Content-Type` associated with every request. This behaviour can be modified with the environment variables:

- `DEFAULT_CONTENT_TYPE`: When the `PUT` request doesn't bear a `Content-Type`, this one will be used. If not specified, this is `text/plain`.
- `FORCED_CONTENT_TYPE`: The specified string will be used as `Content-Type` no matter what is transmitted with the `PUT` request.

## Docker container

    docker run --name apimock -p 8800:8800 -d pierreprinetti/apimock:latest

## Features

It currently supports:
- [x] CORS headers (responses always bear `Allow-Origin: *` and a bunch of authorized headers and methods)
- [x] `OPTIONS`
- [x] `PUT`
- [x] `GET`
- [x] `DELETE`
- [x] `Content-Type` header

What it might support in the future:
- [ ] `POST` to an endpoint with fake ID generator (e.g. `POST` to `example.com/items` would result in the storage of the element in `example.com/items/1`
