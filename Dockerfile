FROM golang:1.21.8-alpine3.18 as builder
WORKDIR /app
COPY . /app
RUN go build -o main .

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main /app/main
CMD ["/app/main"]

