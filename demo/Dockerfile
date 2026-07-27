# Reproducible demo environment for the README recordings.
#
# The container's network namespace only contains the synthetic dev stack
# started by entrypoint.sh, so recordings never expose the host's real
# processes. Build from the repository root:
#
#   docker build -f demo/Dockerfile -t ports-demo .

FROM golang:1.25-bookworm AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 go build -ldflags "-s -w -X main.version=demo" -o /out/ports ./cmd/ports

COPY demo/go.mod ./demo/go.mod
COPY demo/listener ./demo/listener
RUN cd demo && CGO_ENABLED=0 go build -ldflags "-s -w" -o /out/listener ./listener

FROM debian:bookworm-slim
RUN apt-get update \
 && apt-get install -y --no-install-recommends iproute2 procps \
 && rm -rf /var/lib/apt/lists/*

COPY --from=build /out/ports /usr/local/bin/ports
COPY --from=build /out/listener /opt/demo/listener
COPY demo/entrypoint.sh /usr/local/bin/demo-entrypoint

RUN chmod +x /usr/local/bin/demo-entrypoint \
 && printf '%s\n' \
      'PS1="\[\033[38;5;99m\]❯\[\033[0m\] "' \
      'export TERM=xterm-256color' \
    > /root/.bashrc

ENV TERM=xterm-256color
WORKDIR /root
ENTRYPOINT ["/usr/local/bin/demo-entrypoint"]
