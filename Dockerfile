FROM node:22-alpine AS node-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

FROM golang:1.26-alpine AS go-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=node-builder /app/ui/static ./ui/static
RUN CGO_ENABLED=0 go build -o /bin/server ./cmd/web

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=go-builder /bin/server ./server
COPY migrations/ ./migrations/
RUN mkdir -p ./data
EXPOSE 8080
CMD ["./server"]
