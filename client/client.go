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
package client

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	brokerHttp "github.com/trustedanalytics/tap-go-common/http"
	commonLogger "github.com/trustedanalytics/tap-go-common/logger"
)

var (
	logger, _ = commonLogger.InitLogger("client")
)

type TapBlobStoreApi interface {
	StoreBlob(blob_id string, file multipart.File) error
	GetBlob(blob_id string, dest io.Writer) error
	DeleteBlob(blob_id string) (int, error)
}

func NewTapBlobStoreApiWithBasicAuth(address, username, password string) (*TapBlobStoreApiConnector, error) {
	client, _, err := brokerHttp.GetHttpClient()
	if err != nil {
		return nil, err
	}
	return &TapBlobStoreApiConnector{address, username, password, client}, nil
}

type TapBlobStoreApiConnector struct {
	Address  string
	Username string
	Password string
	Client   *http.Client
}

func (c *TapBlobStoreApiConnector) getApiConnector(url string) brokerHttp.ApiConnector {
	return brokerHttp.ApiConnector{
		BasicAuth: &brokerHttp.BasicAuth{User: c.Username, Password: c.Password},
		Client:    c.Client,
		Url:       url,
	}
}

func (c *TapBlobStoreApiConnector) StoreBlob(blobID string, file multipart.File) error {
	connector := c.getApiConnector(fmt.Sprintf("%s/api/v1/blobs", c.Address))

	bodyPipeReader, bodyPipeWriter := io.Pipe()
	defer func() {
		err := bodyPipeReader.Close()
		if err != nil {
			logger.Error("Pipe reader closing error: ", err)
		}
	}()

	go writeBlobAsync(bodyPipeWriter, blobID, file)

	var req *http.Request
	req, err := http.NewRequest("POST", connector.Url, bodyPipeReader)
	if err != nil {
		logger.Error("Creating request failed: ", err)
		return err
	}

	req.Header.Add("Authorization", brokerHttp.GetBasicAuthHeader(connector.BasicAuth))

	logger.Infof("Doing: POST %v ", connector.Url)
	_, err = connector.Client.Do(req)
	if err != nil {
		logger.Error("Make http request POST failed: ", err)
		return err
	}

	return nil
}

func writeBlobAsync(pw *io.PipeWriter, blobID string, blobFile multipart.File) {
	var err error
	defer func() {
		if err != nil {
			err:=pw.CloseWithError(err)
			if err != nil {
				logger.Error("Pipe writer closing with error: ", err)
			}
		} else {
			err:=pw.Close()
			if err != nil {
				logger.Error("Pipe writer closing error: ", err)
			}
		}
	}()

	bodyWriter := multipart.NewWriter(pw)
	defer func(){
		err := bodyWriter.Close()
		if err != nil {
			logger.Error("Body writer closing with error: ", err)
		}
	}()

	err = bodyWriter.WriteField("blob_id", blobID)
	if err != nil {
		logger.Errorf("bodyWriter.WriteField(%v) failed: %v", blobID, err)
		return
	}

	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", "blob.tar.gz")
	if err != nil {
		logger.Errorf("bodyWriter.CreateFormFile(\"uploadfile\", \"blob.tar.gz\") failed: %v", err)
		return
	}

	_, err = io.Copy(fileWriter, blobFile)
	if err != nil {
		logger.Errorf("copying to writer failed: %v", err)
	}
}

func (c *TapBlobStoreApiConnector) GetBlob(blob_id string, dest io.Writer) error {
	connector := c.getApiConnector(fmt.Sprintf("%s/api/v1/blobs/%s", c.Address, blob_id))
	size, err := brokerHttp.DownloadBinary(connector.Url, brokerHttp.GetBasicAuthHeader(connector.BasicAuth), connector.Client, dest)
	if err != nil {
		return err
	}
	logger.Infof("Written %v bytes of binary data to destination", size)
	return err
}

func (c *TapBlobStoreApiConnector) DeleteBlob(blob_id string) (int, error) {
	connector := c.getApiConnector(fmt.Sprintf("%s/api/v1/blobs/%s", c.Address, blob_id))
	status, _, err := brokerHttp.RestDELETE(connector.Url, "", brokerHttp.GetBasicAuthHeader(connector.BasicAuth), connector.Client)
	
	return status, err
}
