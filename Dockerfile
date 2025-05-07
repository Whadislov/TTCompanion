# Step 1 : Build du backend
FROM golang:1.23.8 AS builder
WORKDIR /src

COPY api api
COPY cmd/web/main.go cmd/web/main.go
COPY internal/my_db internal/my_db
COPY internal/my_functions internal/my_functions
COPY internal/my_types internal/my_types
COPY credentials.env ./
COPY config_app.json ./
COPY go.mod go.sum ./

# Compile backend
RUN CGO_ENABLED=0 GOARCH=amd64 go build -o backend ./cmd/web
# Step 2 : launch app

FROM alpine:latest
#FROM nginx:latest
WORKDIR /app


# Install Curl for debugging
RUN apk add --no-cache curl

# Compiled backend
COPY --from=builder /src/backend ./backend

# Copy WASM + config
COPY wasm ./wasm
COPY credentials.env ./
COPY config_app.json ./

# Permissions
RUN chmod +x /app/backend

# Ports
EXPOSE 8080

# Start app
CMD ["/app/backend"]