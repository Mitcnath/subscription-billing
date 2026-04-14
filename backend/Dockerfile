# ---- Stage 1: Builder ----
# Use the official Go Alpine image as the build environment.
# Named "builder" so the second stage can reference it with --from=builder.
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy dependency manifests first and download modules.
# This is a separate step so Docker can cache the module layer —
# it only re-downloads if go.mod or go.sum changes.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code.
COPY . .

# Compile both binaries.
# CGO_ENABLED=0 — disables cgo, producing a fully static binary with no
# dynamic C library dependencies. Required for Alpine (musl libc).
# GOOS=linux   — targets Linux explicitly (important when building from macOS).
# -o <name>    — names the output binary.
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o seed ./cmd/seed

# ---- Stage 2: Runner ----
# Minimal Alpine image — no Go toolchain, no source code, just the binaries.
# This keeps the final image small (typically under 20MB).
FROM alpine:3.20

WORKDIR /app

# Copy only the compiled binaries from the builder stage.
# --from=builder references Stage 1 by its alias.
COPY --from=builder /app/server .
COPY --from=builder /app/seed .

# Document that the server listens on 8080 (does not publish the port itself —
# that is handled in docker-compose.yml).
EXPOSE 8080

# Default command — runs the API server.
CMD ["./server"]