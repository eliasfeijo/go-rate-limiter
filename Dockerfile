# Stage 1: Build the binary
FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -C cmd -o app

# Stage 2: Create the final image
FROM scratch

COPY --from=builder /app/cmd/app /bin/app

CMD ["/bin/app"]
