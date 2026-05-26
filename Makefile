.PHONY: benchmark
benchmark:
	go test -bench=.

.PHONY: build-docker
build-docker:
	docker build --platform linux/amd64 . -t quineglobal/hyper-quineglobal-com

.PHONY: push-docker
push-docker: build-docker
	docker push quineglobal/hyper-quineglobal-com

.PHONY: cover
cover:
	go tool cover -html=cover.out

.PHONY: css
css:
	npm run css

.PHONY: css-watch
css-watch:
	npm run css:watch

.PHONY: dev
dev:
	@trap 'kill 0' SIGINT SIGTERM; npm run css:watch & go run ./cmd/app; wait

.PHONY: install
install:
	npm install

.PHONY: lint
lint:
	golangci-lint run

.PHONY: start
start:
	go run ./cmd/app

.PHONY: test
test:
	go test -coverprofile=cover.out -shuffle on ./...
