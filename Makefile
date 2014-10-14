
APP_FILE=gobackuper
APP_PATH=./$(APP_FILE)

NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m
DEPS=$(go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)

VERSION:=`cat VERSION`

deps:
	@echo "$(OK_COLOR)==> Installing dependencies$(NO_COLOR)"
	@go get -d -v ./...
	@echo $(DEPS) | xargs -n1 go get -d

updatedeps:
	@echo "$(OK_COLOR)==> Updating all dependencies$(NO_COLOR)"
	@go get -d -v -u ./...
	@echo $(DEPS) | xargs -n1 go get -d -u

format:
	@echo "$(OK_COLOR)==> Formatting$(NO_COLOR)"
	go fmt ./...

lint:
	@echo "$(OK_COLOR)==> Linting$(NO_COLOR)"
	golint .

build:
	@echo "$(OK_COLOR)==> Building$(NO_COLOR)"
	@echo "Current Version: $(VERSION)"
	go build -ldflags "-X main.Version $(VERSION)" -o $(APP_PATH)
	#strip $(APP_PATH)

clear:
	@echo "$(OK_COLOR)==> Clearing$(NO_COLOR)"
	rm -f $(APP_PATH)

brew: build

