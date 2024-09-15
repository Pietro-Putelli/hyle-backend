export
LOCAL_BIN := $(CURDIR)/bin
OUT_DIR := $(CURDIR)/out
PATH:=$(LOCAL_BIN):$(PATH)
include .env

# https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md

$(eval export $(shell sed -ne 's/ *#.*$$//; /./ s/=.*$$// p' .env))

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

functions := $(shell find functions -mindepth 1 -type d -exec basename {} \;)

build-AuthTokenPostFun: ## Build AuthTokenPostFun
	@GOOS=linux GOARCH=amd64 go build -o functions/AuthTokenPostFun/bootstrap functions/AuthTokenPostFun/main.go
	cp functions/AuthTokenPostFun/bootstrap $(ARTIFACTS_DIR)/.

build-AuthRefreshPostFun: ## Build AuthRefreshPostFun
	@GOOS=linux GOARCH=amd64 go build -o functions/AuthRefreshPostFun/bootstrap functions/AuthRefreshPostFun/main.go
	cp functions/AuthRefreshPostFun/bootstrap $(ARTIFACTS_DIR)/.

build-CustomAuthorizerFun: ## Build CustomAuthorizerFun
	@GOOS=linux GOARCH=amd64 go build -o functions/CustomAuthorizerFun/bootstrap functions/CustomAuthorizerFun/main.go
	cp functions/CustomAuthorizerFun/bootstrap $(ARTIFACTS_DIR)/.

build-UserSessionPatchFun: ## Build UserSessionPatchFun
	@GOOS=linux GOARCH=amd64 go build -o functions/UserSessionPatchFun/bootstrap functions/UserSessionPatchFun/main.go
	cp functions/UserSessionPatchFun/bootstrap $(ARTIFACTS_DIR)/.

build-UserLogoutDeleteFun: ## Build UserLogoutDeleteFun
	@GOOS=linux GOARCH=amd64 go build -o functions/UserLogoutDeleteFun/bootstrap functions/UserLogoutDeleteFun/main.go
	cp functions/UserLogoutDeleteFun/bootstrap $(ARTIFACTS_DIR)/.

build-UserProfilePutFun: ## Build UserProfilePutFun
	@GOOS=linux GOARCH=amd64 go build -o functions/UserProfilePutFun/bootstrap functions/UserProfilePutFun/main.go
	cp functions/UserProfilePutFun/bootstrap $(ARTIFACTS_DIR)/.

build-UserProfileHealthFun: ## Build UserProfileHealthFun
	@GOOS=linux GOARCH=amd64 go build -o functions/UserProfileHealthFun/bootstrap functions/UserProfileHealthFun/main.go
	cp functions/UserProfileHealthFun/bootstrap $(ARTIFACTS_DIR)/.

build-BookPickPostFun: ## Build BookPickPostFun
	@GOOS=linux GOARCH=amd64 go build -o functions/BookPickPostFun/bootstrap functions/BookPickPostFun/main.go
	cp functions/BookPickPostFun/bootstrap $(ARTIFACTS_DIR)/.

build-BookDeleteFun: ## Build BookDeleteFun
	@GOOS=linux GOARCH=amd64 go build -o functions/BookDeleteFun/bootstrap functions/BookDeleteFun/main.go
	cp functions/BookDeleteFun/bootstrap $(ARTIFACTS_DIR)/.

build-BookPickDeleteFun: ## Build BookPickDeleteFun
	@GOOS=linux GOARCH=amd64 go build -o functions/BookPickDeleteFun/bootstrap functions/BookPickDeleteFun/main.go
	cp functions/BookPickDeleteFun/bootstrap $(ARTIFACTS_DIR)/.

build-BookPutFun: ## Build BookPutFun
	@GOOS=linux GOARCH=amd64 go build -o functions/BookPutFun/bootstrap functions/BookPutFun/main.go
	cp functions/BookPutFun/bootstrap $(ARTIFACTS_DIR)/.

build-BookPicksGetFun: ## Build BookPicksGetFun
	@GOOS=linux GOARCH=amd64 go build -o functions/BookPicksGetFun/bootstrap functions/BookPicksGetFun/main.go
	cp functions/BookPicksGetFun/bootstrap $(ARTIFACTS_DIR)/.

build-BookTopicsGetFun: ## Build BookTopicsGetFun
	@GOOS=linux GOARCH=amd64 go build -o functions/BookTopicsGetFun/bootstrap functions/BookTopicsGetFun/main.go
	cp functions/BookTopicsGetFun/bootstrap $(ARTIFACTS_DIR)/.

