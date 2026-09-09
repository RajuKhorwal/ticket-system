# --- Stage 1: Build the Go binary ---
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy dependency files first so Docker can cache downloaded modules
# between builds (faster rebuilds when only application code changes).
COPY go.mod go.sum ./
RUN go mod download

# Now copy the rest of the source code and build a static binary.
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ticket-system .

# --- Stage 2: Minimal runtime image ---
FROM alpine:latest

WORKDIR /app

# Copy only the compiled binary from the builder stage.
# This keeps the final image small and does not include the Go toolchain.
COPY --from=builder /app/ticket-system .

EXPOSE 8080

CMD ["./ticket-system"]
