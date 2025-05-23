FROM golang:1.22 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

#COPY cmd/server/*.go ./
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main-app cmd/server/main.go

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/main-app ./

ENV DATABASE_URL postgres://postgres:postgres@host.docker.internal:5432/postgres?sslmode=disable

EXPOSE 8080

CMD ["./main-app"]