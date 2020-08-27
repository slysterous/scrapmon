FROM golang:1.13 as builder
WORKDIR /home/print-scrape
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -a -installsuffix cgo -o print-scrape ./cmd/print-scrape/main.go && wget https://github.com/golang-migrate/migrate/releases/download/v4.1.0/migrate.linux-amd64.tar.gz && tar -xvf migrate.linux-amd64.tar.gz && mv migrate.linux-amd64 migrate

FROM alpine

RUN apk add --no-cache tzdata

COPY --from=builder /home/print-scrape/print-scrape .
COPY --from=builder /home/print-scrape/migrate .
COPY --from=builder /home/print-scrape/internal/migrations ./internal/migrations
COPY --from=builder /home/print-scrape/migrate.sh .

CMD [ "sh", "-c",  "/migrate.sh && /print-scrape" ]