build-BookPickPutFun: ## Build BookPickPutFun
	@GOOS=linux GOARCH=amd64 go build -o functions/BookPickPutFun/bootstrap functions/BookPickPutFun/main.go
	cp functions/BookPickPutFun/bootstrap $(ARTIFACTS_DIR)/.

build-BookGetFun: ## Build BookGetFun
	@GOOS=linux GOARCH=amd64 go build -o functions/BookGetFun/bootstrap functions/BookGetFun/main.go
	cp functions/BookGetFun/bootstrap $(ARTIFACTS_DIR)/.

build-BookListFun: ## Build BookListFun
	@GOOS=linux GOARCH=amd64 go build -o functions/BookListFun/bootstrap functions/BookListFun/main.go
	cp functions/BookListFun/bootstrap $(ARTIFACTS_DIR)/.

build-CreatePickKeywordsFun: ## Build CreatePickKeywordsFun
	@GOOS=linux GOARCH=amd64 go build -o functions/CreatePickKeywordsFun/bootstrap functions/CreatePickKeywordsFun/main.go
	cp functions/CreatePickKeywordsFun/bootstrap $(ARTIFACTS_DIR)/.

build-SharpPickFun: ## Build SharpPickFun
	@GOOS=linux GOARCH=amd64 go build -o functions/SharpPickFun/bootstrap functions/SharpPickFun/main.go
	cp functions/SharpPickFun/bootstrap $(ARTIFACTS_DIR)/.

build-KeywordDetailFun: ## Build KeywordDetailFun
	@GOOS=linux GOARCH=amd64 go build -o functions/KeywordDetailFun/bootstrap functions/KeywordDetailFun/main.go
	cp functions/KeywordDetailFun/bootstrap $(ARTIFACTS_DIR)/.

build-SemanticSearchFun: ## Build SemanticSearchFun
	@GOOS=linux GOARCH=amd64 go build -o functions/SemanticSearchFun/bootstrap functions/SemanticSearchFun/main.go
	cp functions/SemanticSearchFun/bootstrap $(ARTIFACTS_DIR)/.

build-BookSavePostFun: ## Build BookSavePostFun
	@GOOS=linux GOARCH=amd64 go build -o functions/BookSavePostFun/bootstrap functions/BookSavePostFun/main.go
	cp functions/BookSavePostFun/bootstrap $(ARTIFACTS_DIR)/.

build-TranslateWordFun: ## Build TranslateWordFun
	@GOOS=linux GOARCH=amd64 go build -o functions/TranslateWordFun/bootstrap functions/TranslateWordFun/main.go
	cp functions/TranslateWordFun/bootstrap $(ARTIFACTS_DIR)/.

build-EligibleUsersForPNFun: ## Build EligibleUsersForPNFun
	@GOOS=linux GOARCH=amd64 go build -o functions/EligibleUsersForPNFun/bootstrap functions/EligibleUsersForPNFun/main.go
	cp functions/EligibleUsersForPNFun/bootstrap $(ARTIFACTS_DIR)/.

build-SendPushNotificationFun: ## Build SendPushNotificationFun
	@GOOS=linux GOARCH=amd64 go build -o functions/SendPushNotificationFun/bootstrap functions/SendPushNotificationFun/main.go
	cp functions/SendPushNotificationFun/bootstrap $(ARTIFACTS_DIR)/.

build: ## Build all functions
	sam build
.PHONY: build

test: setup ## Run tests
	@mkdir -p coverage
	@$(LOCAL_BIN)/ginkgo run -r --randomize-suites --fail-on-pending --cover --coverprofile=coverage/coverage.out --json-report=coverage/coverage-report.json .
.PHONY: setup

setup: ### Install bin dependencies
	GOBIN=$(LOCAL_BIN) go install github.com/onsi/ginkgo/v2/ginkgo
	GOBIN=$(LOCAL_BIN) go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	GOBIN=$(LOCAL_BIN) go install go.uber.org/mock/mockgen@latest

migrate: setup ## Run migrations, take care of .env file
	@echo "Running migrations..."
	@$(LOCAL_BIN)/migrate -verbose -path migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=require" up
	@echo "Migrations done"
.PHONY: migrate

generate: setup ### Generate mocks
	GOBIN=$(LOCAL_BIN) go generate ./...
.PHONY: generate

deploy: ## Deploy to AWS with SAM
	@sam deploy --config-file samconfig.toml --template-file .aws-sam/build/template.yaml --stack-name feynman --capabilities CAPABILITY_IAM --resolve-s3
.PHONY: deploy

update: ### Update dependencies
	go mod tidy
.PHONY: update
