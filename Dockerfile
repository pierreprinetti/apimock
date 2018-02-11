FROM golang:1.9-alpine

EXPOSE 80

WORKDIR $GOPATH/src/github.com/pierreprinetti/apimock

COPY . .

RUN go build -o app

ENTRYPOINT ./app
