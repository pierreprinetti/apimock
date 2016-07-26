# mockapi

This is a minimal fake API server.

It is an in-memory, non-persistent key-value store you can fill with `PUT` requests, where the request path is the key and the request body is the value.
Retrieve the saved value with a subsequent `GET` request at the same endpoint.

