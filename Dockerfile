FROM golang:1.15.2 AS build

WORKDIR /src/
COPY main.go /src/
RUN go build -o /bin/http

FROM ubuntu:latest
COPY --from=build /bin/http /bin/http
EXPOSE 3000
ENTRYPOINT ["/bin/http", "-folder", "/web"]