# ---------- STAGE 1: Build ----------
FROM golang:1.24 AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o taskflow ./main.go

# ---------- STAGE 2: Dev/Test ----------
FROM builder AS dev
CMD ["go", "run", "./main.go"]

# ---------- STAGE 3: Run Tests ----------
FROM builder AS test
RUN go test -v ./...

# ---------- STAGE 4: Production ----------
FROM gcr.io/distroless/base-debian11 AS prod

WORKDIR /

COPY --from=builder /app/taskflow /taskflow

USER nonroot:nonroot
ENTRYPOINT ["/taskflow"]
