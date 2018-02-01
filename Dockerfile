FROM golang:1.9-alpine

EXPOSE 80

RUN apk add --no-cache git
ADD https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

WORKDIR $GOPATH/src/github.com/pierreprinetti/apimock

COPY . .

RUN dep ensure

RUN go build -o app

ENTRYPOINT ./app
