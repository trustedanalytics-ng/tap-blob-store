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
	"errors"
	"github.com/trustedanalytics/tapng-catalog/models"
	"io"
	"os"
)

type ImageGetResponse struct {
	ImageId    string `json:"id"`
	Type       string `json:"type"`
	State      string `json:"state"`
	AuditTrail models.AuditTrail
}

type ImageStatePutRequest struct {
	State string `json:"state"`
}

type BuildImagePostRequest struct {
	ImageId string `json:"id"`
}

var (
	ImagesMap = map[string]string{
		"JAVA":   "tap-base-java:java8-jessie",
		"GO":     "tap-base-binary:binary-jessie",
		"NODEJS": "tap-base-node:node4.4-jessie",
		"PYTHON": "tap-base-python:python2.7-jessie",
	}
)

func GetCatalogAddress() string {
	return os.Getenv("CATALOG_HOST") + ":" + os.Getenv("CATALOG_PORT") + "/api/v1"
}

func GetBlobStoreAddress() string {
	return os.Getenv("BLOB_STORE_HOST") + ":" + os.Getenv("BLOB_STORE_PORT") + "/api/v1"
}

func GetHubAddress() string {
	return os.Getenv("HUB_ADDRESS")
}

func GetDockerApiVersion() string {
	return os.Getenv("DOCKER_API_VERSION")
}

func GetDockerHostAddress() string {
	return os.Getenv("DOCKER_HOST")
}

func StreamToByte(stream io.Reader) ([]byte, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(stream)
	if err != nil {
		return nil, errors.New("Could not read stream into byte array: " + err.Error())
	}
	return buf.Bytes(), nil
}

func StreamToString(stream io.Reader) (string, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(stream)
	if err != nil {
		return "", errors.New("Could not read stream into string: " + err.Error())
	}
	return buf.String(), nil
}
