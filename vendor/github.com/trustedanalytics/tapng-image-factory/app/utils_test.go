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
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gocraft/web"
	"github.com/smartystreets/goconvey/convey"
	"strings"
)

const (
	catalogHost   = "127.0.0.1"
	catalogPort   = "8083"
	imageId       = "25de06b4-ac21-4454-bb82-e72bc05f3a5c"
	imagesPath    = "/images/"
	blobStoreHost = "127.0.0.1"
	blobStorePort = "8084"
	blobId        = "img_" + imageId
	blobsPath     = "/blobs/"
	blob          = "abcdefg"
)

func prepareMocksAndRouter(t *testing.T) (router *web.Router) {
	router = web.New(Context{})
	return router
}

func sendRequest(rType, path string, body []byte, r *web.Router) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(rType, path, bytes.NewReader(body))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}

func marshalToJson(t *testing.T, serviceInstance interface{}) []byte {
	if body, err := json.Marshal(serviceInstance); err != nil {
		t.Errorf(err.Error())
		t.FailNow()
		return nil
	} else {
		return body
	}
}

func assertResponse(rr *httptest.ResponseRecorder, body string, code int) {
	if body != "" {
		convey.So(strings.TrimSpace(string(rr.Body.Bytes())), convey.ShouldContainSubstring, body)
	}
	convey.So(rr.Code, convey.ShouldEqual, code)
}
