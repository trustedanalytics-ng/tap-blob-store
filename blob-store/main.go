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
	"github.com/gocraft/web"
	"github.com/trustedanalytics/blob-store/minio-wrapper"
	"github.com/trustedanalytics/tap-go-common/logger"
	"net/http"
	"os"
)

type Context struct {
	wrappedMinio *miniowrapper.Wrapper
}

var (
	logger = logger_wrapper.InitLogger("main")
	port   = os.Getenv("BLOB_STORE_PORT")
	host   = os.Getenv("BLOB_STORE_HOST")
)

const (
	bucketName = "blobstore"
	URLblobs   = "/api/v1/blobs/"
)

func main() {
	wrappedMinio, err := miniowrapper.CreateWrappedMinio(bucketName)
	if err != nil {
		logger.Fatal(err)
	}

	context := Context{wrappedMinio}
	router := web.New(context)
	router.Post(URLblobs, context.StoreBlob)
	router.Get(URLblobs+":blob_id", context.RetrieveBlob)
	router.Delete(URLblobs+":blob_id", context.RemoveBlob)

	err = http.ListenAndServe(host+":"+port, router)
	if err != nil {
		logger.Critical("Couldn't serve blob store on host:", host, ":", port, " Application will be closed now.")
		logger.Fatal(err)
	}
}
