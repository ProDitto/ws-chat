# Stage 1: Build frontend assets
FROM node:18-alpine AS frontend-builder
WORKDIR /app/frontend
COPY quikchat-frontend/package*.json ./
RUN npm install
COPY quikchat-frontend/ ./
RUN npm run build

# Stage 2: Build backend application
FROM golang:1.24-alpine AS backend-builder
WORKDIR /app
# Copy Go module files
COPY quikchat-backend/go.mod quikchat-backend/go.sum ./
RUN go mod download
# Copy backend source
COPY quikchat-backend/ ./
# Copy frontend assets from the previous stage
COPY --from=frontend-builder /app/frontend/public ./web
COPY --from=frontend-builder /app/frontend ./web
# Build the backend binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /quikchat ./cmd/api

# Stage 3: Final production image
FROM alpine:latest
WORKDIR /root/
# Copy the binary from the builder stage
COPY --from=backend-builder /quikchat .
# Copy the static web assets
COPY --from=backend-builder /app/web ./web
# Copy database migrations
COPY quikchat-backend/migrations ./migrations
# Copy the example env file
COPY quikchat-backend/.env.example .

EXPOSE 8080
CMD ["./quikchat"]

