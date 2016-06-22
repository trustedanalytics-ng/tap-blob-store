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
	"github.com/trustedanalytics/blob-store/minioClient"
)

const (
	ErrKeyNotExist = "The specified key does not exist."
)

func (c *Context) StoreBlob(rw web.ResponseWriter, req *web.Request) {
	blob_id := req.FormValue("blob_id")

	logger.Info("Storing blob -", blob_id)

	blob,err := getBlobFromRequest(rw, req)
	if err != nil {
		switch err.Error() {
		//case Err...:
		//	logger.Notice(err)
		//	http.Error(rw, "...", http.StatusBadRequest)
		default:
			logUnhandledError(rw, err)
		}
		return
		logger.Error(err)
	}

	err = minioClient.StoreInMinio(blob_id, blob)
	if err != nil {
		switch err.Error() {
		//case Err...:
		//	logger.Notice(err)
		//	http.Error(rw, "The specified Blob ID is already in use", http.StatusConflict)
		//case Err...:
		//	logger.Notice(err)
		//	http.Error(rw, "The space allocated for the Blob Store has been exhausted", 507)
		default:
			logUnhandledError(rw, err)
		}
		return
	}

	http.Error(rw, "The blob has been successfully stored", http.StatusCreated)
}

func getBlobFromRequest(w web.ResponseWriter, r *web.Request) (multipart.File, error) {
	logger.Info("Getting blob from request")

	r.ParseMultipartForm(32 << 20)
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

	blob, err := minioClient.RetrieveFromMinio(blob_id)
	if err != nil {
		switch err.Error() {
		case ErrKeyNotExist:
			logger.Notice(err)
			http.Error(rw, "The specified blob does not exist", http.StatusNotFound)
		default:
			logUnhandledError(rw, err)
		}
		return
	}
	blob.Stat()

	err = seekThroughFile(blob)
	if err != nil {
		logUnhandledError(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/octet-stream")
	http.ServeContent(rw, req.Request, blob_id, time.Now(), blob)
}

func seekThroughFile(blob *minio.Object) error {
	stat, err := blob.Stat()
	if err != nil {
		return err
	}

	_, err = blob.Seek(stat.Size, os.SEEK_CUR)
	if err != nil {
		return err
	}

	logger.Info("Seek through file. Size of file -", stat.Size)
	return nil
}

func (c *Context) RemoveBlob(rw web.ResponseWriter, req *web.Request) {
	blob_id := req.PathParams["blob_id"]

	err := minioClient.RemoveFromMinio(blob_id)
	if err != nil {
		switch err.Error() {
		case ErrKeyNotExist:
			logger.Notice(err)
			http.Error(rw, "The specified blob does not exist", http.StatusNotFound)
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

	logger.Error(error_id, err)
	http.Error(rw, fmt.Sprint("Unhandled Exception, errorID = ", error_id), http.StatusInternalServerError)
}