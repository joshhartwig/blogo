# ---- build stage ----
FROM golang:1.26.2 AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
  go build -trimpath -ldflags="-s -w" -o /out/blogo ./cmd/server

# ---- runtime stage ----
FROM gcr.io/distroless/static:nonroot
WORKDIR /app

COPY --from=build /out/blogo /app/blogo

# No need to copy ui/static or ui/templates, as they are embedded in the binary

# Optional: bake default content too (handy for first run)
COPY content /app/content

EXPOSE 4000
USER nonroot:nonroot
ENTRYPOINT ["/app/blogo"]
CMD ["serve", "--config", "/config/config.toml", "--addr", ":4000"]