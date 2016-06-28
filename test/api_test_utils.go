package api_test_utils

import (
	"testing"
	"net/http"
	"net/http/httptest"

	"strings"
	"bytes"
	"encoding/json"

	"github.com/gocraft/web"
	"github.com/smartystreets/goconvey/convey"
	"io"
	"mime/multipart"
	"os"
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

func PrepareForm(blobId, filename string) (bodyBuf *bytes.Buffer, contentType string)  {

	bodyBuf = &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	if blobId != "" {
		bodyWriter.WriteField("blob_id", blobId)
	}

	if(filename != "") {
		fileWriter, _ := bodyWriter.CreateFormFile("uploadfile", filename)
		fh, _ := os.Open(filename)	//TODO: It requires Real File, Mock It !!!
		_, _ = io.Copy(fileWriter, fh)
	}

	contentType = bodyWriter.FormDataContentType()
	bodyWriter.Close()

	return //magic :)
}

func AssertResponse(rr *httptest.ResponseRecorder, body string, code int) {
	if body != "" {
		convey.So(strings.TrimSpace(string(rr.Body.Bytes())), convey.ShouldContainSubstring, body)
	}
	convey.So(rr.Code, convey.ShouldEqual, code)
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
