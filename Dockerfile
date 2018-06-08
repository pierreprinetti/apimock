FROM golang:1.10 AS builder
WORKDIR $GOPATH/src/github.com/pierreprinetti/apimock
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app .

FROM scratch
COPY --from=builder /app ./
EXPOSE 80
ENTRYPOINT ["./app"]
