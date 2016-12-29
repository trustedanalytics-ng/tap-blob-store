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

package api

import (
	"net/http"
	"testing"

	"github.com/gocraft/web"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/trustedanalytics/tap-blob-store/minio-wrapper"
	TestUtils "github.com/trustedanalytics/tap-blob-store/test"
)

const (
	NewBlobID       = "17"
	UnhandledBlobID = "unhandledBlob"
	ExistedBlobID   = "1234"
	NilBlobID       = "nil"

	TestFilePath = "../test/testFile.txt"
)

func prepareMocksAndRouter(t *testing.T) (router *web.Router, c ApiContext) {
	c = ApiContext{WrappedMinio: &miniowrapper.Wrapper{Mc: &MinioClientMock{}, BucketName: ""}}
	router = web.New(c)
	return router, c
}

func TestStoreBlob(t *testing.T) {
	router, context := prepareMocksAndRouter(t)
	router.Post(URLblobs, context.StoreBlob)

	Convey("Test Store Blob", t, func() {
		Convey("Blob ID not specified. Should return error message", func() {
			bodyBuf, contentType := TestUtils.PrepareForm("", "")
			response := TestUtils.SendForm(URLblobs, bodyBuf, contentType, router)
			TestUtils.AssertResponse(response, "The blobID is not specified", http.StatusBadRequest)
		})
		Convey("Blob file not specified. Should return error message", func() {
			bodyBuf, contentType := TestUtils.PrepareForm(NewBlobID, "")
			response := TestUtils.SendForm(URLblobs, bodyBuf, contentType, router)
			TestUtils.AssertResponse(response, "Blob not specified.", http.StatusBadRequest)
		})
		Convey("Blob ID already exists. Should return error message", func() {
			bodyBuf, contentType := TestUtils.PrepareForm(ExistedBlobID, TestFilePath)
			response := TestUtils.SendForm(URLblobs, bodyBuf, contentType, router)
			TestUtils.AssertResponse(response, "The specified Blob ID is already in use", http.StatusConflict)
		})
		Convey("Error! Should return unhandled exception", func() {
			bodyBuf, contentType := TestUtils.PrepareForm(UnhandledBlobID, TestFilePath)
			response := TestUtils.SendForm(URLblobs, bodyBuf, contentType, router)
			TestUtils.AssertResponse(response, "Unhandled Exception, error:", http.StatusInternalServerError)
		})
		Convey("Should return proper response", func() {
			bodyBuf, contentType := TestUtils.PrepareForm(NewBlobID, TestFilePath)
			response := TestUtils.SendForm(URLblobs, bodyBuf, contentType, router)
			TestUtils.AssertResponse(response, "", http.StatusCreated)
		})
	})
}

func TestRetrieveBlob(t *testing.T) {

	oldBlobStat := blobStat
	oldBlobSeek := blobSeek
	oldBlobServe := blobServe

	blobStat = mockBlobStat
	blobSeek = mockBlobSeek
	blobServe = mockBlobServe

	defer func() {
		blobStat = oldBlobStat
		blobSeek = oldBlobSeek
		blobServe = oldBlobServe
	}()

	router, context := prepareMocksAndRouter(t)
	router.Get(URLblobs+":blob_id", context.RetrieveBlob)

	Convey("Test Retrieve Blob", t, func() {
		Convey("Blob ID not existed. Should return error message", func() {
			response := TestUtils.SendRequest("GET", URLblobs+NewBlobID, nil, router)
			TestUtils.AssertResponse(response, "The specified blob does not exist.", http.StatusNotFound)
		})
		Convey("Error! Should return unhandled exception", func() {
			response := TestUtils.SendRequest("GET", URLblobs+UnhandledBlobID, nil, router)
			TestUtils.AssertResponse(response, "Unhandled Exception, error:", http.StatusInternalServerError)
		})
		Convey("Blob ID exist, but Minio contains nil object. Should return error message", func() {
			response := TestUtils.SendRequest("GET", URLblobs+NilBlobID, nil, router)
			TestUtils.AssertResponse(response, "Unhandled Exception, error:", http.StatusInternalServerError)
		})
		Convey("Should return proper response", func() {
			response := TestUtils.SendRequest("GET", URLblobs+ExistedBlobID, nil, router)
			TestUtils.AssertResponse(response, "", http.StatusOK)
		})
	})
}

func TestRemoveBlob(t *testing.T) {
	router, context := prepareMocksAndRouter(t)
	router.Delete(URLblobs+":blob_id", context.RemoveBlob)

	Convey("Test Remove Blob", t, func() {
		Convey("Blob ID not existed. Should return error message", func() {
			response := TestUtils.SendRequest("DELETE", URLblobs+NewBlobID, nil, router)
			TestUtils.AssertResponse(response, "The specified blob does not exist.", http.StatusNotFound)
		})
		Convey("Error! Should return unhandled exception", func() {
			response := TestUtils.SendRequest("DELETE", URLblobs+UnhandledBlobID, nil, router)
			TestUtils.AssertResponse(response, "Unhandled Exception, error:", http.StatusInternalServerError)
		})
		Convey("Should return proper response", func() {
			response := TestUtils.SendRequest("DELETE", URLblobs+ExistedBlobID, nil, router)
			TestUtils.AssertResponse(response, "", http.StatusNoContent)
		})
	})
}
