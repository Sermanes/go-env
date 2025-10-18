DOCKER_REGISTRY = registry.digitalocean.com/computernerd/
PROJECT = PROJECT_NAME
DOCKER_TAG = TAG

.PHONY: install-deps tests format linting ci build docker-buildx docker-push swagger

install-deps:
	go mod init || true
	go mod tidy

tests:
	go clean -testcache && go test -v -failfast ./... -short

format:
	go fmt ./...
