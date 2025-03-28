FROM golang:1.21 AS build

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o Verses

FROM build as test

WORKDIR /app

RUN go test -v ./...

FROM --platform=linux/amd64 debian:stable-slim

RUN apt-get update && apt-get install -y ca-certificates

COPY --from=build /app/Verses /usr/bin/Verses

EXPOSE 8080

CMD ["Verses"]

