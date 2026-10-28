# Stage 1: Build Frontend Assets
FROM node:24-alpine AS frontend-builder
WORKDIR /app
RUN npm install -g pnpm@12.14.0
COPY pnpm-workspace.yaml package.json pnpm-lock.yaml ./
COPY web/package.json ./web/
RUN pnpm install --filter bn-riftstats-web
COPY web/ ./web/
RUN pnpm approve-builds --all && cd web && pnpm exec vite build

# Stage 2: Build Go Backend Binaries
FROM golang:1.24-alpine AS backend-builder
WORKDIR /app
RUN apk add --no-cache git ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY migrations/ ./migrations/
COPY --from=frontend-builder /app/web/dist ./web/dist
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/ingestor ./cmd/ingestor
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/aggregator ./cmd/aggregator

# Stage 3: Minimal Runtime Container
FROM alpine:3.20 AS runner
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=backend-builder /bin/api /bin/ingestor /bin/aggregator /app/
COPY --from=frontend-builder /app/web/dist /app/web/dist
COPY migrations/ /app/migrations/

EXPOSE 8080
ENV PORT=8080
ENV ENV=production
ENV GOVERNOR_MODE=DEV_PROTOTYPE_MODE

ENTRYPOINT ["/app/api"]
