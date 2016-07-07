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
	"github.com/gocraft/web"
	"github.com/trustedanalytics/tapng-go-common/util"
)

type Context struct {
	BlobStoreConnector *Connector
	CatalogConnector   *Connector
	DockerConnector    *DockerClient
}

func (c *Context) SetupContext() {
	c.BlobStoreConnector = NewBlobStoreConnector()
	c.CatalogConnector = NewCatalogConnector()
	dockerClient, err := NewDockerClient()
	if err != nil {
		logger.Panic(err)
	}
	c.DockerConnector = dockerClient
}

func (c *Context) BuildImage(rw web.ResponseWriter, req *web.Request) {
	req_json := BuildImagePostRequest{}
	err := util.ReadJson(req, &req_json)
	if err != nil {
		logger.Error(err.Error())
		rw.WriteHeader(400)
		return
	}
	imgDetails, err := c.CatalogConnector.GetImageDetails(req_json.ImageId)
	if err != nil {
		logger.Error(err.Error())
		rw.WriteHeader(500)
		return
	}
	blobBytes, err := c.BlobStoreConnector.GetImageBlob(imgDetails.ImageId)
	if err != nil {
		logger.Error(err.Error())
		rw.WriteHeader(500)
		return
	}
	err = c.CatalogConnector.UpdateImageState(imgDetails.ImageId, "BUILDING")
	if err != nil {
		logger.Error(err.Error())
		rw.WriteHeader(500)
		return
	}
	tag := GetDockerHostAddress() + "/" + imgDetails.ImageId
	err = c.DockerConnector.CreateImage(bytes.NewReader(blobBytes), imgDetails.Type, tag)
	if err != nil {
		logger.Error(err.Error())
		rw.WriteHeader(500)
		return
	}
	err = c.DockerConnector.PushImage(tag)
	if err != nil {
		logger.Error(err.Error())
		rw.WriteHeader(500)
		return
	}
	err = c.CatalogConnector.UpdateImageState(imgDetails.ImageId, "READY")
	if err != nil {
		logger.Error(err.Error())
		rw.WriteHeader(500)
		return
	}
	err = c.BlobStoreConnector.DeleteImageBlob(imgDetails.ImageId)
	if err != nil {
		logger.Error(err.Error())
		rw.Write([]byte(err.Error()))
		rw.WriteHeader(500)
		return
	}
	rw.WriteHeader(201)
}
