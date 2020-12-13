mod:
	go mod download

gen:
	go generate -v ./...

test:
	go test -v --race ./... -cover -coverprofile=coverage.out

coverage:
	go tool cover -html=coverage.out

lint:
	golangci-lint run ./...

up:
	docker-compose up -d

down:
	docker-compose down

build:
	$(eval GIT_BRANCH=$(shell git rev-parse --short HEAD))
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-X main.version=$(GIT_BRANCH)" -ldflags="-X main.configFile=c.json" -a -installsuffix cgo -o appbin .
	docker build -t testassignment . --no-cache

run: build up