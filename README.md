# scrapmon [![PkgGoDev](https://pkg.go.dev/badge/github.com/slysterous/scrapmon)](https://pkg.go.dev/github.com/slysterous/scrapmon)
[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](https://lbesson.mit-license.org/)
![example workflow](https://github.com/slysterous/scrapmon/actions/workflows/tests.yml/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/slysterous/scrapmon/badge.svg?branch=master)](https://coveralls.io/github/slysterous/scrapmon)
[![Go Report Card](https://goreportcard.com/badge/github.com/slysterous/scrapmon)](https://goreportcard.com/report/github.com/slysterous/scrapmon)

A highly concurrent prntscr scrapper.

## Requirements

docker
docker compose
make 
go version to be defined
go mod to be defined

## Usage
go run --race cmd/scrapmon/main.go start --from=lHB0T --iterations=5000 --workers=16
