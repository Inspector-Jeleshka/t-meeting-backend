# syntax=docker/dockerfile:1
FROM golang:1.25-alpine
WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download -x

COPY . .

RUN go build -o /bin/server ./cmd

# Expose the port that the application listens on.
EXPOSE 33

# What the container should run when it is started.
CMD ["/bin/server"]
