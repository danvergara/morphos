HTMX_VERSION=1.9.6
BOOTSTRAP_VERSION=5.3.2
GO_VERSION=1.21.5

.PHONY: run
## run: Runs the air command.
run:
	MORPHOS_PORT=8080 air -c .air.toml

.PHONY: download-htmx
## download-htmx: Downloads HTMX minified js file
download-htmx:
	curl -o static/htmx.min.js https://unpkg.com/htmx.org@${HTMX_VERSION}/dist/htmx.min.js

.PHONY: download-bootstrap
## download-bootstrap: Downloads Bootstrap minified css/js file
download-bootstrap:
	curl -o static/bootstrap.min.css https://cdn.jsdelivr.net/npm/bootstrap@${BOOTSTRAP_VERSION}/dist/css/bootstrap.min.css
	curl -o static/bootstrap.min.js https://cdn.jsdelivr.net/npm/bootstrap@${BOOTSTRAP_VERSION}/dist/js/bootstrap.bundle.min.js

.PHONY: docker-build
## docker-build: Builds the container image
docker-build:
	docker build --build-arg="GO_VERSION=${GO_VERSION}" -t morphos .

.PHONY: docker-run
## docker-run: Runs the container
docker-run:
	docker run --rm -p 8080:8080 -v /tmp:/tmp morphos

.PHONY: help
## help: Prints this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
