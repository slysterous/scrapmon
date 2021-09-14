# This is for testing and linting purposes
# It is to be used by CI tools.

FROM golang:1.16
WORKDIR /scrapmon
RUN go get golang.org/x/lint/golint
RUN wget https://github.com/golang-migrate/migrate/releases/download/v4.1.0/migrate.linux-amd64.tar.gz && tar -xvf migrate.linux-amd64.tar.gz && mv migrate.linux-amd64 migrate
ARG version=dev
COPY . ./