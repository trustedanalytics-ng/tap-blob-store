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
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/gocraft/web"
	"github.com/minio/minio-go"

	"github.com/trustedanalytics/tap-blob-store/minio-wrapper"
	commonLogger "github.com/trustedanalytics/tap-go-common/logger"
)

const (
	ErrMsgKeyNotExist      = "The specified key does not exist."
	ErrMsgBlobNotSpecified = "http: no such file"
	URLblobs               = "/blobs/"
	defaultMaxMemory       = 32 << 20 // 32 MB
)

var (
	logger, _ = commonLogger.InitLogger("api")
	blobStat  = minioBlobStat
	blobSeek  = minioBlobSeek
	blobServe = minioBlobServe
)

type ApiContext struct {
	WrappedMinio *miniowrapper.Wrapper
}

func NewApiContext(wrapper *miniowrapper.Wrapper) *ApiContext {
	return &ApiContext{wrapper}
}

func minioBlobStat(blob *minio.Object) (minio.ObjectInfo, error) {
	return blob.Stat()
}

func minioBlobSeek(blob *minio.Object, offset int64, whence int) (n int64, err error) {
	return blob.Seek(offset, whence)
}

func minioBlobServe(w http.ResponseWriter, req *http.Request, name string, modtime time.Time, content io.ReadSeeker) {
	http.ServeContent(w, req, name, modtime, content)
}

func RegisterRoutes(router *web.Router, context ApiContext) {
	router.Middleware(context.BasicAuthorizeMiddleware)

	router.Post(URLblobs, context.StoreBlob)
	router.Get(URLblobs+":blob_id", context.RetrieveBlob)
	router.Delete(URLblobs+":blob_id", context.RemoveBlob)
}

func (c *ApiContext) StoreBlob(rw web.ResponseWriter, req *web.Request) {
	blobID := req.FormValue("blob_id")
	if blobID == "" {
		logNoticedError(rw, "The blob_id is not specified.", nil, http.StatusBadRequest)
		return
	}

	logger.Info("Storing blob:", blobID)
	blob, err := getBlobFromRequest(rw, req)
	if err != nil {
		switch err.Error() {
		case ErrMsgBlobNotSpecified:
			logNoticedError(rw, "Blob not specified.", err, http.StatusBadRequest)
		default:
			logUnhandledError(rw, err)
		}
		return
	}
	err = c.WrappedMinio.StoreInMinio(blobID, blob)
	if err != nil {
		switch err.Error() {
		case miniowrapper.ErrKeyAlreadyInUse.Error():
			logNoticedError(rw, "The specified Blob ID is already in use", err, http.StatusConflict)
		//TODO: github.com/minio/minio/api-errors.go there is something like ErrStorageFull
		//case TBD:
		//	logNoticedError(rw, "The space allocated for the Blob Store has been exhausted", err, 507)
		default:
			logUnhandledError(rw, err)
		}
		return
	}

	http.Error(rw, "The blob has been successfully stored", http.StatusCreated)
}

func getBlobFromRequest(w web.ResponseWriter, r *web.Request) (multipart.File, error) {
	logger.Info("Getting blob from request")

	r.ParseMultipartForm(defaultMaxMemory)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		return nil, err
	}

	logger.Info("Request content:", handler.Header)
	return file, nil
}

func (c *ApiContext) RetrieveBlob(rw web.ResponseWriter, req *web.Request) {
	blobID := req.PathParams["blob_id"]
	logger.Info("Retrieving blob:", blobID)

	blob, err := c.WrappedMinio.RetrieveFromMinio(blobID)
	if err != nil {
		switch err.Error() {
		case ErrMsgKeyNotExist:
			logNoticedError(rw, "The specified blob does not exist.", err, http.StatusNotFound)
		default:
			logUnhandledError(rw, err)
		}
		return
	}

	err = seekThroughFile(blob)
	if err != nil {
		logUnhandledError(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/octet-stream")
	blobServe(rw, req.Request, blobID, time.Now(), blob)
}

func seekThroughFile(blob *minio.Object) error {
	stat, err := blobStat(blob)
	if err != nil {
		return err
	}

	_, err = blobSeek(blob, stat.Size, os.SEEK_CUR)
	if err != nil {
		return err
	}

	logger.Debug("Seek through file. Size of file:", stat.Size)
	return nil
}

func (c *ApiContext) RemoveBlob(rw web.ResponseWriter, req *web.Request) {
	blobID := req.PathParams["blob_id"]
	logger.Info("Removing blob:", blobID)

	err := c.WrappedMinio.RemoveFromMinio(blobID)
	if err != nil {
		switch err.Error() {
		case ErrMsgKeyNotExist:
			logNoticedError(rw, "The specified blob does not exist.", err, http.StatusNotFound)
		default:
			logUnhandledError(rw, err)
		}
		return
	}

	http.Error(rw, "The blob has been successfully removed", http.StatusNoContent)
}

func logUnhandledError(rw web.ResponseWriter, err error) {
	logInWrapper(logger.Error, "Unhandled Exception, error:", err)
	http.Error(rw, fmt.Sprint("Unhandled Exception, error:", err), http.StatusInternalServerError)
}

func logNoticedError(rw web.ResponseWriter, message string, err error, statusCode int) {
	logInWrapper(logger.Notice, message, "Error: ", err)
	http.Error(rw, message, statusCode)
}

func logInWrapper(logLevel func(args ...interface{}), args ...interface{}) {
	logger.ExtraCalldepth += 3
	logLevel(args...)
	logger.ExtraCalldepth -= 3
}
