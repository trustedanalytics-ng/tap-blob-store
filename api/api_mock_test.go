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
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/minio/minio-go"
	"github.com/stretchr/testify/mock"

	"github.com/trustedanalytics-ng/tap-blob-store/minio-wrapper"
)

var (
	ErrUnhandledException = errors.New("Unhandled Error")
)

type MinioClientMock struct {
	mock.Mock
}

func (m *MinioClientMock) MakeBucket(bucketName, regionName string) error {
	return ErrUnhandledException
}

func (m *MinioClientMock) BucketExists(bucketName string) error {
	return ErrUnhandledException
}

func (m *MinioClientMock) StatObject(bucketName, objectName string) (minio.ObjectInfo, error) {
	mo := minio.ObjectInfo{}

	switch objectName {
	case NewBlobID:
		return mo, errors.New(ErrMsgKeyNotExist)
	case ExistedBlobID:
		return mo, nil
	default:
		return mo, ErrUnhandledException
	}
}

func (m *MinioClientMock) PutObject(bucketName, objectName string, reader io.Reader, contentType string) (n int64, err error) {
	switch objectName {
	case ExistedBlobID:
		return 0, miniowrapper.ErrKeyAlreadyInUse
	case NewBlobID:
		return 0, nil
	default:
		return 0, ErrUnhandledException
	}
}

func (m *MinioClientMock) GetObject(bucketName, objectName string) (*minio.Object, error) {
	switch objectName {
	case NilBlobID:
		return nil, nil
	case NewBlobID:
		return nil, errors.New(ErrMsgKeyNotExist)
	case ExistedBlobID:
		return &minio.Object{}, nil
	default:
		return nil, ErrUnhandledException
	}
}

func (m *MinioClientMock) RemoveObject(bucketName, objectName string) error {
	switch objectName {
	case NewBlobID:
		return errors.New(ErrMsgKeyNotExist)
	case ExistedBlobID:
		return nil
	default:
		return ErrUnhandledException
	}
}

func mockBlobStat(blob *minio.Object) (minio.ObjectInfo, error) {
	if blob == nil {
		return minio.ObjectInfo{}, errors.New("Object is nil")
	}
	return minio.ObjectInfo{}, nil
}

func mockBlobSeek(blob *minio.Object, offset int64, whence int) (n int64, err error) {
	return 0, nil
}

func mockBlobServe(w http.ResponseWriter, req *http.Request, name string, modtime time.Time, content io.ReadSeeker) {
}
