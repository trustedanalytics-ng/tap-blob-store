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
	"github.com/jarcoal/httpmock"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestGetImageDetails(t *testing.T) {
	originalCatalogHost := os.Getenv("CATALOG_HOST")
	originalCatalogPort := os.Getenv("CATALOG_PORT")
	os.Setenv("CATALOG_HOST", catalogHost)
	os.Setenv("CATALOG_PORT", catalogPort)
	os.Setenv("IMAGE_FACTORY_PORT", "8080")
	defer os.Setenv("CATALOG_HOST", originalCatalogHost)
	defer os.Setenv("CATALOG_PORT", originalCatalogPort)
	httpmock.Activate()
	c := NewCatalogConnector()
	c.Client.Transport = httpmock.DefaultTransport
	defer httpmock.DeactivateAndReset()

	Convey("Test GetImageDetails", t, func() {
		Convey("Should return proper response", func() {
			httpmock.RegisterResponder("GET", GetCatalogAddress()+imagesPath+imageId,
				httpmock.NewStringResponder(200, `{"id":"`+imageId+`"}`))
			res, err := c.GetImageDetails(imageId)
			So(err, ShouldBeNil)
			So(res.ImageId, ShouldEqual, imageId)
		})

		Convey("Should return not found", func() {
			httpmock.RegisterResponder("GET", GetCatalogAddress()+imagesPath+imageId,
				httpmock.NewStringResponder(404, ``))
			res, err := c.GetImageDetails(imageId)
			So(err.Error(), ShouldEqual, "Invalid status: 404")
			So(res.ImageId, ShouldBeEmpty)
		})
	})
}

func TestUpdateImageState(t *testing.T) {
	originalCatalogHost := os.Getenv("CATALOG_HOST")
	originalCatalogPort := os.Getenv("CATALOG_PORT")
	os.Setenv("CATALOG_HOST", catalogHost)
	os.Setenv("CATALOG_PORT", catalogPort)
	os.Setenv("IMAGE_FACTORY_PORT", "8080")
	defer os.Setenv("CATALOG_HOST", originalCatalogHost)
	defer os.Setenv("CATALOG_PORT", originalCatalogPort)
	httpmock.Activate()
	c := NewCatalogConnector()
	c.Client.Transport = httpmock.DefaultTransport
	defer httpmock.DeactivateAndReset()

	Convey("Test UpdateImageState", t, func() {
		Convey("Should return proper response", func() {
			httpmock.RegisterResponder("PATCH", GetCatalogAddress()+imagesPath+imageId,
				httpmock.NewStringResponder(200, ``))
			err := c.UpdateImageState(imageId, "created")
			So(err, ShouldBeNil)

		})

		Convey("Should return not found", func() {
			httpmock.RegisterResponder("PATCH", GetCatalogAddress()+imagesPath+imageId,
				httpmock.NewStringResponder(404, ``))
			err := c.UpdateImageState(imageId, "created")
			So(err.Error(), ShouldEqual, "Invalid status: 404")
		})
	})
}

func TestGetBlob(t *testing.T) {
	originalBlobStoreHost := os.Getenv("BLOB_STORE_HOST")
	originalBlobStorePort := os.Getenv("BLOB_STORE_PORT")
	os.Setenv("BLOB_STORE_HOST", blobStoreHost)
	os.Setenv("BLOB_STORE_PORT", blobStorePort)
	os.Setenv("IMAGE_FACTORY_PORT", "8080")
	defer os.Setenv("BLOB_STORE_HOST", originalBlobStoreHost)
	defer os.Setenv("BLOB_STORE_PORT", originalBlobStorePort)
	httpmock.Activate()
	c := NewBlobStoreConnector()
	c.Client.Transport = httpmock.DefaultTransport
	defer httpmock.DeactivateAndReset()

	Convey("Test GetBlob", t, func() {
		Convey("Should return proper response", func() {
			httpmock.RegisterResponder("GET", GetBlobStoreAddress()+blobsPath+blobId,
				httpmock.NewStringResponder(200, blob))
			res, err := c.GetImageBlob(imageId)
			So(err, ShouldBeNil)
			So(string(res), ShouldEqual, blob)
		})

		Convey("Should return not found", func() {
			httpmock.RegisterResponder("GET", GetBlobStoreAddress()+blobsPath+blobId,
				httpmock.NewStringResponder(404, ""))
			res, err := c.GetImageBlob(imageId)
			So(err.Error(), ShouldEqual, "Invalid status: 404 Response: ")
			So(string(res), ShouldEqual, "")
		})
	})
}

func TestDeleteBlob(t *testing.T) {
	originalBlobStoreHost := os.Getenv("BLOB_STORE_HOST")
	originalBlobStorePort := os.Getenv("BLOB_STORE_PORT")
	os.Setenv("BLOB_STORE_HOST", blobStoreHost)
	os.Setenv("BLOB_STORE_PORT", blobStorePort)
	os.Setenv("IMAGE_FACTORY_PORT", "8080")
	defer os.Setenv("BLOB_STORE_HOST", originalBlobStoreHost)
	defer os.Setenv("BLOB_STORE_PORT", originalBlobStorePort)
	httpmock.Activate()
	c := NewBlobStoreConnector()
	c.Client.Transport = httpmock.DefaultTransport
	defer httpmock.DeactivateAndReset()

	Convey("Test DeleteBlob", t, func() {
		Convey("Should return proper response", func() {
			httpmock.RegisterResponder("DELETE", GetBlobStoreAddress()+blobsPath+blobId,
				httpmock.NewStringResponder(204, ""))
			err := c.DeleteImageBlob(imageId)
			So(err, ShouldBeNil)
		})

		Convey("Should return not found", func() {
			httpmock.RegisterResponder("DELETE", GetBlobStoreAddress()+blobsPath+blobId,
				httpmock.NewStringResponder(404, ""))
			err := c.DeleteImageBlob(imageId)
			So(err.Error(), ShouldEqual, "Invalid status: 404 Response: ")
		})
	})
}
