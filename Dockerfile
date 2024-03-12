
FROM golang:1.22.1-alpine3.19 AS build

WORKDIR /app
COPY . /app

#RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o app -ldflags="-extldflags=-static" ./cmd/cli.go

FROM gcr.io/distroless/static-debian11

WORKDIR /app

COPY --from=build /app/app /app/app
COPY --from=build /app/app /app/app

ENTRYPOINT [ "/app/app" ]
