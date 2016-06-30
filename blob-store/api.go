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
	"time"
	"os"
	"fmt"
	"math/rand"
	"mime/multipart"
	"github.com/minio/minio-go"
	"github.com/gocraft/web"
	"github.com/trustedanalytics/blob-store/minioWrapper"
	"io"
)

const (
	ErrMsgKeyNotExist = "The specified key does not exist."
	ErrMsgBlobNotSpecified = "http: no such file"

	defaultMaxMemory = 32 << 20 // 32 MB
)

var (
	blobStat = minioBlobStat
	blobSeek = minioBlobSeek
	blobServe = minioBlobServe
)

func minioBlobStat(blob *minio.Object) (minio.ObjectInfo, error) {
	return blob.Stat()
}

func minioBlobSeek(blob *minio.Object, offset int64, whence int) (n int64, err error) {
	return blob.Seek(offset, whence)
}

func minioBlobServe(w http.ResponseWriter, req *http.Request, name string, modtime time.Time, content io.ReadSeeker) {
	http.ServeContent(w,req,name,modtime,content)
}


func (c *Context) StoreBlob(rw web.ResponseWriter, req *web.Request) {
	blob_id := req.FormValue("blob_id")
	if(blob_id == "") {
		logNoticedError(rw, "The blob_id is not specified.", nil, http.StatusBadRequest)
		return
	}

	logger.Info("Storing blob -", blob_id)
	blob,err := getBlobFromRequest(rw, req)
	if err != nil {
		switch err.Error() {
		case ErrMsgBlobNotSpecified:
			logNoticedError(rw, "Blob not specified.", err, http.StatusBadRequest)
		default:
			logUnhandledError(rw, err)
		}
		return
	}

	err = c.wrappedMinio.StoreInMinio(blob_id, blob)
	if err != nil {
		switch err.Error() {
		case minioWrapper.ErrKeyAlreadyInUse.Error():
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
	defer file.Close()

	logger.Info("Request content -", handler.Header)
	return file, nil
}

func (c *Context) RetrieveBlob(rw web.ResponseWriter, req *web.Request) {
	blob_id := req.PathParams["blob_id"]
	logger.Info("Retrieving blob -", blob_id)

	blob, err := c.wrappedMinio.RetrieveFromMinio(blob_id)
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
	blobServe(rw, req.Request, blob_id, time.Now(), blob)
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

	logger.Info("Seek through file. Size of file -", stat.Size)
	return nil
}

func (c *Context) RemoveBlob(rw web.ResponseWriter, req *web.Request) {
	blob_id := req.PathParams["blob_id"]
	logger.Info("Removing blob -", blob_id)

	err := c.wrappedMinio.RemoveFromMinio(blob_id)
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
	rand.Seed( time.Now().UnixNano())
	error_id := rand.Intn(999999)

	logInWrapper(logger.Error, "errorID =", error_id, "-", err)
	http.Error(rw, fmt.Sprint("Unhandled Exception, errorID = ", error_id), http.StatusInternalServerError)
}

func logNoticedError(rw web.ResponseWriter, message string, err error, statusCode int) {
	logInWrapper(logger.Notice, message, "Error: ", err)
	http.Error(rw, message, statusCode)
}

func logInWrapper(logLevel func(args ...interface{}), args ...interface{}) {
	logger.ExtraCalldepth += 3
	logLevel(args ...)
	logger.ExtraCalldepth -= 3
}
