GOBIN=$(GOPATH)/bin
APP_DIR_LIST=$(shell go list ./... | grep -v /vendor/)
PROJECT_VERSION=`grep VERSION: manifest.yml | cut -d ":" -f 2 | tr -d "\" "`
COMMIT_COUNT=`git rev-list --count origin/master`
COMMIT_SHA=`git rev-parse HEAD`
VERSION=$(PROJECT_VERSION)

build: verify_gopath
	CGO_ENABLED=0 go install -tags netgo $(APP_DIR_LIST)
	go fmt $(APP_DIR_LIST)

deps_fetch_newest: bin/govendor
	$(GOBIN)/govendor remove +all
	@echo "Update deps used in project to their newest versions"
	$(GOBIN)/govendor fetch -v +external, +missing

deps_update: verify_gopath
	$(MAKE) bin/govendor
	@echo "Update all vendor deps according to their current version in GOPATH"
	$(GOBIN)/govendor remove +all
	$(GOBIN)/govendor update +external
	@echo "Done"

deps_list: bin/govendor
	@echo "Project dependencies list:"
	$(GOBIN)/govendor list

verify_gopath:
	@if [ -z "$(GOPATH)" ] || [ "$(GOPATH)" = "" ]; then\
		echo "GOPATH not set. You need to set GOPATH before run this command";\
		exit 1 ;\
	fi

mock_update: bin/gomock
	$(GOBIN)/mockgen -source=state/state.go -package=state -destination=state/state_mock.go

tests: verify_gopath mock_update
	go test --cover $(APP_DIR_LIST)

bin/gomock: verify_gopath
	go get -v -u github.com/golang/mock/mockgen

bin/govendor: verify_gopath
	go get -v -u github.com/kardianos/govendor
