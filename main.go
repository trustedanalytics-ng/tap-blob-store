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
	"net/http"

	"github.com/gocraft/web"

	"github.com/trustedanalytics/tap-blob-store/api"
	"github.com/trustedanalytics/tap-blob-store/minio-wrapper"
	httpGoCommon "github.com/trustedanalytics/tap-go-common/http"
	commonLogger "github.com/trustedanalytics/tap-go-common/logger"
)

var (
	logger, _ = commonLogger.InitLogger("main")
)

const (
	bucketName = "blobstore"
)

func Healthz(rw web.ResponseWriter, req *web.Request) {
	rw.WriteHeader(http.StatusOK)
	fmt.Fprint(rw, "ok\n")
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

	httpGoCommon.StartServer(router)
}
