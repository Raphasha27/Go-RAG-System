FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/rag-system .

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/rag-system .

USER appuser

ENV DATABASE_URL=""
ENV OPENAI_API_KEY=""

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD ["/app/rag-system", "-health"] || exit 1

ENTRYPOINT ["/app/rag-system"]
