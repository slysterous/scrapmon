help:
	@echo "Usage: make [option]"
	@echo "OPTIONS:"
	@echo "  lint            use to lint files."
	@echo "  fmt             use to gofmt all files excluding vendor."
	@echo "  fmtcheck        use to check gofmt compatibility of files."
	@echo "  ci              use to run CI pipeline (via docker)."
	@echo "  ci-cleanup      to kill & remove all ci containers."
	@echo "  run       use to run the project locally (via docker)."

lint:
	golint -set_exit_status=1 `go list ./...`

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

ci:
	docker-compose down
	docker-compose build
	docker-compose up -d
	docker-compose ps
	docker-compose run print-scrape ./scripts/ci.sh
	docker-compose down

run:
	docker-compose down
	docker-compose -f docker-compose.yml -f docker-compose.local.yml up -d

ci-cleanup:
	docker-compose down
fmt:
	go fmt ./...
