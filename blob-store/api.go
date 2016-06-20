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
)

func (c *Context) StoreBlob(rw web.ResponseWriter, req *web.Request) {

	rw.WriteHeader(201)	//Created - the blob has been successfully stored
	//rw.WriteHeader(409)	//Conflict - the specified Blob ID is already in use
	//rw.WriteHeader(507)	//Insufficient Storage - the space allocated for the Blob Store has been exhausted.
}


func (c *Context) RetrieveBlob(rw web.ResponseWriter, req *web.Request) {
	//blob_id := req.PathParams["blob_id"]

	rw.WriteHeader(200)
	//rw.WriteHeader(404)	//Not found - if the specified Blob ID is not in blobStore
}


func (c *Context) RemoveBlob(rw web.ResponseWriter, req *web.Request) {
	//blob_id := req.PathParams["blob_id"]

	rw.WriteHeader(204)	//No content - the blob has been successfully removed
	//rw.WriteHeader(404)	//Not found - the specified Blob ID is not assigned to any of the blobs stored.
}

