# syntax=docker/dockerfile:1.7
# =============================================================================
# OPTIMISED Dockerfile — multi-stage.  ~850 MB -> ~17 MB  (98% smaller)
# =============================================================================

# ------------------------- Stage 1: compile ---------------------------------
FROM golang:1.22-alpine AS build
WORKDIR /src

COPY go.mod ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .

# CGO_ENABLED=0  -> a fully static binary with no C library dependency, so it
#                   runs on ANY Linux base, even one with no libc at all.
# -trimpath      -> strips build machine paths (reproducible + no info leak)
# -ldflags="-s -w" -> drops the debug symbol table, ~25% smaller binary
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/order-service .

# ------------------------- Stage 2: runtime ---------------------------------
# The Go toolchain, the source, and the module cache all stay behind in stage 1.
# Only the compiled binary crosses over.
FROM alpine:3.20 AS runtime

RUN apk add --no-cache ca-certificates wget \
 && adduser -D -u 10001 app

COPY --from=build /out/order-service /usr/local/bin/order-service

USER app
ENV PORT=3003
EXPOSE 3003

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -qO- "http://127.0.0.1:$PORT/health" >/dev/null || exit 1

ENTRYPOINT ["/usr/local/bin/order-service"]
