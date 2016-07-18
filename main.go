/**
 * Copyright (c) 2016 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"
	"github.com/gocraft/web"
	"github.com/trustedanalytics/tapng-blob-store/api"
	"github.com/trustedanalytics/tapng-blob-store/minio-wrapper"
	httpGoCommon "github.com/trustedanalytics/tapng-go-common/http"
	"github.com/trustedanalytics/tapng-go-common/logger"
	"net/http"
	"os"
)

var (
	logger = logger_wrapper.InitLogger("main")
)

const (
	bucketName = "blobstore"
)

func Healthz(rw web.ResponseWriter, req *web.Request) {
	rw.WriteHeader(http.StatusOK)
	fmt.Fprintf(rw, "ok\n")
}

func main() {
	wrappedMinio, err := miniowrapper.CreateWrappedMinio(bucketName)
	if err != nil {
		logger.Fatal(err)
	}

	context := api.NewApiContext(wrappedMinio)
	router := web.New(*context)
	router.Get("/healthz", Healthz)
	apiRouter := router.Subrouter(*context, "/api")

	v1Router := apiRouter.Subrouter(*context, "/v1")
	api.RegisterRoutes(v1Router, *context)

	v1AliasRouter := apiRouter.Subrouter(*context, "/v1.0")
	api.RegisterRoutes(v1AliasRouter, *context)

	if os.Getenv("BLOB_STORE_SSL_CERT_FILE_LOCATION") != "" {
		httpGoCommon.StartServerTLS(os.Getenv("BLOB_STORE_SSL_CERT_FILE_LOCATION"),
			os.Getenv("BLOB_STORE_SSL_KEY_FILE_LOCATION"), router)
	} else {
		httpGoCommon.StartServer(router)
	}
}
