FROM golang:1.25-trixie AS builder


ENV DEBIAN_FRONTEND=noninteractive
ARG ENVIRONMENT=prod
ENV TZ=Etc/UTC
COPY . .

RUN apt-get update && apt-get install -y --no-install-recommends \
    make npm \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* \
    && npm install -g pnpm@latest-10 \
    && cd client && pnpm install && cd .. \
    && make ENVIRONMENT="$ENVIRONMENT"

# --------------------------------------------------------------#

FROM debian:trixie-slim AS runner

ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=Etc/UTC

RUN apt-get update && apt-get install -y --no-install-recommends \
    chromium ca-certificates \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /go/tgminiapp /go/tgminiapp
ENV CHROME_EXECUTABLE="chromium"


ENTRYPOINT ["./go/tgminiapp"]
