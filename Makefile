export GOBIN=$(GOPATH)/bin
export APP_DIR_LIST=$(shell go list ./... | grep -v /vendor/)
COMMIT_COUNT=`git rev-list --count origin/master`
COMMIT_SHA=`git rev-parse HEAD`
VERSION=0.1.0
all: build

build: bin/all
	@echo "build complete."

bin/all: verify_gopath
	CGO_ENABLED=0 go install -tags netgo $(APP_DIR_LIST)
	go fmt $(APP_DIR_LIST)

verify_gopath:
	@if [ -z "$(GOPATH)" ] || [ "$(GOPATH)" = "" ]; then\
		echo "GOPATH not set. You need to set GOPATH before run this command";\
		exit 1 ;\
	fi

local_bin/all: verify_gopath
	CGO_ENABLED=0 go install -tags local $(APP_DIR_LIST)
	go fmt $(APP_DIR_LIST)

run: local_bin/all
    MINIO_ACCESS_KEY=access_key MINIO_SECRET_KEY=secret_key $(GOPATH)/bin/minio server --address localhost:9001 &
	MINIO_ACCESS_KEY=access_key MINIO_SECRET_KEY=secret_key MINIO_HOST=localhost MINIO_PORT=9001 BLOB_STORE_PORT=8084 BLOB_STORE_HOST=localhost $(GOPATH)/bin/blob-store


pack: build
	mkdir -p build
	cp -Rf $(GOBIN)/blob-store build/blob-store
	cp -Rf $(GOBIN)/minio build/minio
	echo "commit_sha=$(COMMIT_SHA)" > build/build_info.ini
	zip -r -q blob-store-${VERSION}.zip build/blob-store build/minio build/build_info.ini
