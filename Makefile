GOBIN=$(GOPATH)/bin
APP_DIR_LIST=$(shell go list ./... | grep -v /vendor/)
MINIO_IN_LAB_URL=http://rrceph01.sclab.intel.com/dependencies/minio
MINIO_EXT_URL=https://dl.minio.io/server/minio/release/linux-amd64/minio
COMMIT_COUNT=`git rev-list --count origin/master`
COMMIT_SHA=`git rev-parse HEAD`
VERSION=0.1.0
all: build

build: bin/blob-store
	@echo "build complete."

bin/minio: verify_gopath
	@if [ ! -f "$(GOBIN)/minio" ]; then\
		echo "Minio server was not found. It will be downloaded";\
		wget "$(MINIO_IN_LAB_URL)" -O $(GOBIN)/minio || wget "$(MINIO_EXT_URL)" -O $(GOBIN)/minio ;\
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

deps_fetch_specific: bin/govendor
	@if [ "$(DEP_URL)" = "" ]; then\
		echo "DEP_URL not set. Run this comand as follow:";\
		echo " make deps_fetch_specific DEP_URL=github.com/nu7hatch/gouuid";\
	exit 1 ;\
	fi
	@echo "Fetching specific dependency in newest versions"
	$(GOBIN)/govendor fetch -v $(DEP_URL)

deps_update_tapng: verify_gopath
	$(GOBIN)/govendor update github.com/trustedanalytics/...
	rm -Rf vendor/github.com/trustedanalytics/tapng-blob-store
	@echo "Done"

local_bin/minio: bin/minio
	mkdir -p ~/MINIO

local_bin/blob-store: verify_gopath
	CGO_ENABLED=0 go install -tags local $(APP_DIR_LIST)
	go fmt $(APP_DIR_LIST)

run: local_bin/blob-store local_bin/minio
	MINIO_ACCESS_KEY=access_key MINIO_SECRET_KEY=secret_key $(GOBIN)/minio server ~/MINIO --address localhost:9000 &\
	sleep 2 &&\
	MINIO_ACCESS_KEY=access_key MINIO_SECRET_KEY=secret_key MINIO_HOST=localhost MINIO_PORT=9000 PORT=8084 BIND_ADDRESS=localhost $(GOBIN)/tapng-blob-store

build_anywhere: prepare_dirs
	$(eval GOPATH=$(shell cd ./temp; pwd))
	$(eval GOBIN=$(GOPATH)/bin)
	$(eval APP_DIR_LIST=$(shell GOPATH=$(GOPATH) go list ./temp/src/github.com/trustedanalytics/tapng-blob-store/... | grep -v /vendor/))
	GOPATH=$(GOPATH) CGO_ENABLED=0 go install -tags netgo -pkgdir /temp/pkg $(APP_DIR_LIST)
	$(MAKE) bin/minio
	mkdir -p build
	cp -Rf $(GOBIN)/tapng-blob-store build/tapng-blob-store
	cp -Rf $(GOBIN)/minio build/minio
	cp -Rf build/ minio/
	echo "commit_sha=$(COMMIT_SHA)" > build/build_info.ini
	zip -r -q tapng-blob-store-${VERSION}.zip build/tapng-blob-store build/minio build/build_info.ini
	rm -Rf ./temp

prepare_dirs:
	mkdir -p ./temp/src/github.com/trustedanalytics/tapng-blob-store
	$(eval REPOFILES=$(shell pwd)/*)
	ln -sf $(REPOFILES) temp/src/github.com/trustedanalytics/tapng-blob-store

docker_build: build_anywhere
	docker build -t tapng-blob-store .
	docker build -t tapng-blob-store/minio ./minio/

kubernetes_deploy: docker_build
	kubectl create -f service.yaml
	kubectl create -f minio-configmap.yaml
	kubectl create -f minio-secret.yaml
	kubectl create -f deployment.yaml

kubernetes_update: docker_build
	kubectl delete -f deployment.yaml
	kubectl create -f deployment.yaml