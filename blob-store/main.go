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
	"net/http"
	"github.com/gocraft/web"
	"github.com/trustedanalytics/tap-go-common/logger"
	"github.com/trustedanalytics/blob-store/minioWrapper"
)

type Context struct{
	wrappedMinio *minioWrapper.Wrapper
}

var (
	logger = logger_wrapper.InitLogger("main")
)

const (
	port = "8080"
	bucketName = "blobstore"
)

const (
	URLblobs = "/api/v1/blobs/"
)


func main() {
	wrappedMinio, err := minioWrapper.CreateWrappedMinio(bucketName)
	if err != nil {
		logger.Fatal(err)
	}

	context := Context{wrappedMinio}
	router := web.New(context)
	router.Post(URLblobs, context.StoreBlob)
	router.Get(URLblobs + ":blob_id", context.RetrieveBlob)
	router.Delete(URLblobs + ":blob_id", context.RemoveBlob)

	err = http.ListenAndServe("localhost:" + port, router)
	if err != nil {
		logger.Critical("Couldn't serve blob store on port ", port, " Application will be closed now.")
		logger.Fatal(err)
	}
}