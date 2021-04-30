# scrapmon [![PkgGoDev](https://pkg.go.dev/badge/github.com/slysterous/scrapmon)]((https://pkg.go.dev/github.com/slysterous/scrapmon))
[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](https://lbesson.mit-license.org/)
[![Coverage Status](https://coveralls.io/repos/github/slysterous/scrapmon/badge.svg?branch=main)](https://coveralls.io/github/slysterous/scrapmon?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/slysterous/scrapmon)](https://goreportcard.com/report/github.com/slysterous/scrapmon)

A highly concurrent prntscr scrapper.

requirements

docker
docker compose
make 
go version to be defined
go mod to be defined

## Usage
go run --race cmd/scrapmon/main.go start --from=lHB0T --iterations=5000 --workers=16
