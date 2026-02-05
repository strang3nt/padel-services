FROM golang:1.24-trixie AS builder


ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=Etc/UTC
COPY . .

RUN apt-get update && apt-get install -y --no-install-recommends \
    make npm \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* \
    && make

# --------------------------------------------------------------#

FROM debian:trixie-slim AS runner

ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=Etc/UTC

RUN apt-get update && apt-get install -y --no-install-recommends \
    chromium \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /go/tgminiapp /go/tgminiapp
ENV CHROME_EXECUTABLE="chromium"


ENTRYPOINT ["./go/tgminiapp"]
