FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go test -v

RUN go build -o ./out/k-taxes .

FROM alpine:latest

ENV PORT=8080

ENV DATABASE_URL=postgresql://postgres:postgres@postgres:5432/ktaxes?sslmode=disable

ENV ADMIN_USERNAME=adminTax

ENV ADMIN_PASSWORD=admin!

COPY --from=builder /app/out/k-taxes /app/k-taxes

EXPOSE 8080

CMD ["/app/k-taxes"]