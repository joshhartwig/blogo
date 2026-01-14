# ---- build stage ----
FROM golang:1.24.2 AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -trimpath -ldflags="-s -w" -o /out/blogo ./cmd/server

# ---- runtime stage ----
FROM gcr.io/distroless/static:nonroot
WORKDIR /app

COPY --from=build /out/blogo /app/blogo

# Bake themes into the image (recommended)
COPY themes /app/themes

# Optional: bake default content too (handy for first run)
COPY content /app/content

EXPOSE 3999
USER nonroot:nonroot
ENTRYPOINT ["/app/blogo"]