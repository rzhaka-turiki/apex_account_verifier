FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /apex-verifier ./cmd/server


FROM alpine:3.22

WORKDIR /app

COPY --from=builder /apex-verifier .

EXPOSE 50051

CMD ["./apex-verifier"]