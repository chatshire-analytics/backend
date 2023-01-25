PROJECT_NAME=mentat-backend
MODULE_NAME=mentat-backend

.DEFAULT_GOAL := build

.PHONY: build
build:
	@go build ./cmd/mentat-backend/

.PHONY: fmt
fmt:
	@go fmt ./...

.PHONY: test
test:
	@go test -v -coverprofile coverage.out ./...

.PHONY: coverage
coverage:
	@go tool cover -html=coverage.out

.PHONY: get
get:
	@go mod download

.PHONY: docker
docker:
	@docker build -f ./build/package/Dockerfile -t mentat-backend:latest .

.PHONY: deploy
deploy:
	@kubectl apply -f deployments/app-deployment.yaml

.PHONY: undeploy
undeploy:
	@kubectl delete -f deployments/app-deployment.yaml