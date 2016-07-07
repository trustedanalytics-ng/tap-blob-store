GOBIN=$(GOPATH)/bin
APP_DIR_LIST=$(shell go list ./... | grep -v /vendor/)
COMMIT_COUNT=`git rev-list --count origin/master`
COMMIT_SHA=`git rev-parse HEAD`
VERSION=0.1.0
all: build

build: bin/blob-store bin/minio
	@echo "build complete."

bin/minio: verify_gopath
	@if [ ! -f "$(GOBIN)/minio" ]; then\
		echo "Minio server was not found. It will be downloaded";\
		wget https://dl.minio.io/server/minio/release/linux-amd64/minio -O $(GOBIN)/minio ;\
		chmod +x $(GOBIN)/minio ;\
	fi

bin/blob-store: verify_gopath
	CGO_ENABLED=0 go install -tags netgo $(APP_DIR_LIST)
	go fmt $(APP_DIR_LIST)

verify_gopath:
	@if [ -z "$(GOPATH)" ] || [ "$(GOPATH)" = "" ]; then\
		echo "GOPATH not set. You need to set GOPATH before run this command";\
		exit 1 ;\
	fi

local_bin/minio: bin/minio
	mkdir -p ~/MINIO

local_bin/blob-store: verify_gopath
	CGO_ENABLED=0 go install -tags local $(APP_DIR_LIST)
	go fmt $(APP_DIR_LIST)

run: local_bin/blob-store local_bin/minio
	MINIO_ACCESS_KEY=access_key MINIO_SECRET_KEY=secret_key $(GOBIN)/minio server ~/MINIO --address localhost:9001 &\
	sleep 2 &&\

pack: build
	mkdir -p build
	cp -Rf $(GOBIN)/tapng-blob-store build/tapng-blob-store
	cp -Rf $(GOBIN)/minio build/minio
	echo "commit_sha=$(COMMIT_SHA)" > build/build_info.ini
	zip -r -q tapng-blob-store-${VERSION}.zip build/tapng-blob-store build/minio build/build_info.ini

prepare_dirs:
	mkdir -p ./temp/src/github.com/trustedanalytics/tapng-blob-store
	$(eval REPOFILES=$(shell pwd)/*)
	ln -sf $(REPOFILES) temp/src/github.com/trustedanalytics/tapng-blob-store

build_anywhere: prepare_dirs
	$(eval GOPATH=$(shell cd ./temp; pwd))
	$(eval APP_DIR_LIST=$(shell GOPATH=$(GOPATH) go list ./temp/src/github.com/trustedanalytics/tapng-blob-store/... | grep -v /vendor/))
	GOPATH=$(GOPATH) CGO_ENABLED=0 go install -tags netgo $(APP_DIR_LIST)
	rm -Rf application && mkdir application
	cp $(GOPATH)/bin/tapng-blob-store ./application/tapng-blob-store
