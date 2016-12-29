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
package api_test_utils

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/gocraft/web"
	"github.com/smartystreets/goconvey/convey"
)

func SendRequest(rType, path string, body []byte, r *web.Router) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(rType, path, bytes.NewReader(body))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}

func SendForm(path string, body io.Reader, contentType string, r *web.Router) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("POST", path, body)
	req.Header.Set("Content-Type", contentType)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}

func PrepareForm(blobId, filepath string) (bodyBuf *bytes.Buffer, contentType string) {
	bodyBuf = &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	defer bodyWriter.Close()

	if blobId != "" {
		bodyWriter.WriteField("blobID", blobId)
	}

	if filepath != "" {
		fileWriter, _ := bodyWriter.CreateFormFile("uploadfile", filepath)
		fh, err := os.Open(filepath) //TODO: It requires Real File, Mock It !!!
		if err != nil {
			panic(err)
		}
		_, err = io.Copy(fileWriter, fh)
		if err != nil {
			panic(err)
		}
	}

	contentType = bodyWriter.FormDataContentType()
	return
}

func AssertResponse(rr *httptest.ResponseRecorder, body string, code int) {
	if body != "" {
		convey.Convey("Body should be expected", func() {
			convey.So(strings.TrimSpace(string(rr.Body.Bytes())), convey.ShouldContainSubstring, body)
		})
	}
	convey.Convey("Response code should equal "+strconv.Itoa(code), func() {
		convey.So(rr.Code, convey.ShouldEqual, code)
	})
}

func MarshallToJson(t *testing.T, serviceInstance interface{}) []byte {
	if body, err := json.Marshal(serviceInstance); err != nil {
		t.Errorf(err.Error())
		t.FailNow()
		return nil
	} else {
		return body
	}
}
