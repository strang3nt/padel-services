FROM golang:1.25

RUN apt-get update && apt-get install -y \
    chromium \
    wget \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY . .
RUN mkdir cmd/tgminiapp/dist && touch cmd/tgminiapp/dist/placeholder

CMD ["./start.development.sh"]
