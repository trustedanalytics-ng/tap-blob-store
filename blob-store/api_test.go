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
	"github.com/trustedanalytics/blob-store/minioWrapper"

	TestUtils "github.com/trustedanalytics/blob-store/test"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	NewBlobId = "17"
	UnhandledBlobId = "unhandledBlob"
	ExistedBlobId = "1234"
	NilBlobId = "nil"

	TestFileName = "testFile.txt"
)


func prepareMocksAndRouter(t *testing.T) (router *web.Router, c Context) {
	c = Context{&minioWrapper.Wrapper{&MinioClientMock{},""}}
	router = web.New(c)
	return router, c
}

func TestStoreBlob(t *testing.T) {
	router, context := prepareMocksAndRouter(t)
	router.Post(URLblobs, context.StoreBlob)

	Convey("Test Store Blob", t, func() {
		Convey("Blob ID not specified. Should returns error message", func() {
			bodyBuf, contentType := TestUtils.PrepareForm("", "")
			response := TestUtils.SendForm(URLblobs, bodyBuf, contentType, router)
			TestUtils.AssertResponse(response, "The blob_id is not specified", 400)
		})
		Convey("Blob file not specified, . Should returns error message", func() {
			bodyBuf, contentType := TestUtils.PrepareForm(NewBlobId, "")
			response := TestUtils.SendForm(URLblobs, bodyBuf, contentType, router)
			TestUtils.AssertResponse(response, "Blob not specified.", 400)
		})
		Convey("Blob ID already exists. Should returns error message", func() {
			bodyBuf, contentType := TestUtils.PrepareForm(ExistedBlobId, TestFileName)
			response := TestUtils.SendForm(URLblobs, bodyBuf, contentType, router)
			TestUtils.AssertResponse(response, "The specified Blob ID is already in use", 409)
		})
		Convey("Error! Should returns unhandled exception", func() {
			bodyBuf, contentType := TestUtils.PrepareForm(UnhandledBlobId, TestFileName)
			response := TestUtils.SendForm(URLblobs, bodyBuf, contentType, router)
			TestUtils.AssertResponse(response, "Unhandled Exception, errorID = ", 500)
		})
		Convey("Should returns proper response", func() {
			bodyBuf, contentType := TestUtils.PrepareForm(NewBlobId, TestFileName)
			response := TestUtils.SendForm(URLblobs, bodyBuf, contentType, router)
			TestUtils.AssertResponse(response, "", 201)
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

	defer func () {
		blobStat = oldBlobStat
		blobSeek = oldBlobSeek
		blobServe = oldBlobServe
	}()

	router, context := prepareMocksAndRouter(t)
	router.Get(URLblobs + ":blob_id", context.RetrieveBlob)

	Convey("Test Retrieve Blob", t, func() {
		Convey("Blob ID not existed. Should returns error message", func() {
			response := TestUtils.SendRequest("GET", URLblobs + NewBlobId, nil, router)
			TestUtils.AssertResponse(response, "The specified blob does not exist.", 404)
		})
		Convey("Error! Should returns unhandled exception", func() {
			response := TestUtils.SendRequest("GET", URLblobs + UnhandledBlobId, nil, router)
			TestUtils.AssertResponse(response, "Unhandled Exception, errorID = ", 500)
		})
		Convey("Blob ID exist, but Minio contains nil object. Should returns error message", func() {
			response := TestUtils.SendRequest("GET", URLblobs + NilBlobId, nil, router)
			TestUtils.AssertResponse(response, "Unhandled Exception, errorID = ", 500)
		})
		Convey("Should returns proper response", func() {
			response := TestUtils.SendRequest("GET", URLblobs + ExistedBlobId, nil, router)
			TestUtils.AssertResponse(response, "", 200)
		})
	})
}

func TestRemoveBlob(t *testing.T) {
	router, context := prepareMocksAndRouter(t)
	router.Delete(URLblobs + ":blob_id", context.RemoveBlob)

	Convey("Test Remove Blob", t, func() {
		Convey("Blob ID not existed. Should returns error message", func() {
			response := TestUtils.SendRequest("DELETE", URLblobs + NewBlobId, nil, router)
			TestUtils.AssertResponse(response, "The specified blob does not exist.", 404)
		})
		Convey("Error! Should returns unhandled exception", func() {
			response := TestUtils.SendRequest("DELETE", URLblobs + UnhandledBlobId, nil, router)
			TestUtils.AssertResponse(response, "Unhandled Exception, errorID = ", 500)
		})
		Convey("Should returns proper response", func() {
			response := TestUtils.SendRequest("DELETE", URLblobs + ExistedBlobId, nil, router)
			TestUtils.AssertResponse(response, "", 204)
		})
	})
}

