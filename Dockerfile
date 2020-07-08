# Accept the Go version for the image to be set as a build argument.
# Default to Go 1.12
ARG GO_VERSION=1.12

# First stage: build the executable.
FROM golang:${GO_VERSION}-alpine AS builder

# Create the user and group files that will be used in the running container to
# run the process as an unprivileged user.
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /src

# Fetch dependencies first; they are less susceptible to change on every build
# and will therefore be cached for speeding up the next build
COPY ./go.mod ./go.sum ./
RUN go mod download

# Import the code from the context.
COPY ./ ./

# Build the executable to `/app`. Mark the build as statically linked.
RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -o /app .

# Final stage: the running container.
FROM scratch AS final

# Import the user and group files from the first stage.
COPY --from=builder /user/group /user/passwd /etc/

# Import the compiled executable from the first stage.
COPY --from=builder /app /app

ARG BUILD_DATE
ARG VCS_REF
LABEL \
    org.opencontainers.image.created=$BUILD_DATE \
    org.opencontainers.image.authors="https://pierreprinetti.com" \
    org.opencontainers.image.url="https://hub.docker.com/r/pierreprinetti/apimock/" \
    org.opencontainers.image.source="https://github.com/pierreprinetti/apimock" \
    org.opencontainers.image.version=$VCS_REF \
    org.opencontainers.image.revision=$VCS_REF \
    org.opencontainers.image.title="apimock" \
    org.opencontainers.image.description="A mock API server"

# Declare the port on which the webserver will be exposed.
# As we're going to run the executable as an unprivileged user, we can't bind
# to ports below 1024.
EXPOSE 8080

# Perform any further action as an unprivileged user.
USER nobody:nobody

# Run the compiled binary.
ENTRYPOINT ["/app"]
