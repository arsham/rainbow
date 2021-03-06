help: ## Show help messages.
	@grep -E '^[0-9a-zA-Z_-]+:(.*?## .*)?$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

run="."
dir="./..."
short="-short"
flags=""
timeout=40s

TARGET=$(shell git describe --abbrev=0 --tags)
RELEADE_NAME=rainbow
DEPLOY_FOLDER=deploy
CHECKSUM_FILE=CHECKSUM
MAKEFLAGS += -j1

.PHONY: install
install: ## Install the binary.
	@go install -trimpath -ldflags="-s -w"


.PHONY: unittest
unittest: ## Run unit tests in watch mode. You can set: [run, timeout, short, dir, flags]. Example: make unittest flags="-race".
	@echo "running tests on $(run). waiting for changes..."
	@-zsh -c "go test -trimpath --timeout=$(timeout) $(short) $(dir) -run $(run) $(flags); repeat 100 printf '#'; echo"
	@reflex -d none -r "(\.go$$)|(go.mod)" -- zsh -c "go test -trimpath --timeout=$(timeout) $(short) $(dir) -run $(run) $(flags); repeat 100 printf '#'"


.PHONY: ci_tests
ci_tests: ## Run tests for CI.
	go fmt ./...
	go vet ./...
	golangci-lint run ./...
	go test -trimpath --timeout=5m -failfast -v -race -covermode=atomic -coverprofile=coverage.out ./...


.PHONY: dependencies
dependencies: ## Install dependencies requried for development operations.
	@go get -u github.com/cespare/reflex
	@go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.30.0
	@go get -u github.com/git-chglog/git-chglog/cmd/git-chglog
	@go mod tidy


.PHONY: changelog
changelog: ## Update the changelog.
	@git-chglog > CHANGELOG.md
	@echo "Changelog has been updated."


.PHONY: changelog_release
changelog_release: ## Update the changelog with a release tag.
	@git-chglog --next-tag $(tag) > CHANGELOG.md
	@echo "Changelog has been updated."


.PHONY: clean
clean: ## Clean test caches and tidy up modules.
	@go clean -testcache
	@go mod tidy
	@rm -rf $(DEPLOY_FOLDER)


.PHONY: tmpfolder
tmpfolder: ## Create the temporary folder.
	@mkdir -p $(DEPLOY_FOLDER)
	@rm -rf $(DEPLOY_FOLDER)/$(CHECKSUM_FILE) 2> /dev/null

.PHONY: linux
linux: tmpfolder
linux: ## Build for GNU/Linux.
	@GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o $(DEPLOY_FOLDER)/$(RELEADE_NAME) main.go
	@tar -czf $(DEPLOY_FOLDER)/figurine_linux_$(TARGET).tar.gz $(DEPLOY_FOLDER)/$(RELEADE_NAME)
	@cd $(DEPLOY_FOLDER) ; sha256sum figurine_linux_$(TARGET).tar.gz >> $(CHECKSUM_FILE)
	@echo "Linux target:" $(DEPLOY_FOLDER)/figurine_linux_$(TARGET).tar.gz
	@rm $(DEPLOY_FOLDER)/$(RELEADE_NAME)

.PHONY: darwin
darwin: tmpfolder
darwin: ## Build for Mac.
	@GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o $(DEPLOY_FOLDER)/$(RELEADE_NAME) main.go
	@tar -czf $(DEPLOY_FOLDER)/figurine_darwin_$(TARGET).tar.gz $(DEPLOY_FOLDER)/$(RELEADE_NAME)
	@cd $(DEPLOY_FOLDER) ; sha256sum figurine_darwin_$(TARGET).tar.gz >> $(CHECKSUM_FILE)
	@echo "Darwin target:" $(DEPLOY_FOLDER)/figurine_darwin_$(TARGET).tar.gz
	@rm $(DEPLOY_FOLDER)/$(RELEADE_NAME)

.PHONY: windows
windows: tmpfolder
windows: ## Build for windoze.
	@GOOS=windows GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o $(DEPLOY_FOLDER)/$(RELEADE_NAME).exe main.go
	@zip -r $(DEPLOY_FOLDER)/figurine_windows_$(TARGET).zip $(DEPLOY_FOLDER)/$(RELEADE_NAME).exe
	@cd $(DEPLOY_FOLDER) ; sha256sum figurine_windows_$(TARGET).zip >> $(CHECKSUM_FILE)
	@echo "Windows target:" $(DEPLOY_FOLDER)/figurine_windows_$(TARGET).zip
	@rm $(DEPLOY_FOLDER)/$(RELEADE_NAME).exe


.PHONY: release
release: ## Create releases for Linux, Mac, and windoze.
release: tmpfolder linux darwin windows


.PHONY: coverage
coverage: ## Show the test coverage on browser.
	go test -covermode=count -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | tail -n 1
	go tool cover -html=coverage.out
