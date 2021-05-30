FROM golang:1.15.12-alpine3.12 as builder
COPY . /app
WORKDIR /app
RUN go build -o app main.go

FROM alpine
COPY --from=builder /app/app /usr/local/bin/app
ENTRYPOINT ["app"]