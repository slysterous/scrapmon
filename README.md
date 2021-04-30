# scrapmon
A highly concurrent prntscr scrapper.

requirements

docker
docker compose
make 
go version to be defined
go mod to be defined

## Usage
go run --race cmd/scrapmon/main.go start --from=lHB0T --iterations=5000 --workers=16