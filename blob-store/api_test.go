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
	"testing"

	"github.com/gocraft/web"

	TestUtils "github.com/trustedanalytics/blob-store/test"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	URLstoreBlob = "/api/v1/blobs"
	URLretrieveBlob = "/api/v1/blobs/"
	URLremoveBlob = "/api/v1/blobs/"
	blobId = "17"
)

func prepareMocksAndRouter(t *testing.T) (router *web.Router) {
	router = web.New(Context{})
	return router
}

func TestStoreBlob(t *testing.T) {
	router := prepareMocksAndRouter(t)
	router.Post(URLstoreBlob, (*Context).StoreBlob)

	Convey("Test Store Blob", t, func() {
		Convey("Should returns proper response", func() {
			response := TestUtils.SendRequest("POST", URLstoreBlob, nil, router)
			TestUtils.AssertResponse(response, "", 201)
		})
	})
}

func TestRetrieveBlob(t *testing.T) {
	router := prepareMocksAndRouter(t)
	router.Get(URLretrieveBlob + ":blob_id", (*Context).RetrieveBlob)

	Convey("Test Retrieve Blob", t, func() {
		Convey("Should returns proper response", func() {
			response := TestUtils.SendRequest("GET", URLretrieveBlob + blobId, nil, router)
			TestUtils.AssertResponse(response, "", 200)
		})
	})
}

func TestRemoveBlob(t *testing.T) {
	router := prepareMocksAndRouter(t)
	router.Delete(URLremoveBlob + ":blob_id", (*Context).RemoveBlob)

	Convey("Test Remove Blob", t, func() {
		Convey("Should returns proper response", func() {
			response := TestUtils.SendRequest("DELETE", URLremoveBlob + blobId, nil, router)
			TestUtils.AssertResponse(response, "", 204)
		})
	})
}

