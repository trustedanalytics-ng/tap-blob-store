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
	"io"
	"errors"

	"github.com/stretchr/testify/mock"
	"github.com/minio/minio-go"
)


type MinioClientMock struct {
	mock.Mock
}

func (m *MinioClientMock) MakeBucket(bucketName, regionName string) error {
	return nil
}

func (m *MinioClientMock) BucketExists(bucketName string) error {
	return nil
}

func (m *MinioClientMock) StatObject(bucketName, objectName string) (minio.ObjectInfo, error) {
	return minio.ObjectInfo{}, nil
}

func (m *MinioClientMock) PutObject(bucketName, objectName string, reader io.Reader, contentType string) (n int64, err error) {
	return 0, nil
}

func (m *MinioClientMock) GetObject(bucketName, objectName string) (*minio.Object, error) {
	switch objectName {
	case NilBlobId:
		return nil, nil
	case BlobId:
		return &minio.Object{}, nil
	default:
		return nil, errors.New(ErrMsgKeyNotExist)
	}


	return &minio.Object{}, nil
}

func (m *MinioClientMock) RemoveObject(bucketName, objectName string) error {
	if(objectName != BlobId) {
		return errors.New(ErrMsgKeyNotExist)
	}
	return nil
}

