# Builds the frontend and the backend, then ships both in one small image:
# the Go server serves the API and the static frontend from the same origin.

# Node major matches frontend/.nvmrc.
FROM node:24-alpine AS frontend
WORKDIR /app
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# The Go tag must be >= the go directive in backend/go.mod: official images set
# GOTOOLCHAIN=local, so they won't download a newer toolchain.
FROM golang:1.27-alpine AS backend
WORKDIR /src
COPY backend/ ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /server ./cmd/server

# No shell or curl in this image, so no HEALTHCHECK; CI smoke-tests /healthz.
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=backend /server /server
COPY --from=frontend /app/dist /app/web
ENV STATIC_DIR=/app/web PORT=8080
EXPOSE 8080
ENTRYPOINT ["/server"]
