FROM golang:1.22 as BUILDER

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY *.go ./
COPY ./internal ./internal

RUN CGO_ENABLED=0 go build -o bluesnews-full-rss

RUN mkdir -p articles

FROM debian:bullseye-slim

# Install cron
RUN apt-get update && apt-get install -y cron 

WORKDIR /app

RUN mkdir -p articles

# Copy binary from builder stage
COPY --from=builder /app/bluesnews-full-rss /app/bluesnews-full-rss

COPY crontab /etc/cron.d/bluenews-cron

# Apply the cron job
RUN crontab /etc/cron.d/bluenews-cron

# Create log file for cron job
RUN touch /var/log/cron.log

CMD ["sh", "-c", "cron && /app/bluesnews-full-rss"]
