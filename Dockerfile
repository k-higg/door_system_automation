FROM golang:1.25-alpine as builder
LABEL authors="Kenny Higgins"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/goapp ./main.go

from alpine:latest

WORKDIR /app

COPY --from=builder /app/goapp .
COPY --from=builder /app/resources ./resources

EXPOSE 8080

CMD ["./goapp"]