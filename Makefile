HTMX_VERSION=1.9.6
BOOTSTRAP_VERSION=5.3.2

.PHONY: run
## run: Runs the air command.
run:
	MORPHOS_PORT=3000 air -c .air.toml

.PHONY: download-htmx
## download-htmx: Downloads HTMX minified js file
download-htmx:
	curl -o static/htmx.min.js https://unpkg.com/htmx.org@${HTMX_VERSION}/dist/htmx.min.js

.PHONY: download-bootstrap
## download-bootstrap: Downloads Bootstrap minified css/js file
download-bootstrap:
	curl -o static/bootstrap.min.css https://cdn.jsdelivr.net/npm/bootstrap@${BOOTSTRAP_VERSION}/dist/css/bootstrap.min.css
	curl -o static/bootstrap.min.js https://cdn.jsdelivr.net/npm/bootstrap@${BOOTSTRAP_VERSION}/dist/js/bootstrap.bundle.min.js

.PHONY: build
## build: Builds the container image
build:
	docker build -t morphos .

.PHONY: docker-run
## docker-run: Runs the container
docker-run:
	docker run -d -p 8080:8080 morphos

.PHONY: help
## help: Prints this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
