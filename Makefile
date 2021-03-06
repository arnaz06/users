SOURCES := $(shell find . -name '*.go' -type f -not -path './vendor/*'  -not -path '*/mocks/*')
TEST_OPTS := -covermode=atomic $(TEST_OPTS)
GIT_COMMIT := $(shell git rev-parse --short HEAD)

IMAGE_NAME = users

# Database
MYSQL_USER ?= users
MYSQL_PASSWORD ?= users-pass
MYSQL_ADDRESS ?= 127.0.0.1:3306
MYSQL_DATABASE ?= users


# Dependency Management
.PHONY: vendor
vendor: go.mod go.sum
	@GO111MODULE=on go get ./...

# Linter
.PHONY: lint-prepare
lint-prepare:
	@echo "Installing golangci-lint"
	@wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.31.0

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: mockery-prepare
mockery-prepare:
	@echo "Installing mockery"
	@GO111MODULE=off go get -u github.com/vektra/mockery/.../

# Testing
.PHONY: unittest
unittest: vendor
	GO111MODULE=on go test -short $(TEST_OPTS) ./...

.PHONY: test
test: vendor
	GO111MODULE=on go test $(TEST_OPTS) ./...


# Database Migration
.PHONY: migrate-prepare
migrate-prepare:
	@GO111MODULE=off go get -tags 'mysql' -u github.com/golang-migrate/migrate/cmd/migrate

.PHONY: migrate-up
migrate-up:
	@migrate -database "mysql://$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_ADDRESS))/$(MYSQL_DATABASE)" \
	-path=internal/mysql/migrations up

.PHONY: migrate-down
migrate-down:
	@migrate -database "mysql://$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_ADDRESS))/$(MYSQL_DATABASE)" \
	-path=internal/mysql/migrations down

#Build
.PHONY: users
users:
	GO111MODULE=on go build -o users -ldflags="-X 'github.com/arnaz06/users.Version=${GIT_COMMIT}'" github.com/arnaz06/users/cmd/api

.PHONY: docker
docker: vendor $(SOURCES)
	@docker build -t $(IMAGE_NAME) .


# Docker
.PHONY: mysql-up
mysql-up:
	@docker-compose up -d mysql

.PHONY: mysql-down
mysql-down:
	@docker stop users_mysql

.PHONY: run
run:
	@docker-compose up -d

.PHONY: stop
stop:
	@docker-compose down


# Mock
UserRepository: user.go
	@mockery -name=UserRepository

UserService: user.go
	@mockery -name=UserService
