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
	"bytes"
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

func NewTapBlobStoreApiWithSSLAndBasicAuth(address, username, password, certPemFile, keyPemFile, caPemFile string) (*TapBlobStoreApiConnector, error) {
	client, _, err := brokerHttp.GetHttpClientWithCertAndCaFromFile(certPemFile, keyPemFile, caPemFile)
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
		BasicAuth: &brokerHttp.BasicAuth{c.Username, c.Password},
		Client:    c.Client,
		Url:       url,
	}
}

func (c *TapBlobStoreApiConnector) StoreBlob(blob_id string, file multipart.File) error {
	connector := c.getApiConnector(fmt.Sprintf("%s/api/v1/blobs", c.Address))

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	err := bodyWriter.WriteField("blob_id", blob_id)
	if err != nil {
		return err
	}

	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", "blob.tar.gz")
	if err != nil {

		fmt.Println("error writing to buffer")
		return err
	}

	size, err := io.Copy(fileWriter, file)
	if err != nil {
		return err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	var req *http.Request
	req, _ = http.NewRequest("POST", connector.Url, bytes.NewBuffer(bodyBuf.Bytes()))
	req.Header.Add("Authorization", brokerHttp.GetBasicAuthHeader(connector.BasicAuth))
	brokerHttp.SetContentType(req, contentType)

	logger.Infof("Doing: POST %v Sending %v bytes", connector.Url, size)
	_, err = connector.Client.Do(req)
	if err != nil {
		logger.Error("ERROR: Make http request POST", err)
		return err
	}

	return err
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
	if err != nil {
		return status, err
	}

	return status, err
}
