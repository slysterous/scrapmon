FROM golang:1.15 as builder
WORKDIR /home/scrapmon
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -a -installsuffix cgo -o scrapmon ./cmd/scrapmon/main.go && wget https://github.com/golang-migrate/migrate/releases/download/v4.1.0/migrate.linux-amd64.tar.gz && tar -xvf migrate.linux-amd64.tar.gz && mv migrate.linux-amd64 migrate

FROM alpine

RUN apk add --no-cache tzdata

COPY --from=builder /home/scrapmon/scrapmon .
COPY --from=builder /home/scrapmon/scripts/migrate .
COPY --from=builder /home/scrapmon/internal/migrations ./internal/migrations
COPY --from=builder /home/scrapmon/scripts/migrate.sh .

CMD [ "sh", "-c",  "/scripts/migrate.sh && /scrapmon" ]