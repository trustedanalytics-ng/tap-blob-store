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

package miniowrapper

import (
	"errors"
	"io"
	"os"

	"github.com/minio/minio-go"

	"github.com/trustedanalytics/tapng-go-common/logger"
)

type ReducedMinioClient interface {
	BucketExists(bucketName string) error
	MakeBucket(bucketName string, location string) error
	StatObject(bucketName, objectName string) (minio.ObjectInfo, error)
	PutObject(bucketName, objectName string, reader io.Reader, contentType string) (n int64, err error)
	GetObject(bucketName, objectName string) (*minio.Object, error)
	RemoveObject(bucketName, objectName string) error
}

type Wrapper struct {
	Mc         ReducedMinioClient
	BucketName string
}

var (
	logger          = logger_wrapper.InitLogger("minio")
	endpoint        = os.Getenv("MINIO_HOST") + ":" + os.Getenv("MINIO_PORT")
	accessKeyID     = os.Getenv("MINIO_ACCESS_KEY")
	secretAccessKey = os.Getenv("MINIO_SECRET_KEY")
)

const (
	ssl = false
)

const (
	ErrMsgBucketNotExist = "The specified bucket does not exist."
	ErrMsgKeyNotExist    = "The specified key does not exist."
)

var (
	ErrKeyAlreadyInUse = errors.New("The specified key already exists.")
)

func CreateWrappedMinio(bucketName string) (wrap *Wrapper, err error) {
	logger.Info("Starting connection Minio Client to server:", endpoint)
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, ssl)
	if err != nil {
		return nil, err
	}

	wrap = &Wrapper{ReducedMinioClient(*minioClient), bucketName}

	err = wrap.InitMinio()
	if err != nil {
		return nil, err
	}

	logger.Info("Successfully started Minio Client.")
	return wrap, nil
}

func (mw *Wrapper) InitMinio() (err error) {
	logger.Info("Starting Minio Client in bucket:", mw.BucketName)
	err = mw.Mc.BucketExists(mw.BucketName)
	if err != nil {
		switch err.Error() {
		case ErrMsgBucketNotExist:
			err = mw.Mc.MakeBucket(mw.BucketName, "")
			if err != nil {
				return err
			}
			logger.Info("Successfully created bucket:", mw.BucketName)

		default:
			return err
		}
	}
	return nil
}

func (mw *Wrapper) StoreInMinio(objectName string, object io.Reader) (err error) {
	logger.Info("Trying to store object:", objectName, "- in bucket:", mw.BucketName)

	_, err = mw.Mc.StatObject(mw.BucketName, objectName)
	if err == nil {
		return ErrKeyAlreadyInUse
	}
	if err.Error() != ErrMsgKeyNotExist {
		return err
	}

	_, err = mw.Mc.PutObject(mw.BucketName, objectName, object, "application/octet-stream")
	if err != nil {
		return err
	}

	logger.Info("Object successfully stored:", objectName)
	return nil
}

func (mw *Wrapper) RemoveFromMinio(objectName string) (err error) {
	logger.Info("Trying to remove object:", objectName, "- in bucket:", mw.BucketName)

	_, err = mw.Mc.StatObject(mw.BucketName, objectName)
	if err != nil {
		return err
	}

	err = mw.Mc.RemoveObject(mw.BucketName, objectName)
	if err != nil {
		return err
	}

	logger.Info("Object successfully removed:", objectName)
	return nil
}

func (mw *Wrapper) RetrieveFromMinio(objectName string) (blob *minio.Object, err error) {
	logger.Info("Trying to retrieve object:", objectName, "- from bucket:", mw.BucketName)

	blob, err = mw.Mc.GetObject(mw.BucketName, objectName)
	if err != nil {
		return nil, err
	}

	logger.Info("Object successfully retrieved:", objectName)
	return blob, nil
}
