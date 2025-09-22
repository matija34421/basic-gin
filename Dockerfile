FROM golang:1.25-alpine AS build
WORKDIR /src


COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /out/app ./cmd/api

FROM alpine:3.20
RUN adduser -D -g '' app && apk add --no-cache ca-certificates
USER app
WORKDIR /home/app

COPY --from=build /out/app /usr/local/bin/app
EXPOSE 8080
ENTRYPOINT ["app"]
