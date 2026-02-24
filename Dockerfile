# syntax=docker/dockerfile:1
FROM golang:1.25-alpine as build
WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download -x

COPY . .

RUN go build -o /bin/server ./cmd

FROM scratch
WORKDIR /app

COPY --from=build /bin/server /bin/server

# Expose the port that the application listens on.
EXPOSE 33

# What the container should run when it is started.
CMD ["/bin/server"]
