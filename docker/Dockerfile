FROM golang:1.22.6-bookworm AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Cloud Run sets $PORT; your code already reads it.
# Build a static binary (fast cold starts, no shell)
ARG TARGETOS=linux
ARG TARGETARCH=amd64
ENV CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH
RUN go build -o /out/app ./cmd/api 2>/dev/null || go build -o /out/app .

# ---- Runtime stage (distroless, non-root) ----
FROM gcr.io/distroless/static-debian12:latest

# Minimal envs Cloud Run likes; PORT is injected at runtime
ENV PORT=8080
WORKDIR /app
COPY --from=builder /out/app /app/app

USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/app/app"]
