FROM golang:1.13 as builder
WORKDIR /go/src/github.com/slysterous/print-scrape

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -a -installsuffix cgo -o print-scrape ./cmd/print-scrape/main.go

FROM alpine

RUN apk add --no-cache tzdata

COPY --from=builder /go/src/github.com/slysterous/print-scrape .

ENV INPUT ""
ENV OUTPUT ""

CMD ./print-scrape $INPUT $OUTPUT