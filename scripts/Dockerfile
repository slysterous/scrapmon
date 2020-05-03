# This is for testing and linting purposes
# It is to be used by CI tools.

FROM golang:1.13
WORKDIR /go/src/github.com/slysterous/print-scrape
RUN go get golang.org/x/lint/golint
ARG version=dev
COPY . ./